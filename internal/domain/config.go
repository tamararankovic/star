package domain

type ParamSet map[string]string

type NamedParamSet struct {
	Name string
	Set  ParamSet
}

type ConfigBase struct {
	Org       string
	Name      string
	Version   string
	CreatedAt string
	Namespace string
}

type StandaloneConfig struct {
	ConfigBase
	Set ParamSet
}

type ConfigGroup struct {
	ConfigBase
	Sets []NamedParamSet
}

type ConfigStore interface {
	PutStandalone(config *StandaloneConfig) *Error
	GetStandalone(org, name, version, namespace string) (*StandaloneConfig, *Error)
	PutGroup(config *ConfigGroup) *Error
	GetGroup(org, name, version, namespace string) (*ConfigGroup, *Error)
}
