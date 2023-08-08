package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/c12s/config/pkg/config"
	"github.com/c12s/star/internal/domain"
	"github.com/c12s/star/internal/services"
	"github.com/nats-io/nats.go"
	"log"
)

type NatsConfigHandler struct {
	conn       *nats.Conn
	nodeIdRepo domain.NodeIdRepo
	marshaller config.Marshaller
	service    *services.ConfigService
}

func NewNatsConfigHandler(conn *nats.Conn, nodeIdRepo domain.NodeIdRepo, marshaller config.Marshaller, service *services.ConfigService) (*NatsConfigHandler, error) {
	return &NatsConfigHandler{
		conn:       conn,
		nodeIdRepo: nodeIdRepo,
		marshaller: marshaller,
		service:    service,
	}, nil
}

func (n *NatsConfigHandler) Handle(nodeIdCh chan string) (*nats.Subscription, error) {
	var subscription *nats.Subscription

	go func() {
		nodeId, err := n.nodeIdRepo.Get()
		var id string
		if err != nil {
			id = <-nodeIdCh
		} else {
			id = nodeId.Value
		}
		subscription, err = n.conn.Subscribe(n.subject(id), n.handleNewConfig)
		if err != nil {
			log.Println(err)
		}
	}()

	return subscription, nil
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
	req := domain.PutConfigGroupReq{
		Group: domain.ConfigGroup{
			Name:      group.Name,
			Namespace: group.Namespace,
			Configs:   make([]domain.Config, len(group.Configs)),
		},
	}
	for _, config := range group.Configs {
		req.Group.Configs = append(req.Group.Configs, domain.Config{
			Key:   config.Key,
			Value: config.Value,
		})
	}
	_, err = n.service.Put(req)
	if err != nil {
		log.Println(err)
	}
}

func (n *NatsConfigHandler) subject(nodeId string) string {
	return fmt.Sprintf("%s.configs", nodeId)
}
