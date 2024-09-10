package sharedmiddlewareinterface

import (
	"github.com/valyala/fasthttp"
)

type HttpMiddleware interface {
	Middleware(handler fasthttp.RequestHandler) fasthttp.RequestHandler
}
