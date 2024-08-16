package pageconsumer

import (
	"context"
	pagerepositoryinterface "github.com/Borislavv/scrapper/internal/spider/domain/repository/interface"
	pagecomparatorinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/comparator/interface"
	pageconsumerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/consumer/interface"
	scannerdtointerface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/scanner/dto/interface"
	logger "github.com/Borislavv/scrapper/internal/spider/infrastructure/logger/interface"
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

			if dto.Error() != nil {
				c.logger.ErrorMsg(ctx, pageconsumerinterface.ScanURLError.Error(), logger.Fields{
					"url":       dto.URL(),
					"err":       dto.Error(),
					"userAgent": dto.UserAgent(),
				})
				return
			}

			cur := dto.Page()

			prev, found, err := c.repository.FindByURL(ctx, cur.URL)
			if err != nil {
				c.logger.ErrorMsg(ctx, pageconsumerinterface.FindPageError.Error(), logger.Fields{
					"url":       dto.URL(),
					"err":       dto.Error(),
					"userAgent": dto.UserAgent(),
				})
				return
			} else if !found {
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

			if !c.comparator.IsEquals(cur, prev) {
				if err = c.repository.Save(ctx, cur); err != nil {
					c.logger.ErrorMsg(ctx, pageconsumerinterface.SavePageError.Error(), logger.Fields{
						"url":       dto.URL(),
						"err":       err.Error(),
						"userAgent": dto.UserAgent(),
					})
				}
			} else {
				c.logger.InfoMsg(ctx, pageconsumerinterface.PagesAreEqualMsg, logger.Fields{
					"url":       cur.URL,
					"userAgent": dto.UserAgent(),
				})
			}
		}(dto)
	}
}
