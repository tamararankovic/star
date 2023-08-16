package startup

import (
	"context"
	"errors"
	kuiperapi "github.com/c12s/kuiper/pkg/api"
	magnetarapi "github.com/c12s/magnetar/pkg/api"
	oortapi "github.com/c12s/oort/pkg/api"
	"github.com/c12s/star/internal/configs"
	"github.com/c12s/star/internal/domain"
	"github.com/c12s/star/internal/repos"
	"github.com/c12s/star/internal/servers"
	"github.com/c12s/star/internal/services"
	"github.com/c12s/star/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

type app struct {
	config                    *configs.Config
	grpcServer                *grpc.Server
	configAsyncServer         *servers.ConfigAsyncServer
	configGrpcServer          api.StarConfigServer
	configService             *services.ConfigService
	registrationService       *services.RegistrationService
	configClient              *kuiperapi.KuiperAsyncClient
	evaluatorClient           oortapi.OortEvaluatorClient
	registrationClient        *magnetarapi.RegistrationAsyncClient
	nodeIdRepo                domain.NodeIdRepo
	configRepo                domain.ConfigRepo
	shutdownProcesses         []func()
	gracefulShutdownProcesses []func(wg *sync.WaitGroup)
}

func NewAppWithConfig(config *configs.Config) (*app, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	return &app{
		config:                    config,
		shutdownProcesses:         make([]func(), 0),
		gracefulShutdownProcesses: make([]func(wg *sync.WaitGroup), 0),
	}, nil
}

func (a *app) Start() error {
	a.init()

	err := a.startConfigAsyncServer()
	if err != nil {
		return err
	}
	err = a.startGrpcServer()
	if err != nil {
		return err
	}

	return nil
}

func (a *app) GracefulStop(ctx context.Context) {
	// call all shutdown processes after a timeout or graceful shutdown processes completion
	defer a.shutdown()

	// wait for all graceful shutdown processes to complete
	wg := &sync.WaitGroup{}
	wg.Add(len(a.gracefulShutdownProcesses))

	for _, gracefulShutdownProcess := range a.gracefulShutdownProcesses {
		go gracefulShutdownProcess(wg)
	}

	// notify when graceful shutdown processes are done
	gracefulShutdownDone := make(chan struct{})
	go func() {
		wg.Wait()
		gracefulShutdownDone <- struct{}{}
	}()

	// wait for graceful shutdown processes to complete or for ctx timeout
	select {
	case <-ctx.Done():
		log.Println("ctx timeout ... shutting down")
	case <-gracefulShutdownDone:
		log.Println("app gracefully stopped")
	}
}

func (a *app) init() {
	natsConn, err := NewNatsConn(a.config.NatsAddress())
	if err != nil {
		log.Fatalln(err)
	}
	a.shutdownProcesses = append(a.shutdownProcesses, func() {
		log.Println("closing nats conn")
		natsConn.Close()
	})

	a.initNodeIdFsRepo()
	a.initRegistrationClient()
	a.initRegistrationService()

	if !a.registrationService.Registered() {
		err := a.registrationService.Register(a.config.MaxRegistrationRetries())
		if err != nil {
			log.Fatalln(err)
		}
	}

	nodeId, err := a.nodeIdRepo.Get()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("NODE ID: " + nodeId.Value)
	a.initConfigInMemRepo()
	a.initEvaluatorClient()
	a.initConfigService()
	a.initConfigAsyncClient(nodeId.Value)

	a.initConfigAsyncServer()
	a.initConfigGrpcServer()
	a.initGrpcServer()
}

func (a *app) initGrpcServer() {
	if a.configGrpcServer == nil {
		log.Fatalln("config grpc server is nil")
	}
	s := grpc.NewServer()
	api.RegisterStarConfigServer(s, a.configGrpcServer)
	reflection.Register(s)
	a.grpcServer = s
}

func (a *app) initConfigGrpcServer() {
	if a.configService == nil {
		log.Fatalln("config service is nil")
	}
	server, err := servers.NewStarConfigServer(a.configService)
	if err != nil {
		log.Fatalln(err)
	}
	a.configGrpcServer = server
}

func (a *app) initConfigAsyncServer() {
	if a.configClient == nil {
		log.Fatalln("config client is nil")
	}
	if a.configService == nil {
		log.Fatalln("config service is nil")
	}
	server, err := servers.NewConfigAsyncServer(a.configClient, a.configService)
	if err != nil {
		log.Fatalln(err)
	}
	a.configAsyncServer = server
}

func (a *app) initRegistrationService() {
	if a.nodeIdRepo == nil {
		log.Fatalln("nodeid repo is nil")
	}
	if a.registrationClient == nil {
		log.Fatalln("registrator is nil")
	}
	a.registrationService = services.NewRegistrationService(a.registrationClient, a.nodeIdRepo)
}

func (a *app) initConfigService() {
	if a.configRepo == nil {
		log.Fatalln("config repo is nil")
	}
	if a.evaluatorClient == nil {
		log.Fatalln("oort evaluator is nil")
	}
	service, err := services.NewConfigService(a.configRepo, a.evaluatorClient)
	if err != nil {
		log.Fatalln(err)
	}
	a.configService = service
}

func (a *app) initEvaluatorClient() {
	client, err := newOortEvaluatorClient(a.config.OortAddress())
	if err != nil {
		log.Fatalln(err)
	}
	a.evaluatorClient = client
}

func (a *app) initConfigAsyncClient(nodeId string) {
	configClient, err := kuiperapi.NewKuiperAsyncClient(a.config.NatsAddress(), nodeId)
	if err != nil {
		log.Fatalln(err)
	}
	a.configClient = configClient
}

func (a *app) initRegistrationClient() {
	client, err := magnetarapi.NewRegistrationAsyncClient(a.config.NatsAddress())
	if err != nil {
		log.Fatalln(err)
	}
	a.registrationClient = client
}

func (a *app) initConfigInMemRepo() {
	repo, err := repos.NewConfigInMemRepo()
	if err != nil {
		log.Fatalln(err)
	}
	a.configRepo = repo
}

func (a *app) initNodeIdFsRepo() {
	repo, err := repos.NewNodeIdFSRepo(a.config.NodeIdDirPath(), a.config.NodeIdFileName())
	if err != nil {
		log.Fatalln(err)
	}
	a.nodeIdRepo = repo
}

func (a *app) startConfigAsyncServer() error {
	a.configAsyncServer.Serve()
	a.gracefulShutdownProcesses = append(a.gracefulShutdownProcesses, func(wg *sync.WaitGroup) {
		a.configAsyncServer.GracefulStop()
		log.Println("config async server gracefully stopped")
		wg.Done()
	})
	return nil
}

func (a *app) startGrpcServer() error {
	lis, err := net.Listen("tcp", a.config.GrpcServerAddress())
	if err != nil {
		return err
	}
	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := a.grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	a.gracefulShutdownProcesses = append(a.gracefulShutdownProcesses, func(wg *sync.WaitGroup) {
		a.grpcServer.GracefulStop()
		log.Println("grpc server gracefully stopped")
		wg.Done()
	})
	return nil
}

func (a *app) shutdown() {
	for _, shutdownProcess := range a.shutdownProcesses {
		shutdownProcess()
	}
}
