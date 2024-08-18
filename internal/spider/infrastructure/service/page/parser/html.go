package pageparser

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"gitlab.xbet.lan/web-backend/php/spider/internal/shared/domain/entity"
	pageparserinterface "gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/service/page/parser/interface"
	"gitlab.xbet.lan/web-backend/php/spider/internal/spider/infrastructure/logger/interface"
	"net/http"
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
			return nil, p.logger.Error(ctx, pageparserinterface.DOMTreeParseError, logger.Fields{
				"url":       resp.Request.URL.String(),
				"userAgent": resp.Request.UserAgent(),
				"err":       err.Error(),
			})
		}
	}

	page := entity.NewPage(resp.Request.URL.String(), resp.Request.UserAgent())

	html, err := doc.Html()
	if err != nil {
		p.logger.ErrorMsg(ctx, pageparserinterface.QueryHTMLError.Error(), logger.Fields{
			"url":       resp.Request.URL.String(),
			"userAgent": resp.Request.UserAgent(),
			"err":       err.Error(),
			"selector":  "html",
		})
	}
	page.HTML = html

	// query for "title"
	page.Title = doc.Find("title").Text()
	if page.Title == "" {
		p.logger.WarningMsg(ctx, pageparserinterface.EmptyTitleError.Error(), logger.Fields{
			"url":       resp.Request.URL.String(),
			"userAgent": resp.Request.UserAgent(),
			"selector":  "title",
		})
	}

	// query for "description"
	description, descriptionExists := doc.Find("meta[name='description']").Attr("content")
	if !descriptionExists {
		p.logger.WarningMsg(ctx, pageparserinterface.NotExistsDescriptionError.Error(), logger.Fields{
			"url":       resp.Request.URL.String(),
			"userAgent": resp.Request.UserAgent(),
			"selector":  "meta[name='description.content",
		})
	} else if description == "" {
		p.logger.WarningMsg(ctx, pageparserinterface.EmptyDescriptionError.Error(), logger.Fields{
			"url":       resp.Request.URL.String(),
			"userAgent": resp.Request.UserAgent(),
			"selector":  "meta[name='description.content",
		})
	}
	page.Description = description

	// query for "canonical"
	canonical, canonicalExists := doc.Find("link[rel='canonical']").Attr("href")
	if !canonicalExists {
		p.logger.WarningMsg(ctx, pageparserinterface.NotExistsCanonicalError.Error(), logger.Fields{
			"url":       resp.Request.URL.String(),
			"userAgent": resp.Request.UserAgent(),
			"selector":  "link[rel='canonical'].href",
		})
	} else if canonical == "" {
		p.logger.WarningMsg(ctx, pageparserinterface.EmptyCanonicalError.Error(), logger.Fields{
			"url":       resp.Request.URL.String(),
			"userAgent": resp.Request.UserAgent(),
			"selector":  "link[rel='canonical'].href",
		})
	}
	page.Canonical = canonical

	// query for "H1"
	page.H1 = doc.Find("h1").First().Text()
	if page.H1 == "" {
		p.logger.WarningMsg(ctx, pageparserinterface.EmptyH1Error.Error(), logger.Fields{
			"url":       resp.Request.URL.String(),
			"userAgent": resp.Request.UserAgent(),
			"selector":  "h1",
		})
	}

	// query for "plaintext" (content)
	plaintext, err := doc.Find(".seo-text").Html()
	if err != nil {
		p.logger.WarningMsg(ctx, pageparserinterface.QueryPlainTextError.Error(), logger.Fields{
			"url":       resp.Request.URL.String(),
			"userAgent": resp.Request.UserAgent(),
			"err":       err.Error(),
			"selector":  ".seo-text",
		})
	}
	page.PlainText = plaintext

	// query for "hreflang"
	links := doc.Find("link[hreflang]")
	hreflangMap := make(map[string]string, links.Length())
	links.Each(func(index int, item *goquery.Selection) {
		link, linkExists := item.Attr("href")
		if !linkExists {
			p.logger.WarningMsg(ctx, pageparserinterface.NotExistsHrefLangHrefError.Error(), logger.Fields{
				"url":       resp.Request.URL.String(),
				"userAgent": resp.Request.UserAgent(),
				"selector":  "link[hreflang][].href",
			})
			return
		} else if link == "" {
			p.logger.WarningMsg(ctx, pageparserinterface.EmptyHrefLangLinkError.Error(), logger.Fields{
				"url":       resp.Request.URL.String(),
				"userAgent": resp.Request.UserAgent(),
				"selector":  "link[hreflang][].href",
			})
			return
		}

		rel, existsRel := item.Attr("rel")
		if !existsRel {
			p.logger.WarningMsg(ctx, pageparserinterface.NotExistsHrefLangRelError.Error(), logger.Fields{
				"url":       resp.Request.URL.String(),
				"userAgent": resp.Request.UserAgent(),
				"selector":  "link[hreflang][].rel",
			})
			return
		} else if rel == "" {
			p.logger.WarningMsg(ctx, pageparserinterface.EmptyHrefLangRelError.Error(), logger.Fields{
				"url":       resp.Request.URL.String(),
				"userAgent": resp.Request.UserAgent(),
				"selector":  "link[hreflang][].rel",
			})
			return
		}

		language, existsLanguage := item.Attr("hreflang")
		if !existsLanguage {
			p.logger.WarningMsg(ctx, pageparserinterface.NotExistsHrefLangError.Error(), logger.Fields{
				"url":       resp.Request.URL.String(),
				"userAgent": resp.Request.UserAgent(),
				"selector":  "link[hreflang][].hreflang",
			})
			return
		} else if language == "" {
			p.logger.WarningMsg(ctx, pageparserinterface.EmptyHrefLangError.Error(), logger.Fields{
				"url":       resp.Request.URL.String(),
				"userAgent": resp.Request.UserAgent(),
				"selector":  "link[hreflang][].hreflang",
			})
			return
		}

		hreflangMap[rel+"_"+language] = link
	})
	page.HrefLangs = hreflangMap

	// query for "relinkingBlock"
	relinkingMap := make(map[string]map[string]string, 3)
	doc.Find("div.seo-list-section").Each(func(i int, section *goquery.Selection) {
		blockTitle := section.Find(".seo-list-section__title").Text()
		if blockTitle == "" {
			p.logger.WarningMsg(ctx, pageparserinterface.EmptyRelinkingBlockNameError.Error(), logger.Fields{
				"url":       resp.Request.URL.String(),
				"userAgent": resp.Request.UserAgent(),
				"block":     blockTitle,
				"selector":  "div.seo-list-section[].seo-list-section__title",
			})
		}

		relinkingLinks := section.Find("li.seo-list-links__item")

		relinkingMap[blockTitle] = make(map[string]string, relinkingLinks.Length())

		relinkingLinks.Each(func(j int, item *goquery.Selection) {
			link := item.Find("a")
			href, exists := link.Attr("href")
			if !exists {
				p.logger.WarningMsg(ctx, pageparserinterface.NotExistsRelinkingBlockHrefError.Error(), logger.Fields{
					"url":       resp.Request.URL.String(),
					"userAgent": resp.Request.UserAgent(),
					"block":     blockTitle,
					"selector":  "div.seo-list-section[].li.seo-list-links__item[].href",
				})
				return
			} else if href == "" {
				p.logger.WarningMsg(ctx, pageparserinterface.EmptyRelinkingBlockHrefError.Error(), logger.Fields{
					"url":       resp.Request.URL.String(),
					"userAgent": resp.Request.UserAgent(),
					"block":     blockTitle,
					"selector":  "div.seo-list-section[].li.seo-list-links__item[].href",
				})
				return
			}

			anchor := link.Text()
			if anchor == "" {
				p.logger.WarningMsg(ctx, pageparserinterface.EmptyRelinkingBlockAnchorError.Error(), logger.Fields{
					"url":       resp.Request.URL.String(),
					"userAgent": resp.Request.UserAgent(),
					"block":     blockTitle,
					"selector":  "div.seo-list-section[].li.seo-list-links__item[].a",
				})
				return
			}

			relinkingMap[blockTitle][anchor] = href
		})
	})
	page.RelinkingBlock = relinkingMap

	// query for "FAQ"
	questions := doc.Find(`div[class="seo-faq__body"]`).Find(`div[itemprop="mainEntity"]`)
	faqMap := make(map[string]string, questions.Length())
	questions.Each(func(i int, div *goquery.Selection) {
		question := div.Find("button").Find(`span[itemprop="name"]`).Text()
		if question == "" {
			p.logger.WarningMsg(ctx, pageparserinterface.EmptyFAQQuestionError.Error(), logger.Fields{
				"url":       resp.Request.URL.String(),
				"userAgent": resp.Request.UserAgent(),
				"selector":  ".seo-faq__body.div[itemprop='mainEntity'][].button.span[itemprop='name']",
			})
			return
		}

		answer := div.Find(`span[itemprop="acceptedAnswer"]`).Find(`span[itemprop="text"]`).Text()
		if answer == "" {
			p.logger.WarningMsg(ctx, pageparserinterface.EmptyFAQAnswerError.Error(), logger.Fields{
				"url":       resp.Request.URL.String(),
				"userAgent": resp.Request.UserAgent(),
				"selector":  ".seo-faq__body.div[itemprop='mainEntity'][].span[itemprop='acceptedAnswer'].span[itemprop='text']",
			})
			return
		}

		faqMap[question] = answer
	})
	page.FAQ = faqMap

	// setup response headers
	page.Headers = resp.Header

	return page, nil
}
