package servers

import (
	"errors"
	"log"

	kuiperapi "github.com/c12s/kuiper/pkg/api"
	"github.com/c12s/star/internal/domain"
	"github.com/c12s/star/internal/mappers/proto"
)

type ConfigAsyncServer struct {
	client  *kuiperapi.KuiperAsyncClient
	configs domain.ConfigStore
}

func NewConfigAsyncServer(client *kuiperapi.KuiperAsyncClient, configs domain.ConfigStore) (*ConfigAsyncServer, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	return &ConfigAsyncServer{
		client:  client,
		configs: configs,
	}, nil
}

func (c *ConfigAsyncServer) Serve() {
	err := c.client.ReceiveConfig(
		func(cmd *kuiperapi.ApplyStandaloneConfigCommand) {
			config, err := proto.ApplyStandaloneConfigCommandToDomain(cmd)
			if err != nil {
				log.Println(err)
				return
			}
			putErr := c.configs.PutStandalone(config)
			if putErr != nil {
				log.Println(putErr)
			}
			// todo: vrati odgovor o statusu task-a nazad
		},
		func(cmd *kuiperapi.ApplyConfigGroupCommand) {
			config, err := proto.ApplyConfigGroupCommandToDomain(cmd)
			if err != nil {
				log.Println(err)
				return
			}
			putErr := c.configs.PutGroup(config)
			if putErr != nil {
				log.Println(putErr)
			}
			// todo: vrati odgovor o statusu task-a nazad
		})
	if err != nil {
		log.Println(err)
	}
}

func (c *ConfigAsyncServer) GracefulStop() {
	c.client.GracefulStop()
}
