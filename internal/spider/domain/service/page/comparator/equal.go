package pagecomparator

import (
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	fieldenum "github.com/Borislavv/scrapper/internal/shared/domain/enum/field"
	"github.com/Borislavv/scrapper/internal/spider/domain/entity/interface"
	pagecomparatorinterface "github.com/Borislavv/scrapper/internal/spider/domain/service/page/comparator/interface"
	"regexp"
)

const (
	dateRgxPatter   = `(\b\d{4}[-./]\d{2}[-./]\d{2}(?: \d{2}:\d{2})?\b)|(\b\d{2}[-./]\d{2}[-./]\d{4}(?: \d{2}:\d{2})?\b)|(\b\d{2}[-./]\d{2}[-./]\d{2}\b)|(\b\d{4}\b)|(\b\d{2}\b)`
	datePlaceholder = "{{date}}"
)

type Equal struct {
	dateRgx *regexp.Regexp
}

func NewEqual() *Equal {
	return &Equal{dateRgx: regexp.MustCompile(dateRgxPatter)}
}

func (c *Equal) IsEquals(cur, prev entityinterface.Page) (bool, *entity.Reason) {
	fields := make(map[string]*entity.Field)

	c.compareTitles(cur, prev, fields)
	c.compareDescriptions(cur, prev, fields)
	c.compareCanonicals(cur, prev, fields)
	c.compareH1s(cur, prev, fields)
	c.comparePlainTexts(cur, prev, fields)
	c.compareFAQs(cur, prev, fields)
	c.compareHrefLangs(cur, prev, fields)
	c.compareAlternateMedias(cur, prev, fields)
	c.compareRelinkingBlocks(cur, prev, fields)

	return len(fields) == 0, entity.NewReason(prev.GetVersion(), fields)
}

func (c *Equal) compareTitles(cur, prev entityinterface.Page, fields map[string]*entity.Field) {
	curTitle := c.dateRgx.ReplaceAllString(prev.GetTitle(), datePlaceholder)
	prevTitle := c.dateRgx.ReplaceAllString(prev.GetTitle(), datePlaceholder)

	if prevTitle != curTitle {
		fields[fieldenum.Title] = entity.NewReasonField(
			pagecomparatorinterface.TitleDoesNotMatchMsg,
			prev.GetTitle(),
			cur.GetTitle(),
		)
	}
}

func (c *Equal) compareDescriptions(cur, prev entityinterface.Page, fields map[string]*entity.Field) {
	curDescription := c.dateRgx.ReplaceAllString(prev.GetDescription(), datePlaceholder)
	prevDescription := c.dateRgx.ReplaceAllString(prev.GetDescription(), datePlaceholder)

	if prevDescription != curDescription {
		fields[fieldenum.Description] = entity.NewReasonField(
			pagecomparatorinterface.DescriptionDoesNotMatchMsg,
			prev.GetDescription(),
			cur.GetDescription(),
		)
	}
}

func (c *Equal) compareCanonicals(cur, prev entityinterface.Page, fields map[string]*entity.Field) {
	curCanonical := c.dateRgx.ReplaceAllString(prev.GetCanonical(), datePlaceholder)
	prevCanonical := c.dateRgx.ReplaceAllString(prev.GetCanonical(), datePlaceholder)

	if prevCanonical != curCanonical {
		fields[fieldenum.Canonical] = entity.NewReasonField(
			pagecomparatorinterface.CanonicalDoesNotMatchMsg,
			prev.GetCanonical(),
			cur.GetCanonical(),
		)
	}
}

func (c *Equal) compareH1s(cur, prev entityinterface.Page, fields map[string]*entity.Field) {
	curH1 := c.dateRgx.ReplaceAllString(prev.GetH1(), datePlaceholder)
	prevH1 := c.dateRgx.ReplaceAllString(prev.GetH1(), datePlaceholder)

	if prevH1 != curH1 {
		fields[fieldenum.H1] = entity.NewReasonField(
			pagecomparatorinterface.H1DoesNotMatchMsg,
			prev.GetH1(),
			cur.GetH1(),
		)
	}
}

func (c *Equal) comparePlainTexts(cur, prev entityinterface.Page, fields map[string]*entity.Field) {
	curPlainText := c.dateRgx.ReplaceAllString(prev.GetPlainText(), datePlaceholder)
	prevPlainText := c.dateRgx.ReplaceAllString(prev.GetPlainText(), datePlaceholder)

	if prevPlainText != curPlainText {
		fields[fieldenum.PlainText] = entity.NewReasonField(
			pagecomparatorinterface.PlainTextDoesNotMatchMsg,
			prev.GetPlainText(),
			cur.GetPlainText(),
		)
	}
}

func (c *Equal) compareFAQs(cur, prev entityinterface.Page, fields map[string]*entity.Field) {
	// check the previous FAQ block is equals the current
	for prevQuestion, prevAnswer := range prev.GetFAQ() {
		curAnswer, answerExists := cur.GetFAQ()[prevQuestion]
		if !answerExists {
			fields[fieldenum.FAQ+"."+prevQuestion] = entity.NewReasonField(
				pagecomparatorinterface.FAQAnswerDoesNotExistsMsg,
				fieldenum.FAQAnswer+": "+prevAnswer,
				"",
			)
		} else if prevAnswer != curAnswer {
			fields[fieldenum.FAQ+"."+prevQuestion] = entity.NewReasonField(
				pagecomparatorinterface.FAQAnswersDoNotMatchMsg,
				fieldenum.FAQAnswer+": "+prevAnswer,
				fieldenum.FAQAnswer+": "+curAnswer,
			)
		}
	}
	// check wether the current FAQ block is enriched by new elements
	for curQuestion, curAnswer := range cur.GetFAQ() {
		_, prevAnswerExists := prev.GetFAQ()[curQuestion]
		if !prevAnswerExists {
			fields[fieldenum.FAQ+"."+curQuestion] = entity.NewReasonField(
				pagecomparatorinterface.FAQWasEnrichedByNewElementMsg,
				"",
				fieldenum.FAQAnswer+": "+curAnswer,
			)
		}
	}
}

func (c *Equal) compareHrefLangs(cur, prev entityinterface.Page, fields map[string]*entity.Field) {
	// check the previous HrefLang block is equals the current
	for uniq, prevHrefLang := range prev.GetHrefLangs() {
		curHreflang, curHreflangExists := cur.GetHrefLangs()[uniq]
		if !curHreflangExists {
			fields[fieldenum.HrefLangs+"."+uniq] = entity.NewReasonField(
				pagecomparatorinterface.HrefLangLinkDoesNotExistsMsg,
				fieldenum.HrefLangLink+": "+prevHrefLang,
				"",
			)
		} else if prevHrefLang != curHreflang {
			fields[fieldenum.HrefLangs+"."+uniq] = entity.NewReasonField(
				pagecomparatorinterface.HrefLangLinksDoNotMatchMsg,
				fieldenum.HrefLangLink+": "+prevHrefLang,
				fieldenum.HrefLangLink+": "+curHreflang,
			)
		}
	}
	// check wether the current HrefLangs are enriched by new elements
	for uniq, curHrefLang := range cur.GetHrefLangs() {
		_, prevHreflangExists := prev.GetHrefLangs()[uniq]
		if !prevHreflangExists {
			fields[fieldenum.HrefLangs+"."+uniq] = entity.NewReasonField(
				pagecomparatorinterface.HrefLangWasEnrichedByNewElementMsg,
				"",
				fieldenum.HrefLangLink+": "+curHrefLang,
			)
		}
	}
}

func (c *Equal) compareAlternateMedias(cur, prev entityinterface.Page, fields map[string]*entity.Field) {
	// check the previous AlternateMedia block is equals the current
	for uniq, prevAlternateMedia := range prev.GetAlternateMedias() {
		curAlternateMedia, curAlternateMediaExists := cur.GetAlternateMedias()[uniq]
		if !curAlternateMediaExists {
			fields[fieldenum.AlternateMedias+"."+uniq] = entity.NewReasonField(
				pagecomparatorinterface.AlternateMediaLinkDoesNotExistsMsg,
				fieldenum.AlternateMedias+": "+prevAlternateMedia,
				"",
			)
		} else if prevAlternateMedia != curAlternateMedia {
			fields[fieldenum.AlternateMedias+"."+uniq] = entity.NewReasonField(
				pagecomparatorinterface.AlternateMediaLinksDoNotMatchMsg,
				fieldenum.AlternateMedias+": "+prevAlternateMedia,
				fieldenum.AlternateMedias+": "+curAlternateMedia,
			)
		}
	}
	// check wether the current AlternateMedia are enriched by new elements
	for uniq, curAlternateMedia := range cur.GetAlternateMedias() {
		_, prevAlternateMediaExists := prev.GetAlternateMedias()[uniq]
		if !prevAlternateMediaExists {
			fields[fieldenum.AlternateMedias+"."+uniq] = entity.NewReasonField(
				pagecomparatorinterface.AlternateMediaWasEnrichedByNewElementMsg,
				"",
				fieldenum.AlternateMediaLink+": "+curAlternateMedia,
			)
		}
	}
}

func (c *Equal) compareRelinkingBlocks(cur, prev entityinterface.Page, fields map[string]*entity.Field) {
	// check the previous RelinkingBlock block is equals the current
	if len(prev.GetRelinkingBlock()) != len(cur.GetRelinkingBlock()) {
		fields[fieldenum.RelinkingBlock] = entity.NewReasonField(
			pagecomparatorinterface.RelinkingBlockLengthsDoNotMatchMsg,
			len(prev.GetRelinkingBlock()),
			len(cur.GetRelinkingBlock()),
		)
	}
}
