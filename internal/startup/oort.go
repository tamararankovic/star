package startup

import (
	oort "github.com/c12s/oort/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newOortClient(address string) (oort.OortEvaluatorClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return oort.NewOortEvaluatorClient(conn), nil
}
