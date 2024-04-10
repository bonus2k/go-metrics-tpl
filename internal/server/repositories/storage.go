package repositories

import (
	"context"

	"github.com/bonus2k/go-metrics-tpl/internal/models"
)

type Storage interface {
	AddGauge(context.Context, string, float64) error
	GetGauge(context.Context, string) (float64, error)
	AddCounter(context.Context, string, int64) error
	GetCounter(context.Context, string) (int64, error)
	GetAllMetrics(context.Context) ([]Metric, error)
	AddMetrics(context.Context, []models.Metrics) error
	CheckConnection() error
}
