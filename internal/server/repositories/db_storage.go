package repositories

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

var db Storage

type DbStorageImpl struct {
	db *sql.DB
}

func NewDbStorage(connect string) (*Storage, error) {
	if db == nil {
		dataBase, err := sql.Open("pgx", connect)
		if err != nil {
			return nil, err
		}
		db = &DbStorageImpl{db: dataBase}
	}

	return &db, nil
}

func (d *DbStorageImpl) AddGauge(s string, f float64) {
	//TODO implement me
	panic("implement me")
}

func (d *DbStorageImpl) GetGauge(s string) (float64, bool) {
	//TODO implement me
	panic("implement me")
}

func (d *DbStorageImpl) AddCounter(s string, i int64) {
	//TODO implement me
	panic("implement me")
}

func (d *DbStorageImpl) GetCounter(s string) (int64, bool) {
	//TODO implement me
	panic("implement me")
}

func (d *DbStorageImpl) GetAllMetrics() []Metric {
	//TODO implement me
	panic("implement me")
}

func (d *DbStorageImpl) CheckConnection() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()
	return d.db.PingContext(ctx)
}
