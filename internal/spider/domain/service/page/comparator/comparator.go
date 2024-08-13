package pagecomparator

import (
	"github.com/Borislavv/scrapper/internal/spider/domain/entity/interface"
)

type PageComparator struct {
}

func New() *PageComparator {
	return &PageComparator{}
}

func (c *PageComparator) IsEquals(prev, cur entityinterface.Page) bool {
	isEqScalar := prev.GetTitle() == cur.GetTitle() &&
		prev.GetDescription() == cur.GetDescription() &&
		prev.GetCanonical() == cur.GetCanonical() &&
		prev.GetH1() == cur.GetH1() &&
		prev.GetPlainText() == cur.GetPlainText()

	if !isEqScalar {
		return false
	}

	if len(prev.GetFAQ()) != len(cur.GetFAQ()) {
		return false
	}
	for i, prevFAQ := range prev.GetFAQ() {
		if prevFAQ != cur.GetFAQ()[i] {
			return false
		}
	}

	if len(prev.GetHrefLangs()) != len(cur.GetHrefLangs()) {
		return false
	}
	for i, prevHrefLang := range prev.GetHrefLangs() {
		if prevHrefLang != cur.GetHrefLangs()[i] {
			return false
		}
	}

	if len(prev.GetRelinkingBlock()) != len(cur.GetRelinkingBlock()) {
		return false
	}
	for i, prevRelinkingBlock := range prev.GetRelinkingBlock() {
		if prevRelinkingBlock != cur.GetRelinkingBlock()[i] {
			return false
		}
	}

	return true
}
