package build

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"linkShortOzon/build/migrations"
	"linkShortOzon/config"
	errPkg "linkShortOzon/internals/myerror"
)

func CreateConn(configDB config.DatabasePostgres) (*pgxpool.Pool, error) {
	addressPostgres := "postgres://" + configDB.UserName + ":" + configDB.Password +
		"@" + configDB.Host + ":" + configDB.Port + "/" + configDB.SchemaName

	conn, errCreateConn := pgxpool.Connect(context.Background(), addressPostgres)
	if errCreateConn != nil {
		return nil, errCreateConn
	}
	return conn, nil
}

func CreateDB(conn *pgxpool.Pool) error {
	migrator, errNewMigrator := migrations.NewMigrator(conn.Config().ConnString())
	if errNewMigrator != nil {
		return errNewMigrator
	}

	now, exp, _, errMigrInfo := migrator.Info()
	if errMigrInfo != nil {
		return errMigrInfo
	}
	if now < exp {
		errMigrate := migrator.Migrate()
		if errMigrate != nil {
			return errMigrate
		}
	} else {
		return &errPkg.MyErrors{
			Text: errPkg.MMigrateDontNeeded,
		}
	}

	return nil
}
