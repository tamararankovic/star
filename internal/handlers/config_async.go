package handlers

import (
	"errors"
	configapi "github.com/c12s/config/pkg/api"
	"github.com/c12s/star/internal/domain"
	"github.com/c12s/star/internal/mappers/proto"
	"github.com/c12s/star/internal/services"
	"log"
)

type ConfigHandler struct {
	client  *configapi.ConfigClient
	service *services.ConfigService
}

func NewConfigHandler(client *configapi.ConfigClient, service *services.ConfigService) (*ConfigHandler, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	return &ConfigHandler{
		client:  client,
		service: service,
	}, nil
}

func (c *ConfigHandler) Handle() {
	err := c.client.ReceiveConfig(func(group *configapi.ConfigGroup) {
		groupDomain, err := proto.ConfigGroupToDomain(group)
		if err != nil {
			log.Println(err)
			return
		}
		req := domain.PutConfigGroupReq{
			Group: *groupDomain,
		}
		_, err = c.service.Put(req)
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Println(err)
	}
}
