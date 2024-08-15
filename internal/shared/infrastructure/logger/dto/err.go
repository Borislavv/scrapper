package loggerdto

import (
	"context"
	logger "github.com/Borislavv/scrapper/internal/spider/infrastructure/logger/interface"
)

type ErrDto struct {
	Ctx    context.Context
	Level  string
	Err    error
	Fields logger.Fields
}
