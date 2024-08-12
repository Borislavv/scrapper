package pagesaver

import (
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	entityinterface "github.com/Borislavv/scrapper/internal/spider/domain/entity/interface"
	pagerepositoryinterface "github.com/Borislavv/scrapper/internal/spider/domain/repository/interface"
	"log"
)

type PageSaver struct {
	repository pagerepositoryinterface.PageRepository
}

func New(repository pagerepositoryinterface.PageRepository) *PageSaver {
	return &PageSaver{repository: repository}
}

func (s *PageSaver) Save(page entityinterface.Page) error {
	p, ok := page.(*entity.Page)
	if !ok {
		err := errors.New("unable to cast page by interface to pointer, type assertion failed")
		log.Println("PageSaver: " + err.Error())
		return err
	}

	if err := s.repository.Save(p); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
