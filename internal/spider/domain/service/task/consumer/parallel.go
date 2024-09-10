package taskconsumer

import (
	"context"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/enum/ctx"
	logger "github.com/Borislavv/scrapper/internal/shared/domain/service/logger/interface"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	taskconsumerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/consumer/interface"
	taskrunnerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/runner/interface"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
	"net/url"
	"sync"
)

type Parallel struct {
	logger     logger.Logger
	config     spiderconfiginterface.Configurator
	taskRunner taskrunnerinterface.TaskRunner
}

func NewParallel(
	logger logger.Logger,
	config spiderconfiginterface.Configurator,
	taskRunner taskrunnerinterface.TaskRunner,
) *Parallel {
	return &Parallel{logger: logger, config: config, taskRunner: taskRunner}
}

func (c *Parallel) Consume(ctx context.Context, urlsCh <-chan *url.URL) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	limiter := rate.NewLimiter(
		rate.Limit(c.config.GetTasksPerSecondLimit()),
		c.config.GetTasksConcurrencyLimit(),
	)

	for u := range urlsCh {
		if err := limiter.Wait(ctx); err != nil {
			if !(errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)) {
				c.logger.ErrorMsg(ctx, taskconsumerinterface.RateLimiterError.Error(), logger.Fields{
					"tasksPerSecondLimit":  c.config.GetTasksPerSecondLimit(),
					"taskConcurrencyLimit": c.config.GetTasksConcurrencyLimit(),
				})
			}
			return
		}

		// set up a task UUID for propagate into logs
		ctx = context.WithValue(ctx, ctxenum.TaskUUIDKey, uuid.NewString())

		wg.Add(1)
		go func(u *url.URL) {
			defer wg.Done()
			c.taskRunner.Run(ctx, u)
		}(u)
	}
}
