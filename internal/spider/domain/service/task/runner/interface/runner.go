package taskrunnerinterface

import (
	"context"
	"net/url"
	"sync"
)

type TaskRunner interface {
	Run(ctx context.Context, wg *sync.WaitGroup, url url.URL)
}
