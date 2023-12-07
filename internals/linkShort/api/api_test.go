package api

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"linkShortOzon/internals/linkShort/api/mocks"
	errPkg "linkShortOzon/internals/myerror"
	"testing"
)

var CreateLinkShortHandler = []struct {
	testName                string
	inputValueReqId         interface{}
	inputValueUnmarshal     []byte
	out                     []byte
	inputErrorfArgs         []interface{}
	inputErrorfFormat       string
	countErrorf             int
	inputWarnfArgs          []interface{}
	inputWarnfFormat        string
	countWarnf              int
	inputCreateLinkShortApp string
	outCreatLinkShortApp    string
	errCreate               error
	countCreate             int
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
		errCreate:               nil,
		countCreate:             1,
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
		errCreate:               nil,
		countCreate:             1,
	}, {
		testName:            "Error unmarshal ",
		inputValueReqId:     "1",
		inputValueUnmarshal: []byte("{{\"link\":\"www.mail.ru\"}"),
		out:                 []byte(errPkg.ErrUnmarshal),
		inputErrorfArgs:     []interface{}{errPkg.ErrUnmarshal, "invalid character '{' looking for beginning of object key string", 1},
		inputErrorfFormat:   "%s, %s, requestId: %d",
		countErrorf:         1,
		countWarnf:          0,
		errCreate:           nil,
		countCreate:         0,
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
		errCreate:               errors.New(errPkg.LSHCreateLinkShortNotInsertUnique),
		countCreate:             1,
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
		errCreate:               errors.New(errPkg.LSHCreateLinkShortTransactionNotCreate),
		countCreate:             1,
	},
}

func TestCreateLinkShortHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrlApp := gomock.NewController(t)
	defer ctrl.Finish()
	defer ctrlApp.Finish()

	mockMultilogger := mocks.NewMockMultiLoggerInterface(ctrl)
	mockApplication := mocks.NewMockLinkShortAppInterface(ctrlApp)
	for _, tt := range CreateLinkShortHandler {
		ctxIn := fasthttp.RequestCtx{}
		ctxIn.SetUserValue("reqId", tt.inputValueReqId)
		ctxIn.Request.SetBody(tt.inputValueUnmarshal)
		ctxExpected := fasthttp.RequestCtx{}
		ctxExpected.Response.SetBody(tt.out)
		mockMultilogger.
			EXPECT().
			Errorf(tt.inputErrorfFormat, tt.inputErrorfArgs).
			Times(tt.countErrorf)

		mockMultilogger.
			EXPECT().
			Warnf(tt.inputWarnfFormat, tt.inputWarnfArgs).
			Times(tt.countWarnf)

		mockApplication.
			EXPECT().
			CreateLinkShortApp(tt.inputCreateLinkShortApp).
			Return(tt.outCreatLinkShortApp, tt.errCreate).
			Times(tt.countCreate)

		linkShortApi := LinkShortApi{Application: mockApplication, Logger: mockMultilogger}
		t.Run(tt.testName, func(t *testing.T) {
			linkShortApi.CreateLinkShortHandler(&ctxIn)
			//println(string(ctxIn.Response.Body()))
			require.Equal(
				t,
				ctxExpected.Response.Body(),
				ctxIn.Response.Body(),
				fmt.Sprintf("Expected: %v\nbut got: %v", ctxExpected.Response.Body(), ctxIn.Response.Body()),
			)
		})
	}
}
