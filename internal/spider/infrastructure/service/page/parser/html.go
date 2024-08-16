package pageparser

import (
	"context"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/vo"
	pageparserinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/parser/interface"
	"github.com/Borislavv/scrapper/internal/spider/infrastructure/logger/interface"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)

type HTML struct {
	logger logger.Logger
}

// NewHTML is a constructor of HTML parser.
func NewHTML(logger logger.Logger) *HTML {
	return &HTML{logger: logger}
}

// Parse is method which parsing a DOM tree for target meta-tags.
func (p *HTML) Parse(ctx context.Context, resp *http.Response) (*entity.Page, error) {
	// create a new DOM tree from reader
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		if !(errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)) {
			return nil, p.logger.Error(ctx, pageparserinterface.ParseDOMFailed, logger.Fields{
				"url":       resp.Request.URL.String(),
				"userAgent": resp.Request.UserAgent(),
				"err":       err.Error(),
			})
		}
	}

	page := &entity.Page{
		URL:       resp.Request.URL.String(),
		UserAgent: resp.Request.UserAgent(),
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// query for "title"
	page.Title = doc.Find("title").Text()

	// query for "description"
	page.Description, _ = doc.Find("meta[name='description']").Attr("content")

	// query for "canonical"
	page.Canonical, _ = doc.Find("link[rel='canonical']").Attr("href")

	// query for "H1"
	page.H1 = doc.Find("h1").First().Text()

	// query for "plaintext" (content)
	plaintext, err := doc.Find(".seo-text").Html()
	if err != nil {
		p.logger.ErrorMsg(ctx, pageparserinterface.ParsePlainTextFailed.Error(), logger.Fields{
			"url":       resp.Request.URL.String(),
			"userAgent": resp.Request.UserAgent(),
			"err":       err.Error(),
		})
	}
	page.PlainText = plaintext

	// setup response headers
	page.Headers = resp.Header

	return page, nil
}
