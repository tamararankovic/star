package proto

import (
	configapi "github.com/c12s/kuiper/pkg/api"
	"github.com/c12s/star/internal/domain"
	"github.com/c12s/star/pkg/api"
)

func ApplyConfigGroupCommandToDomain(cmd *configapi.ApplyConfigGroupCommand) (*domain.ConfigGroup, error) {
	resp := &domain.ConfigGroup{
		ConfigBase: domain.ConfigBase{
			Org:       cmd.Group.Organization,
			Name:      cmd.Group.Name,
			Version:   cmd.Group.Version,
			CreatedAt: cmd.Group.CreatedAt,
			Namespace: cmd.Namespace,
		},
	}
	for _, paramSet := range cmd.Group.ParamSets {
		set := domain.NamedParamSet{
			Name: paramSet.Name,
		}
		for _, param := range paramSet.ParamSet {
			set.Set[param.Key] = param.Value
		}
		resp.Sets = append(resp.Sets, set)
	}
	return resp, nil
}

func ApplyStandaloneConfigCommandToDomain(cmd *configapi.ApplyStandaloneConfigCommand) (*domain.StandaloneConfig, error) {
	resp := &domain.StandaloneConfig{
		ConfigBase: domain.ConfigBase{
			Org:       cmd.Config.Organization,
			Name:      cmd.Config.Name,
			Version:   cmd.Config.Version,
			CreatedAt: cmd.Config.CreatedAt,
			Namespace: cmd.Namespace,
		},
		Set: make(domain.ParamSet),
	}
	for _, param := range cmd.Config.ParamSet {
		resp.Set[param.Key] = param.Value
	}

	return resp, nil
}

func ConfigGroupFromDomain(domainGroup domain.ConfigGroup) (*api.ConfigGroup, error) {
	group := &api.ConfigGroup{
		Organization: domainGroup.Org,
		Name:         domainGroup.Name,
		Version:      domainGroup.Version,
		CreatedAt:    domainGroup.CreatedAt,
	}
	for _, paramSet := range domainGroup.Sets {
		set := &api.NamedParamSet{
			Name: paramSet.Name,
		}
		for key, value := range paramSet.Set {
			set.ParamSet = append(set.ParamSet, &api.Param{Key: key, Value: value})
		}
		group.ParamSets = append(group.ParamSets, set)
	}
	return group, nil
}

func StandaloneConfigFromDomain(domainConfig domain.StandaloneConfig) (*api.StandaloneConfig, error) {
	config := &api.StandaloneConfig{
		Organization: domainConfig.Org,
		Name:         domainConfig.Name,
		Version:      domainConfig.Version,
		CreatedAt:    domainConfig.CreatedAt,
	}
	for key, value := range domainConfig.Set {
		config.ParamSet = append(config.ParamSet, &api.Param{Key: key, Value: value})
	}
	return config, nil
}
