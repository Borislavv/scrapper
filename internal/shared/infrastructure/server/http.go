package sharedserver

import (
	"context"
	"errors"
	logger "github.com/Borislavv/scrapper/internal/shared/domain/service/logger/interface"
	sharedserverconfiginterface "github.com/Borislavv/scrapper/internal/shared/infrastructure/server/config/interface"
	sharedservercontrollerinterface "github.com/Borislavv/scrapper/internal/shared/infrastructure/server/controller/interface"
	sharedservermiddlewareinterface "github.com/Borislavv/scrapper/internal/shared/infrastructure/server/middleware/interface"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"sync"
)

type HTTP struct {
	ctx    context.Context
	logger logger.Logger
	server *fasthttp.Server
	config sharedserverconfiginterface.Configurator
}

func NewHTTP(
	ctx context.Context,
	logger logger.Logger,
	config sharedserverconfiginterface.Configurator,
	controllers []sharedservercontrollerinterface.HttpController,
	middlewares []sharedservermiddlewareinterface.HttpMiddleware,
) *HTTP {
	s := &HTTP{ctx: ctx, logger: logger, config: config}
	s.initServer(s.buildRouter(controllers), middlewares)
	return s
}

func (s *HTTP) ListenAndServe() {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go s.serve(wg)

	wg.Add(1)
	go s.shutdown(wg)
}

func (s *HTTP) serve(wg *sync.WaitGroup) {
	defer wg.Done()

	port := s.config.GetServerPort()

	s.logger.InfoMsg(s.ctx, s.config.GetServerName()+" http server was started", logger.Fields{"port": port})
	defer s.logger.InfoMsg(s.ctx, s.config.GetServerName()+" http server was stopped", logger.Fields{"port": port})

	if err := s.server.ListenAndServe(port); err != nil {
		s.logger.ErrorMsg(s.ctx, err.Error(), logger.Fields{"port": port})
	}
}

func (s *HTTP) shutdown(wg *sync.WaitGroup) {
	defer wg.Done()

	<-s.ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), s.config.GetServerShutDownTimeout())
	defer cancel()

	if err := s.server.ShutdownWithContext(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			s.logger.ErrorMsg(s.ctx, err.Error(), logger.Fields{"port": s.config.GetServerPort()})
		}
		return
	}
}

func (s *HTTP) buildRouter(controllers []sharedservercontrollerinterface.HttpController) *router.Router {
	r := router.New()
	for _, controller := range controllers {
		controller.AddRoute(r)
	}
	return r
}

func (s *HTTP) initServer(r *router.Router, mdws []sharedservermiddlewareinterface.HttpMiddleware) {
	h := r.Handler

	for i := len(mdws) - 1; i >= 0; i-- {
		h = mdws[i].Middleware(h)
	}

	s.server = &fasthttp.Server{Handler: h}
}
