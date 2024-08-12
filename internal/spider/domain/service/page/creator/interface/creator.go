package pagecreatorinterface

import "github.com/Borislavv/scrapper/internal/shared/domain/entity"

type PageCreator interface {
	Create(data map[string]interface{}) *entity.Page
}
