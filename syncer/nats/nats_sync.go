package nats

import (
	"github.com/c12s/star/syncer"
	nats "github.com/nats-io/nats.go"
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

func (n *NatsSync) Subscribe(topic string, f syncer.Fn) {
	n.natsConnection.Subscribe(topic, func(msg *nats.Msg) {
		f(msg.Data)
	})
}

func (n *NatsSync) Error(topic string, data []byte) {
	n.natsConnection.Publish(topic, data)
	n.natsConnection.Flush()
}
