package taskconsumerinterface

import (
	"context"
	"errors"
	"net/url"
)

var RateLimiterError = errors.New("task consuming failed due to error occurred while limiting rate")

type Consumer interface {
	Consume(ctx context.Context, urlsCh <-chan *url.URL)
}
