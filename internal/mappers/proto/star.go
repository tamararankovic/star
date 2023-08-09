package proto

import (
	configapi "github.com/c12s/config/pkg/api"
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

func ConfigGroupToDomain(group *configapi.ConfigGroup) (*domain.ConfigGroup, error) {
	resp := &domain.ConfigGroup{
		Name:      group.Name,
		Namespace: group.Namespace,
		Configs:   make([]domain.Config, len(group.Configs)),
	}
	for i, config := range group.Configs {
		resp.Configs[i] = domain.Config{
			Key:   config.Key,
			Value: config.Value,
		}
	}
	return resp, nil
}

func ConfigGroupFromDomain(domainGroup domain.ConfigGroup) (*configapi.ConfigGroup, error) {
	group := &configapi.ConfigGroup{
		Name:      domainGroup.Name,
		Namespace: domainGroup.Namespace,
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
