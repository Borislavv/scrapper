package entity

import (
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/vo"
)

type Page struct {
	ID             vo.ID    `bson:",inline"`
	Title          string   `bson:"title"`
	Description    string   `bson:"description"`
	Canonical      string   `bson:"canonical"`
	H1             string   `bson:"H1"`
	PlainText      string   `bson:"plainText"`
	HTML           string   `bson:"html"`
	FAQ            []string `bson:"FAQ"`
	RelinkingBlock []string `bson:"relinkingBlock"`
	HrefLangs      []string `bson:"hrefLangs"`
	logs           []string `bson:"logs"`
	network        []string `bson:"network"`
}

func (p *Page) GetID() vo.ID {
	return p.ID
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

func (p *Page) GetLogs() []string {
	return p.logs
}

func (p *Page) GetNetwork() []string {
	return p.network
}
