package startup

import (
	"github.com/c12s/star/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func startServer(address string, starConfigServer api.StarConfigServer) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	api.RegisterStarConfigServer(s, starConfigServer)
	reflection.Register(s)

	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return s, nil
}
