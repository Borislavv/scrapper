package pagescrapperinterface

import (
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"net/url"
)

type PageScrapper interface {
	Scrape(url *url.URL) (*entity.Page, error)
}
