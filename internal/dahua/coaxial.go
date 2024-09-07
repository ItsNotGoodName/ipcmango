package dahua

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/jmoiron/sqlx"
	"github.com/thejerf/suture/v4"
)

func NewCoaxialWorker(conn Conn, db *sqlx.DB, store *Store) CoaxialWorker {
	return CoaxialWorker{
		conn:  conn,
		db:    db,
		store: store,
	}
}

type CoaxialWorker struct {
	conn  Conn
	db    *sqlx.DB
	store *Store
}

func (w CoaxialWorker) String() string {
	return fmt.Sprintf("dahua.CoaxialWorker(name=%s)", w.conn.Name)
}

func (w CoaxialWorker) Serve(ctx context.Context) error {
	client, err := w.store.GetClient(ctx, w.conn.Key)
	if err != nil {
		return err
	}

	channel := 1

	// Check if device supports coaxial
	caps, err := GetCoaxialCaps(ctx, client.RPC, channel)
	if err != nil {
		return err
	}
	if !(caps.SupportControlSpeaker || caps.SupportControlLight || caps.SupportControlFullcolorLight) {
		return suture.ErrDoNotRestart
	}

	slog.Info("Started service", slog.String("service", w.String()))

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	// Publish initial status
	lastStatus, err := GetCoaxialStatus(ctx, client.RPC, channel)
	if err != nil {
		return err
	}
	HandleCoaxialStatusChange(ctx, w.conn.Key, channel, lastStatus)

	// Get and send status if it changes on an interval
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			status, err := GetCoaxialStatus(ctx, client.RPC, channel)
			if err != nil {
				return err
			}
			if lastStatus.Speaker == status.Speaker && lastStatus.WhiteLight == status.WhiteLight {
				continue
			}
			lastStatus = status

			HandleCoaxialStatusChange(ctx, w.conn.Key, channel, status)
		}
	}
}

func HandleCoaxialStatusChange(ctx context.Context, deviceKey types.Key, channel int, coaxialStatus DeviceCoaxialStatus) {
	bus.Publish(bus.CoaxialStatusUpdated{
		DeviceKey:  deviceKey,
		Channel:    channel,
		WhiteLight: coaxialStatus.WhiteLight,
		Speaker:    coaxialStatus.Speaker,
	})
}
