package proto

import (
	configapi "github.com/c12s/kuiper/pkg/api"
	"github.com/c12s/star/internal/domain"
	"github.com/c12s/star/pkg/api"
)

func ApplyConfigGroupCommandToDomain(config *configapi.ConfigGroup, namespace string) (*domain.ConfigGroup, error) {
	resp := &domain.ConfigGroup{
		ConfigBase: domain.ConfigBase{
			Org:       config.Organization,
			Name:      config.Name,
			Version:   config.Version,
			CreatedAt: config.CreatedAt,
			Namespace: namespace,
		},
	}
	for _, paramSet := range config.ParamSets {
		set := domain.NamedParamSet{
			Name: paramSet.Name,
			Set:  make(domain.ParamSet),
		}
		for _, param := range paramSet.ParamSet {
			set.Set[param.Key] = param.Value
		}
		resp.Sets = append(resp.Sets, set)
	}
	return resp, nil
}

func ApplyStandaloneConfigCommandToDomain(config *configapi.StandaloneConfig, namespace string) (*domain.StandaloneConfig, error) {
	resp := &domain.StandaloneConfig{
		ConfigBase: domain.ConfigBase{
			Org:       config.Organization,
			Name:      config.Name,
			Version:   config.Version,
			CreatedAt: config.CreatedAt,
			Namespace: namespace,
		},
		Set: make(domain.ParamSet),
	}
	for _, param := range config.ParamSet {
		resp.Set[param.Key] = param.Value
	}

	return resp, nil
}

func ConfigGroupFromDomain(domainGroup domain.ConfigGroup) (*api.NodeConfigGroup, error) {
	group := &api.NodeConfigGroup{
		Organization: domainGroup.Org,
		Name:         domainGroup.Name,
		Version:      domainGroup.Version,
		CreatedAt:    domainGroup.CreatedAt,
	}
	for _, paramSet := range domainGroup.Sets {
		set := &api.NodeNamedParamSet{
			Name: paramSet.Name,
		}
		for key, value := range paramSet.Set {
			set.ParamSet = append(set.ParamSet, &api.NodeParam{Key: key, Value: value})
		}
		group.ParamSets = append(group.ParamSets, set)
	}
	return group, nil
}

func StandaloneConfigFromDomain(domainConfig domain.StandaloneConfig) (*api.NodeStandaloneConfig, error) {
	config := &api.NodeStandaloneConfig{
		Organization: domainConfig.Org,
		Name:         domainConfig.Name,
		Version:      domainConfig.Version,
		CreatedAt:    domainConfig.CreatedAt,
	}
	for key, value := range domainConfig.Set {
		config.ParamSet = append(config.ParamSet, &api.NodeParam{Key: key, Value: value})
	}
	return config, nil
}
