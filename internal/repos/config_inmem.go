package repos

import "github.com/c12s/star/internal/domain"

type configInMemRepo struct {
	groups map[string]*domain.ConfigGroup
}

func NewConfigInMemRepo() (domain.ConfigRepo, error) {
	return &configInMemRepo{
		groups: make(map[string]*domain.ConfigGroup),
	}, nil
}

func (c configInMemRepo) Put(group domain.ConfigGroup) error {
	c.groups[group.Id()] = &group
	return nil
}

func (c configInMemRepo) Get(groupId string) (*domain.ConfigGroup, error) {
	group, ok := c.groups[groupId]
	if !ok {
		return nil, domain.ErrNotFound()
	}
	return group, nil
}
