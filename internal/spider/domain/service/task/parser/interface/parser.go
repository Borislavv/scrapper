package taskparserinterface

import (
	"context"
	"errors"
	"net/url"
)

var (
	BuildFilepathError = errors.New("parsing tasks failed due to error occurred while building path to file")
	OpenFileError      = errors.New("parsing tasks failed due to error occurred while opening file")
	ReadFileError      = errors.New("parsing tasks failed due to error occurred while reading URLs")
	ParseURLError      = errors.New("parsing task failed due to error occurred while parsing URL to *url.URL")
)

type TaskParser interface {
	Parse(ctx context.Context) ([]*url.URL, error)
}
