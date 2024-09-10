package entity

import (
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/vo"
	"time"
)

type Page struct {
	vo.ID       `bson:",inline"`
	Version     int    `bson:"version"`
	URL         string `bson:"url"`
	UserAgent   string `bson:"userAgent"`
	Title       string `bson:"title"`
	Description string `bson:"description"`
	Canonical   string `bson:"canonical"`
	H1          string `bson:"h1"`
	PlainText   string `bson:"plainText"`
	HTML        string `bson:"html"`
	// FAQ is a map which contains map[question] = answer
	FAQ map[string]string `bson:"faq"`
	// HrefLangs is a map which contains map[rel + "_" + language] = link
	HrefLangs map[string]string `bson:"hrefLangs"`
	// AlternateMedias is a map which contains map[rel + "_" + media] = link
	AlternateMedias map[string]string `bson:"alternateMedias"`
	// RelinkingBlock is a map of maps which contains map[blockTitle][anchor] = link
	RelinkingBlock map[string]map[string]string `bson:"relinkingBlock"`
	// Headers is a map of slices which contains map[headerName][][headerValue1, headerValue2, ...]
	Headers map[string][]string `bson:"headers"`
	// Reason is a struct which contains the diff. of fields.
	Reason       *Reason `bson:"reason,omitempty"`
	vo.Timestamp `bson:",inline"`
}

func NewPage(url, userAgent string) *Page {
	return &Page{
		URL:       url,
		UserAgent: userAgent,
		Version:   1,
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func (p *Page) GetID() vo.ID {
	return p.ID
}

func (p *Page) GetVersion() int {
	return p.Version
}

func (p *Page) UpVersion(previous *Page) {
	p.Version = previous.GetVersion() + 1
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

func (p *Page) GetFAQ() map[string]string {
	return p.FAQ
}

func (p *Page) GetRelinkingBlock() map[string]map[string]string {
	return p.RelinkingBlock
}

func (p *Page) GetHrefLangs() map[string]string {
	return p.HrefLangs
}

func (p *Page) GetAlternateMedias() map[string]string {
	return p.AlternateMedias
}

func (p *Page) GetHTML() string {
	return p.HTML
}

func (p *Page) GetHeaders() map[string][]string {
	return p.Headers
}

func (p *Page) GetReason() *Reason {
	return p.Reason
}
