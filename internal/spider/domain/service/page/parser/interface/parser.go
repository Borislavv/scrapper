package pageparserinterface

import (
	"context"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"net/http"
)

var (
	ParseDOMFailed       = errors.New("parsing page from HTML failed due to error occurred while building a DOM tree and reading it")
	ParsePlainTextFailed = errors.New("parsing page from HTML failed due to error occurred while querying the plaintext (content)")
)

type PageParser interface {
	Parse(ctx context.Context, resp *http.Response) (*entity.Page, error)
}
