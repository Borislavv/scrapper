package scannerdtointerface

import "github.com/Borislavv/scrapper/internal/shared/domain/entity"

type Result interface {
	URL() string
	Page() *entity.Page
	Error() error
}
