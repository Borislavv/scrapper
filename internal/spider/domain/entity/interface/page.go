package entityinterface

import (
	vointerface "github.com/Borislavv/scrapper/internal/shared/domain/vo/interface"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/vo"
)

type Page interface {
	GetID() vo.ID
	GetURL() string
	GetTitle() string
	GetDescription() string
	GetCanonical() string
	GetH1() string
	GetPlainText() string
	GetFAQ() []string
	GetRelinkingBlock() []string
	GetHrefLangs() []string
	GetHTML() string
	GetHeaders() map[string][]string
	vointerface.Timestamper
}
