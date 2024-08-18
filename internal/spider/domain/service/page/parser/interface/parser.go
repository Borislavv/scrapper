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

	EmptyTitleError                = errors.New("parsing page Title from HTML failed due to 'title' is empty")
	EmptyDescriptionError          = errors.New("parsing page Description from HTML failed due to 'description' is empty")
	EmptyCanonicalError            = errors.New("parsing page Canonical from HTML failed due to 'canonical' is empty")
	EmptyH1Error                   = errors.New("parsing page H1 from HTML failed due to 'h1' is empty")
	EmptyHrefLangLinkError         = errors.New("parsing page HrefLang from HTML failed due to 'href' is empty")
	EmptyHrefLangRelError          = errors.New("parsing page HrefLang from HTML failed due to 'rel' is empty")
	EmptyHrefLangError             = errors.New("parsing page HrefLang from HTML failed due to 'hreflang' is empty")
	EmptyRelinkingBlockHrefError   = errors.New("parsing page RelinkingBlock from HTML failed due to 'href' is empty")
	EmptyRelinkingBlockAnchorError = errors.New("parsing page RelinkingBlock from HTML failed due to 'anchor' is empty")
	EmptyRelinkingBlockNameError   = errors.New("parsing page RelinkingBlock from HTML failed due to block name is empty")
	EmptyFAQQuestionError          = errors.New("parsing page FAQ from HTML failed due to 'question' is empty")
	EmptyFAQAnswerError            = errors.New("parsing page FAQ from HTML failed due to 'answer' is empty")

	NotExistsDescriptionError        = errors.New("parsing page Description from HTML failed due to 'description' element is not exists")
	NotExistsCanonicalError          = errors.New("parsing page Canonical from HTML failed due to 'canonical' element is not exists")
	NotExistsHrefLangHrefError       = errors.New("parsing page HrefLang from HTML failed due to 'href' attribute does not exists")
	NotExistsHrefLangRelError        = errors.New("parsing page HrefLang from HTML failed due to 'rel' attribute does not exists")
	NotExistsHrefLangError           = errors.New("parsing page HrefLang from HTML failed due to 'hreflang' attribute does not exists")
	NotExistsRelinkingBlockHrefError = errors.New("parsing page RelinkingBlock from HTML failed due to 'href' attribute does not exists")
)

type PageParser interface {
	Parse(ctx context.Context, resp *http.Response) (*entity.Page, error)
}
