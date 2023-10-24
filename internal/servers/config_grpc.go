package servers

import (
	"context"
	"errors"
	"github.com/c12s/star/internal/domain"
	"github.com/c12s/star/internal/mappers/proto"
	"github.com/c12s/star/internal/services"
	"github.com/c12s/star/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type starConfigServer struct {
	api.UnimplementedStarConfigServer
	service *services.ConfigService
}

func NewStarConfigServer(service *services.ConfigService) (api.StarConfigServer, error) {
	return &starConfigServer{
		service: service,
	}, nil
}

func (s starConfigServer) GetConfigGroup(ctx context.Context, req *api.GetConfigGroupReq) (*api.GetConfigGroupResp, error) {
	domainReq, err := proto.GetConfigGroupReqToDomain(req)
	if err != nil {
		return nil, err
	}
	domainResp, err := s.service.Get(*domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized()) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		if errors.Is(err, domain.ErrNotFound()) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp, err := proto.GetConfigGroupRespFromDomain(*domainResp)
	return resp, nil
}
