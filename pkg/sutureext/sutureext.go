package sutureext

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/thejerf/suture/v4"
)

func NewSimple(name string) *suture.Supervisor {
	return suture.New("root", suture.Spec{
		EventHook: EventHook(),
	})
}

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

// Service forces the use of the String method
type Service interface {
	String() string
	suture.Service
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
