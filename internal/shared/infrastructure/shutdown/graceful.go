package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Graceful struct {
	wg         *sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
	osSigsCh   chan os.Signal
}

// NewGraceful is a constructor of new Graceful shutdown implementation.
// Accepts main context.Context and context.CancelFunc.
func NewGraceful(ctx context.Context, cancelFunc context.CancelFunc) *Graceful {
	osSigsCh := make(chan os.Signal, 1)
	signal.Notify(osSigsCh, syscall.SIGINT, syscall.SIGTERM)

	wg := &sync.WaitGroup{}

	return &Graceful{
		wg:         wg,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		osSigsCh:   osSigsCh,
	}
}

// Add is a clone of (sync.WaitGroup).Add() method which must be called on main goroutines.
func (g *Graceful) Add(n int) {
	g.wg.Add(n)
}

// Done is a clone of (sync.WaitGroup).Done() method which must be called on closing main goroutines.
func (g *Graceful) Done() {
	g.wg.Done()
}

// ListenCancelAndAwait will catch one of channels (osSigsCh:[syscall.SIGINT, syscall.SIGTERM])
// or ctx.Done() and awaits while all main goroutines will be finished by sync.WaitGroup.
// NOTE: the ListenCancelAndAwait method is a synchronous (blocking) and must not be called from goroutine.
// Also, must be called at the end of the main function.
func (g *Graceful) ListenCancelAndAwait() {
	defer g.wg.Wait()

	select {
	case <-g.ctx.Done():
		g.cancelFunc()
		return
	case <-g.osSigsCh:
		g.cancelFunc()
		return
	}
}
