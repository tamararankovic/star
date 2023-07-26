package configs

import (
	"os"
	"strconv"
)

type Config struct {
	natsAddress                        string
	registrationSubject                string
	registrationReqTimeoutMilliseconds int64
	maxRegistrationRetries             int8
	nodeIdDirPath                      string
	nodeIdFileName                     string
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

func (c *Config) MaxRegistrationRetries() int8 {
	return c.maxRegistrationRetries
}

func (c *Config) NodeIdDirPath() string {
	return c.nodeIdDirPath
}

func (c *Config) NodeIdFileName() string {
	return c.nodeIdFileName
}

func NewFromEnv() (*Config, error) {
	registrationReqTimeoutMilliseconds, err := strconv.Atoi(os.Getenv("REGISTRATION_REQ_TIMEOUT_MILLISECONDS"))
	maxRegistrationRetries, err := strconv.Atoi(os.Getenv("MAX_REGISTRATION_RETRIES"))
	if err != nil {
		return nil, err
	}
	return &Config{
		natsAddress:                        os.Getenv("NATS_ADDRESS"),
		registrationSubject:                os.Getenv("REGISTRATION_SUBJECT"),
		registrationReqTimeoutMilliseconds: int64(registrationReqTimeoutMilliseconds),
		maxRegistrationRetries:             int8(maxRegistrationRetries),
		nodeIdDirPath:                      os.Getenv("NODE_ID_DIR_PATH"),
		nodeIdFileName:                     os.Getenv("NODE_ID_FILE_NAME"),
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
