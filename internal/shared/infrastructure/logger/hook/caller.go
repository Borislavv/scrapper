package loggerhook

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type CallerHook struct {
	Skip int
}

func (h CallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h CallerHook) Fire(entry *logrus.Entry) error {
	pc, file, line, ok := runtime.Caller(h.Skip)
	if !ok {
		return nil
	}

	// Получаем имя функции
	funcName := runtime.FuncForPC(pc).Name()

	// Обрезаем имя функции, чтобы оставить только последнее звено пути
	funcName = funcName[strings.LastIndex(funcName, "/")+1:]

	entry.Data["file"] = file
	entry.Data["line"] = line
	entry.Data["func"] = funcName

	return nil
}
