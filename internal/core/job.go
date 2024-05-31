package core

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/pkg/jobs"
	"github.com/maragudk/goqite"
)

func NewJobClient(queue *goqite.Queue, runner *jobs.Runner) JobClient {
	return JobClient{
		queue:  queue,
		runner: runner,
	}
}

type JobClient struct {
	queue  *goqite.Queue
	runner *jobs.Runner
}

func (s JobClient) Serve(ctx context.Context) error {
	s.runner.Start(ctx)
	return nil
}

func NewJob[T any](client JobClient, job func(ctx context.Context, data T) error) Job[T] {
	action := fmt.Sprintf("%T", *new(T))
	client.runner.Register(action, func(ctx context.Context, m []byte) error {
		var data T
		if err := gob.NewDecoder(bytes.NewReader(m)).Decode(&data); err != nil {
			return err
		}
		return job(ctx, data)
	})
	return Job[T]{
		action: action,
		queue:  client.queue,
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

func (j Job[T]) CreateAndGetID(ctx context.Context, data T) (goqite.ID, error) {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(data); err != nil {
		return "", err
	}
	return jobs.CreateAndGetID(ctx, j.queue, j.action, buffer.Bytes())
}

func (j Job[T]) CreateAndGetIDTx(ctx context.Context, tx *sql.Tx, data T) (goqite.ID, error) {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(data); err != nil {
		return "", err
	}
	return jobs.CreateAndGetIDTx(ctx, tx, j.queue, j.action, buffer.Bytes())
}
