package api

import (
	"context"
	"sync"
)

type Api struct {
}

func NewApp() *Api {
	return &Api{}
}

func (a *Api) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	<-ctx.Done()
}
