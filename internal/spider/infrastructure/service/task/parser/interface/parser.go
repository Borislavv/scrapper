package taskparserinterface

import "net/url"

type TaskParser interface {
	ParseURLs() ([]*url.URL, error)
}
