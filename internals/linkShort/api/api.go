//go:generate mockgen -destination=mocks/api.go -package=mocks linkShortOzon/internals/myerror MultiLoggerInterface
//go:generate mockgen -destination=mocks/apiApplication.go -package=mocks linkShortOzon/internals/linkShort/application LinkShortAppInterface
//go:generate mockgen -destination=mocks/myErrors.go -package=mocks linkShortOzon/internals/myerror CheckErrorInterface
package api

import (
	"encoding/json"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/valyala/fasthttp"
	"linkShortOzon/internals/linkShort"
	"linkShortOzon/internals/linkShort/application"
	errPkg "linkShortOzon/internals/myerror"
	"linkShortOzon/internals/util"
	"net/http"
)

var UnknownReqId = -1

type LinkShortApiInterface interface {
	CreateLinkShortHandler(ctx *fasthttp.RequestCtx)
	TakeLinkFullHandler(ctx *fasthttp.RequestCtx)
}

type LinkShortApi struct {
	Application application.LinkShortAppInterface
	CheckErrors errPkg.CheckErrorInterface
	Logger      errPkg.MultiLoggerInterface
}

func (l *LinkShortApi) CreateLinkShortHandler(ctx *fasthttp.RequestCtx) {
	reqIdCtx := ctx.UserValue("reqId")
	reqId, errConvert := util.InterfaceConvertInt(reqIdCtx)
	if errConvert != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBody([]byte(errConvert.Error()))
		l.Logger.Errorf("%s", errConvert.Error())
	}
	if reqId != errPkg.IntNil {
		l.CheckErrors.SetRequestIdUser(reqId)
	} else {
		l.CheckErrors.SetRequestIdUser(UnknownReqId)
	}

	var linkFullIn linkShort.LinkFull
	errUnmarshal := json.Unmarshal(ctx.Request.Body(), &linkFullIn)
	if errUnmarshal != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBody([]byte(errPkg.ErrUnmarshal))
		l.Logger.Errorf("%s, %s, requestId: %d", errPkg.ErrUnmarshal, errUnmarshal.Error(), reqId)
		return
	}

	linkShortOut, errIn := l.Application.CreateLinkShortApp(linkFullIn.Link)

	errOut, resultOut, codeHTTP := l.CheckErrors.CheckErrorCreateLinkShort(errIn)
	if errOut != nil {
		switch errOut.Error() {
		case errPkg.ErrMarshal:
			ctx.Response.SetStatusCode(codeHTTP)
			ctx.Response.SetBody([]byte(errPkg.ErrMarshal))
			return
		case errPkg.ErrCheck:
			ctx.Response.SetStatusCode(codeHTTP)
			ctx.Response.SetBody(resultOut)
			return
		}
	}

	request, errResponse := json.Marshal(&util.Result{
		Status: http.StatusCreated,
		Body: linkShort.ResponseLinkShort{
			LinkShort: linkShort.LinkShort{
				Link: linkShortOut,
			},
		},
	})
	if errResponse != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBody([]byte(errPkg.ErrEncode))
		l.Logger.Errorf("%s, %s, requestId: %d", errPkg.ErrEncode, errResponse.Error(), reqId)
		return
	}

	ctx.Response.SetBody(request)
	json.NewEncoder(ctx)
	ctx.Response.SetStatusCode(http.StatusOK)
}

func (l *LinkShortApi) TakeLinkFullHandler(ctx *fasthttp.RequestCtx) {
	reqIdCtx := ctx.UserValue("reqId")
	reqId, errConvert := util.InterfaceConvertInt(reqIdCtx)
	if errConvert != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBody([]byte(errConvert.Error()))
		l.Logger.Errorf("%s", errConvert.Error())
	}

	if reqId != errPkg.IntNil {
		l.CheckErrors.SetRequestIdUser(reqId)
	} else {
		l.CheckErrors.SetRequestIdUser(UnknownReqId)
	}

	var linkShortIn linkShort.LinkShort
	errUnmarshal := json.Unmarshal(ctx.Request.Body(), &linkShortIn)
	if errUnmarshal != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBody([]byte(errPkg.ErrUnmarshal))
		l.Logger.Errorf("%s, %s, requestId: %d", errPkg.ErrUnmarshal, errUnmarshal.Error(), reqId)
		return
	}

	linkFullOut, errIn := l.Application.TakeLinkFullApp(linkShortIn.Link)

	errOut, resultOut, codeHTTP := l.CheckErrors.CheckErrorTakeLinkFull(errIn)
	if errOut != nil {
		switch errOut.Error() {
		case errPkg.ErrMarshal:
			ctx.Response.SetStatusCode(codeHTTP)
			ctx.Response.SetBody([]byte(errPkg.ErrMarshal))
			return
		case errPkg.ErrCheck:
			ctx.Response.SetStatusCode(codeHTTP)
			ctx.Response.SetBody(resultOut)
			return
		}
	}

	request, errResponse := json.Marshal(&util.Result{
		Status: http.StatusCreated,
		Body: linkShort.ResponseLinkFull{
			LinkShort: linkShort.LinkFull{
				Link: linkFullOut,
			},
		},
	})
	if errResponse != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Response.SetBody([]byte(errPkg.ErrEncode))
		l.Logger.Errorf("%s, %s, requestId: %d", errPkg.ErrEncode, errResponse.Error(), reqId)
		return
	}

	ctx.Response.SetBody(request)
	json.NewEncoder(ctx)
	ctx.Response.SetStatusCode(http.StatusOK)

}
