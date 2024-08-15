package loggerdto

import (
	"context"
	logger "github.com/Borislavv/scrapper/internal/spider/infrastructure/logger/interface"
)

type MsgDto struct {
	Ctx    context.Context
	Level  string
	Msg    string
	Fields logger.Fields
}
