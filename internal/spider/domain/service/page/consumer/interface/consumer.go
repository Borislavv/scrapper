package pageconsumerinterface

import (
	"context"
	"errors"
	scannerdtointerface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/scanner/dto/interface"
)

var (
	ScanURLError  = errors.New("page consume failed due to error occurred while scanning url")
	FindPageError = errors.New("page consume failed due to error occurred while searching page by repository")
	SavePageError = errors.New("page consume failed due to error occurred while saving page by repository")
)

var (
	PagesAreEqualMsg        = "pages are equal"
	PageSavedAtFirstTimeMsg = "page saved at first time"
)

type Consumer interface {
	Consume(ctx context.Context, resultCh <-chan scannerdtointerface.Result)
}
