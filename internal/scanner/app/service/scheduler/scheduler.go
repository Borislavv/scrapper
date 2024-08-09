package scheduler

import (
	"context"
	"time"
)

type Scheduler struct {
}

func New() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Manage(ctx context.Context) <-chan time.Time {
	tCh := make(chan time.Time, 1)
	tCh <- time.Now()
	close(tCh)
	return tCh

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()
	return ticker.C
}
