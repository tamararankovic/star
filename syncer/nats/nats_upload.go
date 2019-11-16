package nats

import (
	fPb "github.com/c12s/scheme/flusher"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
)

type NatsUploader struct {
	natsConnection *nats.Conn
	nodeId         string
	topic          string
	errtopic       string
}

func NewNatsUploader(address, nodeid, topic, errtopic string) (*NatsUploader, error) {
	natsConnection, err := nats.Connect(address)
	if err != nil {
		return nil, err
	}

	return &NatsUploader{
		natsConnection: natsConnection,
		nodeId:         nodeid,
		topic:          topic,
		errtopic:       errtopic,
	}, nil
}

func (n *NatsUploader) SetNodeId(nodeid string) {
	n.nodeId = nodeid
}

func (n *NatsUploader) Upload(data *fPb.Update) {
	state, err := proto.Marshal(data)
	if err != nil {
		return
	}
	n.natsConnection.Publish(n.topic, state)
}

func (n *NatsUploader) Error(topic string, data []byte) {
	n.natsConnection.Publish(topic, data)
	n.natsConnection.Flush()
}

func (n *NatsUploader) NodeId() string {
	return n.nodeId
}
