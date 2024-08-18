package pagerepositoryinterface

import (
	"context"
	"errors"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
)

var (
	FindByURLError       = errors.New("searching page by URL failed due to error occurred while executing query")
	UpdateError          = errors.New("updating page failed due to error occurred while executing query")
	UpdateMarshalError   = errors.New("updating page failed due to error occurred while marshaling page into []byte")
	UpdateUnmarshalError = errors.New("updating page failed due to error occurred while unmarshaling []byte into bson.M")
	SaveError            = errors.New("saving page failed due to error occurred while executing query")
	InsertedIDCastError  = errors.New("casting inserted id failed to primitive.ObjectID")
)

type PageRepository interface {
	FindOneLatestByURL(ctx context.Context, url string) (page *entity.Page, found bool, err error)
	Update(ctx context.Context, page *entity.Page) error
	Save(ctx context.Context, page *entity.Page) error
}
