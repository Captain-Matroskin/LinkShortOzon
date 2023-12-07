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

var oldCreateLinkShortHandler = []struct {
	testName            string
	inputValueReqId     interface{}
	inputValueUnmarshal []byte
	out                 []byte
	//errorf
	inputErrorfArgs   []interface{}
	inputErrorfFormat string
	countErrorf       int
	//warnf
	inputWarnfArgs   []interface{}
	inputWarnfFormat string
	countWarnf       int
	//app
	inputCreateLinkShortApp string
	outCreatLinkShortApp    string
	//CheckErrors
	errCreateLSHApp      error
	countCreateLSHApp    int
	inputCheckError      error
	outCheckErrTypeError error
	//outCheckErrResultString string
	outCheckErrResultBytes []byte
	outCheckErrCodeHTTP    error
}{
	{
		testName:                "Successful CreateLinkShort handler",
		inputValueReqId:         10,
		inputValueUnmarshal:     []byte("{\"link\":\"www.mail.ru\"}"),
		out:                     []byte("{\"status\":201,\"body\":{\"link_short\":{\"link\":\"hf89h4qwer\"}}}"),
		countErrorf:             0,
		countWarnf:              0,
		inputCreateLinkShortApp: "www.mail.ru",
		outCreatLinkShortApp:    "hf89h4qwer",
		errCreateLSHApp:         nil,
		countCreateLSHApp:       1,
	},
	{
		testName:                "Error reqId ",
		inputValueReqId:         nil,
		inputValueUnmarshal:     []byte("{\"link\":\"www.mail.ru\"}"),
		out:                     []byte("{\"status\":201,\"body\":{\"link_short\":{\"link\":\"hf89h4qwer\"}}}"),
		inputErrorfArgs:         []interface{}{"expected type string or int"},
		inputErrorfFormat:       "%s",
		countErrorf:             1,
		countWarnf:              0,
		inputCreateLinkShortApp: "www.mail.ru",
		outCreatLinkShortApp:    "hf89h4qwer",
		errCreateLSHApp:         nil,
		countCreateLSHApp:       1,
	}, {
		testName:            "Error unmarshal ",
		inputValueReqId:     "1",
		inputValueUnmarshal: []byte("{{\"link\":\"www.mail.ru\"}"),
		out:                 []byte(errPkg.ErrUnmarshal),
		inputErrorfArgs:     []interface{}{errPkg.ErrUnmarshal, "invalid character '{' looking for beginning of object key string", 1},
		inputErrorfFormat:   "%s, %s, requestId: %d",
		countErrorf:         1,
		countWarnf:          0,
		errCreateLSHApp:     nil,
		countCreateLSHApp:   0,
	}, {
		testName:                "Error checkError warnf ",
		inputValueReqId:         "1",
		inputValueUnmarshal:     []byte("{\"link\":\"www.mail.ru\"}"),
		out:                     []byte("{\"status\":409,\"explain\":\"link is not unique CreateLinkShortPostgres\"}"),
		countErrorf:             0,
		inputWarnfArgs:          []interface{}{"link is not unique CreateLinkShortPostgres", 1},
		inputWarnfFormat:        "%s, requestId: %d",
		countWarnf:              1,
		inputCreateLinkShortApp: "www.mail.ru",
		outCreatLinkShortApp:    "",
		errCreateLSHApp:         errors.New(errPkg.LSHCreateLinkShortNotInsertUnique),
		countCreateLSHApp:       1,
	}, {
		testName:                "Error checkError errorf ",
		inputValueReqId:         "1",
		inputValueUnmarshal:     []byte("{\"link\":\"www.mail.ru\"}"),
		out:                     []byte("{\"status\":500,\"explain\":\"" + errPkg.ErrDB + "\"}"),
		inputErrorfArgs:         []interface{}{"transaction Create Link Short not create CreateLinkShortPostgres", 1},
		inputErrorfFormat:       "%s, requestId: %d",
		countErrorf:             1,
		countWarnf:              0,
		inputCreateLinkShortApp: "www.mail.ru",
		outCreatLinkShortApp:    "",
		errCreateLSHApp:         errors.New(errPkg.LSHCreateLinkShortTransactionNotCreate),
		countCreateLSHApp:       1,
	},
}
