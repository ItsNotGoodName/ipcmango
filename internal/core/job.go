package core

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"

	"github.com/maragudk/goqite"
	"github.com/maragudk/goqite/jobs"
)

func NewJobService(runner *jobs.Runner) JobService {
	return JobService{
		runner: runner,
	}
}

type JobService struct {
	runner *jobs.Runner
}

func (s JobService) Serve(ctx context.Context) error {
	s.runner.Start(ctx)
	return nil
}

func NewJob[T any](queue *goqite.Queue, runner *jobs.Runner, job func(ctx context.Context, data T) error) Job[T] {
	action := fmt.Sprintf("%T", *new(T))
	runner.Register(action, func(ctx context.Context, m []byte) error {
		var data T
		if err := gob.NewDecoder(bytes.NewReader(m)).Decode(&data); err != nil {
			return err
		}
		return job(ctx, data)
	})
	return Job[T]{
		action: action,
		queue:  queue,
	}
}

type Job[T any] struct {
	action string
	queue  *goqite.Queue
}

func (j Job[T]) Create(ctx context.Context, data T) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(data); err != nil {
		return err
	}
	return jobs.Create(ctx, j.queue, j.action, buffer.Bytes())
}

func (j Job[T]) CreateTx(ctx context.Context, tx *sql.Tx, data T) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(data); err != nil {
		return err
	}
	return jobs.CreateTx(ctx, tx, j.queue, j.action, buffer.Bytes())
}
