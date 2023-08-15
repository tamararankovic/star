package domain

import (
	"errors"
	"fmt"
)

type Config struct {
	Key   string
	Value string
}

type ConfigGroup struct {
	Name      string
	Namespace string
	Configs   []Config
}

func (c *ConfigGroup) Id() string {
	return fmt.Sprintf("%s/%s", c.Namespace, c.Name)
}

type ConfigRepo interface {
	Put(group ConfigGroup) error
	Get(groupId string) (*ConfigGroup, error)
}

type PutConfigGroupReq struct {
	Group ConfigGroup
}

type PutConfigGroupResp struct {
}

type GetConfigGroupReq struct {
	GroupName string
	Namespace string
	SubId     string
	SubKind   string
}

type GetConfigGroupResp struct {
	Group ConfigGroup
}

var (
	errUnauthorized = errors.New("unauthorized")
	errNotFound     = errors.New("not found")
)

func ErrUnauthorized() error {
	return errUnauthorized
}

func ErrNotFound() error {
	return errNotFound
}
