package untils

import (
	"context"
	"time"
)

func NewTicker(ctx context.Context, interval time.Duration) (ch <-chan time.Time, cancelFn func()) {
	ctx, cancel := context.WithCancel(ctx)

	tickCh := make(chan time.Time, 1)
	tickCh <- time.Now()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				tickCh <- t
			}
		}
	}()

	return tickCh, func() { cancel() }
}
