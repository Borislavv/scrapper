package entityinterface

import "github.com/Borislavv/scrapper/internal/shared/infrastructure/vo"

type Page interface {
	GetID() vo.ID
	GetTitle() string
	GetDescription() string
	GetCanonical() string
	GetH1() string
	GetPlainText() string
	GetFAQ() []string
	GetRelinkingBlock() []string
	GetHrefLangs() []string
	GetHTML() string
	GetLogs() []string
	GetNetwork() []string
}
