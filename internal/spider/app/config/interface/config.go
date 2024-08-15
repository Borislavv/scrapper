package spiderconfiginterface

import (
	"time"
)

type Configurator interface {
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
