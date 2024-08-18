package pageproviderinterface

import (
	"context"
	scannerdtointerface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/page/scanner/dto/interface"
	"net/url"
)

type Provider interface {
	Provide(ctx context.Context, url *url.URL, resultCh chan<- scannerdtointerface.Result)
}
