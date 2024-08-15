package pageparserinterface

import (
	"context"
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"net/http"
)

type PageParser interface {
	Parse(ctx context.Context, resp *http.Response) (*entity.Page, error)
}
