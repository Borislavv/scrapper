package pagefinder

import (
	"context"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	pagerepositoryinterface "github.com/Borislavv/scrapper/internal/spider/domain/repository/interface"
	"net/url"
)

type PageFinder struct {
	repository pagerepositoryinterface.PageRepository
}

func New(repository pagerepositoryinterface.PageRepository) *PageFinder {
	return &PageFinder{repository: repository}
}

func (f *PageFinder) FindByURL(ctx context.Context, url url.URL) (*entity.Page, error) {
	return f.repository.FindByURL(ctx, url)
}
