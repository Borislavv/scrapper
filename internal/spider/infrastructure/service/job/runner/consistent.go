package jobrunner

import (
	"context"
	spiderinterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	taskproviderinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/provider/interface"
	taskrunnerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/runner/interface"
	"log"
	"sync"
	"time"
)

type Consistent struct {
	config       spiderinterface.Config
	taskRunner   taskrunnerinterface.TaskRunner
	taskProvider taskproviderinterface.TaskProvider
}

func New(
	config spiderinterface.Config,
	taskRunner taskrunnerinterface.TaskRunner,
	taskProvider taskproviderinterface.TaskProvider,
) *Consistent {
	return &Consistent{
		config:       config,
		taskRunner:   taskRunner,
		taskProvider: taskProvider,
	}
}

func (s *Consistent) Run(ctx context.Context) {
	defer log.Println("Consistent: closed")
	log.Println("Consistent: started")

	// close the job 5 minutes before next will be running
	jobCtx, jobCancel := context.WithTimeout(ctx, s.config.GetJobsFrequency()-(time.Minute*5))
	defer jobCancel()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	//limiter := rate.NewLimiter(
	//	rate.Limit(s.config.GetTasksPerSecondLimit()),
	//	s.config.GetTasksConcurrencyLimit(),
	//)

	for url := range s.taskProvider.Provide(jobCtx) {
		//if err := limiter.Wait(ctx); err != nil {
		//	if !(errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)) {
		//		log.Println("Consistent: " + err.Error())
		//	}
		//	return
		//}

		log.Println("processing url: " + url.String())

		wg.Add(1)
		go s.taskRunner.Run(jobCtx, wg, url)
	}
}
