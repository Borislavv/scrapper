package pagecomparatorinterface

import (
	"github.com/Borislavv/scrapper/internal/spider/domain/entity/interface"
)

type PageComparator interface {
	IsEquals(prev, cur entityinterface.Page) bool
}
