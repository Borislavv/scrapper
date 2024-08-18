package pagecomparatorinterface

import (
	"gitlab.xbet.lan/web-backend/php/spider/internal/spider/domain/entity/interface"
)

type PageComparator interface {
	IsEquals(prev, cur entityinterface.Page) bool
}
