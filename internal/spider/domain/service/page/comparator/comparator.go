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
	return false
}
