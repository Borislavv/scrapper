package pagecomparator

import (
	"gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/entity/interface"
)

type Equal struct {
}

func NewEqual() *Equal {
	return &Equal{}
}

func (c *Equal) IsEquals(prev, cur entityinterface.Page) bool {
	isEqScalar := prev.GetTitle() == cur.GetTitle() &&
		prev.GetDescription() == cur.GetDescription() &&
		prev.GetCanonical() == cur.GetCanonical() &&
		prev.GetH1() == cur.GetH1() &&
		prev.GetPlainText() == cur.GetPlainText()

	if !isEqScalar {
		return false
	}

	// compare number of questions
	if len(prev.GetFAQ()) != len(cur.GetFAQ()) {
		return false
	}
	for prevQuestion, prevAnswer := range prev.GetFAQ() {
		curAnswer, answerExists := cur.GetFAQ()[prevQuestion]
		if !answerExists {
			return false
		}
		if prevAnswer != curAnswer {
			return false
		}
	}

	// compare number of links
	if len(prev.GetHrefLangs()) != len(cur.GetHrefLangs()) {
		return false
	}
	for uniq, prevHrefLang := range prev.GetHrefLangs() {
		curHreflang, curHreflangExists := cur.GetHrefLangs()[uniq]
		if !curHreflangExists {
			return false
		}
		if prevHrefLang != curHreflang {
			return false
		}
	}

	// compare number of blocks
	if len(prev.GetRelinkingBlock()) != len(cur.GetRelinkingBlock()) {
		return false
	}
	for prevTitle, prevBlockMap := range prev.GetRelinkingBlock() {
		curBlockMap, curBlockExists := cur.GetRelinkingBlock()[prevTitle]
		if !curBlockExists {
			return false
		}
		if len(prevBlockMap) != len(curBlockMap) {
			return false
		}

		for prevAnchor, prevLink := range prevBlockMap {
			curLink, curLinkExists := curBlockMap[prevAnchor]
			if !curLinkExists {
				return false
			}
			if prevLink != curLink {
				return false
			}
		}
	}

	return true
}
