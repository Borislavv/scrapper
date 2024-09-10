package entityinterface

import (
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	vointerface "github.com/Borislavv/scrapper/internal/shared/domain/vo/interface"
)

type Page interface {
	vointerface.Identifier
	GetVersion() int
	UpVersion(previous *entity.Page)
	GetURL() string
	GetTitle() string
	GetDescription() string
	GetCanonical() string
	GetH1() string
	GetPlainText() string
	GetHTML() string
	GetFAQ() map[string]string
	GetHrefLangs() map[string]string
	GetAlternateMedias() map[string]string
	GetHeaders() map[string][]string
	GetRelinkingBlock() map[string]map[string]string
	GetReason() *entity.Reason
	vointerface.Timestamper
}
