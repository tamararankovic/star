package proto

import (
	configapi "github.com/c12s/kuiper/pkg/api"
	"github.com/c12s/star/internal/domain"
	"github.com/c12s/star/pkg/api"
)

func GetConfigGroupReqToDomain(req *api.GetConfigGroupReq) (*domain.GetConfigGroupReq, error) {
	return &domain.GetConfigGroupReq{
		GroupId: req.GroupId,
		SubId:   req.SubId,
		SubKind: req.SubKind,
	}, nil
}

func GetConfigGroupRespFromDomain(domainResp domain.GetConfigGroupResp) (*api.GetConfigGroupResp, error) {
	group, err := ConfigGroupFromDomain(domainResp.Group)
	if err != nil {
		return nil, err
	}
	return &api.GetConfigGroupResp{
		Group: group,
	}, nil
}

func ApplyConfigCommandToDomain(cmd *configapi.ApplyConfigCommand) (*domain.PutConfigGroupReq, error) {
	resp := &domain.PutConfigGroupReq{
		Group: domain.ConfigGroup{
			Id:      cmd.Id,
			Configs: make([]domain.Config, len(cmd.Configs)),
		},
	}
	for i, config := range cmd.Configs {
		resp.Group.Configs[i] = domain.Config{
			Key:   config.Key,
			Value: config.Value,
		}
	}
	return resp, nil
}

func ConfigGroupFromDomain(domainGroup domain.ConfigGroup) (*api.NodeConfigGroup, error) {
	group := &api.NodeConfigGroup{
		Id: domainGroup.Id,
	}
	for _, domainConfig := range domainGroup.Configs {
		config := &api.NodeConfig{
			Key:   domainConfig.Key,
			Value: domainConfig.Value,
		}
		group.Configs = append(group.Configs, config)
	}
	return group, nil
}
