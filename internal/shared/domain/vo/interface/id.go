package vointerface

import "github.com/Borislavv/scrapper/internal/shared/infrastructure/vo"

type Identifier interface {
	GetID() vo.ID
}
