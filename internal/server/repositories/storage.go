package repositories

import (
	"fmt"
	"strings"
)

var mem MemStorage

type MemStorageImpl struct {
	gauge   map[string]float64
	counter map[string]int64
}

type Metric struct {
	Name  string
	Value string
}

func NewMemStorage() MemStorage {
	if mem == nil {
		mem = &MemStorageImpl{gauge: make(map[string]float64), counter: make(map[string]int64)}
	}
	return mem
}

func (ms *MemStorageImpl) AddGauge(name string, value float64) {
	trimName := strings.TrimSpace(name)
	ms.gauge[trimName] = value
}

func (ms *MemStorageImpl) GetGauge(name string) (float64, bool) {
	trimName := strings.TrimSpace(name)
	f, ok := ms.gauge[trimName]
	return f, ok
}

func (ms *MemStorageImpl) AddCounter(name string, value int64) {
	trimName := strings.TrimSpace(name)
	int64s, found := ms.counter[trimName]
	if found {
		ms.counter[trimName] = int64s + value
	} else {
		ms.counter[trimName] = value
	}
}

func (ms *MemStorageImpl) GetCounter(name string) (int64, bool) {
	trimName := strings.TrimSpace(name)
	int64s, ok := ms.counter[trimName]
	return int64s, ok
}

func (ms *MemStorageImpl) GetAllMetrics() []Metric {
	var metrics []Metric
	for k, v := range ms.gauge {
		metrics = append(metrics, Metric{Name: k, Value: fmt.Sprintf("%v", v)})
	}
	for k, v := range ms.counter {
		metrics = append(metrics, Metric{Name: k, Value: fmt.Sprintf("%v", v)})
	}
	return metrics
}

type MemStorage interface {
	AddGauge(string, float64)
	GetGauge(string) (float64, bool)
	AddCounter(string, int64)
	GetCounter(string) (int64, bool)
	GetAllMetrics() []Metric
}
