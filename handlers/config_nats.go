package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/config/pkg/config"
	"github.com/c12s/star/domain"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type NatsConfigHandler struct {
	conn       *nats.Conn
	nodeIdRepo domain.NodeIdRepo
	marshaller config.Marshaller
}

func NewNatsConfigHandler(conn *nats.Conn, nodeIdRepo domain.NodeIdRepo, marshaller config.Marshaller) (*NatsConfigHandler, error) {
	return &NatsConfigHandler{
		conn:       conn,
		nodeIdRepo: nodeIdRepo,
		marshaller: marshaller,
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
	group, err := n.marshaller.UnmarshalConfigGroup(msg.Data)
	if err != nil {
		log.Println(err)
	}
	// todo: remove this later
	jsonGroup, err := json.Marshal(group)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(jsonGroup))
}

func (n *NatsConfigHandler) subject(nodeId string) string {
	return fmt.Sprintf("%s.config", nodeId)
}
