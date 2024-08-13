package spiderconfig

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	// URLsFilepath is an additional path to file with URLs.
	URLsFilepath string `envconfig:"URLS_FILEPATH" default:"/public/data/urls.csv"`
	// JobFrequency is an interval between jobs running (interval between spider execution).
	JobFrequency time.Duration `envconfig:"JOB_FREQUENCY" default:"1h"`
	// TasksPerSecondLimit indicates how many tasks will be processed ber one second.
	TasksPerSecondLimit int `envconfig:"TASKS_PER_SECOND_LIMIT" default:"10"`
	// TasksConcurrencyLimit indicates how many tasks will be processed at the same time.
	TasksConcurrencyLimit int `envconfig:"TASKS_CONCURRENCY_LIMIT" default:"10"`
	// TimeoutPerURL is a timeout per request.
	TimeoutPerURL time.Duration `envconfig:"TIMEOUT_PER_URL" default:"5s"`
	// MongoHost is a host from docker-compose (mongodb container name).
	MongoHost string `envconfig:"MONGO_HOST" default:"mongodb"`
	// MongoPort is an exposed port of mongodb.
	MongoPort int `envconfig:"MONGO_PORT" default:"27017"`
	// MongoLogin is a login for simple auth.
	MongoLogin string `envconfig:"MONGO_LOGIN" default:"seo"`
	// MongoPassword is a password for simple auth.
	MongoPassword string `envconfig:"MONGO_PASSWORD" default:"seo"`
	// MongoDatabase is a name of target mongodb.
	MongoDatabase string `envconfig:"MONGO_DATABASE" default:"seo"`
	// MongoPagesCollection is a name of page entities collection.
	MongoPagesCollection string `envconfig:"MONGO_PAGES_COLLECTION" default:"pages"`
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

func (c *Config) GetMongoPagesCollection() string {
	return c.MongoPagesCollection
}

func (c *Config) GetMongoRequestTimeout() time.Duration {
	return c.MongoRequestTimeout
}
