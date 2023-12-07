package api

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"linkShortOzon/internals/linkShort/api/mocks"
	errPkg "linkShortOzon/internals/myerror"
	"net/http"
	"testing"
)

type testApiCreateLinkShortHandler struct {
	testName      string
	reqId         interface{}
	reqIdInt      int
	body          []byte
	logger        testLogger
	App           testAppCreateLinkShort
	checkError    testCheckErrorCreateLSH
	countSetReqId int
	outBody       []byte
}

type testLogger struct {
	errorf loggerErrorf
}

//type loggerWarnf struct {
//	args   []interface{}
//	format string
//	count  int
//}

type loggerErrorf struct {
	format string
	args   []interface{}
	count  int
}

type testAppCreateLinkShort struct {
	in        string
	outResult string
	outErr    error
	count     int
}

type testCheckErrorCreateLSH struct {
	outErr      error
	outResult   []byte
	outCodeHTTP int
	count       int
}

var createLinkShortHandler = []testApiCreateLinkShortHandler{
	{
		testName: "CreateLinkShort handler: successful",
		reqId:    10,
		reqIdInt: 10,
		body:     []byte("{\"link\":\"www.site.ru\"}"),
		logger: testLogger{
			errorf: loggerErrorf{count: 0},
		},
		App: testAppCreateLinkShort{
			in:        "www.site.ru",
			outResult: "hf89h4qwer",
			outErr:    nil,
			count:     1,
		},
		checkError: testCheckErrorCreateLSH{
			outErr:      nil,
			outResult:   nil,
			outCodeHTTP: errPkg.IntNil,
			count:       1,
		},
		countSetReqId: 1,
		outBody:       []byte("{\"status\":201,\"body\":{\"link_short\":{\"link\":\"hf89h4qwer\"}}}"),
	},
	{
		testName: "CreateLinkShort handler: (error) reqId",
		reqId:    nil,
		reqIdInt: 0,
		body:     []byte("{\"link\":\"www.site.ru\"}"),
		logger: testLogger{errorf: loggerErrorf{
			format: "%s",
			args:   []interface{}{"expected type string or int"},
			count:  1,
		}},
		App: testAppCreateLinkShort{
			in:        "www.site.ru",
			outResult: "hf89h4qwer",
			outErr:    nil,
			count:     1,
		},
		checkError: testCheckErrorCreateLSH{
			outErr:      nil,
			outResult:   nil,
			outCodeHTTP: errPkg.IntNil,
			count:       1,
		},
		countSetReqId: 1,
		outBody:       []byte("{\"status\":201,\"body\":{\"link_short\":{\"link\":\"hf89h4qwer\"}}}"),
	},
	{
		testName: "CreateLinkShort handler: (error) unmarshall body",
		reqId:    10,
		reqIdInt: 10,
		body:     []byte("{{{\"link\":\"www.site.ru\"}"),
		logger: testLogger{
			errorf: loggerErrorf{
				format: "%s, %s, requestId: %d",
				args: []interface{}{errPkg.ErrUnmarshal,
					"invalid character '{' looking for beginning of object key string", 10},
				count: 1,
			},
		},
		App: testAppCreateLinkShort{
			count: 0,
		},
		checkError: testCheckErrorCreateLSH{
			count: 0,
		},
		countSetReqId: 1,
		outBody:       []byte(errPkg.ErrUnmarshal),
	},
	{
		testName: "CreateLinkShort handler: (error) CheckErrorCreateLinkShort - errMarshall",
		reqId:    10,
		reqIdInt: 10,
		body:     []byte("{\"link\":\"www.site.ru\"}"),
		logger: testLogger{
			errorf: loggerErrorf{count: 0},
		},
		App: testAppCreateLinkShort{
			in:        "www.site.ru",
			outResult: "",
			outErr:    errors.New(errPkg.LSHCreateLinkShortAppNotGenerate),
			count:     1,
		},
		checkError: testCheckErrorCreateLSH{
			outErr:      errors.New(errPkg.ErrMarshal),
			outResult:   nil,
			outCodeHTTP: http.StatusInternalServerError,
			count:       1,
		},
		countSetReqId: 1,
		outBody:       []byte(errPkg.ErrMarshal),
	},
	{
		testName: "CreateLinkShort handler: (error) CheckErrorCreateLinkShort - errCheck",
		reqId:    10,
		reqIdInt: 10,
		body:     []byte("{\"link\":\"www.site.ru\"}"),
		logger: testLogger{
			errorf: loggerErrorf{count: 0},
		},
		App: testAppCreateLinkShort{
			in:        "www.site.ru",
			outResult: "",
			outErr:    errors.New(errPkg.LSHCreateLinkShortAppNotGenerate),
			count:     1,
		},
		checkError: testCheckErrorCreateLSH{
			outErr:      errors.New(errPkg.ErrCheck),
			outResult:   nil,
			outCodeHTTP: http.StatusInternalServerError,
			count:       1,
		},
		countSetReqId: 1,
		outBody:       nil,
	},
}

func TestCreateLinkShortHandler(t *testing.T) {
	ctrlMultiLogger := gomock.NewController(t)
	ctrlApp := gomock.NewController(t)
	ctrlCheckError := gomock.NewController(t)
	defer ctrlMultiLogger.Finish()
	defer ctrlApp.Finish()
	defer ctrlCheckError.Finish()

	mockMultiLogger := mocks.NewMockMultiLoggerInterface(ctrlMultiLogger)
	mockApplication := mocks.NewMockLinkShortAppInterface(ctrlApp)
	mockCheckError := mocks.NewMockCheckErrorInterface(ctrlCheckError)

	for _, curTest := range createLinkShortHandler {
		ctxIn := fasthttp.RequestCtx{}
		ctxIn.SetUserValue("reqId", curTest.reqId)
		ctxIn.Request.SetBody(curTest.body)
		ctxExpected := fasthttp.RequestCtx{}
		ctxExpected.Response.SetBody(curTest.outBody)
		mockMultiLogger.
			EXPECT().
			Errorf(curTest.logger.errorf.format, curTest.logger.errorf.args).
			Times(curTest.logger.errorf.count)

		mockApplication.
			EXPECT().
			CreateLinkShortApp(curTest.App.in).
			Return(curTest.App.outResult, curTest.App.outErr).
			Times(curTest.App.count)

		if curTest.reqIdInt != errPkg.IntNil {
			mockCheckError.
				EXPECT().
				SetRequestIdUser(curTest.reqId).
				Times(curTest.countSetReqId)
		} else {
			mockCheckError.
				EXPECT().
				SetRequestIdUser(UnknownReqId).
				Times(curTest.countSetReqId)
		}

		mockCheckError.
			EXPECT().
			CheckErrorCreateLinkShort(curTest.App.outErr).
			Return(curTest.checkError.outErr, curTest.checkError.outResult, curTest.checkError.outCodeHTTP).
			Times(curTest.checkError.count)

		linkShortApi := LinkShortApi{Application: mockApplication, Logger: mockMultiLogger, CheckErrors: mockCheckError}
		t.Run(curTest.testName, func(t *testing.T) {
			linkShortApi.CreateLinkShortHandler(&ctxIn)
			require.Equal(
				t,
				ctxExpected.Response.Body(),
				ctxIn.Response.Body(),
				fmt.Sprintf("Expected: %s\nbut got: %s", string(ctxExpected.Response.Body()), string(ctxIn.Response.Body())),
			)
		})
	}
}

type testApiTakeLinkFullHandler struct {
	testName      string
	reqId         interface{}
	reqIdInt      int
	body          []byte
	logger        testLogger
	App           testAppTakeLinkFull
	checkError    testCheckErrorTakeLinkFull
	countSetReqId int
	outBody       []byte
}

type testAppTakeLinkFull struct {
	in        string
	outResult string
	outErr    error
	count     int
}

type testCheckErrorTakeLinkFull struct {
	outErr      error
	outResult   []byte
	outCodeHTTP int
	count       int
}

var takeLinkFullHandler = []testApiTakeLinkFullHandler{
	{
		testName: "TakeLinkFullHandler handler: successful",
		reqId:    10,
		reqIdInt: 10,
		body:     []byte("{\"link\":\"ozon.click.ru/_FeLIUZ33Y\"}"),
		logger: testLogger{
			errorf: loggerErrorf{count: 0},
		},
		App: testAppTakeLinkFull{
			in:        "ozon.click.ru/_FeLIUZ33Y",
			outResult: "www.site.ru",
			outErr:    nil,
			count:     1,
		},
		checkError: testCheckErrorTakeLinkFull{
			outErr:      nil,
			outResult:   nil,
			outCodeHTTP: errPkg.IntNil,
			count:       1,
		},
		countSetReqId: 1,
		outBody:       []byte("{\"status\":201,\"body\":{\"link_full\":{\"link\":\"www.site.ru\"}}}"),
	},
	{
		testName: "TakeLinkFullHandler handler: (error) reqId",
		reqId:    nil,
		reqIdInt: 0,
		body:     []byte("{\"link\":\"ozon.click.ru/_FeLIUZ33Y\"}"),
		logger: testLogger{errorf: loggerErrorf{
			format: "%s",
			args:   []interface{}{"expected type string or int"},
			count:  1,
		}},
		App: testAppTakeLinkFull{
			in:        "ozon.click.ru/_FeLIUZ33Y",
			outResult: "www.site.ru",
			outErr:    nil,
			count:     1,
		},
		checkError: testCheckErrorTakeLinkFull{
			outErr:      nil,
			outResult:   nil,
			outCodeHTTP: errPkg.IntNil,
			count:       1,
		},
		countSetReqId: 1,
		outBody:       []byte("{\"status\":201,\"body\":{\"link_full\":{\"link\":\"www.site.ru\"}}}"),
	},
	{
		testName: "TakeLinkFullHandler handler: (error) unmarshall body",
		reqId:    10,
		reqIdInt: 10,
		body:     []byte("{{{\"link\":\"ozon.click.ru/_FeLIUZ33Y\"}"),
		logger: testLogger{
			errorf: loggerErrorf{
				format: "%s, %s, requestId: %d",
				args: []interface{}{errPkg.ErrUnmarshal,
					"invalid character '{' looking for beginning of object key string", 10},
				count: 1,
			},
		},
		App: testAppTakeLinkFull{
			count: 0,
		},
		checkError: testCheckErrorTakeLinkFull{
			count: 0,
		},
		countSetReqId: 1,
		outBody:       []byte(errPkg.ErrUnmarshal),
	},
	{
		testName: "TakeLinkFullHandler handler: (error) CheckErrorTakeLinkFull - errMarshall",
		reqId:    10,
		reqIdInt: 10,
		body:     []byte("{\"link\":\"ozon.click.ru/_FeLIUZ33Y\"}"),
		logger: testLogger{
			errorf: loggerErrorf{count: 0},
		},
		App: testAppTakeLinkFull{
			in:        "ozon.click.ru/_FeLIUZ33Y",
			outResult: "",
			outErr:    errors.New(errPkg.LSHCreateLinkShortAppNotGenerate),
			count:     1,
		},
		checkError: testCheckErrorTakeLinkFull{
			outErr:      errors.New(errPkg.ErrMarshal),
			outResult:   nil,
			outCodeHTTP: http.StatusInternalServerError,
			count:       1,
		},
		countSetReqId: 1,
		outBody:       []byte(errPkg.ErrMarshal),
	},
	{
		testName: "TakeLinkFullHandler handler: (error) CheckErrorTakeLinkFull - errCheck",
		reqId:    10,
		reqIdInt: 10,
		body:     []byte("{\"link\":\"ozon.click.ru/_FeLIUZ33Y\"}"),
		logger: testLogger{
			errorf: loggerErrorf{count: 0},
		},
		App: testAppTakeLinkFull{
			in:        "ozon.click.ru/_FeLIUZ33Y",
			outResult: "",
			outErr:    errors.New(errPkg.LSHCreateLinkShortAppNotGenerate),
			count:     1,
		},
		checkError: testCheckErrorTakeLinkFull{
			outErr:      errors.New(errPkg.ErrCheck),
			outResult:   nil,
			outCodeHTTP: http.StatusInternalServerError,
			count:       1,
		},
		countSetReqId: 1,
		outBody:       nil,
	},
}

func TestTakeLinkFullHandler(t *testing.T) {
	ctrlMultiLogger := gomock.NewController(t)
	ctrlApp := gomock.NewController(t)
	ctrlCheckError := gomock.NewController(t)
	defer ctrlMultiLogger.Finish()
	defer ctrlApp.Finish()
	defer ctrlCheckError.Finish()

	mockMultiLogger := mocks.NewMockMultiLoggerInterface(ctrlMultiLogger)
	mockApplication := mocks.NewMockLinkShortAppInterface(ctrlApp)
	mockCheckError := mocks.NewMockCheckErrorInterface(ctrlCheckError)

	for _, curTest := range takeLinkFullHandler {
		ctxIn := fasthttp.RequestCtx{}
		ctxIn.SetUserValue("reqId", curTest.reqId)
		ctxIn.Request.SetBody(curTest.body)
		ctxExpected := fasthttp.RequestCtx{}
		ctxExpected.Response.SetBody(curTest.outBody)
		mockMultiLogger.
			EXPECT().
			Errorf(curTest.logger.errorf.format, curTest.logger.errorf.args).
			Times(curTest.logger.errorf.count)

		mockApplication.
			EXPECT().
			TakeLinkFullApp(curTest.App.in).
			Return(curTest.App.outResult, curTest.App.outErr).
			Times(curTest.App.count)

		if curTest.reqIdInt != errPkg.IntNil {
			mockCheckError.
				EXPECT().
				SetRequestIdUser(curTest.reqId).
				Times(curTest.countSetReqId)
		} else {
			mockCheckError.
				EXPECT().
				SetRequestIdUser(UnknownReqId).
				Times(curTest.countSetReqId)
		}

		mockCheckError.
			EXPECT().
			CheckErrorTakeLinkFull(curTest.App.outErr).
			Return(curTest.checkError.outErr, curTest.checkError.outResult, curTest.checkError.outCodeHTTP).
			Times(curTest.checkError.count)

		linkShortApi := LinkShortApi{Application: mockApplication, Logger: mockMultiLogger, CheckErrors: mockCheckError}
		t.Run(curTest.testName, func(t *testing.T) {
			linkShortApi.TakeLinkFullHandler(&ctxIn)
			require.Equal(
				t,
				ctxExpected.Response.Body(),
				ctxIn.Response.Body(),
				fmt.Sprintf("Expected: %s\nbut got: %s", string(ctxExpected.Response.Body()), string(ctxIn.Response.Body())),
			)
		})
	}
}
