package liveness

import (
	"context"
	livenessinterface "github.com/Borislavv/scrapper/internal/shared/infrastructure/liveness/interface"
	"sync/atomic"
)

type Prob struct {
	ctx                context.Context
	livenessAskCh      chan struct{}
	livenessResponseCh chan bool
	isClosed           atomic.Bool
}

func NewProbe(ctx context.Context) *Prob {
	return &Prob{
		ctx:                ctx,
		livenessAskCh:      make(chan struct{}),
		livenessResponseCh: make(chan bool),
		isClosed:           atomic.Bool{},
	}
}

func (p *Prob) Watch(s livenessinterface.Service) (cancel func()) {
	go func() {
		defer func() {
			p.isClosed.Store(true)
			close(p.livenessResponseCh)
		}()
		for range p.livenessAskCh {
			p.livenessResponseCh <- s.IsAlive()
		}
	}()

	return func() {
		close(p.livenessAskCh)
	}
}

func (p *Prob) IsAlive() bool {
	if !p.isClosed.Load() {
		p.livenessAskCh <- struct{}{}
		return <-p.livenessResponseCh
	}
	return false
}
