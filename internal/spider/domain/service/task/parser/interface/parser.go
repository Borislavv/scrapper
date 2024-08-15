package taskparserinterface

import "net/url"

type TaskParser interface {
	Parse() ([]*url.URL, error)
}
