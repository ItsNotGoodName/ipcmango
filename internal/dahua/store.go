package dahua

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
)

func NewStore() *Store {
	return &Store{
		mu:             sync.Mutex{},
		clients:        make(map[string]Client),
		deletedClients: []string{},
	}
}

// Store handles creating and caching clients.
type Store struct {
	mu             sync.Mutex
	clients        map[string]Client
	deletedClients []string
}

func (*Store) String() string {
	return "dahua.Store"
}

// Close closes all clients.
func (s *Store) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

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

// GetClient retrieves a cached client or creates a new one.
func (s *Store) GetClient(ctx context.Context, conn Conn) (Client, error) {
	s.mu.Lock()
	if slices.Contains(s.deletedClients, conn.Key.UUID) {
		s.mu.Unlock()
		return Client{}, fmt.Errorf("client deleted")
	}

	client, ok := s.clients[conn.Key.UUID]
	if !ok {
		// Not found

		client = NewClient(conn)
		s.clients[conn.Key.UUID] = client
	} else if !client.Conn.EQ(conn) && client.Conn.UpdatedAt.Before(conn.UpdatedAt) {
		// Found but not equal and old

		err := client.CloseNoWait(ctx)
		if err != nil {
			slog.Error("Failed to close Client connection", slog.String("uuid", client.Conn.Key.UUID))
		}

		client = NewClient(conn)
		s.clients[conn.Key.UUID] = client
	}
	s.mu.Unlock()

	return client, nil
}

// DeleteClient removes a client and prevents it from being created again.
func (s *Store) DeleteClient(ctx context.Context, deviceUUID string) error {
	s.mu.Lock()
	client, found := s.clients[deviceUUID]
	if found {
		delete(s.clients, deviceUUID)
	}
	s.deletedClients = append(s.deletedClients, deviceUUID)
	s.mu.Unlock()

	if !found {
		return nil
	}

	return client.Close(ctx)
}

func (s *Store) Register() *Store {
	bus.Subscribe(s.String(), func(ctx context.Context, event bus.DeviceDeleted) error {
		return s.DeleteClient(ctx, event.DeviceKey.UUID)
	})
	return s
}
