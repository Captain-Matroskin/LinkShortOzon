package api

import (
	"github.com/valyala/fasthttp"
	errPkg "linkShortOzon/internals/myerror"
	"math"
)

type MiddlewareApiInterface interface {
	LogURL(h fasthttp.RequestHandler) fasthttp.RequestHandler
}

type MiddlewareApi struct {
	ReqId  int
	Logger errPkg.MultiLoggerInterface
}

// middleware, даем request id каждому запросу для логов и в дальнейшем, чтобы можно было отслеживать
func (m *MiddlewareApi) LogURL(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		if m.ReqId == math.MaxInt {
			m.ReqId = 1
		}
		m.ReqId++
		m.Logger.Infof("Method: %s, URL: %s, requestId: %d", string(ctx.Method()), ctx.URI(), m.ReqId)
		ctx.SetUserValue("reqId", m.ReqId)

		h(ctx)
	})
}
