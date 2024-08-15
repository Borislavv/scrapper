package pagerepositoryinterface

import (
	"context"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
)

var (
	NotFoundError = errors.New("page not found")
)

type PageRepository interface {
	FindByURL(ctx context.Context, url string) (page *entity.Page, found bool, err error)
	Save(ctx context.Context, page *entity.Page) error
}
