// Code generated by generate-bus.go; DO NOT EDIT.
package event

import (
	"context"
	"errors"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog/log"
)

func busLogError(err error) {
	if err != nil {
		log.Err(err).Str("package", "event").Msg("Failed to emit event")
	}
}

func NewBus() *Bus {
	return &Bus{
		lock: core.NewLockStore[string](),
		ServiceContext: sutureext.NewServiceContext("event.Bus"),
	}
}

type Bus struct {
	sutureext.ServiceContext
	lock *core.LockStore[string]
	onDahuaDeviceCreated []func(ctx context.Context, event DahuaDeviceCreated) error
	onDahuaDeviceUpdated []func(ctx context.Context, event DahuaDeviceUpdated) error
	onDahuaDeviceDeleted []func(ctx context.Context, event DahuaDeviceDeleted) error
	onDahuaEvent []func(ctx context.Context, event DahuaEvent) error
	onDahuaEventWorkerConnecting []func(ctx context.Context, event DahuaEventWorkerConnecting) error
	onDahuaEventWorkerConnect []func(ctx context.Context, event DahuaEventWorkerConnect) error
	onDahuaEventWorkerDisconnect []func(ctx context.Context, event DahuaEventWorkerDisconnect) error
	onDahuaCoaxialStatus []func(ctx context.Context, event DahuaCoaxialStatus) error
}

func (b *Bus) Register(pub pubsub.Pub) (*Bus) {
	b.OnDahuaDeviceCreated(func(ctx context.Context, evt DahuaDeviceCreated) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnDahuaDeviceUpdated(func(ctx context.Context, evt DahuaDeviceUpdated) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnDahuaDeviceDeleted(func(ctx context.Context, evt DahuaDeviceDeleted) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnDahuaEvent(func(ctx context.Context, evt DahuaEvent) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnDahuaEventWorkerConnecting(func(ctx context.Context, evt DahuaEventWorkerConnecting) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnDahuaEventWorkerConnect(func(ctx context.Context, evt DahuaEventWorkerConnect) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnDahuaEventWorkerDisconnect(func(ctx context.Context, evt DahuaEventWorkerDisconnect) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnDahuaCoaxialStatus(func(ctx context.Context, evt DahuaCoaxialStatus) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	return b
}


func (b *Bus) OnDahuaDeviceCreated(h func(ctx context.Context, evt DahuaDeviceCreated) error) {
	b.onDahuaDeviceCreated = append(b.onDahuaDeviceCreated, h)
}

func (b *Bus) OnDahuaDeviceUpdated(h func(ctx context.Context, evt DahuaDeviceUpdated) error) {
	b.onDahuaDeviceUpdated = append(b.onDahuaDeviceUpdated, h)
}

func (b *Bus) OnDahuaDeviceDeleted(h func(ctx context.Context, evt DahuaDeviceDeleted) error) {
	b.onDahuaDeviceDeleted = append(b.onDahuaDeviceDeleted, h)
}

func (b *Bus) OnDahuaEvent(h func(ctx context.Context, evt DahuaEvent) error) {
	b.onDahuaEvent = append(b.onDahuaEvent, h)
}

func (b *Bus) OnDahuaEventWorkerConnecting(h func(ctx context.Context, evt DahuaEventWorkerConnecting) error) {
	b.onDahuaEventWorkerConnecting = append(b.onDahuaEventWorkerConnecting, h)
}

func (b *Bus) OnDahuaEventWorkerConnect(h func(ctx context.Context, evt DahuaEventWorkerConnect) error) {
	b.onDahuaEventWorkerConnect = append(b.onDahuaEventWorkerConnect, h)
}

func (b *Bus) OnDahuaEventWorkerDisconnect(h func(ctx context.Context, evt DahuaEventWorkerDisconnect) error) {
	b.onDahuaEventWorkerDisconnect = append(b.onDahuaEventWorkerDisconnect, h)
}

func (b *Bus) OnDahuaCoaxialStatus(h func(ctx context.Context, evt DahuaCoaxialStatus) error) {
	b.onDahuaCoaxialStatus = append(b.onDahuaCoaxialStatus, h)
}



func (b *Bus) DahuaDeviceCreated(evt DahuaDeviceCreated) {
	ctx := b.Context()
	unlock, err := b.lock.Lock(ctx, evt.EventTopic())
	if err != nil{
		busLogError(err)
		return
	}
	for _, v := range b.onDahuaDeviceCreated {
		busLogError(v(ctx, evt))
	}
	unlock()
}

func (b *Bus) DahuaDeviceUpdated(evt DahuaDeviceUpdated) {
	ctx := b.Context()
	unlock, err := b.lock.Lock(ctx, evt.EventTopic())
	if err != nil{
		busLogError(err)
		return
	}
	for _, v := range b.onDahuaDeviceUpdated {
		busLogError(v(ctx, evt))
	}
	unlock()
}

func (b *Bus) DahuaDeviceDeleted(evt DahuaDeviceDeleted) {
	ctx := b.Context()
	unlock, err := b.lock.Lock(ctx, evt.EventTopic())
	if err != nil{
		busLogError(err)
		return
	}
	for _, v := range b.onDahuaDeviceDeleted {
		busLogError(v(ctx, evt))
	}
	unlock()
}

func (b *Bus) DahuaEvent(evt DahuaEvent) {
	ctx := b.Context()
	unlock, err := b.lock.Lock(ctx, evt.EventTopic())
	if err != nil{
		busLogError(err)
		return
	}
	for _, v := range b.onDahuaEvent {
		busLogError(v(ctx, evt))
	}
	unlock()
}

func (b *Bus) DahuaEventWorkerConnecting(evt DahuaEventWorkerConnecting) {
	ctx := b.Context()
	unlock, err := b.lock.Lock(ctx, evt.EventTopic())
	if err != nil{
		busLogError(err)
		return
	}
	for _, v := range b.onDahuaEventWorkerConnecting {
		busLogError(v(ctx, evt))
	}
	unlock()
}

func (b *Bus) DahuaEventWorkerConnect(evt DahuaEventWorkerConnect) {
	ctx := b.Context()
	unlock, err := b.lock.Lock(ctx, evt.EventTopic())
	if err != nil{
		busLogError(err)
		return
	}
	for _, v := range b.onDahuaEventWorkerConnect {
		busLogError(v(ctx, evt))
	}
	unlock()
}

func (b *Bus) DahuaEventWorkerDisconnect(evt DahuaEventWorkerDisconnect) {
	ctx := b.Context()
	unlock, err := b.lock.Lock(ctx, evt.EventTopic())
	if err != nil{
		busLogError(err)
		return
	}
	for _, v := range b.onDahuaEventWorkerDisconnect {
		busLogError(v(ctx, evt))
	}
	unlock()
}

func (b *Bus) DahuaCoaxialStatus(evt DahuaCoaxialStatus) {
	ctx := b.Context()
	unlock, err := b.lock.Lock(ctx, evt.EventTopic())
	if err != nil{
		busLogError(err)
		return
	}
	for _, v := range b.onDahuaCoaxialStatus {
		busLogError(v(ctx, evt))
	}
	unlock()
}

