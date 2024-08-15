package taskrunnerinterface

import (
	"context"
	"net/url"
)

type TaskRunner interface {
	Run(ctx context.Context, url *url.URL)
}
