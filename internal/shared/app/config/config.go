package sharedconfig

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	// LoggerLevel: info, debug, warning, error, fatal, panic.
	LoggerLevel string `envconfig:"LOGGER_LEVEL"  default:"debug"`
	// LoggerLevel: /dev/null/, stdout, or path to file (logs will store in the {projectRoot}/var/log dir.).
	LoggerOutput string `envconfig:"LOGGER_OUTPUT" default:"stdout"`
	// LoggerFormatter: text, json.
	LoggerFormatter string `envconfig:"LOGGER_FORMAT" default:"json"`
	// LoggerReportCaller determines whether the call function and code line writes into log.
	LoggerReportCaller bool `envconfig:"LOGGER_REPORT_CALLER" default:"true"`
	// LoggerContextExtraFields determines which fields must be extract from context.Context and passed into log record.
	LoggerContextExtraFields []string `envconfig:"LOGGER_CONTEXT_EXTRA_FIELD" default:"JOB_ID,TASK_ID"`
	// MongoHost is a host from docker-compose (mongodb container name).
	MongoHost string `envconfig:"MONGO_HOST" default:"mongodb"`
	// MongoPort is an exposed port of mongodb.
	MongoPort int `envconfig:"MONGO_PORT" default:"27017"`
	// MongoLogin is a login for simple auth.
	MongoLogin string `envconfig:"MONGO_LOGIN" default:"useo"`
	// MongoPassword is a password for simple auth.
	MongoPassword string `envconfig:"MONGO_PASSWORD" default:"pseo"`
	// MongoDatabase is a name of target mongodb.
	MongoDatabase string `envconfig:"MONGO_DATABASE" default:"seo"`
	// MongoRequestTimeout is a timeout per request to mongodb.
	MongoRequestTimeout time.Duration `envconfig:"MONGO_REQUEST_TIMEOUT" default:"5s"`
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

func (c *Config) GetLoggerReportCaller() bool {
	return c.LoggerReportCaller
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
