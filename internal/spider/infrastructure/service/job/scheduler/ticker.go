package jobscheduler

import (
	"context"
	"gitlab.xbet.lan/web-backend/php/spider/internal/shared/infrastructure/util"
	spiderinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/app/config/interface"
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
