// Package repositories реализует сохрание и получение метрик из БД
package repositories

import (
	"context"

	"github.com/bonus2k/go-metrics-tpl/internal/models"
)

type Storage interface {
	// AddGauge сохранить метрику Gauge из хранилища
	AddGauge(context.Context, string, float64) error
	// GetGauge получить из хранилища метрику Gauge
	GetGauge(context.Context, string) (float64, error)
	// AddCounter сохранить метрику Counter в хранилища
	AddCounter(context.Context, string, int64) error
	// GetCounter получить из хранилища метрику Counter
	GetCounter(context.Context, string) (int64, error)
	// GetAllMetrics получить из хранилища метрики
	GetAllMetrics(context.Context) ([]Metric, error)
	// AddMetrics сохранить в хранилища пакет метрик Gauge
	AddMetrics(context.Context, []models.Metrics) error
	// CheckConnection проверить соединение с хранилища
	CheckConnection() error
}
