package pageprovider

import (
	"context"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	scannerdtointerface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/scanner/dto/interface"
	pagescannerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/scanner/interface"
	"net/url"
	"sync"
)

type Chan struct {
	config  spiderconfiginterface.Configurator
	scanner pagescannerinterface.PageScanner
}

func NewChan(config spiderconfiginterface.Configurator, scanner pagescannerinterface.PageScanner) *Chan {
	return &Chan{config: config, scanner: scanner}
}

func (p *Chan) Provide(ctx context.Context, url *url.URL) (resultCh chan scannerdtointerface.Result) {
	resultCh = make(chan scannerdtointerface.Result, 1)
	defer close(resultCh)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for _, userAgent := range p.config.GetUserAgents() {
		wg.Add(1)
		go p.scanner.Scan(ctx, wg, url, userAgent, resultCh, p.config.GetRequestRetries())
	}

	return resultCh
}
