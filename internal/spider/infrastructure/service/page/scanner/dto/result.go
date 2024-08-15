package scannerdto

import (
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
)

type Result struct {
	page *entity.Page
	url  string
	err  error
}

func NewResult(page *entity.Page, url string, err error) *Result {
	return &Result{page: page, url: url, err: err}
}

func (r *Result) URL() string {
	return r.url
}

func (r *Result) Page() *entity.Page {
	return r.page
}

func (r *Result) Error() error {
	return r.err
}
