package main

import (
	"context"
	"gitlab.xbet.lan/web-backend/php/spider/internal/shared/infrastructure/shutdown"
	spider "gitlab.xbet.lan/web-backend/php/spider/internal/spider/app"
	"log"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	gsh := shutdown.NewGraceful(cancel)

	s, err := spider.New(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Start()
	}()

	gsh.ListenAndCancel()

	log.Println("successfully shut down, exit")
}
