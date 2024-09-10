package pagecomparatorinterface

import (
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"github.com/Borislavv/scrapper/internal/spider/domain/entity/interface"
)

var (
	TitleDoesNotMatchMsg       = "title doesn't match"
	DescriptionDoesNotMatchMsg = "description doesn't match"
	CanonicalDoesNotMatchMsg   = "canonical doesn't match"
	H1DoesNotMatchMsg          = "H1 doesn't match"
	PlainTextDoesNotMatchMsg   = "plain text doesn't match"

	FAQAnswerDoesNotExistsMsg     = "FAQ answer doesn't exists"
	FAQAnswersDoNotMatchMsg       = "FAQ answers don't match"
	FAQWasEnrichedByNewElementMsg = "FAQ was enriched by new element"

	HrefLangLinkDoesNotExistsMsg       = "href lang link doesn't exists"
	HrefLangLinksDoNotMatchMsg         = "href lang links don't match"
	HrefLangWasEnrichedByNewElementMsg = "href lang was enriched by new element"

	AlternateMediaLinkDoesNotExistsMsg       = "alternate media link doesn't exists"
	AlternateMediaLinksDoNotMatchMsg         = "alternate media links don't match"
	AlternateMediaWasEnrichedByNewElementMsg = "alternate media was enriched by new element"

	RelinkingBlockLengthsDoNotMatchMsg = "relinking block lengths don't match"
)

type PageComparator interface {
	IsEquals(cur, prev entityinterface.Page) (bool, *entity.Reason)
}
