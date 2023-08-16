package configs

import (
	"os"
	"strconv"
)

type Config struct {
	natsAddress                        string
	registrationReqTimeoutMilliseconds int64
	maxRegistrationRetries             int8
	nodeIdDirPath                      string
	nodeIdFileName                     string
	grpcServerAddress                  string
	oortAddress                        string
}

func (c *Config) NatsAddress() string {
	return c.natsAddress
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

func (c *Config) GrpcServerAddress() string {
	return c.grpcServerAddress
}

func (c *Config) OortAddress() string {
	return c.oortAddress
}

func NewFromEnv() (*Config, error) {
	registrationReqTimeoutMilliseconds, err := strconv.Atoi(os.Getenv("REGISTRATION_REQ_TIMEOUT_MILLISECONDS"))
	maxRegistrationRetries, err := strconv.Atoi(os.Getenv("MAX_REGISTRATION_RETRIES"))
	if err != nil {
		return nil, err
	}
	return &Config{
		natsAddress:                        os.Getenv("NATS_ADDRESS"),
		registrationReqTimeoutMilliseconds: int64(registrationReqTimeoutMilliseconds),
		maxRegistrationRetries:             int8(maxRegistrationRetries),
		nodeIdDirPath:                      os.Getenv("NODE_ID_DIR_PATH"),
		nodeIdFileName:                     os.Getenv("NODE_ID_FILE_NAME"),
		grpcServerAddress:                  os.Getenv("STAR_ADDRESS"),
		oortAddress:                        os.Getenv("OORT_ADDRESS"),
	}, nil
}
