package main

import (
	"context"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/shutdown"
	spider "github.com/Borislavv/scrapper/internal/spider/app"
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
