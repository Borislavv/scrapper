package loggerdto

import (
	"context"
	logger "github.com/Borislavv/scrapper/internal/shared/domain/service/logger/interface"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type MsgDto struct {
	Ctx    context.Context
	Fields logger.Fields
	Level  string
	Msg    string
	File   string
	Func   string
	Line   int
}

func NewMsg(ctx context.Context, level string, msg string, fields logger.Fields) *MsgDto {
	var fn string
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		fn = runtime.FuncForPC(pc).Name()
		fn = fn[strings.LastIndex(fn, "/")+1:]
	}

	return &MsgDto{
		Ctx:    ctx,
		Level:  level,
		Msg:    msg,
		Fields: fields,
		File:   file,
		Func:   fn,
		Line:   line,
	}
}

func (m *MsgDto) CallerFields() logrus.Fields {
	return logrus.Fields{
		"file": m.File,
		"func": m.Func,
		"line": m.Line,
	}
}
