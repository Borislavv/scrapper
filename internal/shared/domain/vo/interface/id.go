package vointerface

import "gitlab.xbet.lan/web-backend/php/spider/internal/shared/infrastructure/vo"

type Identifier interface {
	GetID() vo.ID
}
