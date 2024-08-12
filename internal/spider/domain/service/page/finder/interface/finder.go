package pagefinderinterface

import (
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"net/url"
)

type PageFinder interface {
	FindByURL(url url.URL) (*entity.Page, error)
}
