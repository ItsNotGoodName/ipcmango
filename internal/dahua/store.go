package dahua

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/ptz"
	"github.com/jmoiron/sqlx"
	"github.com/k0kubun/pp/v3"
)

func NewStoreClient(conn Conn) StoreClient {
	urL, _ := url.Parse("http://" + conn.IP)

	httpClient := http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 5*time.Second)
			},
		},
	}

	clientRPC := dahuarpc.NewClient(&httpClient, urL, conn.Username, conn.Password)
	clientPTZ := ptz.NewClient(clientRPC)
	clientCGI := dahuacgi.NewClient(httpClient, urL, conn.Username, conn.Password)
	clientFile := dahuarpc.NewFileClient(&httpClient, 10)

	var ref int32 = 1

	return StoreClient{
		ref:  &ref,
		Conn: conn,
		URL:  urL,
		RPC:  clientRPC,
		PTZ:  clientPTZ,
		CGI:  clientCGI,
		File: clientFile,
	}
}

type StoreClient struct {
	ref  *int32
	Conn Conn
	URL  *url.URL
	RPC  dahuarpc.Client
	PTZ  ptz.Client
	CGI  dahuacgi.Client
	File dahuarpc.FileClient
}

func (c StoreClient) close(ctx context.Context) error {
	c.File.Close()
	return c.RPC.Close(ctx)
}

func (c StoreClient) Release() error {
	return c.ReleaseContext(context.Background())
}

func (c StoreClient) ReleaseContext(ctx context.Context) error {
	counter := atomic.AddInt32(c.ref, -1)
	if counter == 0 {
		return c.close(ctx)
	}
	return nil
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db:        db,
		clientsMu: sync.Mutex{},
		clients:   make(map[int64]Client),
	}
}

// Store handles creating and caching clients.
type Store struct {
	db *sqlx.DB

	clientsMu sync.Mutex
	clients   map[int64]Client
}

func (*Store) String() string {
	return "dahua.Store"
}

// Close closes all clients.
func (s *Store) Close() {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	wg := sync.WaitGroup{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, client := range s.clients {
		wg.Add(1)
		go func(client Client) {
			defer wg.Done()
			err := client.Close(ctx)
			if err != nil {
				slog.Error("Failed to close Client connection", slog.String("uuid", client.Conn.Key.UUID))
			}
		}(client)
	}

	wg.Wait()
}

func (s *Store) Serve(ctx context.Context) error {
	<-ctx.Done()
	s.Close()
	return ctx.Err()
}

func (s *Store) GetClient(ctx context.Context, deviceKey types.Key) (Client, error) {
	s.clientsMu.Lock()
	var conn Conn
	err := s.db.GetContext(ctx, &conn, `
		SELECT id, uuid, name, ip, username, password
		FROM dahua_devices WHERE id = ? OR uuid = ?
	`, deviceKey.ID, deviceKey.UUID)
	if err != nil {
		pp.Println(err, deviceKey)
		s.clientsMu.Unlock()
		return Client{}, err
	}

	client, ok := s.clients[conn.Key.ID]
	if !ok {
		// Not found

		client = NewClient(conn)
		s.clients[conn.Key.ID] = client
	} else if !client.Conn.EQ(conn) {
		// Found but not equal

		go func(client Client) {
			err := client.Close(context.Background())
			if err != nil {
				slog.Error("Failed to close Client connection", slog.String("name", client.Conn.Name))
			}
		}(client)

		client = NewClient(conn)
		s.clients[conn.Key.ID] = client
	}
	s.clientsMu.Unlock()

	return client, nil
}

func (s *Store) DeleteClient(ctx context.Context, deviceID int64) error {
	s.clientsMu.Lock()
	client, found := s.clients[deviceID]
	if found {
		delete(s.clients, deviceID)
	}
	s.clientsMu.Unlock()

	if !found {
		return nil
	}

	return client.Close(ctx)
}

func (s *Store) Register() *Store {
	bus.Subscribe(s.String(), func(ctx context.Context, event bus.DeviceDeleted) error {
		return s.DeleteClient(ctx, event.DeviceKey.ID)
	})
	return s
}
