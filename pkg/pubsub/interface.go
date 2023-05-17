package pubsub

import (
	"context"

	"github.com/nats-io/nats.go"
)

// IPubSub is a unified interface for Publishing and subscribing to events
type IEvent interface {
	// Publish is used to publish events to a topic
	Publish(ctx context.Context, topic string, data any) error
	// Subscribe is used to subscribe to event from a topic
	Subscribe(ctx context.Context, topic string, h nats.MsgHandler) error
	// Closes an event connection
	Close()
}

type IPublish interface {
	// Publish is used to publish events to a topic
	Publish(ctx context.Context, topic string, data any) error
	// Subscribe is used to subscribe to event from a topic
	Subscribe(ctx context.Context, topic string, h nats.MsgHandler) error
	// Closes an event connection
	Close()
}
