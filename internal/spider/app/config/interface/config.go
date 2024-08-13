package spiderconfiginterface

import (
	"time"
)

type Config interface {
	GetURLsFilepath() string
	GetJobsFrequency() time.Duration
	GetTasksPerSecondLimit() int
	GetTasksConcurrencyLimit() int
	GetTimeoutPerURL() time.Duration
	GetMongoHost() string
	GetMongoPort() int
	GetMongoLogin() string
	GetMongoPassword() string
	GetMongoDatabase() string
	GetMongoPagesCollection() string
	GetMongoRequestTimeout() time.Duration
}
