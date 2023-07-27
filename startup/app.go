package startup

import (
	"github.com/c12s/magnetar/pkg/marshallers/proto"
	"github.com/c12s/star/apis"
	"github.com/c12s/star/configs"
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
	rs := services.NewRegistrationService(registrationAPI, nodeIdRepo, config.MaxRegistrationRetries())
	if !rs.Registered() {
		return rs.Register()
	}
	// todo: if registered, subscribe to a subject for receiving config
	return nil
}
