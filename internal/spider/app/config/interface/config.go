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
}
