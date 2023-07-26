package configs

import (
	"os"
	"strconv"
)

type Config struct {
	natsAddress                        string
	registrationSubject                string
	registrationReqTimeoutMilliseconds int64
}

func (c *Config) NatsAddress() string {
	return c.natsAddress
}

func (c *Config) RegistrationSubject() string {
	return c.registrationSubject
}

func (c *Config) RegistrationReqTimeoutMilliseconds() int64 {
	return c.registrationReqTimeoutMilliseconds
}

func NewFromEnv() (*Config, error) {
	registrationReqTimeoutMilliseconds, err := strconv.Atoi(os.Getenv("REGISTRATION_REQ_TIMEOUT_MILLISECONDS"))
	if err != nil {
		return nil, err
	}
	return &Config{
		natsAddress:                        os.Getenv("NATS_ADDRESS"),
		registrationSubject:                os.Getenv("REGISTRATION_SUBJECT"),
		registrationReqTimeoutMilliseconds: int64(registrationReqTimeoutMilliseconds),
	}, nil
}

//type yamlConfig struct {
//	NatsAddress                        string
//	RegistrationSubject                string
//	RegistrationReqTimeoutMilliseconds int64
//}
//
//func NewFromYamlFile(filepath string) (*Config, error) {
//	fileContents, err := os.ReadFile(filepath)
//	if err != nil {
//		return nil, err
//	}
//	yamlConfig := yamlConfig{}
//	err = yaml.Unmarshal(fileContents, &yamlConfig)
//	if err != nil {
//		return nil, err
//	}
//	return yamlConfig.toConfig(), nil
//}
//
//func (yc yamlConfig) toConfig() *Config {
//	return &Config{
//		natsAddress:         yc.NatsAddress,
//		registrationSubject: yc.RegistrationSubject,
//	}
//}
