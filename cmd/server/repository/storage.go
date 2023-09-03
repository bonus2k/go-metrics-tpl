package repository

import "strings"

type MemStorageImpl struct {
	gauge   map[string]float64
	counter map[string][]int64
}

func NewMemStorage() MemStorage {
	return &MemStorageImpl{gauge: make(map[string]float64), counter: make(map[string][]int64)}
}

func (ms *MemStorageImpl) AddGauge(name string, value float64) {
	trimName := strings.TrimSpace(name)
	ms.gauge[trimName] = value
}

func (ms *MemStorageImpl) GetGauge(name string) float64 {
	trimName := strings.TrimSpace(name)
	return ms.gauge[trimName]
}

func (ms *MemStorageImpl) AddCounter(name string, value int64) {
	trimName := strings.TrimSpace(name)
	int64s, found := ms.counter[trimName]
	if found {
		ms.counter[trimName] = append(int64s, value)
	} else {
		ms.counter[trimName] = []int64{value}
	}
}

func (ms *MemStorageImpl) GetCounter(name string) []int64 {
	trimName := strings.TrimSpace(name)
	return ms.counter[trimName]
}

type MemStorage interface {
	AddGauge(string, float64)
	GetGauge(string) float64
	AddCounter(string, int64)
	GetCounter(string) []int64
}
