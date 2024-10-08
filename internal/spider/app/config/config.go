package spiderconfig

import (
	"github.com/kelseyhightower/envconfig"
	"strings"
	"time"
)

type UserAgents []string

func (c *UserAgents) Decode(value string) error {
	*c = strings.Split(value, "|")
	return nil
}

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

	// URLsFilepath is an additional path to file with URLs.
	URLsFilepath string `envconfig:"URLS_FILEPATH" default:"/public/data/urls_test.csv"`
	// JobFrequency is an interval between jobs running (interval between spider execution).
	JobFrequency time.Duration `envconfig:"JOB_FREQUENCY" default:"1h"`
	// TasksPerSecondLimit indicates how many tasks will be processed ber one second.
	TasksPerSecondLimit int `envconfig:"TASKS_PER_SECOND_LIMIT" default:"10"`
	// TasksConcurrencyLimit indicates how many tasks will be processed at the same time.
	TasksConcurrencyLimit int `envconfig:"TASKS_CONCURRENCY_LIMIT" default:"10"`
	// TimeoutPerURL is a timeout per request.
	TimeoutPerURL time.Duration `envconfig:"TIMEOUT_PER_URL" default:"1m"`
	UserAgents    UserAgents    `envconfig:"USER_AGENTS" default:"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.5993.70 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"`
	//RequestRetries is a value that is used to determine the number of repetitions of the request if an error occurs.
	RequestRetries int `envconfig:"REQUEST_RETRIES" default:"3"`
	// MongoPagesCollection is a name of page entities collection.
	MongoPagesCollection string `envconfig:"MONGO_PAGES_COLLECTION" default:"pages"`
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

func (c *Config) GetURLsFilepath() string {
	return c.URLsFilepath
}

func (c *Config) GetJobsFrequency() time.Duration {
	return c.JobFrequency
}

func (c *Config) GetTasksPerSecondLimit() int {
	return c.TasksPerSecondLimit
}

func (c *Config) GetTasksConcurrencyLimit() int {
	return c.TasksConcurrencyLimit
}

func (c *Config) GetTimeoutPerURL() time.Duration {
	return c.TimeoutPerURL
}

func (c *Config) GetUserAgents() []string {
	return c.UserAgents
}

func (c *Config) GetRequestRetries() int {
	return c.RequestRetries
}

func (c *Config) GetMongoPagesCollection() string {
	return c.MongoPagesCollection
}

func (c *Config) GetMongoRequestTimeout() time.Duration {
	return c.MongoRequestTimeout
}
