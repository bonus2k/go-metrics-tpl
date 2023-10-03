package controllers

import (
	"encoding/json"
	"github.com/bonus2k/go-metrics-tpl/internal/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/models"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func SaveMetric(w http.ResponseWriter, r *http.Request) {
	if MemStorage == nil {
		panic("storage not initialized")
	}
	logger.Log.Debug("decoding request")
	var metric models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metric); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch strings.ToLower(metric.MType) {
	case "gauge":
		MemStorage.AddGauge(metric.ID, *metric.Value)
		w.WriteHeader(http.StatusOK)
		return
	case "counter":
		MemStorage.AddCounter(metric.ID, *metric.Delta)
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	if MemStorage == nil {
		panic("storage not initialized")
	}
	logger.Log.Debug("decoding request")
	var metric models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metric); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if len(metric.ID) == 0 || len(metric.MType) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch strings.ToLower(metric.MType) {
	case "gauge":
		if gauge, ok := MemStorage.GetGauge(metric.ID); ok {
			metric.Value = &gauge
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case "counter":
		if counter, ok := MemStorage.GetCounter(metric.ID); ok {
			metric.Delta = &counter
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
	w.WriteHeader(http.StatusOK)
	return
}

func AllMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := MemStorage.GetAllMetrics()
	marshal, _ := json.Marshal(metrics)
	_, err := w.Write(marshal)
	if err != nil {
		logger.Log.Error("[AllMetrics]", zap.Error(err))
	}
}
