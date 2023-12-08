package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"linkShortOzon/internals/linkShort/api/mocks"
	errPkg "linkShortOzon/internals/myerror"
	"linkShortOzon/internals/proto"
	"net/http"
	"testing"
)

type testApiCreateLinkShortGrpc struct {
	testName     string
	inData       *proto.LinkFull
	inCtx        context.Context
	App          testAppCreateLinkShort
	checkError   testCheckErrorCreateLSHGrpc
	dataExpected *proto.ResultLinkShort
}

type testCheckErrorCreateLSHGrpc struct {
	outErr      error
	outResult   string
	outCodeHTTP int
	count       int
}

var createLinkShortGrpc = []testApiCreateLinkShortGrpc{
	{
		testName: "CreateLinkShort Grpc: successful",
		inData:   &proto.LinkFull{LinkFull: "www.site.ru"},
		inCtx:    context.Background(),
		App: testAppCreateLinkShort{
			outErr:    nil,
			outResult: "hf89h4qwer",
			count:     1,
		},
		checkError: testCheckErrorCreateLSHGrpc{
			outErr:      nil,
			outResult:   "",
			outCodeHTTP: errPkg.IntNil,
			count:       1,
		},
		dataExpected: &proto.ResultLinkShort{
			StatusCode: http.StatusOK,
			Body:       &proto.LinkShort{LinkShort: "hf89h4qwer"},
			Error:      "",
		},
	},
	{
		testName: "CreateLinkShort Grpc: checkError - errCheck",
		inData:   &proto.LinkFull{LinkFull: "www.site.ru"},
		inCtx:    context.Background(),
		App: testAppCreateLinkShort{
			outErr:    errors.New(errPkg.LSHCreateLinkShortNotInsertUnique),
			outResult: "",
			count:     1,
		},
		checkError: testCheckErrorCreateLSHGrpc{
			outErr:      errors.New(errPkg.ErrCheck),
			outResult:   errPkg.ErrCheck,
			outCodeHTTP: http.StatusConflict,
			count:       1,
		},
		dataExpected: &proto.ResultLinkShort{
			StatusCode: http.StatusConflict,
			Error:      errPkg.ErrCheck,
		},
	},
}

func TestCreateLinkShortGrpc(t *testing.T) {
	ctrlApp := gomock.NewController(t)
	ctrlCheckError := gomock.NewController(t)
	defer ctrlApp.Finish()
	defer ctrlCheckError.Finish()

	mockApplication := mocks.NewMockLinkShortAppInterface(ctrlApp)
	mockCheckError := mocks.NewMockCheckErrorInterface(ctrlCheckError)

	for _, curTest := range createLinkShortGrpc {
		mockApplication.
			EXPECT().
			CreateLinkShortApp(curTest.inData.LinkFull).
			Return(curTest.App.outResult, curTest.App.outErr).
			Times(curTest.App.count)

		mockCheckError.
			EXPECT().
			CheckErrorCreateLinkShortGrpc(curTest.App.outErr).
			Return(curTest.checkError.outErr, curTest.checkError.outResult, curTest.checkError.outCodeHTTP).
			Times(curTest.checkError.count)

		linkShortApi := LinkShortManager{Application: mockApplication, CheckErrors: mockCheckError}
		t.Run(curTest.testName, func(t *testing.T) {
			outData, _ := linkShortApi.CreateLinkShort(curTest.inCtx, curTest.inData)
			require.Equal(
				t,
				curTest.dataExpected,
				outData,
				fmt.Sprintf("Expected: %s\nbut got: %s", curTest.dataExpected, outData),
			)
		})
	}
}

type testApiTakeLinkFullGrpc struct {
	testName     string
	inData       *proto.LinkShort
	inCtx        context.Context
	App          testAppTakeLinkFull
	checkError   testCheckErrorCreateLSHGrpc
	dataExpected *proto.ResultLinkFull
}

var takeLinkFullGrpc = []testApiTakeLinkFullGrpc{
	{
		testName: "TakeLinkFull Grpc: successful",
		inData:   &proto.LinkShort{LinkShort: "hf89h4qwer"},
		inCtx:    context.Background(),
		App: testAppTakeLinkFull{
			outErr:    nil,
			outResult: "www.site.ru",
			count:     1,
		},
		checkError: testCheckErrorCreateLSHGrpc{
			outErr:      nil,
			outResult:   "",
			outCodeHTTP: errPkg.IntNil,
			count:       1,
		},
		dataExpected: &proto.ResultLinkFull{
			StatusCode: http.StatusOK,
			Body:       &proto.LinkFull{LinkFull: "www.site.ru"},
			Error:      "",
		},
	},
	{
		testName: "CreateLinkShort Grpc: checkError - errCheck",
		inData:   &proto.LinkShort{LinkShort: "hf89h4qwer"},
		inCtx:    context.Background(),
		App: testAppTakeLinkFull{
			outErr:    errors.New(errPkg.LSHCreateLinkShortNotInsertUnique),
			outResult: "",
			count:     1,
		},
		checkError: testCheckErrorCreateLSHGrpc{
			outErr:      errors.New(errPkg.ErrCheck),
			outResult:   errPkg.ErrCheck,
			outCodeHTTP: http.StatusConflict,
			count:       1,
		},
		dataExpected: &proto.ResultLinkFull{
			StatusCode: http.StatusConflict,
			Error:      errPkg.ErrCheck,
		},
	},
}

func TestTakeLinkFullGrpc(t *testing.T) {
	ctrlApp := gomock.NewController(t)
	ctrlCheckError := gomock.NewController(t)
	defer ctrlApp.Finish()
	defer ctrlCheckError.Finish()

	mockApplication := mocks.NewMockLinkShortAppInterface(ctrlApp)
	mockCheckError := mocks.NewMockCheckErrorInterface(ctrlCheckError)

	for _, curTest := range takeLinkFullGrpc {
		mockApplication.
			EXPECT().
			TakeLinkFullApp(curTest.inData.LinkShort).
			Return(curTest.App.outResult, curTest.App.outErr).
			Times(curTest.App.count)

		mockCheckError.
			EXPECT().
			CheckErrorTakeLinkFullGrpc(curTest.App.outErr).
			Return(curTest.checkError.outErr, curTest.checkError.outResult, curTest.checkError.outCodeHTTP).
			Times(curTest.checkError.count)

		linkShortApi := LinkShortManager{Application: mockApplication, CheckErrors: mockCheckError}
		t.Run(curTest.testName, func(t *testing.T) {
			outData, _ := linkShortApi.TakeLinkFull(curTest.inCtx, curTest.inData)
			require.Equal(
				t,
				curTest.dataExpected,
				outData,
				fmt.Sprintf("Expected: %s\nbut got: %s", curTest.dataExpected, outData),
			)
		})
	}
}
