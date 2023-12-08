package application

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"linkShortOzon/internals/linkShort/application/mocks"
	"testing"
)

type takeLinkFullApp struct {
	testName         string
	linkFull         string
	linkShort        string
	outErr           error
	count            int
	linkFullExpected string
	errorExpected    error
}

var createLSHApp = []takeLinkFullApp{
	{
		testName:         "akeLinkFullPostgres: successful",
		linkFull:         "www.site.ru",
		linkShort:        "ozon.click.ru/_FeLIUZ33Y",
		outErr:           nil,
		count:            1,
		linkFullExpected: "www.site.ru",
		errorExpected:    nil,
	},
}

func TestTakeLinkFullPostgres(t *testing.T) {
	ctrlWrapper := gomock.NewController(t)

	defer ctrlWrapper.Finish()

	mockWrapper := mocks.NewMockLinkShortWrapperInterface(ctrlWrapper)
	for _, curTest := range createLSHApp {
		LSHApp := &LinkShortApp{
			Wrapper: mockWrapper,
		}

		mockWrapper.
			EXPECT().
			TakeLinkFull(curTest.linkShort).
			Return(curTest.linkFull, curTest.outErr).
			Times(curTest.count)

		t.Run(curTest.testName, func(t *testing.T) {
			linkFull, errTakeLSH := LSHApp.TakeLinkFullApp(curTest.linkShort)
			if errTakeLSH != nil && curTest.errorExpected != nil {
				require.Equal(
					t,
					curTest.errorExpected.Error(),
					errTakeLSH.Error(),
					fmt.Sprintf("Expected: %s\nbut got: %s", curTest.errorExpected, errTakeLSH.Error()),
				)
			} else {
				require.Equal(
					t,
					curTest.errorExpected,
					errTakeLSH,
					fmt.Sprintf("Expected: %v\nbut got: %v", curTest.errorExpected, errTakeLSH),
				)
			}
			require.Equal(
				t,
				curTest.linkFullExpected,
				linkFull,
				fmt.Sprintf("Expected: %s\nbut got: %s", curTest.linkFullExpected, linkFull),
			)
		})
	}
}
