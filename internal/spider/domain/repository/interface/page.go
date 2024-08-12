package pagerepositoryinterface

import (
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"net/url"
)

var (
	NotFoundError = errors.New("page not found")
)

type PageRepository interface {
	FindByURL(url url.URL) (*entity.Page, error)
	Save(page *entity.Page) error
}
