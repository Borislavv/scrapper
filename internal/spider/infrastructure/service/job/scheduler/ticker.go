package jobscheduler

import (
	"context"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/util"
	spiderinterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
)

type Ticker struct {
	config spiderinterface.Configurator
}

func NewTicker(config spiderinterface.Configurator) *Ticker {
	return &Ticker{config: config}
}

func (s *Ticker) Manage(ctx context.Context) <-chan struct{} {
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
