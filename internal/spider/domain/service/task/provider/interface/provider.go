package taskproviderinterface

import (
	"context"
	"net/url"
)

type TaskProvider interface {
	Provide(ctx context.Context) <-chan url.URL
}
