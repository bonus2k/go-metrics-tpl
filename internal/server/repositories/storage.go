package repositories

import (
	"context"
)

type Storage interface {
	AddGauge(context.Context, string, float64) error
	GetGauge(context.Context, string) (float64, error)
	AddCounter(context.Context, string, int64) error
	GetCounter(context.Context, string) (int64, error)
	GetAllMetrics(context.Context) ([]Metric, error)
	CheckConnection() error
}
