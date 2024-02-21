package dahua

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuaevents"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type WorkerFactory = func(ctx context.Context, super *suture.Supervisor, device Conn) ([]suture.ServiceToken, error)

// WorkerManager manages devices workers.
type WorkerManager struct {
	super   *suture.Supervisor
	factory WorkerFactory

	workersMu sync.Mutex
	workers   map[int64]workerData
}

func (*WorkerManager) String() string {
	return "dahua.WorkerManager"
}

type workerData struct {
	conn   Conn
	tokens []suture.ServiceToken
}

func NewWorkerManager(super *suture.Supervisor, factory WorkerFactory) *WorkerManager {
	return &WorkerManager{
		super:     super,
		factory:   factory,
		workersMu: sync.Mutex{},
		workers:   make(map[int64]workerData),
	}
}

func (m *WorkerManager) Upsert(ctx context.Context, conn Conn) error {
	m.workersMu.Lock()
	defer m.workersMu.Unlock()

	worker, found := m.workers[conn.ID]
	if found {
		if worker.conn.EQ(conn) {
			return nil
		}

		for _, st := range worker.tokens {
			m.super.Remove(st)
		}
	}

	tokens, err := m.factory(ctx, m.super, conn)
	if err != nil {
		return err
	}

	m.workers[conn.ID] = workerData{
		conn:   conn,
		tokens: tokens,
	}

	return nil
}

func (m *WorkerManager) Delete(id int64) error {
	m.workersMu.Lock()
	worker, found := m.workers[id]
	if !found {
		m.workersMu.Unlock()
		return nil
	}

	for _, token := range worker.tokens {
		m.super.Remove(token)
	}
	delete(m.workers, id)
	m.workersMu.Unlock()
	return nil
}

func (m *WorkerManager) Register(bus *event.Bus, db sqlite.DB) *WorkerManager {
	bus.OnEvent(m.String(), func(ctx context.Context, evt event.Event) error {
		switch evt.Event.Action {
		case event.ActionDahuaDeviceCreated, event.ActionDahuaDeviceUpdated:
			deviceID := event.DataAsInt64(evt.Event)

			conn, err := GetConn(ctx, db, deviceID)
			if err != nil {
				if repo.IsNotFound(err) {
					return m.Delete(deviceID)
				}
				return err
			}

			return m.Upsert(ctx, conn)
		case event.ActionDahuaDeviceDeleted:
			deviceID := event.DataAsInt64(evt.Event)

			return m.Delete(deviceID)
		}
		return nil
	})
	return m
}

func (m *WorkerManager) Bootstrap(ctx context.Context, db sqlite.DB, store *Store) error {
	conns, err := ListConn(ctx, db)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := m.Upsert(ctx, conn); err != nil {
			return err
		}
	}

	return err
}

func DefaultWorkerFactory(bus *event.Bus, pub *pubsub.Pub, db sqlite.DB, store *Store, scanLockStore ScanLockStore, hooks DefaultEventHooks) WorkerFactory {
	return func(ctx context.Context, super *suture.Supervisor, device Conn) ([]suture.ServiceToken, error) {
		var tokens []suture.ServiceToken

		{
			worker := NewEventWorker(device, hooks)
			token := super.Add(worker)
			tokens = append(tokens, token)
		}

		{
			worker := NewCoaxialWorker(bus, db, store, device.ID)
			token := super.Add(worker)
			tokens = append(tokens, token)
		}

		{
			worker := NewQuickScanWorker(pub, db, store, scanLockStore, device.ID)
			token := super.Add(worker)
			tokens = append(tokens, token)
		}

		return tokens, nil
	}
}

type EventHooks interface {
	Connecting(ctx context.Context, deviceID int64)
	Connect(ctx context.Context, deviceID int64)
	Disconnect(ctx context.Context, deviceID int64, err error)
	Event(ctx context.Context, deviceID int64, event dahuacgi.Event)
}

func NewEventWorker(device Conn, hooks EventHooks) EventWorker {
	return EventWorker{
		device: device,
		hooks:  hooks,
	}
}

// EventWorker subscribes to events.
type EventWorker struct {
	device Conn
	hooks  EventHooks
}

func (w EventWorker) String() string {
	return fmt.Sprintf("dahua.EventWorker(id=%d)", w.device.ID)
}

func (w EventWorker) Serve(ctx context.Context) error {
	w.hooks.Connecting(ctx, w.device.ID)
	err := w.serve(ctx)
	w.hooks.Disconnect(context.Background(), w.device.ID, err)
	return sutureext.SanitizeError(ctx, err)
}

func (w EventWorker) serve(ctx context.Context) error {
	c := dahuacgi.NewClient(http.Client{}, w.device.URL, w.device.Username, w.device.Password)

	manager, err := dahuacgi.EventManagerGet(ctx, c, 0)
	if err != nil {
		var httpErr dahuacgi.HTTPError
		if errors.As(err, &httpErr) && slices.Contains([]int{
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
		}, httpErr.StatusCode) {
			log.Err(err).Str("service", w.String()).Msg("Failed to get EventManager")
			return errors.Join(suture.ErrDoNotRestart, err)
		}

		return err
	}
	defer manager.Close()

	w.hooks.Connect(ctx, w.device.ID)

	for reader := manager.Reader(); ; {
		if err := reader.Poll(); err != nil {
			return err
		}

		rawEvent, err := reader.ReadEvent()
		if err != nil {
			return err
		}

		w.hooks.Event(ctx, w.device.ID, rawEvent)
	}
}

func NewCoaxialWorker(bus *event.Bus, db sqlite.DB, store *Store, deviceID int64) CoaxialWorker {
	return CoaxialWorker{
		bus:      bus,
		db:       db,
		store:    store,
		deviceID: deviceID,
	}
}

// CoaxialWorker publishes coaxial status to the bus.
type CoaxialWorker struct {
	bus      *event.Bus
	db       sqlite.DB
	store    *Store
	deviceID int64
}

func (w CoaxialWorker) String() string {
	return fmt.Sprintf("dahua.CoaxialWorker(id=%d)", w.deviceID)
}

func (w CoaxialWorker) Serve(ctx context.Context) error {
	return sutureext.SanitizeError(ctx, w.serve(ctx))
}

func (w CoaxialWorker) serve(ctx context.Context) error {
	client, err := w.store.GetClient(ctx, w.deviceID)
	if err != nil {
		if repo.IsNotFound(err) {
			return suture.ErrDoNotRestart
		}
		return err
	}

	channel := 1

	// Does this device support coaxial?
	caps, err := GetCoaxialCaps(ctx, client.RPC, channel)
	if err != nil {
		return err
	}
	if !(caps.SupportControlSpeaker || caps.SupportControlLight || caps.SupportControlFullcolorLight) {
		return suture.ErrDoNotRestart
	}

	// Get and send initial coaxial status
	coaxialStatus, err := GetCoaxialStatus(ctx, client.RPC, channel)
	if err != nil {
		return err
	}
	w.bus.DahuaCoaxialStatus(event.DahuaCoaxialStatus{
		DeviceID:      w.deviceID,
		Channel:       channel,
		CoaxialStatus: coaxialStatus,
	})

	t := time.NewTicker(1 * time.Second)

	// Get and send coaxial status if it changes on an interval
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}

		s, err := GetCoaxialStatus(ctx, client.RPC, channel)
		if err != nil {
			return err
		}
		if coaxialStatus.Speaker == s.Speaker && coaxialStatus.WhiteLight == s.WhiteLight {
			continue
		}
		coaxialStatus = s

		w.bus.DahuaCoaxialStatus(event.DahuaCoaxialStatus{
			DeviceID:      w.deviceID,
			Channel:       channel,
			CoaxialStatus: coaxialStatus,
		})
	}
}

func NewQuickScanWorker(pub *pubsub.Pub, db sqlite.DB, store *Store, scanLockStore ScanLockStore, deviceID int64) QuickScanWorker {
	return QuickScanWorker{
		pub:           pub,
		db:            db,
		store:         store,
		scanLockStore: scanLockStore,
		deviceID:      deviceID,
	}
}

type QuickScanWorker struct {
	pub           *pubsub.Pub
	db            sqlite.DB
	store         *Store
	scanLockStore ScanLockStore
	deviceID      int64
}

func (w QuickScanWorker) String() string {
	return fmt.Sprintf("dahua.QuickScanWorker(id=%d)", w.deviceID)
}

func (w QuickScanWorker) Serve(ctx context.Context) error {
	return sutureext.SanitizeError(ctx, w.serve(ctx))
}

func (w QuickScanWorker) serve(ctx context.Context) error {
	quickScanC := make(chan struct{}, 1)

	sub, err := w.pub.
		Subscribe(event.DahuaEvent{}).
		Function(ctx, func(ctx context.Context, evt pubsub.Event) error {
			switch e := evt.(type) {
			// case event.DahuaQuickScanQueue:
			// 	if !(e.DeviceID == 0 || e.DeviceID == w.deviceID) {
			// 		return nil
			// 	}
			case event.DahuaEvent:
				if e.Event.DeviceID != w.deviceID || e.Event.Code != dahuaevents.CodeNewFile {
					return nil
				}
			default:
				return nil
			}

			select {
			case quickScanC <- struct{}{}:
			default:
			}

			return nil
		})
	if err != nil {
		return err
	}
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-quickScanC:
			if err := w.scan(ctx); err != nil {
				return err
			}
		}
	}
}

func (w QuickScanWorker) scan(ctx context.Context) error {
	unlock, err := w.scanLockStore.Lock(ctx, w.deviceID)
	if err != nil {
		return err
	}
	defer unlock()

	client, err := w.store.GetClient(ctx, w.deviceID)
	if err != nil {
		return err
	}

	return Scan(ctx, w.db, client.RPC, client.Conn, models.DahuaScanTypeQuick)
}
