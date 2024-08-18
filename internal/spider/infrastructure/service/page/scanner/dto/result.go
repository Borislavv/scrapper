package scannerdto

import (
	"gitlab.xbet.lan/web-backend/php/spider/internal/shared/domain/entity"
)

type Result struct {
	page      *entity.Page
	url       string
	userAgent string
	err       error
}

func NewResult(page *entity.Page, url, userAgent string, err error) *Result {
	return &Result{page: page, url: url, userAgent: userAgent, err: err}
}

func (r *Result) URL() string {
	return r.url
}

func (r *Result) UserAgent() string {
	return r.userAgent
}

func (r *Result) Page() *entity.Page {
	return r.page
}

func (r *Result) Error() error {
	return r.err
}
