package jobrunner

import (
	"context"
	"github.com/google/uuid"
	"gitlab.xbet.lan/web-backend/php/spider/internal/shared/domain/enum/ctx"
	spiderinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/app/config/interface"
	taskconsumerinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/task/consumer/interface"
	taskproviderinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/task/provider/interface"
	logger "gitlab.xbet.lan/web-backend/php/spider/internal/spider/infrastructure/logger/interface"
)

type Runner struct {
	logger       logger.Logger
	config       spiderinterface.Configurator
	taskConsumer taskconsumerinterface.Consumer
	taskProvider taskproviderinterface.TaskProvider
}

// New is a constructor of Runner job runner.
func New(
	logger logger.Logger,
	config spiderinterface.Configurator,
	taskConsumer taskconsumerinterface.Consumer,
	taskProvider taskproviderinterface.TaskProvider,
) *Runner {
	return &Runner{
		logger:       logger,
		config:       config,
		taskProvider: taskProvider,
		taskConsumer: taskConsumer,
	}
}

// Run is a method which executes a job consistently and limiting rate.
func (s *Runner) Run(ctx context.Context) {
	// set up a job UUID for propagate into logs
	ctx = context.WithValue(ctx, ctxenum.JobUUIDKey, uuid.NewString())

	s.logger.InfoMsg(ctx, "started a new job", nil)

	// close the job before a new job will be running
	ctx, cancel := context.WithTimeout(ctx, s.config.GetJobsFrequency()-(s.config.GetJobsFrequency()/10))
	defer cancel()

	s.taskConsumer.Consume(ctx, s.taskProvider.Provide(ctx))

	s.logger.InfoMsg(ctx, "job finished", nil)
}
