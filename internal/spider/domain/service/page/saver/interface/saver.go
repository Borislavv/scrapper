package pagesaverinterface

import entityinterface "github.com/Borislavv/scrapper/internal/spider/domain/entity/interface"

type PageSaver interface {
	Save(page entityinterface.Page) error
}
