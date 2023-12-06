//go:generate mockgen -destination=mocks/orm.go -package=mocks linkShortOzon/internals/linkShort/orm/orm.go LinkShortWrapperInterface,ConnectionPostgresInterface,TransactionInterface,ConnectionRedisInterface

package orm

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	errPkg "linkShortOzon/internals/myerror"
)

type LinkShortWrapperInterface interface {
	CreateLinkShortPostgres(linkFull string, linkShort string) error
	TakeLinkFullPostgres(linkShort string) (string, error)
	CreateLinkShortRedis(linkFull string, linkShort string) error
	TakeLinkFullRedis(linkShort string) (string, error)
	CreateLinkShort(linkFull string, linkShort string) error
	TakeLinkFull(linkShort string) (string, error)
}

type ConnectionPostgresInterface interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type TransactionInterface interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	LargeObjects() pgx.LargeObjects
	Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error)
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	Conn() *pgx.Conn
}

type ConnectionRedisInterface interface {
	Close() error
	Err() error
	Do(commandName string, args ...interface{}) (reply interface{}, err error)
	Send(commandName string, args ...interface{}) error
	Flush() error
	Receive() (reply interface{}, err error)
}

type LinkShortWrapper struct {
	ConnPostgres ConnectionPostgresInterface
	ConnRedis    ConnectionRedisInterface
}

func (w *LinkShortWrapper) CreateLinkShort(linkFull string, linkShort string) error {
	if w.ConnPostgres != nil {
		return w.CreateLinkShortPostgres(linkFull, linkShort)
	}
	if w.ConnRedis != nil {
		return w.CreateLinkShortRedis(linkFull, linkShort)
	}

	return &errPkg.MyErrors{
		Text: errPkg.LSHCreateLinkShortNilConn,
	}
}

func (w *LinkShortWrapper) TakeLinkFull(linkShort string) (string, error) {
	if w.ConnPostgres != nil {
		return w.TakeLinkFullPostgres(linkShort)
	}
	if w.ConnRedis != nil {
		return w.TakeLinkFullRedis(linkShort)
	}

	return "", &errPkg.MyErrors{
		Text: errPkg.LSHTakeLinkShortNilConn,
	}
}

func (w *LinkShortWrapper) CreateLinkShortRedis(linkFull string, linkShort string) error {
	full, _ := w.TakeLinkFullRedis(linkFull)
	if full != "" {
		return &errPkg.MyErrors{
			Text: errPkg.LSHCreateLinkShortExistsRedis,
		}
	}

	_, errLinkFull := redis.String(w.ConnRedis.Do("SET", linkFull, linkShort, "EX", 86400))
	if errLinkFull != nil {
		return &errPkg.MyErrors{
			Text: errPkg.LSHCreateLinkShortNotSetFullLinkRedis,
		}
	}

	_, errLinkShort := redis.String(w.ConnRedis.Do("SET", linkShort, linkFull, "EX", 86400))
	if errLinkShort != nil {
		return &errPkg.MyErrors{
			Text: errPkg.LSHCreateLinkShortNotSetShortLinkRedis,
		}
	}

	return nil
}

func (w *LinkShortWrapper) TakeLinkFullRedis(linkShort string) (string, error) {
	resultGet, errGet := redis.Bytes(w.ConnRedis.Do("GET", linkShort))
	if errGet != nil {
		return "", &errPkg.MyErrors{
			Text: errPkg.LSHTakeLinkShortNotFoundRedis,
		}
	}

	return string(resultGet), nil
}

func (w *LinkShortWrapper) CreateLinkShortPostgres(linkFull string, linkShort string) error {
	contextTransaction := context.Background()
	tx, errBeginConn := w.ConnPostgres.Begin(contextTransaction)
	if errBeginConn != nil {
		return &errPkg.MyErrors{
			Text: errPkg.LSHCreateLinkShortTransactionNotCreate,
		}
	}

	defer tx.Rollback(contextTransaction)

	_, errExecTx := tx.Exec(contextTransaction,
		"INSERT INTO public.link (link, link_short) VALUES ($1, $2)", linkFull, linkShort)
	if errExecTx != nil {
		switch errExecTx.Error() {
		case errPkg.LSHCreateLinkShortNotInsertUniqueDB:
			return &errPkg.MyErrors{
				Text: errPkg.LSHCreateLinkShortNotInsertUnique,
			}
		default:
			return &errPkg.MyErrors{
				Text: errPkg.LSHCreateLinkShortNotInsert,
			}
		}
	}

	errCommitTx := tx.Commit(contextTransaction)
	if errCommitTx != nil {
		return &errPkg.MyErrors{
			Text: errPkg.LSHCreateLinkShortNotCommit,
		}
	}

	return nil
}

func (w *LinkShortWrapper) TakeLinkFullPostgres(linkShort string) (string, error) {
	contextTransaction := context.Background()
	tx, errBeginConn := w.ConnPostgres.Begin(contextTransaction)
	if errBeginConn != nil {
		return "", &errPkg.MyErrors{
			Text: errPkg.LSHTakeLinkShortTransactionNotCreate,
		}
	}

	defer tx.Rollback(contextTransaction)

	var linkFull string
	errQueryRow := tx.QueryRow(contextTransaction,
		"SELECT link FROM public.link WHERE link_short = $1",
		linkShort).Scan(&linkFull)
	if errQueryRow != nil {
		if errQueryRow == pgx.ErrNoRows {
			return "", &errPkg.MyErrors{
				Text: errPkg.LSHTakeLinkShortNotFound,
			}
		}
		return "", &errPkg.MyErrors{
			Text: errPkg.LSHTakeLinkShortNotScan,
		}
	}

	errCommitTx := tx.Commit(contextTransaction)
	if errCommitTx != nil {
		return "", &errPkg.MyErrors{
			Text: errPkg.LSHTakeLinkShortNotCommit,
		}
	}

	return linkFull, nil
}
