package main

import (
	"context"
	"github.com/Borislavv/scrapper/internal/spider/app"
	"github.com/Borislavv/scrapper/pkg/shared/shutdown"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	gsh := shutdown.NewGraceful(cancel)

	wg.Add(1)
	go spider.New().Run(ctx, wg)

	gsh.ListenAndCancel()
	wg.Wait()
}
