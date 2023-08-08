package startup

import (
	configProto "github.com/c12s/config/pkg/proto"
	magnetarProto "github.com/c12s/magnetar/pkg/proto"
	"github.com/c12s/star/internal/apis"
	"github.com/c12s/star/internal/configs"
	"github.com/c12s/star/internal/handlers"
	"github.com/c12s/star/internal/repos"
	"github.com/c12s/star/internal/services"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func StartApp(config *configs.Config) error {
	natsConn, err := NewNatsConn(config.NatsAddress())
	if err != nil {
		return err
	}
	registrationMarshaller := magnetarProto.NewMarshaller()
	registrationAPI := apis.NewNatsRegistrationAPI(natsConn, config.RegistrationSubject(), config.RegistrationReqTimeoutMilliseconds(), registrationMarshaller)
	nodeIdRepo, err := repos.NewNodeIdFSRepo(config.NodeIdDirPath(), config.NodeIdFileName())
	if err != nil {
		return err
	}
	var nodeIdChan chan string
	rs := services.NewRegistrationService(registrationAPI, nodeIdRepo, nodeIdChan, config.MaxRegistrationRetries())

	configRepo, err := repos.NewConfigInMemRepo()
	if err != nil {
		return err
	}
	oortClient, err := newOortClient(config.OortAddress())
	configService, err := services.NewConfigService(configRepo, oortClient)
	configMarshaller, err := configProto.NewMarshaller()
	configNatsHandler, err := handlers.NewNatsConfigHandler(natsConn, nodeIdRepo, configMarshaller, configService)
	configGrpcServer, err := handlers.NewStarConfigServer(configService)
	if err != nil {
		return err
	}
	subscription, err := configNatsHandler.Handle(nodeIdChan)
	server, err := startServer(config.GrpcServerAddress(), configGrpcServer)
	if err != nil {
		return err
	}
	if !rs.Registered() {
		return rs.Register()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGINT)
	<-quit

	server.GracefulStop()
	err = subscription.Unsubscribe()
	if err != nil {
		log.Println(err)
	}
	natsConn.Close()

	return nil
}
