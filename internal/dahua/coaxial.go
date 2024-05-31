package dahua

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
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
	return sutureext.SanitizeError(ctx, w.serve(ctx))
}

func (w CoaxialWorker) serve(ctx context.Context) error {
	client, err := w.store.GetClient(ctx, w.conn)
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

	publish := func(v DeviceCoaxialStatus) {
		bus.Publish(bus.CoaxialStatusUpdated{
			DeviceKey:  w.conn.Key,
			Channel:    channel,
			WhiteLight: v.WhiteLight,
			Speaker:    v.Speaker,
		})
	}

	// Get and publish initial coaxial status
	coaxialStatus, err := GetCoaxialStatus(ctx, client.RPC, channel)
	if err != nil {
		return err
	}
	publish(coaxialStatus)

	t := time.NewTicker(1 * time.Second)

	// Get and send coaxial status if it changes on an interval
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}

		v, err := GetCoaxialStatus(ctx, client.RPC, channel)
		if err != nil {
			return err
		}
		if coaxialStatus.Speaker == v.Speaker && coaxialStatus.WhiteLight == v.WhiteLight {
			continue
		}
		coaxialStatus = v

		publish(v)
	}
}
