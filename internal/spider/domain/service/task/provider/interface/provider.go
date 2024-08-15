package taskproviderinterface

import (
	"context"
	"errors"
	"net/url"
)

var ParseTasksError = errors.New("initialization error due to error occurred while parsing tasks")

type TaskProvider interface {
	Provide(ctx context.Context) chan *url.URL
}
