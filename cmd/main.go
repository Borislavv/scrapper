package main

import (
	"context"
	spider "github.com/Borislavv/scrapper/internal/spider/app"
	"log"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	//gsh := shutdown.NewGraceful(cancel)

	s, err := spider.New(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	wg.Add(1)
	go s.Run(wg)

	//gsh.ListenAndCancel()

	//log.Println("shutting down...")

	wg.Wait()

	cancel()

	log.Println("successfully shut down, exit")
}
