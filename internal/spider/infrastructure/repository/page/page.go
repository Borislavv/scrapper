package pagerepository

import (
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"net/url"
)

type PageRepository struct {
}

func New() *PageRepository {
	return &PageRepository{}
}

func (r *PageRepository) FindByURL(url url.URL) (*entity.Page, error) {
	return nil, nil
}

func (r *PageRepository) Save(page *entity.Page) error {
	return nil
}
