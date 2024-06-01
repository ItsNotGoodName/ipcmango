package dahua

import (
	"context"
	"log/slog"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/jmoiron/sqlx"
	"github.com/thejerf/suture/v4"
)

func NewDefaultServiceFactory(db *sqlx.DB, dahuaStore *Store) DefaultWorkerFactory {
	return DefaultWorkerFactory{
		db:         db,
		dahuaStore: dahuaStore,
	}
}

type DefaultWorkerFactory struct {
	db         *sqlx.DB
	dahuaStore *Store
}

func (d DefaultWorkerFactory) build(device Device) []sutureext.Service {
	conn := NewConn(device)
	return []sutureext.Service{
		NewEventWorker(conn, d.db),
		NewCoaxialWorker(conn, d.db, d.dahuaStore),
	}
}

func (f DefaultWorkerFactory) List(ctx context.Context) ([]Worker, error) {
	var devices []Device
	err := f.db.SelectContext(ctx, &devices, `
		SELECT * FROM dahua_devices
	`)
	if err != nil {
		return nil, err
	}

	var services []Worker
	for _, device := range devices {
		for _, service := range f.build(device) {
			services = append(services, Worker{
				DeviceKey: device.Key,
				Service:   service,
			})
		}
	}

	return services, nil
}

func (f DefaultWorkerFactory) One(ctx context.Context, deviceKey types.Key) ([]Worker, error) {
	var devices []Device
	err := f.db.SelectContext(ctx, &devices, `
		SELECT * FROM dahua_devices WHERE id = ?
	`, deviceKey.ID)
	if err != nil {
		return nil, err
	}

	var services []Worker
	for _, device := range devices {
		for _, service := range f.build(device) {
			services = append(services, Worker{
				DeviceKey: device.Key,
				Service:   service,
			})
		}
	}

	return services, nil
}

type Worker struct {
	DeviceKey types.Key
	Service   sutureext.Service
}

type WorkerFactory interface {
	List(ctx context.Context) ([]Worker, error)
	One(ctx context.Context, deviceKey types.Key) ([]Worker, error)
}

func NewWorkerManager(super *suture.Supervisor, serviceFactory WorkerFactory) *WorkerManager {
	return &WorkerManager{
		super:          super,
		serviceFactory: serviceFactory,
		servicesMu:     sync.Mutex{},
		services:       make(map[string][]suture.ServiceToken),
	}
}

type WorkerManager struct {
	super          *suture.Supervisor
	serviceFactory WorkerFactory

	servicesMu sync.Mutex
	services   map[string][]suture.ServiceToken
}

func (m *WorkerManager) String() string {
	return "dahua.EventManager"
}

func (m *WorkerManager) Serve(ctx context.Context) error {
	slog.Info("Started service", slog.String("service", m.String()))

	if err := m.Start(ctx); err != nil {
		return err
	}

	<-ctx.Done()
	m.Close()
	return ctx.Err()
}

func (m *WorkerManager) Close() {
	m.servicesMu.Lock()
	for _, tokens := range m.services {
		for _, token := range tokens {
			m.super.Remove(token)
		}
	}
	clear(m.services)
	m.servicesMu.Unlock()
}

func (m *WorkerManager) Start(ctx context.Context) error {
	m.servicesMu.Lock()
	defer m.servicesMu.Unlock()

	services, err := m.serviceFactory.List(ctx)
	if err != nil {
		return err
	}

	for _, service := range services {
		m.services[service.DeviceKey.UUID] = append(m.services[service.DeviceKey.UUID], sutureext.Add(m.super, service.Service))
	}

	return nil
}

func (m *WorkerManager) Refresh(ctx context.Context, deviceKey types.Key) error {
	m.servicesMu.Lock()
	defer m.servicesMu.Unlock()

	tokens, ok := m.services[deviceKey.UUID]
	if ok {
		for _, token := range tokens {
			m.super.Remove(token)
		}
	}

	services, err := m.serviceFactory.One(ctx, deviceKey)
	if err != nil {
		return err
	}
	if len(services) == 0 {
		return nil
	}

	for _, service := range services {
		m.services[service.DeviceKey.UUID] = append(m.services[service.DeviceKey.UUID], sutureext.Add(m.super, service.Service))
	}

	return nil
}

func (m *WorkerManager) Register() *WorkerManager {
	bus.Subscribe(m.String(), func(ctx context.Context, event bus.DeviceCreated) error {
		return m.Refresh(ctx, event.DeviceKey)
	})
	bus.Subscribe(m.String(), func(ctx context.Context, event bus.DeviceUpdated) error {
		return m.Refresh(ctx, event.DeviceKey)
	})
	bus.Subscribe(m.String(), func(ctx context.Context, event bus.DeviceDeleted) error {
		return m.Refresh(ctx, event.DeviceKey)
	})
	return m
}
