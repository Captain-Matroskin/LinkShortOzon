package api

import (
	"github.com/valyala/fasthttp"
	"linkShortOzon/internals/linkShort/application"
	errPkg "linkShortOzon/internals/myerror"
)

type LinkShortApiInterface interface {
	CreateLinkShortHandler(ctx *fasthttp.RequestCtx)
	TakeLinkShortHandler(ctx *fasthttp.RequestCtx)
}

type LinkShortApi struct {
	Application application.LinkShortAppInterface
	Logger      errPkg.MultiLoggerInterface
}

func (l *LinkShortApi) CreateLinkShortHandler(ctx *fasthttp.RequestCtx) {

}

func (l *LinkShortApi) TakeLinkShortHandler(ctx *fasthttp.RequestCtx) {

}
