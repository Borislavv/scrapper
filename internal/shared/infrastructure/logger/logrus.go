package logger

import (
	"context"
	sharedconfiginterface "github.com/Borislavv/scrapper/internal/shared/app/config/interface"
	loggerinterface "github.com/Borislavv/scrapper/internal/shared/domain/service/logger/interface"
	sharedloggerinterface "github.com/Borislavv/scrapper/internal/shared/infrastructure/logger/config/interface"
	loggerdto "github.com/Borislavv/scrapper/internal/shared/infrastructure/logger/dto"
	loggerenum "github.com/Borislavv/scrapper/internal/shared/infrastructure/logger/enum"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/util"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

type Logrus struct {
	logger             *logrus.Logger
	contextExtraFields []string
	msgCh              chan *loggerdto.MsgDto
	errCh              chan *loggerdto.ErrDto
}

func NewOutput(cfg sharedconfiginterface.Configurator) (out *os.File, cancel func(), err error) {
	out, err = output(cfg.GetLoggerOutput())
	if err != nil {
		return nil, nil, err
	}
	return out, func() {
		_ = out.Close()
	}, err
}

func NewLogrus(
	cfg sharedloggerinterface.Configurator,
	output loggerinterface.Outputer,
) (
	logger *Logrus,
	cancel loggerinterface.CancelFunc,
	err error,
) {
	l := &Logrus{logger: logrus.New(), contextExtraFields: cfg.GetLoggerContextExtraFields()}

	l.logger.SetLevel(l.getLevel(cfg.GetLoggerLevel()))
	l.logger.SetFormatter(l.getFormat(cfg.GetLoggerFormatter()))
	l.logger.SetOutput(output)

	l.msgCh = make(chan *loggerdto.MsgDto, 1)
	l.errCh = make(chan *loggerdto.ErrDto, 1)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go l.handleErrors(wg)
	go l.handleMessages(wg)

	return l, func() {
		close(l.msgCh)
		close(l.errCh)
		wg.Wait()
	}, nil
}

func (l *Logrus) handleErrors(wg *sync.WaitGroup) {
	defer wg.Done()
	for err := range l.errCh {
		l.logger.
			WithFields(err.Fields).
			WithFields(l.fieldsFromContext(err.Ctx)).
			WithFields(err.CallerFields()).
			Log(l.getLevel(err.Level), err.Err.Error())
	}
}

func (l *Logrus) handleMessages(wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range l.msgCh {
		l.logger.
			WithFields(msg.Fields).
			WithFields(l.fieldsFromContext(msg.Ctx)).
			WithFields(msg.CallerFields()).
			Log(l.getLevel(msg.Level), msg.Msg)
	}
}

func (l *Logrus) DebugMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.msgCh <- loggerdto.NewMsg(ctx, loggerenum.DebugLvl, msg, fields)
}

func (l *Logrus) InfoMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.msgCh <- loggerdto.NewMsg(ctx, loggerenum.InfoLvl, msg, fields)
}

func (l *Logrus) WarningMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.msgCh <- loggerdto.NewMsg(ctx, loggerenum.WarningLvl, msg, fields)
}

func (l *Logrus) ErrorMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.msgCh <- loggerdto.NewMsg(ctx, loggerenum.ErrorLvl, msg, fields)
}

func (l *Logrus) FatalMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.msgCh <- loggerdto.NewMsg(ctx, loggerenum.FatalLvl, msg, fields)
}

func (l *Logrus) PanicMsg(ctx context.Context, msg string, fields loggerinterface.Fields) {
	l.msgCh <- loggerdto.NewMsg(ctx, loggerenum.PanicLvl, msg, fields)
}

func (l *Logrus) Debug(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.errCh <- loggerdto.NewErr(ctx, loggerenum.DebugLvl, err, fields)
	return err
}

func (l *Logrus) Info(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.errCh <- loggerdto.NewErr(ctx, loggerenum.InfoLvl, err, fields)
	return err
}

func (l *Logrus) Warning(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.errCh <- loggerdto.NewErr(ctx, loggerenum.WarningLvl, err, fields)
	return err
}

func (l *Logrus) Error(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.errCh <- loggerdto.NewErr(ctx, loggerenum.ErrorLvl, err, fields)
	return err
}

func (l *Logrus) Fatal(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.errCh <- loggerdto.NewErr(ctx, loggerenum.FatalLvl, err, fields)
	return err
}

func (l *Logrus) Panic(ctx context.Context, err error, fields loggerinterface.Fields) error {
	l.errCh <- loggerdto.NewErr(ctx, loggerenum.PanicLvl, err, fields)
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

func output(output string) (*os.File, error) {
	if output == loggerenum.Stdout {
		return os.Stdout, nil
	} else if output == loggerenum.Stderr {
		return os.Stderr, nil
	}

	path := ""
	if output == "" {
		path = loggerenum.DevNull
	} else {
		fpath, err := util.LogsPath(output)
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
	case loggerenum.InfoLvl:
		return logrus.InfoLevel
	case loggerenum.WarningLvl:
		return logrus.WarnLevel
	case loggerenum.ErrorLvl:
		return logrus.ErrorLevel
	case loggerenum.FatalLvl:
		return logrus.FatalLevel
	case loggerenum.PanicLvl:
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}

func (l *Logrus) getFormat(formatter string) logrus.Formatter {
	switch formatter {
	case loggerenum.TextFormat:
		return &logrus.TextFormatter{}
	default:
		return &logrus.JSONFormatter{}
	}
}
