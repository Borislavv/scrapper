package main

import (
	"context"
	spider "github.com/Borislavv/scrapper/internal/spider/app"
	"log"
	"sync"
)

func main() {
	log.Println("starting the spider...")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	//gsh := shutdown.NewGraceful(cancel)

	wg.Add(1)
	go spider.New(ctx).Run(wg)

	//gsh.ListenAndCancel()

	//log.Println("shutting down...")

	wg.Wait()

	cancel()

	log.Println("successfully shut down")
}
