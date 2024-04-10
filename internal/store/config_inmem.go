package store

import (
	"fmt"

	"github.com/c12s/star/internal/domain"
)

type configInMemStore struct {
	groups     map[string]*domain.ConfigGroup
	standalone map[string]*domain.StandaloneConfig
}

func NewConfigInMemStore() (domain.ConfigStore, error) {
	return &configInMemStore{
		groups:     make(map[string]*domain.ConfigGroup),
		standalone: make(map[string]*domain.StandaloneConfig),
	}, nil
}

func (r *configInMemStore) GetGroup(org string, name string, version string, namespace string) (*domain.ConfigGroup, *domain.Error) {
	config, ok := r.groups[genKey(org, name, version, namespace)]
	if !ok {
		return nil, domain.NewError(domain.ErrTypeNotFound, fmt.Sprintf("config group (org: %s, name: %s, version: %s) not found in namespace %s", org, name, version, namespace))
	}
	return config, nil
}

func (r *configInMemStore) GetStandalone(org string, name string, version string, namespace string) (*domain.StandaloneConfig, *domain.Error) {
	config, ok := r.standalone[genKey(org, name, version, namespace)]
	if !ok {
		return nil, domain.NewError(domain.ErrTypeNotFound, fmt.Sprintf("standalone config (org: %s, name: %s, version: %s) not found in namespace %s", org, name, version, namespace))
	}
	return config, nil
}

func (r *configInMemStore) PutGroup(config *domain.ConfigGroup) *domain.Error {
	r.groups[genKey(config.Org, config.Name, config.Version, config.Namespace)] = config
	return nil
}

func (r *configInMemStore) PutStandalone(config *domain.StandaloneConfig) *domain.Error {
	r.standalone[genKey(config.Org, config.Name, config.Version, config.Namespace)] = config
	return nil
}

func genKey(org string, name string, version string, namespace string) string {
	return fmt.Sprintf("%s/%s/%s/%s", namespace, org, name, version)
}
