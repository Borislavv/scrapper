package entity

import (
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/vo"
)

type Page struct {
	vo.ID          `bson:",inline"`
	URL            string              `bson:"url"`
	Title          string              `bson:"title"`
	Description    string              `bson:"description"`
	Canonical      string              `bson:"canonical"`
	H1             string              `bson:"h1"`
	PlainText      string              `bson:"plainText"`
	HTML           string              `bson:"html"`
	FAQ            []string            `bson:"faq"`
	RelinkingBlock []string            `bson:"relinkingBlock"`
	HrefLangs      []string            `bson:"hrefLangs"`
	Headers        map[string][]string `bson:"headers"`
	vo.Timestamp   `bson:",inline"`
}

func (p *Page) GetID() vo.ID {
	return p.ID
}

func (p *Page) GetURL() string {
	return p.URL
}

func (p *Page) GetTitle() string {
	return p.Title
}

func (p *Page) GetDescription() string {
	return p.Description
}

func (p *Page) GetCanonical() string {
	return p.Canonical
}

func (p *Page) GetH1() string {
	return p.H1
}

func (p *Page) GetPlainText() string {
	return p.PlainText
}

func (p *Page) GetFAQ() []string {
	return p.FAQ
}

func (p *Page) GetRelinkingBlock() []string {
	return p.RelinkingBlock
}

func (p *Page) GetHrefLangs() []string {
	return p.HrefLangs
}

func (p *Page) GetHTML() string {
	return p.HTML
}

func (p *Page) GetHeaders() map[string][]string {
	return p.Headers
}
