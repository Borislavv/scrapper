package taskprovider

import (
	"context"
	spiderinterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	"net/url"
	"runtime"
)

type TaskProvider struct {
	config spiderinterface.Config
}

func New(config spiderinterface.Config) *TaskProvider {
	return &TaskProvider{config: config}
}

// Provide is a method that provides tasks to perform (manages the rate of goroutines creation).
func (p *TaskProvider) Provide(ctx context.Context) <-chan url.URL {
	ch := make(chan url.URL, runtime.NumCPU())

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			case ch <- url.URL{}:
			}
		}
	}()

	return ch
}
