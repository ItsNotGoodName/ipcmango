// Code generated by generate-bus.go; DO NOT EDIT.
package core

import (
	"context"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog/log"
)

func busLogError(err error) bool {
	if err != nil {
		log.Err(err).Str("package", "core").Msg("Failed to handle event")
		return true
	}
	return false
}

func NewBus() *Bus {
	return &Bus{
		ServiceContext: sutureext.NewServiceContext("core.Bus"),
	}
}

type Bus struct {
	sutureext.ServiceContext
	onEventDahuaCameraCreated []func(ctx context.Context, event models.EventDahuaCameraCreated) error
	onEventDahuaCameraUpdated []func(ctx context.Context, event models.EventDahuaCameraUpdated) error
	onEventDahuaCameraDeleted []func(ctx context.Context, event models.EventDahuaCameraDeleted) error
	onEventDahuaCameraEvent []func(ctx context.Context, event models.EventDahuaCameraEvent) error
	onEventDahuaEventWorkerConnecting []func(ctx context.Context, event models.EventDahuaEventWorkerConnecting) error
	onEventDahuaEventWorkerConnect []func(ctx context.Context, event models.EventDahuaEventWorkerConnect) error
	onEventDahuaEventWorkerDisconnect []func(ctx context.Context, event models.EventDahuaEventWorkerDisconnect) error
	onEventDahuaCoaxialStatus []func(ctx context.Context, event models.EventDahuaCoaxialStatus) error
}

func (b *Bus) Register(pub pubsub.Pub) {
	b.OnEventDahuaCameraCreated(func(ctx context.Context, event models.EventDahuaCameraCreated) error {
		return pub.Publish(ctx, event)
	})
	b.OnEventDahuaCameraUpdated(func(ctx context.Context, event models.EventDahuaCameraUpdated) error {
		return pub.Publish(ctx, event)
	})
	b.OnEventDahuaCameraDeleted(func(ctx context.Context, event models.EventDahuaCameraDeleted) error {
		return pub.Publish(ctx, event)
	})
	b.OnEventDahuaCameraEvent(func(ctx context.Context, event models.EventDahuaCameraEvent) error {
		return pub.Publish(ctx, event)
	})
	b.OnEventDahuaEventWorkerConnecting(func(ctx context.Context, event models.EventDahuaEventWorkerConnecting) error {
		return pub.Publish(ctx, event)
	})
	b.OnEventDahuaEventWorkerConnect(func(ctx context.Context, event models.EventDahuaEventWorkerConnect) error {
		return pub.Publish(ctx, event)
	})
	b.OnEventDahuaEventWorkerDisconnect(func(ctx context.Context, event models.EventDahuaEventWorkerDisconnect) error {
		return pub.Publish(ctx, event)
	})
	b.OnEventDahuaCoaxialStatus(func(ctx context.Context, event models.EventDahuaCoaxialStatus) error {
		return pub.Publish(ctx, event)
	})
}


func (b *Bus) OnEventDahuaCameraCreated(h func(ctx context.Context, event models.EventDahuaCameraCreated) error) {
	b.onEventDahuaCameraCreated = append(b.onEventDahuaCameraCreated, h)
}

func (b *Bus) OnEventDahuaCameraUpdated(h func(ctx context.Context, event models.EventDahuaCameraUpdated) error) {
	b.onEventDahuaCameraUpdated = append(b.onEventDahuaCameraUpdated, h)
}

func (b *Bus) OnEventDahuaCameraDeleted(h func(ctx context.Context, event models.EventDahuaCameraDeleted) error) {
	b.onEventDahuaCameraDeleted = append(b.onEventDahuaCameraDeleted, h)
}

func (b *Bus) OnEventDahuaCameraEvent(h func(ctx context.Context, event models.EventDahuaCameraEvent) error) {
	b.onEventDahuaCameraEvent = append(b.onEventDahuaCameraEvent, h)
}

func (b *Bus) OnEventDahuaEventWorkerConnecting(h func(ctx context.Context, event models.EventDahuaEventWorkerConnecting) error) {
	b.onEventDahuaEventWorkerConnecting = append(b.onEventDahuaEventWorkerConnecting, h)
}

func (b *Bus) OnEventDahuaEventWorkerConnect(h func(ctx context.Context, event models.EventDahuaEventWorkerConnect) error) {
	b.onEventDahuaEventWorkerConnect = append(b.onEventDahuaEventWorkerConnect, h)
}

func (b *Bus) OnEventDahuaEventWorkerDisconnect(h func(ctx context.Context, event models.EventDahuaEventWorkerDisconnect) error) {
	b.onEventDahuaEventWorkerDisconnect = append(b.onEventDahuaEventWorkerDisconnect, h)
}

func (b *Bus) OnEventDahuaCoaxialStatus(h func(ctx context.Context, event models.EventDahuaCoaxialStatus) error) {
	b.onEventDahuaCoaxialStatus = append(b.onEventDahuaCoaxialStatus, h)
}



func (b *Bus) EventDahuaCameraCreated(event models.EventDahuaCameraCreated) {
	for _, v := range b.onEventDahuaCameraCreated {
		busLogError(v(b.Context(), event))
	}
}

func (b *Bus) EventDahuaCameraUpdated(event models.EventDahuaCameraUpdated) {
	for _, v := range b.onEventDahuaCameraUpdated {
		busLogError(v(b.Context(), event))
	}
}

func (b *Bus) EventDahuaCameraDeleted(event models.EventDahuaCameraDeleted) {
	for _, v := range b.onEventDahuaCameraDeleted {
		busLogError(v(b.Context(), event))
	}
}

func (b *Bus) EventDahuaCameraEvent(event models.EventDahuaCameraEvent) {
	for _, v := range b.onEventDahuaCameraEvent {
		busLogError(v(b.Context(), event))
	}
}

func (b *Bus) EventDahuaEventWorkerConnecting(event models.EventDahuaEventWorkerConnecting) {
	for _, v := range b.onEventDahuaEventWorkerConnecting {
		busLogError(v(b.Context(), event))
	}
}

func (b *Bus) EventDahuaEventWorkerConnect(event models.EventDahuaEventWorkerConnect) {
	for _, v := range b.onEventDahuaEventWorkerConnect {
		busLogError(v(b.Context(), event))
	}
}

func (b *Bus) EventDahuaEventWorkerDisconnect(event models.EventDahuaEventWorkerDisconnect) {
	for _, v := range b.onEventDahuaEventWorkerDisconnect {
		busLogError(v(b.Context(), event))
	}
}

func (b *Bus) EventDahuaCoaxialStatus(event models.EventDahuaCoaxialStatus) {
	for _, v := range b.onEventDahuaCoaxialStatus {
		busLogError(v(b.Context(), event))
	}
}

