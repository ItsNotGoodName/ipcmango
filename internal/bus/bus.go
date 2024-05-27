package bus

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

var _ctx = context.Background()

func SetContext(ctx context.Context) {
	_ctx = ctx
}

var subs = make(map[string][]func(ctx context.Context, T any))

func Subscribe[T any](name string, fn func(ctx context.Context, event T) error) {
	topic := fmt.Sprintf("%T", *new(T))
	subs[topic] = append(subs[topic], func(ctx context.Context, event any) {
		if err := fn(ctx, event.(T)); err != nil {
			slog.Error("Failed to handle event", "package", "bus", "name", name, "error", err)
		}
	})
}

func Publish[T any](event T) {
	for _, fn := range subs[fmt.Sprintf("%T", event)] {
		fn(_ctx, event)
	}
}

func NewHub[T any]() *Hub[T] {
	return &Hub[T]{
		mu:   sync.Mutex{},
		subs: make(map[*chan T]struct{}),
	}
}

type Hub[T any] struct {
	mu   sync.Mutex
	subs map[*chan T]struct{}
}

func (h *Hub[T]) Broadcast(ctx context.Context, event T) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for sub := range h.subs {
		select {
		case <-ctx.Done():
		case *sub <- event:
		}
	}

	return nil
}

func (h *Hub[T]) Register() *Hub[T] {
	Subscribe("bus.Hub", h.Broadcast)
	return h
}

func (h *Hub[T]) Subscribe(ctx context.Context) (<-chan T, func()) {
	h.mu.Lock()
	c := make(chan T)

	key := &c
	h.subs[key] = struct{}{}
	h.mu.Unlock()

	return c, func() {
		h.mu.Lock()
		delete(h.subs, key)
		h.mu.Unlock()
	}
}
