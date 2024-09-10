package sharedconfig

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	// LoggerLevel: info, debug, warning, error, fatal, panic.
	LoggerLevel string `envconfig:"LOGGER_LEVEL"  default:"debug"`
	// LoggerLevel: /dev/null/, stdout, stderr, or filename (logs will store in the {projectRoot}/var/log dir.).
	LoggerOutput string `envconfig:"LOGGER_OUTPUT" default:"stdout"`
	// LoggerFormatter: text, json.
	LoggerFormatter string `envconfig:"LOGGER_FORMAT" default:"json"`
	// LoggerContextExtraFields determines which fields must be extract from
	// context.Context and passed into log record (see more into ctxenum package).
	LoggerContextExtraFields []string `envconfig:"LOGGER_CONTEXT_EXTRA_FIELD" default:"jobId,taskId"`

	// MongoHost is a host from docker-compose (mongodb container name).
	MongoHost string `envconfig:"MONGO_HOST" default:"mongodb"`
	// MongoPort is an exposed port of mongodb.
	MongoPort int `envconfig:"MONGO_PORT" default:"27017"`
	// MongoLogin is a login for simple auth.
	MongoLogin string `envconfig:"MONGO_LOGIN" default:"spider"`
	// MongoPassword is a password for simple auth.
	MongoPassword string `envconfig:"MONGO_PASSWORD" default:"spider"`
	// MongoDatabase is a name of target mongodb.
	MongoDatabase string `envconfig:"MONGO_DATABASE" default:"spider"`
	// MongoRequestTimeout is a timeout per request to mongodb.
	MongoRequestTimeout time.Duration `envconfig:"MONGO_REQUEST_TIMEOUT" default:"5s"`

	// InfrastructureServerName is a name of the shared server.
	InfrastructureServerName string `envconfig:"INFRASTRUCTURE_SERVER_NAME" default:"infrastructure"`
	// InfrastructureServerPort is a port for shared server (endpoints like a /probe for k8s).
	InfrastructureServerPort string `envconfig:"INFRASTRUCTURE_SERVER_PORT" default:":8000"`
	// InfrastructureServerShutDownTimeout is a duration value before the server will be closed forcefully.
	InfrastructureServerShutDownTimeout time.Duration `envconfig:"INFRASTRUCTURE_SERVER_SHUTDOWN_TIMEOUT" default:"5s"`
	// InfrastructureServerRequestTimeout is a timeout value for close request forcefully.
	InfrastructureServerRequestTimeout time.Duration `envconfig:"INFRASTRUCTURE_SERVER_REQUEST_TIMEOUT" default:"1m"`
}

func Load() (*Config, error) {
	cfg := new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) GetLoggerLevel() string {
	return c.LoggerLevel
}

func (c *Config) GetLoggerOutput() string {
	return c.LoggerOutput
}

func (c *Config) GetLoggerFormatter() string {
	return c.LoggerFormatter
}

func (c *Config) GetLoggerContextExtraFields() []string {
	return c.LoggerContextExtraFields
}

func (c *Config) GetMongoHost() string {
	return c.MongoHost
}

func (c *Config) GetMongoPort() int {
	return c.MongoPort
}

func (c *Config) GetMongoLogin() string {
	return c.MongoLogin
}

func (c *Config) GetMongoPassword() string {
	return c.MongoPassword
}

func (c *Config) GetMongoDatabase() string {
	return c.MongoDatabase
}

func (c *Config) GetServerName() string {
	return c.InfrastructureServerName
}

func (c *Config) GetServerPort() string {
	return c.InfrastructureServerPort
}

func (c *Config) GetServerShutDownTimeout() time.Duration {
	return c.InfrastructureServerShutDownTimeout
}

func (c *Config) GetServerRequestTimeout() time.Duration {
	return c.InfrastructureServerRequestTimeout
}
