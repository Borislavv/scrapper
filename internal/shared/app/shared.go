package shared

import (
	"context"
	sharedconfig "github.com/Borislavv/scrapper/internal/shared/app/config"
	sharedconfiginterface "github.com/Borislavv/scrapper/internal/shared/app/config/interface"
	loggerinterface "github.com/Borislavv/scrapper/internal/shared/domain/service/logger/interface"
	sharedapicontroller "github.com/Borislavv/scrapper/internal/shared/infrastructure/api/controller"
	liveness "github.com/Borislavv/scrapper/internal/shared/infrastructure/liveness/interface"
	"github.com/Borislavv/scrapper/internal/shared/infrastructure/logger"
	sharedserver "github.com/Borislavv/scrapper/internal/shared/infrastructure/server"
	sharedcontrollerinterface "github.com/Borislavv/scrapper/internal/shared/infrastructure/server/controller/interface"
	sharedserverinterface "github.com/Borislavv/scrapper/internal/shared/infrastructure/server/interface"
	sharedmiddleware "github.com/Borislavv/scrapper/internal/shared/infrastructure/server/middleware"
	sharedmiddlewareinterface "github.com/Borislavv/scrapper/internal/shared/infrastructure/server/middleware/interface"
)

type App struct {
	ctx          context.Context
	cancel       context.CancelFunc
	logger       loggerinterface.Logger
	loggerCancel loggerinterface.CancelFunc
	server       sharedserverinterface.Server
}

func New(ctx context.Context, output loggerinterface.Outputer, liveness liveness.Prober) (*App, error) {
	a := new(App)

	ctx, cancel := context.WithCancel(ctx)
	a.ctx = ctx
	a.cancel = cancel

	cfg, err := sharedconfig.Load()
	if err != nil {
		return nil, err
	}

	lgr, lgrCancel, err := logger.NewLogrus(cfg, output)
	if err != nil {
		return nil, err
	}
	a.logger = lgr
	a.loggerCancel = lgrCancel

	a.server = sharedserver.
		NewHTTP(
			ctx,
			lgr,
			cfg,
			controllers(ctx, lgr, liveness),
			middlewares(ctx, cfg),
		)

	return a, nil
}

func (a *App) Start() {
	defer a.Stop()
	a.logger.InfoMsg(a.ctx, "starting the shared application", nil)

	a.server.ListenAndServe()
}

func (a *App) Stop() {
	a.logger.InfoMsg(a.ctx, "closing the shared application", nil)

	a.cancel()
	a.loggerCancel()
}

// controllers returns a slice of server.HttpController[s] for http server (handlers).
func controllers(
	ctx context.Context,
	logger loggerinterface.Logger,
	liveness liveness.Prober,
) []sharedcontrollerinterface.HttpController {
	return []sharedcontrollerinterface.HttpController{
		sharedapicontroller.NewK8SProbe(ctx, logger, liveness),
	}
}

// middlewares returns a slice of server.HttpMiddleware[s] which will executes in reverse order before handling request.
func middlewares(
	ctx context.Context,
	config sharedconfiginterface.Configurator,
) []sharedmiddlewareinterface.HttpMiddleware {
	return []sharedmiddlewareinterface.HttpMiddleware{
		/** exec 1st. */ sharedmiddleware.NewInitCtxMiddleware(ctx, config),
		/** exec 2nd. */ sharedmiddleware.NewApplicationJsonMiddleware(),
	}
}
