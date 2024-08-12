package pagecreator

import "github.com/Borislavv/scrapper/internal/shared/domain/entity"

type PageCreator struct {
}

func New() *PageCreator {
	return &PageCreator{}
}

func (c *PageCreator) Create(data map[string]interface{}) *entity.Page {
	return &entity.Page{}
}
