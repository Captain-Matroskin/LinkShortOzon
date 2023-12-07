package api

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"linkShortOzon/internals/linkShort/api/mocks"
	errPkg "linkShortOzon/internals/myerror"
	"testing"
)

type testApiCreateLinkShortHandler struct {
	testName      string
	reqId         int
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
	args   []interface{}
	format string
	count  int
}

type testAppCreateLinkShort struct {
	in        string
	outResult string
	outErr    error
	count     int
}

type testCheckErrorCreateLSH struct {
	inError     error
	outErr      error
	outResult   []byte
	outCodeHTTP int
	count       int
}

var createLinkShortHandler = []testApiCreateLinkShortHandler{
	{
		testName: "Successful CreateLinkShort handler",
		reqId:    10,
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
			inError:     nil,
			outErr:      nil,
			outResult:   nil,
			outCodeHTTP: errPkg.IntNil,
			count:       1,
		},
		countSetReqId: 1,
		outBody:       []byte("{\"status\":201,\"body\":{\"link_short\":{\"link\":\"hf89h4qwer\"}}}"),
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

		mockCheckError.
			EXPECT().
			SetRequestIdUser(curTest.reqId).
			Times(curTest.countSetReqId)

		mockCheckError.
			EXPECT().
			CheckErrorCreateLinkShort(curTest.checkError.inError).
			Return(curTest.checkError.outErr, curTest.checkError.outResult, curTest.checkError.outCodeHTTP).
			Times(curTest.checkError.count)

		linkShortApi := LinkShortApi{Application: mockApplication, Logger: mockMultiLogger, CheckErrors: mockCheckError}
		t.Run(curTest.testName, func(t *testing.T) {
			linkShortApi.CreateLinkShortHandler(&ctxIn)
			//println(string(ctxIn.Response.Body()))
			require.Equal(
				t,
				ctxExpected.Response.Body(),
				ctxIn.Response.Body(),
				fmt.Sprintf("Expected: %s\nbut got: %s", string(ctxExpected.Response.Body()), string(ctxIn.Response.Body())),
			)
		})
	}
}
