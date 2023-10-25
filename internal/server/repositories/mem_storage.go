package repositories

import (
	"errors"
	"fmt"
	"strings"
)

var mem Storage
var sync bool

type MemStorageImpl struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

type Metric struct {
	Name  string
	Value string
}

func NewMemStorage(syncSave bool) *Storage {
	if mem == nil {
		sync = syncSave
		mem = &MemStorageImpl{Gauge: make(map[string]float64), Counter: make(map[string]int64)}
	}
	return &mem
}

func (ms *MemStorageImpl) AddGauge(name string, value float64) {
	trimName := strings.TrimSpace(name)
	ms.Gauge[trimName] = value
	if sync {
		fileService.Save()
	}

}

func (ms *MemStorageImpl) GetGauge(name string) (float64, bool) {
	trimName := strings.TrimSpace(name)
	f, ok := ms.Gauge[trimName]
	return f, ok
}

func (ms *MemStorageImpl) AddCounter(name string, value int64) {
	trimName := strings.TrimSpace(name)
	int64s, found := ms.Counter[trimName]
	if found {
		ms.Counter[trimName] = int64s + value
	} else {
		ms.Counter[trimName] = value
	}
	if sync {
		fileService.Save()
	}
}

func (ms *MemStorageImpl) GetCounter(name string) (int64, bool) {
	trimName := strings.TrimSpace(name)
	int64s, ok := ms.Counter[trimName]
	return int64s, ok
}

func (ms *MemStorageImpl) GetAllMetrics() []Metric {
	var metrics []Metric
	for k, v := range ms.Gauge {
		metrics = append(metrics, Metric{Name: k, Value: fmt.Sprintf("%v", v)})
	}
	for k, v := range ms.Counter {
		metrics = append(metrics, Metric{Name: k, Value: fmt.Sprintf("%v", v)})
	}
	return metrics
}

func (ms *MemStorageImpl) CheckConnection() error {
	if mem == nil {
		return errors.New("storage of mem don't initialized")
	}
	return nil
}
