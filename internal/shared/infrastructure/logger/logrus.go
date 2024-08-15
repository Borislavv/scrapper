package logger

import (
	"context"
	sharedconfiginterface "github.com/Borislavv/scrapper/internal/shared/app/config/interface"
	loggerdto "github.com/Borislavv/scrapper/internal/shared/infrastructure/logger/dto"
	loggerhook "github.com/Borislavv/scrapper/internal/shared/infrastructure/logger/hook"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/util"
	loggerinterface "github.com/Borislavv/scrapper/internal/spider/infrastructure/logger/interface"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

type Logrus struct {
	wg                 *sync.WaitGroup
	logger             *logrus.Logger
	contextExtraFields []string
	msgCh              chan loggerdto.MsgDto
	errCh              chan loggerdto.ErrDto
}

func NewLogrus(cfg sharedconfiginterface.Configurator) (logger *Logrus, cancel loggerinterface.CancelFunc, err error) {
	lgrs := logrus.New()

	if cfg.GetLoggerReportCaller() {
		lgrs.AddHook(loggerhook.CallerHook{Skip: 7})
	}

	l := &Logrus{logger: lgrs, contextExtraFields: cfg.GetLoggerContextExtraFields()}

	l.logger.SetLevel(l.getLevel(cfg.GetLoggerLevel()))
	l.logger.SetFormatter(l.getFormat(cfg.GetLoggerFormatter()))

	output, err := l.getOutput(cfg.GetLoggerOutput())
	if err != nil {
		return nil, nil, err
	}
	l.logger.SetOutput(output)

	l.msgCh = make(chan loggerdto.MsgDto, 1)
	l.errCh = make(chan loggerdto.ErrDto, 1)

	l.wg = &sync.WaitGroup{}
	l.wg.Add(2)
	go l.handleErrors()
	go l.handleMessages()

	return l, func() {
		close(l.msgCh)
		close(l.errCh)
		l.wg.Wait()
		_ = output.Close()
	}, nil
}

func (l *Logrus) handleErrors() {
	defer l.wg.Done()
	for err := range l.errCh {
		l.logger.WithFields(l.fieldsFromContext(err.Ctx)).WithFields(err.Fields).Log(l.getLevel(err.Level), err.Err.Error())
	}
}

func (l *Logrus) handleMessages() {
	defer l.wg.Done()
	for msg := range l.msgCh {
		l.logger.WithFields(l.fieldsFromContext(msg.Ctx)).WithFields(msg.Fields).Log(l.getLevel(msg.Level), msg.Msg)
	}
}

func (l *Logrus) DebugMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.msgCh <- loggerdto.MsgDto{
		Ctx:    ctx,
		Level:  "debug",
		Msg:    msg,
		Fields: fields,
	}
}

func (l *Logrus) InfoMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.msgCh <- loggerdto.MsgDto{
		Ctx:    ctx,
		Level:  "info",
		Msg:    msg,
		Fields: fields,
	}
}

func (l *Logrus) WarningMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.msgCh <- loggerdto.MsgDto{
		Ctx:    ctx,
		Level:  "warning",
		Msg:    msg,
		Fields: fields,
	}
}

func (l *Logrus) ErrorMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.msgCh <- loggerdto.MsgDto{
		Ctx:    ctx,
		Level:  "error",
		Msg:    msg,
		Fields: fields,
	}
}

func (l *Logrus) FatalMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.msgCh <- loggerdto.MsgDto{
		Ctx:    ctx,
		Level:  "fatal",
		Msg:    msg,
		Fields: fields,
	}
}

func (l *Logrus) PanicMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.msgCh <- loggerdto.MsgDto{
		Ctx:    ctx,
		Level:  "panic",
		Msg:    msg,
		Fields: fields,
	}
}

func (l *Logrus) Debug(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.errCh <- loggerdto.ErrDto{
		Ctx:    ctx,
		Level:  "debug",
		Err:    err,
		Fields: fields,
	}
	return err
}

func (l *Logrus) Info(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.errCh <- loggerdto.ErrDto{
		Ctx:    ctx,
		Level:  "info",
		Err:    err,
		Fields: fields,
	}
	return err
}

func (l *Logrus) Warning(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.errCh <- loggerdto.ErrDto{
		Ctx:    ctx,
		Level:  "warning",
		Err:    err,
		Fields: fields,
	}
	return err
}

func (l *Logrus) Error(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.errCh <- loggerdto.ErrDto{
		Ctx:    ctx,
		Level:  "error",
		Err:    err,
		Fields: fields,
	}
	return err
}

func (l *Logrus) Fatal(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.errCh <- loggerdto.ErrDto{
		Ctx:    ctx,
		Level:  "fatal",
		Err:    err,
		Fields: fields,
	}
	return err
}

func (l *Logrus) Panic(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.errCh <- loggerdto.ErrDto{
		Ctx:    ctx,
		Level:  "panic",
		Err:    err,
		Fields: fields,
	}
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
