package orm

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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
	return nil
}

func (w *LinkShortWrapper) TakeLinkFull(linkShort string) (string, error) {
	return "", nil
}

func (w *LinkShortWrapper) CreateLinkShortRedis(linkFull string, linkShort string) error {
	return nil
}

func (w *LinkShortWrapper) TakeLinkFullRedis(linkShort string) (string, error) {
	return "", nil
}

func (w *LinkShortWrapper) CreateLinkShortPostgres(linkFull string, linkShort string) error {
	return nil
}

func (w *LinkShortWrapper) TakeLinkFullPostgres(linkShort string) (string, error) {
	return "", nil
}
