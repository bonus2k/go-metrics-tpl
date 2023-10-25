package repositories

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

var db Storage

type DBStorageImpl struct {
	db *sql.DB
}

func NewDBStorage(connect string) (*Storage, error) {
	if db == nil {
		dataBase, err := sql.Open("pgx", connect)
		if err != nil {
			return nil, err
		}
		db = &DBStorageImpl{db: dataBase}
	}

	return &db, nil
}

func (d *DBStorageImpl) AddGauge(s string, f float64) {
	//TODO implement me
	panic("implement me")
}

func (d *DBStorageImpl) GetGauge(s string) (float64, bool) {
	//TODO implement me
	panic("implement me")
}

func (d *DBStorageImpl) AddCounter(s string, i int64) {
	//TODO implement me
	panic("implement me")
}

func (d *DBStorageImpl) GetCounter(s string) (int64, bool) {
	//TODO implement me
	panic("implement me")
}

func (d *DBStorageImpl) GetAllMetrics() []Metric {
	//TODO implement me
	panic("implement me")
}

func (d *DBStorageImpl) CheckConnection() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()
	return d.db.PingContext(ctx)
}
