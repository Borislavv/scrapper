package entityinterface

import "github.com/Borislavv/scrapper/internal/shared/domain/entity"

type Reason interface {
	GetComparedVersion() int
	GetFields() map[string]*entity.Field
}
