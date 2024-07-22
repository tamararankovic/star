package servers

import (
	"errors"
	"log"

	meridianapi "github.com/c12s/meridian/pkg/api"
	"github.com/c12s/star/internal/services"
)

type AppConfigAsyncServer struct {
	client *meridianapi.MeridianAsyncClient
	serf   *services.SerfAgent
	nodeId string
}

func NewAppConfigAsyncServer(client *meridianapi.MeridianAsyncClient, serf *services.SerfAgent, nodeId string) (*AppConfigAsyncServer, error) {
	if client == nil {
		return nil, errors.New("client is nil while initializing app config async server")
	}
	return &AppConfigAsyncServer{
		client: client,
		serf:   serf,
		nodeId: nodeId,
	}, nil
}

func (c *AppConfigAsyncServer) Serve() {
	err := c.client.ReceiveConfig(func(orgId, namespaceName, appName, seccompProfile, strategy string, quotas map[string]float64) error {
		log.Printf("Organization: %s\n", orgId)
		log.Printf("Namespace: %s\n", namespaceName)
		log.Printf("Application: %s\n", appName)
		log.Printf("Resource quotas: %s\n", orgId)
		log.Printf("Seccomp profile: %s\n", seccompProfile)
		for resource, quota := range quotas {
			log.Printf("\t%s: %f\n", resource, quota)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}

func (c *AppConfigAsyncServer) GracefulStop() {
	c.client.GracefulStop()
}
