package pageparserinterface

import (
	"github.com/Borislavv/scrapper/internal/shared/domain/entity"
	"net/http"
)

type PageParser interface {
	Parse(resp *http.Response) (*entity.Page, error)
}
