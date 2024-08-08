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
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()
	return ticker.C
}
