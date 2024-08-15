package pagescannerinterface

import (
	"context"
	scannerdtointerface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/scanner/dto/interface"
	"net/url"
	"sync"
)

type PageScanner interface {
	Scan(
		ctx context.Context,
		wg *sync.WaitGroup,
		url *url.URL,
		userAgent string,
		resultCh chan<- scannerdtointerface.Result,
		retries int,
	)
}
