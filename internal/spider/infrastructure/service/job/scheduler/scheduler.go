package jobscheduler

import (
	"context"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/util"
	spiderinterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
)

type JobScheduler struct {
	config spiderinterface.Config
}

func New(config spiderinterface.Config) *JobScheduler {
	return &JobScheduler{config: config}
}

func (s *JobScheduler) Manage(ctx context.Context) <-chan struct{} {
	ticker, cancel := util.NewTicker(ctx, s.config.GetJobsFrequency())

	runJobCh := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				cancel()
				close(runJobCh)
				return
			case <-ticker:
				runJobCh <- struct{}{}
			}
		}
	}()

	return runJobCh
}
