package pagescannerinterface

import (
	"context"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"net/url"
	"sync"
)

type PageScanner interface {
	Scan(
		ctx context.Context,
		wg *sync.WaitGroup,
		url *url.URL,
		userAgent string,
		pagesCh chan<- *entity.Page,
		errsCh chan<- error,
		retries int,
	)
}
