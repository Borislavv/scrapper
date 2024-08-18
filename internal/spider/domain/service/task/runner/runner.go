package taskrunner

import (
	"context"
	spiderconfiginterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/app/config/interface"
	pageconsumerinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/page/consumer/interface"
	pageproviderinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/page/provider/interface"
	scannerdtointerface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/page/scanner/dto/interface"
	"net/url"
	"sync"
)

type Runner struct {
	config       spiderconfiginterface.Configurator
	pageProvider pageproviderinterface.Provider
	pageConsumer pageconsumerinterface.Consumer
}

// New is a constructor of Runner task runner.
func New(
	config spiderconfiginterface.Configurator,
	provider pageproviderinterface.Provider,
	consumer pageconsumerinterface.Consumer,
) *Runner {
	return &Runner{
		config:       config,
		pageProvider: provider,
		pageConsumer: consumer,
	}
}

// Run is a method which provides tasks and consume results.
func (r *Runner) Run(ctx context.Context, url *url.URL) {
	resultCh := make(chan scannerdtointerface.Result, 1)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		r.pageProvider.Provide(ctx, url, resultCh)
		close(resultCh)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		r.pageConsumer.Consume(ctx, resultCh)
	}()
}
