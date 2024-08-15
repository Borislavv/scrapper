package taskconsumer

import (
	"context"
	taskrunnerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/runner/interface"
	"net/url"
	"sync"
)

type Parallel struct {
	taskRunner taskrunnerinterface.TaskRunner
}

func NewParallel(taskRunner taskrunnerinterface.TaskRunner) *Parallel {
	return &Parallel{taskRunner: taskRunner}
}

func (c *Parallel) Consume(ctx context.Context, urlsCh <-chan *url.URL) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	//limiter := rate.NewLimiter(
	//	rate.Limit(s.config.GetTasksPerSecondLimit()),
	//	s.config.GetTasksConcurrencyLimit(),
	//)

	for u := range urlsCh {
		//if err := limiter.Wait(ctx); err != nil {
		//	if !(errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)) {
		//		log.Println("Runner: " + err.Error())
		//	}
		//	return
		//}

		wg.Add(1)
		go func(u *url.URL) {
			defer wg.Done()
			c.taskRunner.Run(ctx, u)
		}(u)
	}
}
