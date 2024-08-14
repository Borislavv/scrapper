package jobrunner

import (
	"context"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/util"
	spiderinterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	taskproviderinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/provider/interface"
	taskrunnerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/runner/interface"
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

	//limiter := rate.NewLimiter(
	//	rate.Limit(s.config.GetTasksPerSecondLimit()),
	//	s.config.GetTasksConcurrencyLimit(),
	//)

	ticker, cancel := util.NewTicker(ctx, time.Second*10)
	defer cancel()

	for url := range s.taskProvider.Provide(ctx) {
		<-ticker

		log.Println("processing url: " + url.String())

		wg.Add(1)
		go s.taskRunner.Run(ctx, wg, url)
	}
}
