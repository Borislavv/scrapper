package pageconsumer

import (
	"context"
	pagerepositoryinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/repository/interface"
	pagecomparatorinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/page/comparator/interface"
	pageconsumerinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/page/consumer/interface"
	scannerdtointerface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/page/scanner/dto/interface"
	logger "gitlab.xbet.lan/web-backend/php/spider/internal/spider/infrastructure/logger/interface"
	"sync"
)

type Parallel struct {
	logger     logger.Logger
	repository pagerepositoryinterface.PageRepository
	comparator pagecomparatorinterface.PageComparator
}

// NewParallel is a constructor of Parallel page consumer.
func NewParallel(
	logger logger.Logger,
	repository pagerepositoryinterface.PageRepository,
	comparator pagecomparatorinterface.PageComparator,
) *Parallel {
	return &Parallel{
		logger:     logger,
		comparator: comparator,
		repository: repository,
	}
}

// Consume is a method which handles pages from channel in parallel.
func (c *Parallel) Consume(ctx context.Context, resultCh <-chan scannerdtointerface.Result) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for dto := range resultCh {
		wg.Add(1)
		go func(dto scannerdtointerface.Result) {
			defer wg.Done()

			if dto.Error() != nil { // check the scanner error for current page
				c.logger.ErrorMsg(ctx, pageconsumerinterface.ScanURLError.Error(), logger.Fields{
					"url":       dto.URL(),
					"err":       dto.Error(),
					"userAgent": dto.UserAgent(),
				})
				return
			}

			// extract current page if errors are not exist
			cur := dto.Page()

			// search a previous version of page
			prev, found, err := c.repository.FindOneLatestByURL(ctx, cur.URL)
			if err != nil {
				c.logger.ErrorMsg(ctx, pageconsumerinterface.FindPageError.Error(), logger.Fields{
					"url":       dto.URL(),
					"err":       dto.Error(),
					"userAgent": dto.UserAgent(),
				})
				return
			} else if !found {
				// page was not found by url, it's mean the page is new (save it at first time)
				if err = c.repository.Save(ctx, cur); err != nil {
					c.logger.ErrorMsg(ctx, pageconsumerinterface.SavePageError.Error(), logger.Fields{
						"url":       dto.URL(),
						"err":       err.Error(),
						"userAgent": dto.UserAgent(),
					})
				} else {
					c.logger.InfoMsg(ctx, pageconsumerinterface.PageSavedAtFirstTimeMsg, logger.Fields{
						"url":       cur.URL,
						"userAgent": dto.UserAgent(),
					})
				}
				return
			}

			// check the previous version as actual
			if !c.comparator.IsEquals(cur, prev) {
				// set up new version for current page
				cur.UpVersion(prev)

				// save blinking page
				if err = c.repository.Save(ctx, cur); err != nil {
					c.logger.ErrorMsg(ctx, pageconsumerinterface.SavePageError.Error(), logger.Fields{
						"url":       dto.URL(),
						"err":       err.Error(),
						"userAgent": dto.UserAgent(),
						"version":   cur.Version,
					})
				} else {
					c.logger.InfoMsg(ctx, pageconsumerinterface.BlinkingPageSavedMsg, logger.Fields{
						"url":       cur.URL,
						"userAgent": dto.UserAgent(),
						"version":   cur.Version,
					})
				}
			} else {
				// update the previous page (set up new page.UpdatedAt value)
				if err = c.repository.Update(ctx, prev); err != nil {
					c.logger.ErrorMsg(ctx, pageconsumerinterface.SavePageError.Error(), logger.Fields{
						"url":       dto.URL(),
						"err":       err.Error(),
						"userAgent": dto.UserAgent(),
						"version":   prev.Version,
					})
				}

				c.logger.InfoMsg(ctx, pageconsumerinterface.PagesAreEqualMsg, logger.Fields{
					"url":       cur.URL,
					"userAgent": dto.UserAgent(),
					"version":   cur.Version,
				})
			}
		}(dto)
	}
}
