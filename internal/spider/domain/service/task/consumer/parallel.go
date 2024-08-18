package taskconsumer

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.xbet.lan/web-backend/php/spider/internal/shared/domain/enum/ctx"
	spiderconfiginterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/app/config/interface"
	taskconsumerinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/task/consumer/interface"
	taskrunnerinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/task/runner/interface"
	logger "gitlab.xbet.lan/web-backend/php/spider/internal/spider/infrastructure/logger/interface"
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
