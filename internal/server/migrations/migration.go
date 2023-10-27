package migrations

import (
	"context"
	"database/sql"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

const migrationsPath = "file://migrations"

func Start(connect string) error {
	dataBase, err := sql.Open("pgx", connect)
	if err != nil {
		logger.Log.Error("can't open connection to db", zap.Error(err))
		return err
	}

	defer func() {
		if dataBase != nil {
			_ = dataBase.Close()
		}
	}()

	if err := establishConnection(dataBase); err != nil {
		logger.Log.Error("can't connected to db", zap.Error(err))
		return err
	}

	err = migrateSQL(dataBase)

	if err != nil {
		panic(err)
	} else {
		logger.Log.Info("migration successfully finished")
	}

	return nil
}

func establishConnection(db *sql.DB) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()
	return db.PingContext(ctx)
}

func migrateSQL(db *sql.DB) error {
	driver, err := mpgx.WithInstance(db, &mpgx.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/server/migrations/sql",
		"pgx", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
