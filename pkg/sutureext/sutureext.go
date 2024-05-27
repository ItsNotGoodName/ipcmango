package sutureext

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/thejerf/suture/v4"
)

func EventHook() suture.EventHook {
	return func(ei suture.Event) {
		switch e := ei.(type) {
		case suture.EventStopTimeout:
			slog.Info("Service failed to terminate in a timely manner", slog.String("supervisor", e.SupervisorName), slog.String("service", e.ServiceName))
		case suture.EventServicePanic:
			slog.Warn("Caught a service panic, which shouldn't happen")
			slog.Info(e.Stacktrace, slog.String("panic", e.PanicMsg))
		case suture.EventServiceTerminate:
			slog.Error("Service failed", slog.Any("error", e.Err), slog.String("supervisor", e.SupervisorName), slog.String("service", e.ServiceName))
			b, _ := json.Marshal(e)
			slog.Debug(string(b))
		case suture.EventBackoff:
			slog.Debug("Too many service failures - entering the backoff state", slog.String("supervisor", e.SupervisorName))
		case suture.EventResume:
			slog.Debug("Exiting backoff state", slog.String("supervisor", e.SupervisorName))
		default:
			slog.Warn("Unknown suture supervisor event type", "type", int(e.Type()))
			b, _ := json.Marshal(e)
			slog.Info(string(b))
		}
	}
}

type ServiceFunc struct {
	name string
	fn   func(ctx context.Context) error
}

func NewServiceFunc(name string, fn func(ctx context.Context) error) ServiceFunc {
	return ServiceFunc{
		name: name,
		fn:   fn,
	}
}

func (s ServiceFunc) String() string {
	return s.name
}

func (s ServiceFunc) Serve(ctx context.Context) error {
	return s.fn(ctx)
}

func OneShotFunc(fn func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		err := fn(ctx)
		if err != nil {
			return errors.Join(suture.ErrTerminateSupervisorTree, err)
		}

		return suture.ErrDoNotRestart
	}
}

func NewServiceContext(name string) ServiceContext {
	return ServiceContext{
		name:  name,
		doneC: make(chan struct{}),
		ctxC:  make(chan context.Context),
	}
}

type ServiceContext struct {
	name  string
	doneC chan struct{}
	ctxC  chan context.Context
}

func (b ServiceContext) String() string {
	return b.name
}

func (b ServiceContext) Serve(ctx context.Context) error {
	select {
	case <-b.doneC:
		return suture.ErrDoNotRestart
	default:
	}
	defer close(b.doneC)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case b.ctxC <- ctx:
		}
	}
}

func (b ServiceContext) Context() context.Context {
	select {
	case <-b.doneC:
		return context.Background()
	case ctx := <-b.ctxC:
		return ctx
	}
}

// SanitizeError prevents the error from being interpreted as a context error unless it
// really is a context error because suture kills the service when it sees a context error.
func SanitizeError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
		return err
	}

	var newErrs [3]error

	if errors.Is(err, suture.ErrDoNotRestart) {
		newErrs[0] = suture.ErrDoNotRestart
	}

	if errors.Is(err, suture.ErrTerminateSupervisorTree) {
		newErrs[1] = suture.ErrTerminateSupervisorTree
	}

	newErrs[2] = errors.New(err.Error())

	return errors.Join(newErrs[:]...)
}
