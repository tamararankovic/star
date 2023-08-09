package startup

import (
	configapi "github.com/c12s/config/pkg/api"
	magnetarapi "github.com/c12s/magnetar/pkg/api"
	"github.com/c12s/magnetar/pkg/messaging"
	"github.com/c12s/magnetar/pkg/messaging/nats"

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

	nodeIdRepo, err := repos.NewNodeIdFSRepo(config.NodeIdDirPath(), config.NodeIdFileName())
	if err != nil {
		return err
	}
	configRepo, err := repos.NewConfigInMemRepo()
	if err != nil {
		return err
	}

	regReqPublisher, err := nats.NewPublisher(natsConn)
	if err != nil {
		return err
	}
	regRespSubscriberFactory := func(subject string) messaging.Subscriber {
		subscriber, _ := nats.NewSubscriber(natsConn, subject, "")
		return subscriber
	}

	oortClient, err := newOortClient(config.OortAddress())
	if err != nil {
		return err
	}
	registrationClient, err := magnetarapi.NewAsyncRegistrationClient(regReqPublisher, regRespSubscriberFactory)
	if err != nil {
		return err
	}

	configService, err := services.NewConfigService(configRepo, oortClient)
	if err != nil {
		return err
	}

	var nodeIdChan chan string
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGINT)

	go func() {
		nodeId := <-nodeIdChan

		configSubscriber, err := nats.NewSubscriber(natsConn, configapi.Subject(nodeId), nodeId)
		if err != nil {
			log.Println(err)
			return
		}
		configClient, err := configapi.NewConfigClient(configSubscriber)
		if err != nil {
			log.Println(err)
			return
		}

		configNatsHandler, err := handlers.NewConfigHandler(configClient, configService)
		if err != nil {
			log.Println(err)
			return
		}
		configNatsHandler.Handle()

		<-quit

		err = configSubscriber.Unsubscribe()
		if err != nil {
			log.Println(err)
		}
	}()

	rs := services.NewRegistrationService(registrationClient, nodeIdRepo, nodeIdChan)

	configGrpcServer, err := handlers.NewStarConfigServer(configService)
	if err != nil {
		return err
	}
	server, err := startServer(config.GrpcServerAddress(), configGrpcServer)
	if err != nil {
		return err
	}

	if !rs.Registered() {
		err = rs.Register(config.MaxRegistrationRetries())
		if err != nil {
			log.Println(err)
		}
	}

	<-quit

	server.GracefulStop()
	if err != nil {
		log.Println(err)
	}
	natsConn.Close()

	return nil
}
