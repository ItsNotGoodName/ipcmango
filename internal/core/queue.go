package core

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/maragudk/goqite"
	"github.com/maragudk/goqite/jobs"
)

func NewJobBuilder[T any](action string) JobBuilder[T] {
	return JobBuilder[T]{
		Action: action,
	}
}

type JobBuilder[T any] struct {
	Action string
	data   T
}

func (e JobBuilder[T]) Create(ctx context.Context, q *goqite.Queue, data T) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return jobs.Create(ctx, q, e.Action, b)
}

func (e JobBuilder[T]) CreateTx(ctx context.Context, tx *sql.Tx, q *goqite.Queue, data T) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return jobs.CreateTx(ctx, tx, q, e.Action, b)
}

func (e JobBuilder[T]) Register(ctx context.Context, r *jobs.Runner, job func(ctx context.Context, data T) error) {
	r.Register(e.Action, func(ctx context.Context, m []byte) error {
		var data T
		if err := json.Unmarshal(m, &data); err != nil {
			return err
		}
		return job(ctx, data)
	})
}
