package build

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"linkShortOzon/config"
	errPkg "linkShortOzon/internals/myerror"
	"strings"
)

func CreateConn(configDB config.DatabasePostgres) (*pgxpool.Pool, error) {
	addressPostgres := "postgres://" + configDB.UserName + ":" + configDB.Password +
		"@" + configDB.Host + ":" + configDB.Port + "/" + configDB.SchemaName

	conn, errCreateConn := pgxpool.Connect(context.Background(), addressPostgres)
	if errCreateConn != nil {
		return nil, &errPkg.MyErrors{
			Text: errPkg.MCreateDBNotConnect,
		}
	}
	return conn, nil
}

func CreateDB(conn *pgxpool.Pool) error {
	contextTransaction := context.Background()
	tx, errTransaction := conn.Begin(contextTransaction)
	if errTransaction != nil {
		return &errPkg.MyErrors{
			Text: errPkg.MCreateDBTransactionNotCreate,
		}
	}

	defer tx.Rollback(contextTransaction)

	file, errRead := ioutil.ReadFile("./build/postgresql/createtables.sql")
	if errRead != nil {
		return &errPkg.MyErrors{
			Text: errPkg.MCreateDBFileNotFound,
		}
	}

	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		_, errTransaction = tx.Exec(context.Background(), request)
		if errTransaction != nil {
			return &errPkg.MyErrors{
				Text: errPkg.MCreateDBFileNotCreate,
			}
		}
	}

	errCommit := tx.Commit(contextTransaction)
	if errCommit != nil {
		return &errPkg.MyErrors{
			Text: errPkg.MCreateDBNotCommit,
		}
	}

	return nil
}
