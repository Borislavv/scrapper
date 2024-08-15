package pageparser

import (
	"context"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/vo"
	pagescannerinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/scanner/interface"
	loggerinterface "github.com/Borislavv/scrapper/internal/spider/infrastructure/logger/interface"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)

type HTML struct {
	logger loggerinterface.Logger
}

// NewHTML is a constructor of HTML parser.
func NewHTML(logger loggerinterface.Logger) *HTML {
	return &HTML{logger: logger}
}

// Parse is method which parsing a DOM tree for target meta-tags.
func (p *HTML) Parse(ctx context.Context, resp *http.Response) (*entity.Page, error) {
	// create a new DOM tree from reader
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		if !(errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)) {
			return nil, p.logger.Error(ctx, pagescannerinterface.ParserError, nil)
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

	// query for "content"
	page.PlainText = doc.Find("seo-text").First().Text()

	// setup response headers
	page.Headers = resp.Header

	return page, nil
}
