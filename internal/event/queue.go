package event

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
)

func NewQueue(db sqlite.DB, bus *Bus) Queue {
	q := Queue{
		db:    db,
		bus:   bus,
		check: make(chan struct{}, 1),
	}

	q.bus.OnEventQueued(q.String(), func(ctx context.Context, evt EventQueued) error {
		select {
		case q.check <- struct{}{}:
		default:
		}
		return nil
	})

	return q
}

type Queue struct {
	db    sqlite.DB
	bus   *Bus
	check chan struct{}
}

func (Queue) String() string {
	return "event.Queue"
}

func (q Queue) Serve(ctx context.Context) error {
	cursor, err := q.db.C().GetEventCursor(ctx)
	if err != nil && !core.IsNotFound(err) {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-q.check:
			for {
				event, err := q.db.C().NextEventByCursor(ctx, cursor)
				if err != nil {
					if core.IsNotFound(err) {
						break
					}
					return err
				}
				cursor = event.ID

				q.bus.Event(Event{
					Event: event,
				})
			}
		}
	}
}
