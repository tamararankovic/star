package domain

import (
	"errors"
)

type Config struct {
	Key   string
	Value string
}

type ConfigGroup struct {
	Id      string
	Configs []Config
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
	GroupId string
	SubId   string
	SubKind string
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
