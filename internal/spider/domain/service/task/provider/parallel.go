package provider

import (
	"context"
	"net/url"

	spiderinterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	taskparserinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/parser/interface"
	taskproviderinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/task/provider/interface"
	logger "github.com/Borislavv/scrapper/internal/spider/infrastructure/logger/interface"
)

type Parallel struct {
	urls   []*url.URL
	logger logger.Logger
	config spiderinterface.Configurator
	parser taskparserinterface.TaskParser
}

// NewParallel is a constructor of Parallel task provider.
func NewParallel(
	ctx context.Context,
	logger logger.Logger,
	config spiderinterface.Configurator,
	parser taskparserinterface.TaskParser,
) (*Parallel, error) {
	ch := &Parallel{config: config, logger: logger, parser: parser}
	if err := ch.parseTasks(ctx); err != nil {
		return nil, err
	}
	return ch, nil
}

// Provide is a method that sends tasks to perform.
func (p *Parallel) Provide(ctx context.Context) chan *url.URL {
	urlsCh := make(chan *url.URL, 1)

	go func() {
		defer close(urlsCh)
		for _, u := range p.urls {
			select {
			case <-ctx.Done():
				return
			case urlsCh <- u:
			}
		}
	}()

	return urlsCh
}

func (p *Parallel) parseTasks(ctx context.Context) error {
	urls, err := p.parser.Parse(ctx)
	if err != nil {
		return p.logger.Fatal(ctx, taskproviderinterface.ParseTasksError, logger.Fields{"err": err.Error()})
	}
	p.urls = urls
	return nil
}
