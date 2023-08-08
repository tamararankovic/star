package handlers

import (
	"context"
	"errors"
	configProto "github.com/c12s/config/pkg/proto"
	"github.com/c12s/star/internal/domain"
	"github.com/c12s/star/internal/services"
	"github.com/c12s/star/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type starConfigServer struct {
	proto.UnimplementedStarConfigServer
	service *services.ConfigService
}

func NewStarConfigServer(service *services.ConfigService) (proto.StarConfigServer, error) {
	return &starConfigServer{
		service: service,
	}, nil
}

func (s starConfigServer) GetConfigGroup(ctx context.Context, req *proto.GetConfigGroupReq) (*proto.GetConfigGroupResp, error) {
	domainReq := domain.GetConfigGroupReq{
		GroupName: req.GroupName,
		Namespace: req.Namespace,
		SubId:     req.Identity.Id,
		SubKind:   req.Identity.Kind,
	}
	domainResp, err := s.service.Get(domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized()) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		if errors.Is(err, domain.ErrNotFound()) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &proto.GetConfigGroupResp{
		Group: &configProto.ConfigGroup{
			Name:      domainResp.Group.Name,
			Namespace: domainResp.Group.Namespace,
		},
	}
	for _, domainConfig := range domainResp.Group.Configs {
		config := &configProto.Config{
			Key:   domainConfig.Key,
			Value: domainConfig.Value,
		}
		resp.Group.Configs = append(resp.Group.Configs, config)
	}
	return resp, nil
}
