package pubsub

import (
	"context"
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

var _ IEvent = (*Nats)(nil)

type Nats struct {
	natsConn *nats.Conn
}

func NewNatsFromCredential(url, pathToCredsFile string) *Nats {
	nc, err := nats.Connect(
		url,
		nats.UserCredentials(pathToCredsFile),
	)
	if err != nil {
		fmt.Println("error establishing nats connection -" + err.Error())
		os.Exit(0)
	}

	fmt.Println("nats connected")
	return &Nats{natsConn: nc}
}

func (e *Nats) Publish(ctx context.Context, topic string, data any) error {
	ec, err := nats.NewEncodedConn(e.natsConn, nats.JSON_ENCODER)
	if err != nil {
		fmt.Println("nats: Publisher Encoder Error -" + err.Error())
		return err
	}

	if err := ec.Publish(topic, data); err != nil {
		fmt.Println("nats: Publisher Encoder Error -" + err.Error())
		return err
	}

	return nil
}

func (e *Nats) Subscribe(ctx context.Context, topic string, h nats.MsgHandler) error {
	_, err := e.natsConn.Subscribe(topic, h)
	if err != nil {
		fmt.Println("nats: Subscriber Error - " + err.Error())
		return err
	}
	return nil
}

func (e *Nats) Close() {
	if e.natsConn != nil {
		if err := e.natsConn.Drain(); err != nil {
			fmt.Println("nats:  error draining nats - " + err.Error())
		}
		e.natsConn.Close()
	}
}
