package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type Graceful struct {
	osSigsCh   chan os.Signal
	cancelFunc context.CancelFunc
}

func NewGraceful(cancelFunc context.CancelFunc) *Graceful {
	return &Graceful{cancelFunc: cancelFunc}
}

func (g *Graceful) ListenAndCancel() {
	g.osSigsCh = make(chan os.Signal, 1)
	signal.Notify(g.osSigsCh, syscall.SIGINT, syscall.SIGTERM)
	<-g.osSigsCh
	g.cancelFunc()
}
