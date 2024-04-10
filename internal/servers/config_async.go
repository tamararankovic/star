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
		func(protoConfig *kuiperapi.StandaloneConfig, namespace string) error {
			config, err := proto.ApplyStandaloneConfigCommandToDomain(protoConfig, namespace)
			if err != nil {
				return err
			}
			putErr := c.configs.PutStandalone(config)
			if putErr != nil {
				return errors.New(putErr.Message())
			}
			return nil
		},
		func(protoConfig *kuiperapi.ConfigGroup, namespace string) error {
			config, err := proto.ApplyConfigGroupCommandToDomain(protoConfig, namespace)
			if err != nil {
				return err
			}
			putErr := c.configs.PutGroup(config)
			if putErr != nil {
				return errors.New(putErr.Message())
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func (c *ConfigAsyncServer) GracefulStop() {
	c.client.GracefulStop()
}
