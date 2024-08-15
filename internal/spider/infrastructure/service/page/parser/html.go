package pageparser

import (
	"context"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

type HTMLParser struct {
}

func NewHTML() *HTMLParser {
	return &HTMLParser{}
}

// Parse is method which parsing a DOM tree for target meta-tags.
func (p *HTMLParser) Parse(resp *http.Response) (*entity.Page, error) {
	// create a new DOM tree from reader
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		if !(errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)) {
			return nil, err
		}
	}

	page := &entity.Page{}

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
