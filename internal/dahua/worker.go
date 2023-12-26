package dahua

import (
	"context"
	"fmt"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/thejerf/suture/v4"
)

func DefaultWorkerBuilder(hooks DefaultEventHooks, bus *core.Bus, store *Store, db repo.DB) WorkerBuilder {
	return func(ctx context.Context, super *suture.Supervisor, device models.DahuaConn) ([]suture.ServiceToken, error) {
		var tokens []suture.ServiceToken

		{
			worker := NewEventWorker(device, hooks)
			token := super.Add(worker)
			tokens = append(tokens, token)
		}

		{
			worker := NewCoaxialWorker(bus, device.ID, store, db)
			token := super.Add(worker)
			tokens = append(tokens, token)
		}

		return tokens, nil
	}
}

type WorkerBuilder = func(ctx context.Context, super *suture.Supervisor, device models.DahuaConn) ([]suture.ServiceToken, error)

type WorkerStore struct {
	super *suture.Supervisor
	build WorkerBuilder

	workersMu sync.Mutex
	workers   map[int64]workerData
}

type workerData struct {
	device models.DahuaConn
	tokens []suture.ServiceToken
}

func NewWorkerStore(super *suture.Supervisor, build WorkerBuilder) *WorkerStore {
	return &WorkerStore{
		super:     super,
		build:     build,
		workersMu: sync.Mutex{},
		workers:   make(map[int64]workerData),
	}
}

func (s *WorkerStore) create(ctx context.Context, device models.DahuaConn) error {
	tokens, err := s.build(ctx, s.super, device)
	if err != nil {
		return err
	}

	s.workers[device.ID] = workerData{
		device: device,
		tokens: tokens,
	}

	return nil
}

func (s *WorkerStore) Create(ctx context.Context, device models.DahuaConn) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	_, found := s.workers[device.ID]
	if found {
		return fmt.Errorf("workers already exists for device by ID: %d", device.ID)
	}

	return s.create(ctx, device)
}

func (s *WorkerStore) Update(ctx context.Context, device models.DahuaConn) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	worker, found := s.workers[device.ID]
	if !found {
		return fmt.Errorf("workers not found for device by ID: %d", device.ID)
	}
	if ConnEqual(worker.device, device) {
		return nil
	}

	for _, st := range worker.tokens {
		s.super.Remove(st)
	}

	return s.create(ctx, device)
}

func (s *WorkerStore) Delete(id int64) error {
	s.workersMu.Lock()
	worker, found := s.workers[id]
	if !found {
		s.workersMu.Unlock()
		return fmt.Errorf("workers not found for device by ID: %d", id)
	}

	for _, token := range worker.tokens {
		s.super.Remove(token)
	}
	delete(s.workers, id)
	s.workersMu.Unlock()
	return nil
}

func (s *WorkerStore) Register(bus *core.Bus) *WorkerStore {
	bus.OnEventDahuaDeviceCreated(func(ctx context.Context, evt models.EventDahuaDeviceCreated) error {
		return s.Create(ctx, evt.Device.DahuaConn)
	})
	bus.OnEventDahuaDeviceUpdated(func(ctx context.Context, evt models.EventDahuaDeviceUpdated) error {
		return s.Update(ctx, evt.Device.DahuaConn)
	})
	bus.OnEventDahuaDeviceDeleted(func(ctx context.Context, evt models.EventDahuaDeviceDeleted) error {
		return s.Delete(evt.DeviceID)
	})
	return s
}

func (s *WorkerStore) Bootstrap(ctx context.Context, db repo.DB, store *Store) error {
	dbDevices, err := db.ListDahuaDevice(ctx)
	if err != nil {
		return err
	}
	var conns []models.DahuaConn
	for _, v := range dbDevices {
		conns = append(conns, v.Convert().DahuaConn)
	}

	clients := store.ClientList(ctx, conns)
	for _, conn := range clients {
		if err := s.Create(ctx, conn.Conn); err != nil {
			return err
		}
	}

	return err
}
