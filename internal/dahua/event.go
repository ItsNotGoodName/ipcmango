package dahua

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/jmoiron/sqlx"
	"github.com/thejerf/suture/v4"
)

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

	connURL, _ := url.Parse("http://" + w.conn.IP)
	c := dahuacgi.NewClient(http.Client{}, connURL, w.conn.Username, w.conn.Password)

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
