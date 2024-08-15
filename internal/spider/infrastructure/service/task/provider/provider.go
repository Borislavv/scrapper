package taskprovider

import (
	"context"
	spiderinterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	"github.com/Borislavv/scrapper/internal/spider/domain/service/task/parser/interface"
	"net/url"
	"runtime"
)

type TaskProvider struct {
	URLs   []*url.URL
	config spiderinterface.Config
}

func New(config spiderinterface.Config, parser taskparserinterface.TaskParser) *TaskProvider {
	URLs, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	return &TaskProvider{config: config, URLs: URLs}
}

// Provide is a method that sends tasks to perform.
func (p *TaskProvider) Provide(ctx context.Context) <-chan *url.URL {
	URLsCh := make(chan *url.URL, runtime.NumCPU())

	go func() {
		defer close(URLsCh)
		for _, URL := range p.URLs {
			select {
			case <-ctx.Done():
				return
			case URLsCh <- URL:
			}
		}
	}()

	return URLsCh
}
