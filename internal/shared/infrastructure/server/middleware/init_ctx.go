package sharedmiddleware

import (
	"context"
	sharedserverconfiginterface "github.com/Borislavv/scrapper/internal/shared/infrastructure/server/config/interface"
	sharedrequestctxenum "github.com/Borislavv/scrapper/internal/shared/infrastructure/server/request/ctx"
	"github.com/valyala/fasthttp"
)

type InitCtxMiddleware struct {
	ctx    context.Context
	config sharedserverconfiginterface.Configurator
}

func NewInitCtxMiddleware(ctx context.Context, config sharedserverconfiginterface.Configurator) *InitCtxMiddleware {
	return &InitCtxMiddleware{ctx: ctx, config: config}
}

func (m *InitCtxMiddleware) Middleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		reqCtx, reqCtxCancel := context.WithTimeout(m.ctx, m.config.GetServerRequestTimeout())
		_ = reqCtxCancel

		ctx.SetUserValue(sharedrequestctxenum.CtxKey, reqCtx)
		ctx.SetUserValue(sharedrequestctxenum.CtxCancelKey, reqCtxCancel)

		next(ctx)
	}
}
