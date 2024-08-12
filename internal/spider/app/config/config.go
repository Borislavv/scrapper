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
