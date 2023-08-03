package startup

import (
	"github.com/c12s/magnetar/pkg/proto"
	"github.com/c12s/star/apis"
	"github.com/c12s/star/configs"
	"github.com/c12s/star/handlers"
	"github.com/c12s/star/repos"
	"github.com/c12s/star/services"
)

func StartApp(config *configs.Config) error {
	natsConn, err := NewNatsConn(config.NatsAddress())
	if err != nil {
		return err
	}
	marshaller := proto.NewMarshaller()
	registrationAPI := apis.NewNatsRegistrationAPI(natsConn, config.RegistrationSubject(), config.RegistrationReqTimeoutMilliseconds(), marshaller)
	nodeIdRepo, err := repos.NewNodeIdFSRepo(config.NodeIdDirPath(), config.NodeIdFileName())
	if err != nil {
		return err
	}
	var nodeIdChan chan string
	rs := services.NewRegistrationService(registrationAPI, nodeIdRepo, nodeIdChan, config.MaxRegistrationRetries())
	configHandler, err := handlers.NewNatsConfigHandler(natsConn, nodeIdRepo)
	if err != nil {
		return err
	}
	subscriptionClosedCh, err := configHandler.Handle(nodeIdChan)
	if err != nil {
		return err
	}
	if !rs.Registered() {
		return rs.Register()
	}

	<-subscriptionClosedCh

	return nil
}
