package pagescanner

import (
	"context"
	"errors"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	pageparserinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/parser/interface"
	scannerdtointerface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/scanner/dto/interface"
	"github.com/Borislavv/scrapper/internal/spider/infrastructure/logger/interface"
	scannerdto "github.com/Borislavv/scrapper/internal/spider/infrastructure/service/page/scanner/dto"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type HTTP struct {
	config         spiderconfiginterface.Configurator
	parser         pageparserinterface.PageParser
	logger         logger.Logger
	httpClientPool *sync.Pool
}

// NewHTTP is a constructor of HTTP scanner.
func NewHTTP(
	config spiderconfiginterface.Configurator,
	parser pageparserinterface.PageParser,
	logger logger.Logger,
) *HTTP {
	return &HTTP{
		config: config,
		parser: parser,
		logger: logger,
		httpClientPool: &sync.Pool{
			New: func() any {
				return &http.Client{
					Timeout: time.Second * 60,
					Transport: &http.Transport{
						MaxIdleConns:        100,
						MaxIdleConnsPerHost: 1,
					},
				}
			},
		},
	}
}

// Scan is method which request target page and parse it into *entity.Page struct (recursive: depends on retries arg.).
func (s *HTTP) Scan(
	ctx context.Context,
	wg *sync.WaitGroup,
	url *url.URL,
	userAgent string,
	resultCh chan<- scannerdtointerface.Result,
	retries int,
) {
	defer wg.Done()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg = &sync.WaitGroup{}
	defer wg.Wait()

	req, err := s.prepareRequest(ctx, url, userAgent, resultCh)
	if err != nil {
		return
	}

	wg.Add(1)
	go s.scan(ctx, wg, req, resultCh, retries, cancel)
}

// scan is method which request target page and parse it into *entity.Page struct (recursive: depends on retries arg.).
func (s *HTTP) scan(
	ctx context.Context,
	wg *sync.WaitGroup,
	req *http.Request,
	resultCh chan<- scannerdtointerface.Result,
	retries int,
	cancel context.CancelFunc,
) {
	defer wg.Done()

	client := s.httpClientPool.Get().(*http.Client)
	defer s.httpClientPool.Put(client)

	resp, err := client.Do(req)
	if err != nil {
		s.logger.ErrorMsg(ctx, "scanning page failed due to request error occurred", logger.Fields{
			"url":       req.URL.String(),
			"userAgent": req.UserAgent(),
			"err":       err.Error(),
		})
		resultCh <- scannerdto.NewResult(nil, req.URL.String(), err)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		if retries == 0 { // retries are exceeded, we go out with logging error
			err = errors.New("non-positive status code received, all retries exceeded")
			s.logger.ErrorMsg(ctx, "scanning page failed due to "+err.Error(), logger.Fields{
				"url":        req.URL.String(),
				"userAgent":  req.UserAgent(),
				"statusCode": resp.StatusCode,
				"retries":    retries,
			})
			resultCh <- scannerdto.NewResult(nil, req.URL.String(), err)
			return
		}

		wg.Add(1) // run a new attempt for scanning
		go s.scan(ctx, wg, req, resultCh, retries-1, cancel)
	}

	page, err := s.parser.Parse(ctx, resp)
	if err != nil {
		s.logger.ErrorMsg(ctx, "scanning page failed due to parser error occurred", logger.Fields{
			"url":        req.URL.String(),
			"userAgent":  req.UserAgent(),
			"statusCode": resp.StatusCode,
			"retries":    retries,
			"err":        err.Error(),
		})
		resultCh <- scannerdto.NewResult(nil, req.URL.String(), err)
		return
	}

	resultCh <- scannerdto.NewResult(page, req.URL.String(), nil)

	cancel()
}

// prepareRequest is method which build a request with context and target user-agent.
func (s *HTTP) prepareRequest(
	ctx context.Context,
	url *url.URL,
	userAgent string,
	resultCh chan<- scannerdtointerface.Result,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		s.logger.ErrorMsg(ctx, "scanning page failed due to preparing request error occurred", logger.Fields{
			"url":       url.String(),
			"userAgent": userAgent,
			"err":       err.Error(),
		})
		resultCh <- scannerdto.NewResult(nil, url.String(), err)
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	return req, nil
}
