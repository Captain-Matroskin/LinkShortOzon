package myerror

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"linkShortOzon/internals/myerror/mocks"
	"testing"
)

type testCheckError struct {
	testName    string
	inErr       error
	log         testLogger
	reqId       int
	expectedErr error
	//expectedRes  string
	//expectedCode int
}

type testLogger struct {
	errorf logger
	warnf  logger
}

type logger struct {
	format string
	args   []interface{}
	count  int
}

var testCheckErrorCreateLinkShort = []testCheckError{
	{
		testName: "CheckErrorCreateLinkShort: successful",
		inErr:    nil,
		log: testLogger{warnf: logger{
			format: "",
			args:   []interface{}{},
			count:  0,
		},
			errorf: logger{
				format: "",
				args:   []interface{}{},
				count:  0,
			},
		},
		reqId:       5,
		expectedErr: nil,
	},
	{
		testName: "CheckErrorCreateLinkShort: error CreateLinkShortNotInsertUnique",
		inErr:    errors.New(LSHCreateLinkShortNotInsertUnique),
		log: testLogger{warnf: logger{
			format: "%s, requestId: %d",
			args:   []interface{}{LSHCreateLinkShortNotInsertUnique, 5},
			count:  1,
		},
			errorf: logger{
				format: "",
				args:   []interface{}{},
				count:  0,
			},
		},
		reqId:       5,
		expectedErr: errors.New(ErrCheck),
	},
	{
		testName: "CheckErrorCreateLinkShort: error default",
		inErr:    errors.New(ErrMarshal),
		log: testLogger{warnf: logger{
			format: "",
			args:   []interface{}{},
			count:  0,
		},
			errorf: logger{
				format: "%s, requestId: %d",
				args:   []interface{}{ErrMarshal, 5},
				count:  1,
			},
		},
		reqId:       5,
		expectedErr: errors.New(ErrCheck),
	},
}

func TestCheckErrorCreateLinkShort(t *testing.T) {
	ctrlMultiLogger := gomock.NewController(t)
	defer ctrlMultiLogger.Finish()

	mockMultiLogger := mocks.NewMockMultiLoggerInterface(ctrlMultiLogger)
	for _, curTest := range testCheckErrorCreateLinkShort {
		checkErr := CheckError{Logger: mockMultiLogger, RequestId: curTest.reqId}
		mockMultiLogger.
			EXPECT().
			Warnf(curTest.log.warnf.format, curTest.log.warnf.args).
			Times(curTest.log.warnf.count)

		mockMultiLogger.
			EXPECT().
			Errorf(curTest.log.errorf.format, curTest.log.errorf.args).
			Times(curTest.log.errorf.count)

		t.Run(curTest.testName, func(t *testing.T) {
			errCheck, _, _ := checkErr.CheckErrorCreateLinkShort(curTest.inErr)
			if errCheck != nil && curTest.expectedErr != nil {
				require.Equal(
					t,
					curTest.expectedErr.Error(),
					errCheck.Error(),
					fmt.Sprintf("Expected: %s\nbut got: %s", curTest.expectedErr.Error(), errCheck.Error()),
				)
			} else {
				require.Equal(
					t,
					curTest.expectedErr,
					errCheck,
					fmt.Sprintf("Expected: %v\nbut got: %v", curTest.expectedErr, errCheck),
				)
			}
		})
	}
}

var testCheckErrorTakeLinkFull = []testCheckError{
	{
		testName: "CheckErrorTakeLinkFull: successful",
		inErr:    nil,
		log: testLogger{warnf: logger{
			format: "",
			args:   []interface{}{},
			count:  0,
		},
			errorf: logger{
				format: "",
				args:   []interface{}{},
				count:  0,
			},
		},
		reqId:       5,
		expectedErr: nil,
	},
	{
		testName: "CheckErrorTakeLinkFull: error CreateLinkShortNotInsertUnique",
		inErr:    errors.New(LSHTakeLinkShortNotFound),
		log: testLogger{warnf: logger{
			format: "%s, requestId: %d",
			args:   []interface{}{LSHTakeLinkShortNotFound, 5},
			count:  1,
		},
			errorf: logger{
				format: "",
				args:   []interface{}{},
				count:  0,
			},
		},
		reqId:       5,
		expectedErr: errors.New(ErrCheck),
	},
	{
		testName: "CheckErrorTakeLinkFull: error default",
		inErr:    errors.New(ErrMarshal),
		log: testLogger{warnf: logger{
			format: "",
			args:   []interface{}{},
			count:  0,
		},
			errorf: logger{
				format: "%s, requestId: %d",
				args:   []interface{}{ErrMarshal, 5},
				count:  1,
			},
		},
		reqId:       5,
		expectedErr: errors.New(ErrCheck),
	},
}

func TestCheckErrorTakeLinkFull(t *testing.T) {
	ctrlMultiLogger := gomock.NewController(t)
	defer ctrlMultiLogger.Finish()

	mockMultiLogger := mocks.NewMockMultiLoggerInterface(ctrlMultiLogger)
	for _, curTest := range testCheckErrorTakeLinkFull {
		checkErr := CheckError{Logger: mockMultiLogger, RequestId: curTest.reqId}
		mockMultiLogger.
			EXPECT().
			Warnf(curTest.log.warnf.format, curTest.log.warnf.args).
			Times(curTest.log.warnf.count)

		mockMultiLogger.
			EXPECT().
			Errorf(curTest.log.errorf.format, curTest.log.errorf.args).
			Times(curTest.log.errorf.count)

		t.Run(curTest.testName, func(t *testing.T) {
			errCheck, _, _ := checkErr.CheckErrorTakeLinkFull(curTest.inErr)
			if errCheck != nil && curTest.expectedErr != nil {
				require.Equal(
					t,
					curTest.expectedErr.Error(),
					errCheck.Error(),
					fmt.Sprintf("Expected: %s\nbut got: %s", curTest.expectedErr.Error(), errCheck.Error()),
				)
			} else {
				require.Equal(
					t,
					curTest.expectedErr,
					errCheck,
					fmt.Sprintf("Expected: %v\nbut got: %v", curTest.expectedErr, errCheck),
				)
			}
		})
	}
}

var testCheckErrorCreateLinkShortGrpc = []testCheckError{
	{
		testName: "CheckErrorCreateLinkShortGrpc: successful",
		inErr:    nil,
		log: testLogger{warnf: logger{
			format: "",
			args:   []interface{}{},
			count:  0,
		},
			errorf: logger{
				format: "",
				args:   []interface{}{},
				count:  0,
			},
		},
		reqId:       5,
		expectedErr: nil,
	},
	{
		testName: "CheckErrorCreateLinkShortGrpc: error CreateLinkShortNotInsertUnique",
		inErr:    errors.New(LSHCreateLinkShortNotInsertUnique),
		log: testLogger{warnf: logger{
			format: "%s, requestId: %d",
			args:   []interface{}{LSHCreateLinkShortNotInsertUnique, 5},
			count:  1,
		},
			errorf: logger{
				format: "",
				args:   []interface{}{},
				count:  0,
			},
		},
		reqId:       5,
		expectedErr: errors.New(ErrCheck),
	},
	{
		testName: "CheckErrorCreateLinkShortGrpc: error default",
		inErr:    errors.New(ErrMarshal),
		log: testLogger{warnf: logger{
			format: "",
			args:   []interface{}{},
			count:  0,
		},
			errorf: logger{
				format: "%s, requestId: %d",
				args:   []interface{}{ErrMarshal, 5},
				count:  1,
			},
		},
		reqId:       5,
		expectedErr: errors.New(ErrInternal),
	},
}

func TestCheckErrorCreateLinkShortGrpc(t *testing.T) {
	ctrlMultiLogger := gomock.NewController(t)
	defer ctrlMultiLogger.Finish()

	mockMultiLogger := mocks.NewMockMultiLoggerInterface(ctrlMultiLogger)
	for _, curTest := range testCheckErrorCreateLinkShortGrpc {
		checkErr := CheckError{Logger: mockMultiLogger, RequestId: curTest.reqId}
		mockMultiLogger.
			EXPECT().
			Warnf(curTest.log.warnf.format, curTest.log.warnf.args).
			Times(curTest.log.warnf.count)

		mockMultiLogger.
			EXPECT().
			Errorf(curTest.log.errorf.format, curTest.log.errorf.args).
			Times(curTest.log.errorf.count)

		t.Run(curTest.testName, func(t *testing.T) {
			errCheck, _, _ := checkErr.CheckErrorCreateLinkShortGrpc(curTest.inErr)
			if errCheck != nil && curTest.expectedErr != nil {
				require.Equal(
					t,
					curTest.expectedErr.Error(),
					errCheck.Error(),
					fmt.Sprintf("Expected: %s\nbut got: %s", curTest.expectedErr.Error(), errCheck.Error()),
				)
			} else {
				require.Equal(
					t,
					curTest.expectedErr,
					errCheck,
					fmt.Sprintf("Expected: %v\nbut got: %v", curTest.expectedErr, errCheck),
				)
			}
		})
	}
}

var testCheckErrorTakeLinkFullGrpc = []testCheckError{
	{
		testName: "CheckErrorTakeLinkFullGrpc: successful",
		inErr:    nil,
		log: testLogger{warnf: logger{
			format: "",
			args:   []interface{}{},
			count:  0,
		},
			errorf: logger{
				format: "",
				args:   []interface{}{},
				count:  0,
			},
		},
		reqId:       5,
		expectedErr: nil,
	},
	{
		testName: "CheckErrorTakeLinkFullGrpc: error LSHTakeLinkShortNotFound",
		inErr:    errors.New(LSHTakeLinkShortNotFound),
		log: testLogger{warnf: logger{
			format: "%s, requestId: %d",
			args:   []interface{}{LSHTakeLinkShortNotFound, 5},
			count:  1,
		},
			errorf: logger{
				format: "",
				args:   []interface{}{},
				count:  0,
			},
		},
		reqId:       5,
		expectedErr: errors.New(ErrCheck),
	},
	{
		testName: "CheckErrorTakeLinkFullGrpc: error default",
		inErr:    errors.New(ErrMarshal),
		log: testLogger{warnf: logger{
			format: "",
			args:   []interface{}{},
			count:  0,
		},
			errorf: logger{
				format: "%s, requestId: %d",
				args:   []interface{}{ErrMarshal, 5},
				count:  1,
			},
		},
		reqId:       5,
		expectedErr: errors.New(ErrInternal),
	},
}

func TestCheckErrorTakeLinkFullGrpc(t *testing.T) {
	ctrlMultiLogger := gomock.NewController(t)
	defer ctrlMultiLogger.Finish()

	mockMultiLogger := mocks.NewMockMultiLoggerInterface(ctrlMultiLogger)
	for _, curTest := range testCheckErrorTakeLinkFullGrpc {
		checkErr := CheckError{Logger: mockMultiLogger, RequestId: curTest.reqId}
		mockMultiLogger.
			EXPECT().
			Warnf(curTest.log.warnf.format, curTest.log.warnf.args).
			Times(curTest.log.warnf.count)

		mockMultiLogger.
			EXPECT().
			Errorf(curTest.log.errorf.format, curTest.log.errorf.args).
			Times(curTest.log.errorf.count)

		t.Run(curTest.testName, func(t *testing.T) {
			errCheck, _, _ := checkErr.CheckErrorTakeLinkFullGrpc(curTest.inErr)
			if errCheck != nil && curTest.expectedErr != nil {
				require.Equal(
					t,
					curTest.expectedErr.Error(),
					errCheck.Error(),
					fmt.Sprintf("Expected: %s\nbut got: %s", curTest.expectedErr.Error(), errCheck.Error()),
				)
			} else {
				require.Equal(
					t,
					curTest.expectedErr,
					errCheck,
					fmt.Sprintf("Expected: %v\nbut got: %v", curTest.expectedErr, errCheck),
				)
			}
		})
	}
}
