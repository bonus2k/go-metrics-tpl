// Package migrations реализует создание таблиц в БД для работы сервиса Server
package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const migrationsPath = "file://migrations"

func Start(connect string) error {
	dataBase, err := sql.Open("pgx", connect)

	if err != nil {
		logger.Log.Error("can't open connection to db", err)
		return fmt.Errorf("can't open connection to db %w", err)
	}

	defer func() {
		if dataBase != nil {
			_ = dataBase.Close()
		}
	}()

	f := func() error {
		if err := establishConnection(dataBase); err != nil {

			return err
		}
		return nil
	}
	err = utils.RetryAfterError(f)
	if err != nil {
		logger.Log.Error("can't connected to db", err)
		return fmt.Errorf("can't connected to db %w", err)
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
	defer func() {
		err = driver.Close()
		if err != nil {
			logger.Log.Error("migrateSQL", err)
		}
	}()
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/server/migrations/sql",
		"pgx", driver)
	defer func() {
		sourceErr, databaseErr := m.Close()
		if databaseErr != nil {
			logger.Log.Error("databaseErr", databaseErr)
		}
		if sourceErr != nil {
			logger.Log.Error("sourceErr", sourceErr)
		}
	}()
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
