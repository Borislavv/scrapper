package pagescannerinterface

import (
	"context"
	"errors"
	scannerdtointerface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/page/scanner/dto/interface"
	"net/url"
)

var (
	RequestError               = errors.New("scanning page failed due to error occurred while request execution")
	NonPositiveStatusCodeError = errors.New("scanning page failed due to received a non-positive status code , all retries exceeded")
	ParserError                = errors.New("scanning page failed due to error occurred while parsing parser page")
	PrepareRequestError        = errors.New("scanning page failed due to error occurred while preparing request")
)

type PageScanner interface {
	Scan(
		ctx context.Context,
		url *url.URL,
		userAgent string,
		resultCh chan<- scannerdtointerface.Result,
		retries int,
	)
}
