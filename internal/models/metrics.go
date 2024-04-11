// Package models описывает схуме JSON для передачи метрик на сервер
package models

import "strconv"

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// ConvertGaugeToMetrics преобразовывает Map с метрика Gauge в список Metrics для отправик на сервис Server
func ConvertGaugeToMetrics(metrics *map[string]string) ([]Metrics, error) {
	listM := make([]Metrics, 0)
	for k, v := range *metrics {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		m := Metrics{ID: k, Value: &value, MType: "gauge"}
		listM = append(listM, m)
	}
	return listM, nil
}
