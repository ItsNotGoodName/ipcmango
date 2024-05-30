package core

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"

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
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(data); err != nil {
		return err
	}
	return jobs.Create(ctx, q, e.Action, buffer.Bytes())
}

func (e JobBuilder[T]) CreateTx(ctx context.Context, tx *sql.Tx, q *goqite.Queue, data T) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(data); err != nil {
		return err
	}
	return jobs.CreateTx(ctx, tx, q, e.Action, buffer.Bytes())
}

func (e JobBuilder[T]) Register(ctx context.Context, r *jobs.Runner, job func(ctx context.Context, data T) error) {
	r.Register(e.Action, func(ctx context.Context, m []byte) error {
		var data T
		if err := gob.NewDecoder(bytes.NewReader(m)).Decode(&data); err != nil {
			return err
		}
		return job(ctx, data)
	})
}
