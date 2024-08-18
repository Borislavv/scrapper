package entityinterface

import (
	"gitlab.xbet.lan/web-backend/php/spider/internal/shared/domain/entity"
	vointerface "gitlab.xbet.lan/web-backend/php/spider/internal/shared/domain/vo/interface"
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
	GetHeaders() map[string][]string
	GetRelinkingBlock() map[string]map[string]string
	vointerface.Timestamper
}
