package dahua

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/jmoiron/sqlx"
)

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
		SELECT id, uuid, name, ip, username, password, updated_at
		FROM dahua_devices WHERE id = ? OR uuid = ? LIMIT 1
	`, deviceKey.ID, deviceKey.UUID)
	if err != nil {
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

		err := client.CloseNoWait(ctx)
		if err != nil {
			slog.Error("Failed to close Client connection", slog.String("uuid", client.Conn.Key.UUID))
		}

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
