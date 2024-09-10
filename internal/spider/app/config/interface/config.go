package spiderconfiginterface

import (
	"time"
)

type Configurator interface {
	GetLoggerLevel() string
	GetLoggerOutput() string
	GetLoggerFormatter() string
	GetLoggerContextExtraFields() []string

	GetMongoHost() string
	GetMongoPort() int
	GetMongoLogin() string
	GetMongoPassword() string
	GetMongoDatabase() string

	GetURLsFilepath() string
	GetJobsFrequency() time.Duration
	GetTasksPerSecondLimit() int
	GetTasksConcurrencyLimit() int
	GetTimeoutPerURL() time.Duration
	GetUserAgents() []string
	GetRequestRetries() int
	GetMongoPagesCollection() string
	GetMongoRequestTimeout() time.Duration
}
