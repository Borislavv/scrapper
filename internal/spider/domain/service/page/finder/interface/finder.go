package pagefinderinterface

import (
	"context"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"net/url"
)

type PageFinder interface {
	FindByURL(ctx context.Context, url *url.URL) (*entity.Page, error)
}
