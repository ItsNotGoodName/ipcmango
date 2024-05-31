package dahua

import (
	"context"
	"log/slog"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/jmoiron/sqlx"
	"github.com/thejerf/suture/v4"
)

func NewDefaultServiceFactory(db *sqlx.DB, dahuaStore *Store) DefaultServiceFactory {
	return DefaultServiceFactory{
		db:         db,
		dahuaStore: dahuaStore,
	}
}

type DefaultServiceFactory struct {
	db         *sqlx.DB
	dahuaStore *Store
}

func (d DefaultServiceFactory) build(device Device) []suture.Service {
	conn := NewConn(device)
	return []suture.Service{
		NewEventWorker(conn, d.db),
		NewCoaxialWorker(conn, d.db, d.dahuaStore),
	}
}

func (d DefaultServiceFactory) List(ctx context.Context) ([]Service, error) {
	var devices []Device
	err := d.db.SelectContext(ctx, &devices, `
		SELECT * FROM dahua_devices
	`)
	if err != nil {
		return nil, err
	}

	var services []Service
	for _, device := range devices {
		for _, service := range d.build(device) {
			services = append(services, Service{
				DeviceKey: device.Key,
				Service:   service,
			})
		}
	}

	return services, nil
}

func (d DefaultServiceFactory) One(ctx context.Context, deviceKey types.Key) ([]Service, error) {
	var devices []Device
	err := d.db.SelectContext(ctx, &devices, `
		SELECT * FROM dahua_devices WHERE id = ?
	`, deviceKey.ID)
	if err != nil {
		return nil, err
	}

	var services []Service
	for _, device := range devices {
		for _, service := range d.build(device) {
			services = append(services, Service{
				DeviceKey: device.Key,
				Service:   service,
			})
		}
	}

	return services, nil
}

type Service struct {
	DeviceKey types.Key
	Service   suture.Service
}

type ServiceFactory interface {
	List(ctx context.Context) ([]Service, error)
	One(ctx context.Context, deviceKey types.Key) ([]Service, error)
}

func NewServiceManager(super *suture.Supervisor, serviceFactory ServiceFactory) *ManagerManager {
	return &ManagerManager{
		super:          super,
		serviceFactory: serviceFactory,
		servicesMu:     sync.Mutex{},
		services:       make(map[string][]suture.ServiceToken),
	}
}

type ManagerManager struct {
	super          *suture.Supervisor
	serviceFactory ServiceFactory

	servicesMu sync.Mutex
	services   map[string][]suture.ServiceToken
}

func (m *ManagerManager) String() string {
	return "dahua.EventManager"
}

func (m *ManagerManager) Serve(ctx context.Context) error {
	slog.Info("Started service", slog.String("service", m.String()))

	if err := m.Start(ctx); err != nil {
		return err
	}

	<-ctx.Done()
	m.Close()
	return ctx.Err()
}

func (m *ManagerManager) Close() {
	m.servicesMu.Lock()
	for _, tokens := range m.services {
		for _, token := range tokens {
			m.super.Remove(token)
		}
	}
	clear(m.services)
	m.servicesMu.Unlock()
}

func (m *ManagerManager) Start(ctx context.Context) error {
	m.servicesMu.Lock()
	defer m.servicesMu.Unlock()

	services, err := m.serviceFactory.List(ctx)
	if err != nil {
		return err
	}

	for _, service := range services {
		m.services[service.DeviceKey.UUID] = append(m.services[service.DeviceKey.UUID], m.super.Add(service.Service))
	}

	return nil
}

func (m *ManagerManager) Refresh(ctx context.Context, deviceKey types.Key) error {
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
		m.services[service.DeviceKey.UUID] = append(m.services[service.DeviceKey.UUID], m.super.Add(service.Service))
	}

	return nil
}

func (m *ManagerManager) Register() *ManagerManager {
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
