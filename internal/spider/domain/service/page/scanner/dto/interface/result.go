package scannerdtointerface

import "gitlab.xbet.lan/web-backend/php/spider/internal/shared/domain/entity"

type Result interface {
	URL() string
	UserAgent() string
	Page() *entity.Page
	Error() error
}
