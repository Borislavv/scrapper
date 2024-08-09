package main

import (
	"context"
	"github.com/Borislavv/scrapper/cmd/api"
	"github.com/Borislavv/scrapper/cmd/scanner"
	"github.com/Borislavv/scrapper/pkg/shared/shutdown"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	gsh := shutdown.NewGraceful(cancel)

	wg.Add(2)
	go api.NewApp().Run(ctx, wg)
	go scanner.NewApp().Run(ctx, wg)

	gsh.ListenAndCancel()
	wg.Wait()
}
