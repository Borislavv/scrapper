package pagesaverinterface

import (
	"context"
	entityinterface "github.com/Borislavv/scrapper/internal/spider/domain/entity/interface"
)

type PageSaver interface {
	Save(ctx context.Context, page entityinterface.Page) error
}
