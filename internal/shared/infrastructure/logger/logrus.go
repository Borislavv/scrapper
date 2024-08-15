package logger

import (
	"context"
	sharedconfiginterface "github.com/Borislavv/scrapper/internal/shared/app/config/interface"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/util"
	loggerinterface "github.com/Borislavv/scrapper/internal/spider/infrastructure/logger/interface"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Logrus struct {
	logger             *logrus.Logger
	contextExtraFields []string
}

func NewLogrus(cfg sharedconfiginterface.Configurator) (logger *Logrus, cancelFunc func(), err error) {
	l := &Logrus{logger: logrus.New(), contextExtraFields: cfg.GetLoggerContextExtraFields()}

	l.logger.SetLevel(l.getLevel(cfg.GetLoggerLevel()))
	l.logger.SetFormatter(l.getFormat(cfg.GetLoggerFormatter()))
	l.logger.SetReportCaller(cfg.GetLoggerReportCaller())

	output, err := l.getOutput(cfg.GetLoggerOutput())
	if err != nil {
		return nil, nil, err
	}
	l.logger.SetOutput(output)

	return logger, func() { _ = output.Close() }, nil
}

func (l *Logrus) DebugMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.logger.WithFields(l.fieldsFromContext(ctx)).WithFields(fields).Debug(msg)
}

func (l *Logrus) InfoMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.logger.WithFields(l.fieldsFromContext(ctx)).WithFields(fields).Info(msg)
}

func (l *Logrus) WarningMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.logger.WithFields(l.fieldsFromContext(ctx)).WithFields(fields).Warning(msg)
}

func (l *Logrus) ErrorMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.logger.WithFields(l.fieldsFromContext(ctx)).WithFields(fields).Error(msg)
}

func (l *Logrus) FatalMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.logger.WithFields(l.fieldsFromContext(ctx)).WithFields(fields).Fatal(msg)
}

func (l *Logrus) PanicMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.logger.WithFields(l.fieldsFromContext(ctx)).WithFields(fields).Panic(msg)
}

func (l *Logrus) Debug(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.logger.WithFields(l.fieldsFromContext(ctx)).WithFields(fields).Debug(err.Error())
	return err
}

func (l *Logrus) Info(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.logger.WithFields(l.fieldsFromContext(ctx)).WithFields(fields).Info(err.Error())
	return err
}

func (l *Logrus) Warning(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.logger.WithFields(l.fieldsFromContext(ctx)).WithFields(fields).Warning(err.Error())
	return err
}

func (l *Logrus) Error(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.logger.WithFields(l.fieldsFromContext(ctx)).WithFields(fields).Error(err.Error())
	return err
}

func (l *Logrus) Fatal(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.logger.WithFields(l.fieldsFromContext(ctx)).WithFields(fields).Fatal(err.Error())
	return err
}

func (l *Logrus) Panic(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.logger.WithFields(l.fieldsFromContext(ctx)).WithFields(fields).Panic(err.Error())
	return err
}

func (l *Logrus) fieldsFromContext(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{}

	for _, field := range l.contextExtraFields {
		if value := ctx.Value(field); value != nil {
			fields[field] = value
		}
	}

	return fields
}

func (l *Logrus) getOutput(output string) (*os.File, error) {
	if output == "stdout" {
		return os.Stdout, nil
	}

	path := ""
	if output == "" {
		path = "/dev/null"
	} else {
		fpath, err := util.Path(output)
		if err != nil {
			return nil, err
		}
		path = fpath
	}

	if _, err := os.ReadDir(filepath.Dir(path)); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (l *Logrus) getLevel(level string) logrus.Level {
	switch level {
	case "info":
		return logrus.InfoLevel
	case "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}

func (l *Logrus) getFormat(formatter string) logrus.Formatter {
	switch formatter {
	case "text":
		return &logrus.TextFormatter{}
	default:
		return &logrus.JSONFormatter{}
	}
}
