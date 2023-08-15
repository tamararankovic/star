package proto

import (
	configapi "github.com/c12s/kuiper/pkg/api"
	"github.com/c12s/star/internal/domain"
	"github.com/c12s/star/pkg/api"
)

func GetConfigGroupReqToDomain(req *api.GetConfigGroupReq) (*domain.GetConfigGroupReq, error) {
	return &domain.GetConfigGroupReq{
		GroupName: req.GroupName,
		Namespace: req.Namespace,
		SubId:     req.Identity.Id,
		SubKind:   req.Identity.Kind,
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
			Name:      cmd.Group.Name,
			Namespace: cmd.Namespace,
			Configs:   make([]domain.Config, len(cmd.Group.Configs)),
		},
	}
	for i, config := range cmd.Group.Configs {
		resp.Group.Configs[i] = domain.Config{
			Key:   config.Key,
			Value: config.Value,
		}
	}
	return resp, nil
}

func ConfigGroupFromDomain(domainGroup domain.ConfigGroup) (*configapi.ConfigGroup, error) {
	group := &configapi.ConfigGroup{
		Name: domainGroup.Name,
	}
	for _, domainConfig := range domainGroup.Configs {
		config := &configapi.Config{
			Key:   domainConfig.Key,
			Value: domainConfig.Value,
		}
		group.Configs = append(group.Configs, config)
	}
	return group, nil
}
