package pagescanner

import (
	"context"
	"errors"
	"fmt"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	spiderconfiginterface "github.com/Borislavv/scrapper/internal/spider/app/config/interface"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type HTTPScanner struct {
	config         spiderconfiginterface.Config
	httpClientPool *sync.Pool
	parser
}

func NewHTTP(config spiderconfiginterface.Config) *HTTPScanner {
	return &HTTPScanner{
		config: config,
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

func (s *HTTPScanner) Scan(
	ctx context.Context,
	wg *sync.WaitGroup,
	url *url.URL,
	userAgent string,
	pagesCh chan<- *entity.Page,
	errsCh chan<- error,
	retries int,
) {
	defer wg.Done()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg = &sync.WaitGroup{}
	defer wg.Wait()

	req, err := s.prepareRequest(ctx, url, userAgent)
	if err != nil {
		s.error(err, errsCh)
		return
	}

	wg.Add(1)
	go s.scan(ctx, wg, req, pagesCh, errsCh, retries, cancel)
}

func (s *HTTPScanner) scan(
	ctx context.Context,
	wg *sync.WaitGroup,
	req *http.Request,
	pagesCh chan<- *entity.Page,
	errsCh chan<- error,
	retries int,
	cancel context.CancelFunc,
) {
	defer wg.Done()

	// get a client from pool
	httpClient := s.httpClientPool.Get().(*http.Client)
	defer s.httpClientPool.Put(httpClient)

	// execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		s.error(err, errsCh)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	// check whether a request success
	if resp.StatusCode != http.StatusOK {
		if retries == 0 { // retries are exceeded, we go out with logging error
			s.error(fmt.Errorf("non-positive status code %d received for url: %s "+
				"after %d retries", resp.StatusCode, req.URL, retries), errsCh)
			return
		}

		wg.Add(1) // run a new attempt for scanning
		go s.scan(ctx, wg, req, pagesCh, errsCh, retries-1, cancel)
	}

	// extract page from raw response
	page, err := s.parseResponse(resp)
	if err != nil {
		s.error(err, errsCh)
		return
	}

	pagesCh <- page

	cancel()
}

func (s *HTTPScanner) prepareRequest(ctx context.Context, url *url.URL, userAgent string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

// error is method which escapes the context errors.
func (s *HTTPScanner) error(err error, errsCh chan<- error) {
	if !(errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)) {
		errsCh <- err
	}
}
