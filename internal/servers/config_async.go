package servers

import (
	"errors"
	kuiperapi "github.com/c12s/kuiper/pkg/api"
	"github.com/c12s/star/internal/mappers/proto"
	"github.com/c12s/star/internal/services"
	"log"
)

type ConfigAsyncServer struct {
	client  *kuiperapi.KuiperAsyncClient
	service *services.ConfigService
}

func NewConfigAsyncServer(client *kuiperapi.KuiperAsyncClient, service *services.ConfigService) (*ConfigAsyncServer, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	return &ConfigAsyncServer{
		client:  client,
		service: service,
	}, nil
}

func (c *ConfigAsyncServer) Serve() {
	err := c.client.ReceiveConfig(func(cmd *kuiperapi.ApplyConfigCommand) {
		req, err := proto.ApplyConfigCommandToDomain(cmd)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = c.service.Put(*req)
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Println(err)
	}
}

func (c *ConfigAsyncServer) GracefulStop() {
	c.client.GracefulStop()
}
