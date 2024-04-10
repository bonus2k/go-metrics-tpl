package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	m "github.com/bonus2k/go-metrics-tpl/internal/models"
)

var mem Storage
var sync bool

type MemStorageImpl struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func (ms *MemStorageImpl) AddMetrics(ctx context.Context, metrics []m.Metrics) error {
	for _, v := range metrics {
		switch v.MType {
		case "counter":
			err := ms.AddCounter(context.TODO(), v.ID, *v.Delta)
			if err != nil {
				return err
			}
		case "gauge":
			err := ms.AddGauge(context.TODO(), v.ID, *v.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewMemStorage(syncSave bool) *Storage {
	if mem == nil {
		sync = syncSave
		mem = &MemStorageImpl{Gauge: make(map[string]float64), Counter: make(map[string]int64)}
	}
	return &mem
}

func (ms *MemStorageImpl) AddGauge(ctx context.Context, name string, value float64) error {
	trimName := strings.TrimSpace(name)
	ms.Gauge[trimName] = value
	if sync {
		return fileService.Save()
	}
	return nil
}

func (ms *MemStorageImpl) GetGauge(ctx context.Context, name string) (float64, error) {
	trimName := strings.TrimSpace(name)
	f, ok := ms.Gauge[trimName]
	if !ok {
		return f, errors.New("gauge not found")
	}
	return f, nil
}

func (ms *MemStorageImpl) AddCounter(ctx context.Context, name string, value int64) error {
	trimName := strings.TrimSpace(name)
	int64s, found := ms.Counter[trimName]
	if found {
		ms.Counter[trimName] = int64s + value
	} else {
		ms.Counter[trimName] = value
	}
	if sync {
		return fileService.Save()
	}
	return nil
}

func (ms *MemStorageImpl) GetCounter(ctx context.Context, name string) (int64, error) {
	trimName := strings.TrimSpace(name)
	int64s, ok := ms.Counter[trimName]
	if !ok {
		return int64s, errors.New("counter not found")
	}
	return int64s, nil
}

func (ms *MemStorageImpl) GetAllMetrics(ctx context.Context) ([]Metric, error) {
	var metrics []Metric
	for k, v := range ms.Gauge {
		metrics = append(metrics, Metric{Name: k, Value: fmt.Sprintf("%v", v)})
	}
	for k, v := range ms.Counter {
		metrics = append(metrics, Metric{Name: k, Value: fmt.Sprintf("%v", v)})
	}
	return metrics, nil
}

func (ms *MemStorageImpl) CheckConnection() error {
	return errors.New("connection don't initialized")
}
