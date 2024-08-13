package pagerepositoryinterface

import (
	"context"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"net/url"
)

var (
	NotFoundError = errors.New("page not found")
	InsertError   = errors.New("error occurred wile inserting page")
)

type PageRepository interface {
	FindByURL(ctx context.Context, url url.URL) (*entity.Page, error)
	Save(ctx context.Context, page *entity.Page) error
}
