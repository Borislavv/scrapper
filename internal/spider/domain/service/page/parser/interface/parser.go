package pageparserinterface

import (
	"context"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"net/http"
)

var (
	DOMTreeParseError = errors.New("parsing page from HTML failed due to error occurred while building a DOM tree and reading it")

	QueryHTMLError      = errors.New("parsing page HTML failed due to error occurred while executing query")
	QueryPlainTextError = errors.New("parsing page PlainText from HTML failed due to error occurred while executing query")

	EmptyHrefLangLinkError         = errors.New("parsing page HrefLang from HTML failed due to 'href' is empty")
	EmptyHrefLangRelError          = errors.New("parsing page HrefLang from HTML failed due to 'rel' is empty")
	EmptyHrefLangError             = errors.New("parsing page HrefLang from HTML failed due to 'hreflang' is empty")
	EmptyRelinkingBlockHrefError   = errors.New("parsing page RelinkingBlock from HTML failed due to 'href' is empty")
	EmptyRelinkingBlockAnchorError = errors.New("parsing page RelinkingBlock from HTML failed due to 'anchor' is empty")
	EmptyFAQQuestionError          = errors.New("parsing page FAQ from HTML failed due to 'question' is empty")
	EmptyFAQAnswerError            = errors.New("parsing page FAQ from HTML failed due to 'answer' is empty")
	EmptyAlternateMediaError       = errors.New("parsing page AlternateMedia from HTML failed due to 'hreflang' is empty")
	EmptyAlternateMediaRelError    = errors.New("parsing page AlternateMedia from HTML failed due to 'rel' is empty")
	EmptyAlternateMediaLinkError   = errors.New("parsing page AlternateMedia from HTML failed due to 'href' is empty")

	NotExistsDescriptionError        = errors.New("parsing page Description from HTML failed due to 'description' element is not exists")
	NotExistsCanonicalError          = errors.New("parsing page Canonical from HTML failed due to 'canonical' element is not exists")
	NotExistsHrefLangHrefError       = errors.New("parsing page HrefLang from HTML failed due to 'href' attribute does not exists")
	NotExistsHrefLangRelError        = errors.New("parsing page HrefLang from HTML failed due to 'rel' attribute does not exists")
	NotExistsHrefLangError           = errors.New("parsing page HrefLang from HTML failed due to 'hreflang' attribute does not exists")
	NotExistsRelinkingBlockHrefError = errors.New("parsing page RelinkingBlock from HTML failed due to 'href' attribute does not exists")
	NotExistsAlternateMediaError     = errors.New("parsing page AlternateMedia from HTML failed due to 'media' attribute does not exists")
	NotExistsAlternateMediaHrefError = errors.New("parsing page AlternateMedia from HTML failed due to 'href' attribute does not exists")
	NotExistsAlternateMediaRelError  = errors.New("parsing page AlternateMedia from HTML failed due to 'rel' attribute does not exists")
)

type PageParser interface {
	Parse(ctx context.Context, resp *http.Response) (*entity.Page, error)
}
