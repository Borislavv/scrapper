package taskconsumerinterface

import (
	"context"
	"net/url"
)

type Consumer interface {
	Consume(ctx context.Context, urlsCh <-chan *url.URL)
}
