package handlers

import (
	"fmt"
	"github.com/c12s/star/domain"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type NatsConfigHandler struct {
	conn       *nats.Conn
	nodeIdRepo domain.NodeIdRepo
}

func NewNatsConfigHandler(conn *nats.Conn, nodeIdRepo domain.NodeIdRepo) (*NatsConfigHandler, error) {
	return &NatsConfigHandler{
		conn:       conn,
		nodeIdRepo: nodeIdRepo,
	}, nil
}

func (n *NatsConfigHandler) Handle(nodeIdCh chan string) (chan bool, error) {
	subscriptionClosedCh := make(chan bool)

	go func() {
		nodeId, err := n.nodeIdRepo.Get()
		var id string
		if err != nil {
			id = <-nodeIdCh
		} else {
			id = nodeId.Value
		}
		subscription, err := n.conn.Subscribe(n.subject(id), n.handleNewConfig)
		if err != nil {
			log.Println(err)
		}

		for subscription.IsValid() {
			time.Sleep(1000 * time.Millisecond)
		}
		subscriptionClosedCh <- true
	}()

	return subscriptionClosedCh, nil
}

func (n *NatsConfigHandler) handleNewConfig(msg *nats.Msg) {
	log.Println("stiglo")
}

func (n *NatsConfigHandler) subject(nodeId string) string {
	return fmt.Sprintf("%s.config", nodeId)
}
