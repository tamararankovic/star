package nats

import (
	"github.com/c12s/star/flusher"
	"github.com/nats-io/go-nats"
)

type NatsSync struct {
	natsConnection *nats.Conn
}

func NewNatsSync(address string) (*NatsSync, error) {
	natsConnection, err := nats.Connect(address)
	if err != nil {
		return nil, err
	}

	return &NatsSync{
		natsConnection: natsConnection,
	}, nil
}

func (n *NatsSync) Subscribe(topic string, f flusher.Fn) {
	n.natsConnection.Subscribe(topic, func(msg *nats.Msg) {
		f(msg.Data)
	})
}
