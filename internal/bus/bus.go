package bus

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"
)

var _ctx = context.Background()

func SetContext(ctx context.Context) {
	_ctx = ctx
}

var (
	subMu       = sync.RWMutex{}
	subNextID   = 0
	subHandlers = make(map[string][]handler)
)

type handler struct {
	ID int
	FN func(ctx context.Context, T any)
}

func subscribe[T any](name string, fn func(ctx context.Context, event T) error) (string, int) {
	// Get topic and handler ID
	topic := fmt.Sprintf("%T", *new(T))
	id := subNextID

	subNextID++

	// Add handler to topic
	subHandlers[topic] = append(subHandlers[topic], handler{
		ID: id,
		FN: func(ctx context.Context, event any) {
			if err := fn(ctx, event.(T)); err != nil {
				slog.Error("Failed to handle event", "package", "bus", "name", name, "error", err)
			}
		},
	})

	return topic, id
}

func unsubscribe(topic string, id int) {
	// Get handlers for topic
	handlers := subHandlers[topic]
	handlersLength := len(handlers)

	if handlersLength < 2 {
		// Remove all handlers
		subHandlers[topic] = []handler{}
	} else {
		// Replace handler with last handler and shrink slice by 1
		idx := slices.IndexFunc(handlers, func(s handler) bool { return s.ID == id })
		handlers[idx] = handlers[handlersLength-1]
		subHandlers[topic] = handlers[:handlersLength-1]
	}
}

func Subscribe[T any](name string, fn func(ctx context.Context, event T) error) func() {
	subMu.Lock()
	topic, id := subscribe(name, fn)
	subMu.Unlock()

	return func() {
		subMu.Lock()
		unsubscribe(topic, id)
		subMu.Unlock()
	}
}

func SubscribeChannel[T any]() (<-chan T, func()) {
	c := make(chan T)
	subMu.Lock()
	topic, id := subscribe("bus.SubscribeChannel", func(ctx context.Context, event T) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case c <- event:
			return nil
		}
	})
	subMu.Unlock()

	return c, func() {
		subMu.Lock()
		unsubscribe(topic, id)
		subMu.Unlock()
	}
}

func Publish[T any](event T) {
	subMu.RLock()
	topic := fmt.Sprintf("%T", event)
	for _, sub := range subHandlers[topic] {
		sub.FN(_ctx, event)
	}
	subMu.RUnlock()
}
