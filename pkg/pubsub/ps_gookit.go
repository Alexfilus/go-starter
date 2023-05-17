package pubsub

import (
	"context"

	goKitEvent "github.com/gookit/event"
)

// var _ IEvent = (*GoKitEvent)(nil)
type GoKitEvent struct{}

func NewGokitEvent() *GoKitEvent {
	return &GoKitEvent{}
}

func (e *GoKitEvent) Publish(ctx context.Context, topic string, data any) error {
	goKitEvent.Fire(topic, goKitEvent.M{"0": data})
	return nil
}

func (e *GoKitEvent) Subscribe(ctx context.Context, topic string, handler any) error {
	listener := mustBeGookitSubHandler(handler)
	goKitEvent.On(topic, listener, goKitEvent.Normal)
	return nil
}

func (e *GoKitEvent) Close() {
	panic("not implemented. pls implement me")
}

// checks that the subscription Handler TYPE is of same with GOOKIT Sub Handler
//
// if they are same it cast it, returning a concrete type
// if different it panics to warn the developer
func mustBeGookitSubHandler(h any) goKitEvent.Listener {
	v, ok := h.(goKitEvent.Listener)
	if !ok {
		panic("invalid GOOKIT subscription handler")
	}
	return v
}
