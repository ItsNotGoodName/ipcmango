package dahua

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/jmoiron/sqlx"
	"github.com/thejerf/suture/v4"
)

func NewEventManager(super *suture.Supervisor, db *sqlx.DB) *EventManager {
	return &EventManager{
		super:      super,
		db:         db,
		servicesMu: sync.Mutex{},
		services:   make(map[string]eventService),
	}
}

type eventService struct {
	Token  suture.ServiceToken
	Worker EventWorker
}

type EventManager struct {
	super *suture.Supervisor
	db    *sqlx.DB

	servicesMu sync.Mutex
	services   map[string]eventService
}

func (m *EventManager) String() string {
	return "dahua.EventManager"
}

func (m *EventManager) Serve(ctx context.Context) error {
	slog.Info("Started service", slog.String("service", m.String()))

	if err := m.Start(); err != nil {
		return err
	}

	<-ctx.Done()
	m.Close()
	return ctx.Err()
}

func (m *EventManager) Close() {
	m.servicesMu.Lock()
	for _, service := range m.services {
		m.super.Remove(service.Token)
	}
	clear(m.services)
	m.servicesMu.Unlock()
}

func (m *EventManager) Start() error {
	m.servicesMu.Lock()
	defer m.servicesMu.Unlock()

	rows, err := m.db.QueryxContext(context.Background(), `
		SELECT * FROM dahua_devices
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var device Device
		if err := rows.StructScan(&device); err != nil {
			return err
		}

		worker := NewEventWorker(NewConn(device), m.db)
		token := m.super.Add(worker)
		m.services[device.UUID] = eventService{
			Token:  token,
			Worker: worker,
		}
	}

	return nil
}

func (m *EventManager) Refresh(ctx context.Context, uuid string) error {
	m.servicesMu.Lock()
	defer m.servicesMu.Unlock()

	service, ok := m.services[uuid]
	if ok {
		m.super.Remove(service.Token)
	}

	var device Device
	err := m.db.GetContext(ctx, &device, `
		SELECT * FROM dahua_devices WHERE uuid = ?
	`, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	worker := NewEventWorker(NewConn(device), m.db)
	token := m.super.Add(worker)
	m.services[device.UUID] = eventService{
		Token:  token,
		Worker: worker,
	}

	return nil
}

func (m *EventManager) Register() *EventManager {
	bus.Subscribe(m.String(), func(ctx context.Context, event bus.DeviceCreated) error {
		return m.Refresh(ctx, event.DeviceKey.UUID)
	})
	bus.Subscribe(m.String(), func(ctx context.Context, event bus.DeviceUpdated) error {
		return m.Refresh(ctx, event.DeviceKey.UUID)
	})
	bus.Subscribe(m.String(), func(ctx context.Context, event bus.DeviceDeleted) error {
		return m.Refresh(ctx, event.DeviceKey.UUID)
	})
	return m
}

func NewEventWorker(conn Conn, db *sqlx.DB) EventWorker {
	return EventWorker{
		conn: conn,
		db:   db,
	}
}

type EventWorker struct {
	conn Conn
	db   *sqlx.DB
}

func (w EventWorker) String() string {
	return fmt.Sprintf("dahua.EventWorker(name=%s)", w.conn.Name)
}

func (w EventWorker) Serve(ctx context.Context) error {
	return sutureext.SanitizeError(ctx, w.serve(ctx))
}

func (w EventWorker) serve(ctx context.Context) error {
	slog.Info("Started service", slog.String("service", w.String()))
	defer slog.Info("Stopped service", slog.String("service", w.String()))

	c := dahuacgi.NewClient(http.Client{}, w.conn.URL, w.conn.Username, w.conn.Password)

	manager, err := dahuacgi.EventManagerGet(ctx, c, 0)
	if err != nil {
		var httpErr dahuacgi.HTTPError
		if errors.As(err, &httpErr) && slices.Contains([]int{
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
		}, httpErr.StatusCode) {
			slog.Error("Failed to get EventManager", slog.Any("error", err), slog.String("service", w.String()))
			return errors.Join(suture.ErrDoNotRestart, err)
		}

		return err
	}
	defer manager.Close()

	for reader := manager.Reader(); ; {
		if err := reader.Poll(); err != nil {
			return err
		}

		event, err := reader.ReadEvent()
		if err != nil {
			return err
		}

		if err := HandleEvent(ctx, w.db, w.conn.Key, event); err != nil {
			return err
		}
	}
}
