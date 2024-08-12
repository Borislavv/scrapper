package jobschedulerinterface

import (
	"context"
)

type JobScheduler interface {
	Manage(ctx context.Context) <-chan struct{}
}
