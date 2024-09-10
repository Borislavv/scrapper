package sharedapicontroller

import (
	"context"
	"encoding/json"
	"github.com/Borislavv/scrapper/internal/shared/domain/service/logger/interface"
	liveness "github.com/Borislavv/scrapper/internal/shared/infrastructure/liveness/interface"
	sharedrequestctxenum "github.com/Borislavv/scrapper/internal/shared/infrastructure/server/request/ctx"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

const K8SProbeGetPath = "/k8s/probe"

type K8SProbe struct {
	ctx      context.Context
	logger   logger.Logger
	liveness liveness.Prober
}

func NewK8SProbe(ctx context.Context, logger logger.Logger, liveness liveness.Prober) *K8SProbe {
	return &K8SProbe{ctx: ctx, logger: logger, liveness: liveness}
}

func (c *K8SProbe) Probe(ctx *fasthttp.RequestCtx) {
	reqCtx, ok := ctx.UserValue(sharedrequestctxenum.CtxKey).(context.Context)
	if !ok {
		c.logger.ErrorMsg(c.ctx, "context.Context is not exists into the fasthttp.RequestCtx "+
			"(unable to handle request)", nil)
		return
	}

	isAlive := c.liveness.IsAlive()

	resp := make(map[string]map[string]bool, 1)
	resp["data"] = make(map[string]bool, 1)
	resp["data"]["success"] = isAlive

	b, err := json.Marshal(resp)
	if err != nil {
		c.logger.ErrorMsg(reqCtx, "unable to handle request,"+
			" error occurred while marshaling data into []byte", nil)
		return
	}

	if _, err = ctx.Write(b); err != nil {
		c.logger.ErrorMsg(reqCtx, "unable to handle request,"+
			" error occurred while writing data into *fasthttp.RequestCtx", nil)
		return
	}

	if !isAlive {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
	}
}

func (c *K8SProbe) AddRoute(router *router.Router) {
	router.GET(K8SProbeGetPath, c.Probe)
}
