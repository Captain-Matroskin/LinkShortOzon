package orm

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
	"linkShortOzon/internals/linkShort/orm/mocks"
	errPkg "linkShortOzon/internals/myerror"
	"testing"
	"time"
)

type createLinkShortPostgres struct {
	testName      string
	inLinkFull    string
	inLinkShort   string
	connBegin     connPostgresBegin
	transExec     transactionExec
	transCommit   transactionCommit
	countRollback int
	errorExpected error
}

type connPostgresBegin struct {
	outError error
	count    int
}

type transactionExec struct {
	outError error
	count    int
}

type transactionCommit struct {
	outError error
	count    int
}

var createLinkShort = []createLinkShortPostgres{
	{
		testName:      "CreateLinkShort orm: successful postgres",
		inLinkFull:    "www.site.ru",
		connBegin:     connPostgresBegin{outError: nil, count: 1},
		transExec:     transactionExec{outError: nil, count: 1},
		transCommit:   transactionCommit{outError: nil, count: 1},
		countRollback: 1,
		errorExpected: nil,
	},
	{
		testName:   "CreateLinkShort orm: (error) TransactionNotCreate",
		inLinkFull: "www.site.ru",
		connBegin: connPostgresBegin{
			outError: errors.New(errPkg.LSHCreateLinkShortTransactionNotCreate),
			count:    1,
		},
		transExec:     transactionExec{outError: nil, count: 0},
		transCommit:   transactionCommit{outError: nil, count: 0},
		countRollback: 0,
		errorExpected: errors.New(errPkg.LSHCreateLinkShortTransactionNotCreate),
	},
	{
		testName:      "CreateLinkShort orm: (error) NotInsertUnique",
		inLinkFull:    "www.site.ru",
		connBegin:     connPostgresBegin{outError: nil, count: 1},
		transExec:     transactionExec{outError: errors.New(errPkg.LSHCreateLinkShortNotInsertUniqueDB), count: 1},
		transCommit:   transactionCommit{outError: nil, count: 0},
		countRollback: 1,
		errorExpected: errors.New(errPkg.LSHCreateLinkShortNotInsertUnique),
	},
	{
		testName:      "CreateLinkShort orm: (error) LSHCreateLinkShortNotInsert",
		inLinkFull:    "www.site.ru",
		connBegin:     connPostgresBegin{outError: nil, count: 1},
		transExec:     transactionExec{outError: errors.New(errPkg.LSHCreateLinkShortNotInsert), count: 1},
		transCommit:   transactionCommit{outError: nil, count: 0},
		countRollback: 1,
		errorExpected: errors.New(errPkg.LSHCreateLinkShortNotInsert),
	},
	{
		testName:      "CreateLinkShort orm: successful postgres",
		inLinkFull:    "www.site.ru",
		connBegin:     connPostgresBegin{outError: nil, count: 1},
		transExec:     transactionExec{outError: nil, count: 1},
		transCommit:   transactionCommit{outError: errors.New(errPkg.LSHCreateLinkShortNotCommit), count: 1},
		countRollback: 1,
		errorExpected: errors.New(errPkg.LSHCreateLinkShortNotCommit),
	},
}

func TestCreateLinkShortPostgres(t *testing.T) {
	ctrlWrapper := gomock.NewController(t)
	ctrlPostgresConn := gomock.NewController(t)
	ctrlTransaction := gomock.NewController(t)

	defer ctrlWrapper.Finish()
	defer ctrlPostgresConn.Finish()
	defer ctrlTransaction.Finish()

	mockPostgresConn := mocks.NewMockConnectionPostgresInterface(ctrlPostgresConn)
	mockTransaction := mocks.NewMockTransactionInterface(ctrlTransaction)
	for _, curTest := range createLinkShort {
		wrapperOrm := &LinkShortWrapper{
			ConnPostgres: mockPostgresConn,
		}
		ctx := context.Background()

		mockPostgresConn.
			EXPECT().
			Begin(ctx).
			Return(mockTransaction, curTest.connBegin.outError).
			Times(curTest.connBegin.count)

		var outExec pgconn.CommandTag
		mockTransaction.
			EXPECT().
			Exec(ctx, "INSERT INTO public.link (link, link_short) VALUES ($1, $2)", curTest.inLinkFull, curTest.inLinkShort).
			Return(outExec, curTest.transExec.outError).
			Times(curTest.transExec.count)

		mockTransaction.
			EXPECT().
			Commit(ctx).
			Return(curTest.transCommit.outError).
			Times(curTest.transCommit.count)

		var errRollback error
		errRollback = nil
		mockTransaction.
			EXPECT().
			Rollback(ctx).
			Return(errRollback).
			Times(curTest.countRollback)

		t.Run(curTest.testName, func(t *testing.T) {
			errCreateLSH := wrapperOrm.CreateLinkShortPostgres(curTest.inLinkFull, curTest.inLinkShort)
			if errCreateLSH != nil && curTest.errorExpected != nil {
				require.Equal(
					t,
					curTest.errorExpected.Error(),
					errCreateLSH.Error(),
					fmt.Sprintf("Expected: %s\nbut got: %s", curTest.errorExpected, errCreateLSH),
				)
			} else {
				require.Equal(
					t,
					curTest.errorExpected,
					errCreateLSH,
					fmt.Sprintf("Expected: %s\nbut got: %s", curTest.errorExpected, errCreateLSH),
				)
			}
		})
	}

}

type Row struct {
	row    []interface{}
	outErr error
}

func (r *Row) Scan(dest ...interface{}) error {
	if r.outErr != nil {
		return r.outErr
	}
	for i := range dest {
		if r.row[i] == nil {
			dest[i] = nil
			continue
		}
		switch dest[i].(type) {
		case *int:
			*dest[i].(*int) = r.row[i].(int)
		case *string:
			*dest[i].(*string) = r.row[i].(string)
		case **string:
			t := r.row[i].(string)
			*dest[i].(**string) = &t
		case *float32:
			*dest[i].(*float32) = float32(r.row[i].(float64))
		case **int32:
			t := int32(r.row[i].(int))
			*dest[i].(**int32) = &t
		case *time.Time:
			*dest[i].(*time.Time) = r.row[i].(time.Time)
		case *bool:
			*dest[i].(*bool) = r.row[i].(bool)
		default:
			dest[i] = nil
		}
	}
	return nil
}

type takeLinkFullPostgres struct {
	testName         string
	inLinkShort      string
	connBegin        connPostgresBegin
	transQueryRow    transQueryRow
	row              Row
	transCommit      transactionCommit
	countRollback    int
	linkFullExpected string
	errorExpected    error
}

type transQueryRow struct {
	outError error
	count    int
}

var takeLinkFullPst = []takeLinkFullPostgres{
	{
		testName:         "takeLinkFullPostgres orm: successful postgres",
		inLinkShort:      "ozon.click.ru/_FeLIUZ33Y",
		connBegin:        connPostgresBegin{outError: nil, count: 1},
		transQueryRow:    transQueryRow{outError: nil, count: 1},
		row:              Row{row: []interface{}{"www.site.ru"}, outErr: nil},
		transCommit:      transactionCommit{outError: nil, count: 1},
		countRollback:    1,
		linkFullExpected: "www.site.ru",
		errorExpected:    nil,
	},
	{
		testName:         "takeLinkFullPostgres orm: (error) transactionNotCreate",
		inLinkShort:      "ozon.click.ru/_FeLIUZ33Y",
		connBegin:        connPostgresBegin{outError: errors.New(errPkg.LSHTakeLinkShortTransactionNotCreate), count: 1},
		transQueryRow:    transQueryRow{outError: nil, count: 0},
		row:              Row{row: []interface{}{"www.site.ru"}, outErr: nil},
		transCommit:      transactionCommit{outError: nil, count: 0},
		countRollback:    0,
		linkFullExpected: "",
		errorExpected:    errors.New(errPkg.LSHTakeLinkShortTransactionNotCreate),
	},
	{
		testName:         "takeLinkFullPostgres orm:  (error) errQueryRow - no rows",
		inLinkShort:      "ozon.click.ru/_FeLIUZ33Y",
		connBegin:        connPostgresBegin{outError: nil, count: 1},
		transQueryRow:    transQueryRow{outError: nil, count: 1},
		row:              Row{outErr: pgx.ErrNoRows},
		transCommit:      transactionCommit{outError: nil, count: 0},
		countRollback:    1,
		linkFullExpected: "",
		errorExpected:    errors.New(errPkg.LSHTakeLinkShortNotFound),
	},
	{
		testName:         "takeLinkFullPostgres orm:  (error) errQueryRow - default",
		inLinkShort:      "ozon.click.ru/_FeLIUZ33Y",
		connBegin:        connPostgresBegin{outError: nil, count: 1},
		transQueryRow:    transQueryRow{outError: nil, count: 1},
		row:              Row{outErr: errors.New(errPkg.LSHTakeLinkShortNotScan)},
		transCommit:      transactionCommit{outError: nil, count: 0},
		countRollback:    1,
		linkFullExpected: "",
		errorExpected:    errors.New(errPkg.LSHTakeLinkShortNotScan),
	},
	{
		testName:         "takeLinkFullPostgres orm: (error) commit",
		inLinkShort:      "ozon.click.ru/_FeLIUZ33Y",
		connBegin:        connPostgresBegin{outError: nil, count: 1},
		transQueryRow:    transQueryRow{outError: nil, count: 1},
		row:              Row{row: []interface{}{"www.site.ru"}, outErr: nil},
		transCommit:      transactionCommit{outError: errors.New(errPkg.LSHTakeLinkShortNotCommit), count: 1},
		countRollback:    1,
		linkFullExpected: "",
		errorExpected:    errors.New(errPkg.LSHTakeLinkShortNotCommit),
	},
}

func TestTakeLinkFullPostgres(t *testing.T) {
	ctrlWrapper := gomock.NewController(t)
	ctrlPostgresConn := gomock.NewController(t)
	ctrlTransaction := gomock.NewController(t)

	defer ctrlWrapper.Finish()
	defer ctrlPostgresConn.Finish()
	defer ctrlTransaction.Finish()

	mockPostgresConn := mocks.NewMockConnectionPostgresInterface(ctrlPostgresConn)
	mockTransaction := mocks.NewMockTransactionInterface(ctrlTransaction)
	for _, curTest := range takeLinkFullPst {
		wrapperOrm := &LinkShortWrapper{
			ConnPostgres: mockPostgresConn,
		}
		ctx := context.Background()

		mockPostgresConn.
			EXPECT().
			Begin(ctx).
			Return(mockTransaction, curTest.connBegin.outError).
			Times(curTest.connBegin.count)

		mockTransaction.
			EXPECT().
			QueryRow(ctx, "SELECT link FROM public.link WHERE link_short = $1", curTest.inLinkShort).
			Return(&curTest.row).
			Times(curTest.transQueryRow.count)

		mockTransaction.
			EXPECT().
			Commit(ctx).
			Return(curTest.transCommit.outError).
			Times(curTest.transCommit.count)

		var errRollback error
		errRollback = nil
		mockTransaction.
			EXPECT().
			Rollback(ctx).
			Return(errRollback).
			Times(curTest.countRollback)

		t.Run(curTest.testName, func(t *testing.T) {
			linkFull, errTakeLSH := wrapperOrm.TakeLinkFullPostgres(curTest.inLinkShort)
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
