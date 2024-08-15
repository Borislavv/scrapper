package jobrunner

import (
	"context"
	spiderinterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	taskconsumerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/consumer/interface"
	taskproviderinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/provider/interface"
	"time"
)

type Runner struct {
	config       spiderinterface.Configurator
	taskConsumer taskconsumerinterface.Consumer
	taskProvider taskproviderinterface.TaskProvider
}

// New is a constructor of Runner job runner.
func New(
	config spiderinterface.Configurator,
	taskConsumer taskconsumerinterface.Consumer,
	taskProvider taskproviderinterface.TaskProvider,
) *Runner {
	return &Runner{
		config:       config,
		taskProvider: taskProvider,
		taskConsumer: taskConsumer,
	}
}

// Run is a method which executes a job consistently and limiting rate.
func (s *Runner) Run(ctx context.Context) {
	// close the job 5 minutes before next will be running
	ctx, cancel := context.WithTimeout(ctx, s.config.GetJobsFrequency()-(time.Minute*5))
	defer cancel()

	s.taskConsumer.Consume(ctx, s.taskProvider.Provide(ctx))
}
