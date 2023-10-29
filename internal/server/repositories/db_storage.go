package repositories

import (
	"context"
	"database/sql"
	m "github.com/bonus2k/go-metrics-tpl/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

var db Storage

const t = 10 * time.Second

type DBStorageImpl struct {
	db *sql.DB
}

func (d *DBStorageImpl) AddMetrics(ctx context.Context, metrics []m.Metrics) error {
	context, cancelFunc := context.WithTimeout(ctx, t)
	defer cancelFunc()
	tx, err := d.db.BeginTx(context, nil)
	defer tx.Commit()
	if err != nil {
		return err
	}
	stmtG, err := tx.PrepareContext(context, "INSERT INTO gauge (name, value) VALUES($1,$2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value;")
	if err != nil {
		return err
	}
	stmtC, err := tx.PrepareContext(context, "INSERT INTO count (name, value) VALUES($1,$2) ON CONFLICT (name) DO UPDATE SET value = count.value + EXCLUDED.value;")
	if err != nil {
		return err
	}

	for _, v := range metrics {
		switch v.MType {
		case "counter":
			_, err := stmtC.ExecContext(context, v.ID, v.Delta)
			if err != nil {
				tx.Rollback()
				return err
			}
		case "gauge":
			_, err := stmtG.ExecContext(context, v.ID, v.Value)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}

func (d *DBStorageImpl) AddGauge(ctx context.Context, s string, f float64) error {
	timeout, cancelFunc := context.WithTimeout(ctx, t)
	defer cancelFunc()
	_, err := d.db.ExecContext(
		timeout,
		"INSERT INTO gauge (name, value) VALUES($1,$2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value;",
		s, f,
	)
	return err
}

func (d *DBStorageImpl) GetGauge(ctx context.Context, s string) (float64, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, t)
	defer cancelFunc()
	var v float64
	row := d.db.QueryRowContext(
		timeout,
		"SELECT value FROM gauge WHERE name = $1",
		s,
	)
	err := row.Scan(&v)
	return v, err
}

func (d *DBStorageImpl) AddCounter(ctx context.Context, s string, i int64) error {
	timeout, cancelFunc := context.WithTimeout(ctx, t)
	defer cancelFunc()
	_, err := d.db.ExecContext(
		timeout,
		"INSERT INTO count (name, value) VALUES($1,$2) ON CONFLICT (name) DO UPDATE SET value = count.value + EXCLUDED.value;",
		s, i,
	)
	return err
}

func (d *DBStorageImpl) GetCounter(ctx context.Context, s string) (int64, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, t)
	defer cancelFunc()
	var v int64
	row := d.db.QueryRowContext(
		timeout,
		"SELECT value FROM count WHERE name = $1",
		s,
	)
	err := row.Scan(&v)
	return v, err
}

func (d *DBStorageImpl) GetAllMetrics(ctx context.Context) ([]Metric, error) {
	context, cancelFunc := context.WithTimeout(ctx, t)
	defer cancelFunc()
	stmtG, err := d.db.PrepareContext(context, "SELECT name, value FROM gauge")
	if err != nil {
		return nil, err
	}
	defer stmtG.Close()
	stmtC, err := d.db.PrepareContext(context, "SELECT name, value FROM count")
	if err != nil {
		return nil, err
	}
	defer stmtC.Close()

	rowsGauge, err := stmtG.Query()
	if err != nil {
		return nil, err
	}
	rowsCount, err := stmtC.Query()
	if err != nil {
		return nil, err
	}
	gauges, err := funcName(rowsGauge)
	if err != nil {
		return nil, err
	}
	counts, err := funcName(rowsCount)
	if err != nil {
		return nil, err
	}
	return append(gauges, counts...), nil
}

func funcName(rows *sql.Rows) ([]Metric, error) {
	metrics := make([]Metric, 0)
	for rows.Next() {
		metric := Metric{}
		err := rows.Scan(&metric.Name, &metric.Value)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
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

func (d *DBStorageImpl) CheckConnection() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()
	return d.db.PingContext(ctx)
}
