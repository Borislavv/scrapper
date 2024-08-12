package jobrunner

import (
	"context"
	spiderinterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	taskproviderinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/provider/interface"
	taskrunnerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/runner/interface"
	"golang.org/x/time/rate"
	"log"
	"sync"
	"time"
)

type JobRunner struct {
	config       spiderinterface.Config
	taskRunner   taskrunnerinterface.TaskRunner
	taskProvider taskproviderinterface.TaskProvider
}

func New(
	config spiderinterface.Config,
	taskRunner taskrunnerinterface.TaskRunner,
	taskProvider taskproviderinterface.TaskProvider,
) *JobRunner {
	return &JobRunner{
		config:       config,
		taskRunner:   taskRunner,
		taskProvider: taskProvider,
	}
}

func (s *JobRunner) Run(ctx context.Context) {
	// close the job 5 minutes before next will be running
	ctx, cancel := context.WithTimeout(ctx, s.config.GetJobsFrequency()-(time.Minute*5))
	defer cancel()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	limiter := rate.NewLimiter(rate.Limit(s.config.GetTasksPerSecondLimit()), s.config.GetTasksConcurrencyLimit())

	for url := range s.taskProvider.Provide(ctx) {
		if err := limiter.Wait(nil); err != nil {
			log.Println("JobRunner: " + err.Error())
			return
		}

		wg.Add(1)
		if err := s.taskRunner.Run(ctx, wg, url); err != nil {
			log.Println("JobRunner: " + err.Error())
			continue
		}
	}
}
