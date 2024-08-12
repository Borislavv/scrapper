package pagescrapperinterface

import (
	"context"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"net/url"
)

type PageScrapper interface {
	Scrape(ctx context.Context, url url.URL) (*entity.Page, error)
}
