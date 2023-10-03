package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
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
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch strings.ToLower(metric.MType) {
	case "gauge":
		logger.Log.Info("save", zap.Any("gauge", metric))
		MemStorage.AddGauge(metric.ID, *metric.Value)
		w.WriteHeader(http.StatusOK)
		return
	case "counter":
		logger.Log.Info("save", zap.Any("counter", metric))
		MemStorage.AddCounter(metric.ID, *metric.Delta)
		w.WriteHeader(http.StatusOK)
		return
	default:
		logger.Log.Info("default", zap.Any("metric", metric))
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
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
}

func CounterPage(w http.ResponseWriter, r *http.Request) {
	if MemStorage == nil {
		panic("storage not initialized")
	}

	w.Header().Set("Content-Type", "text/plain")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if num, err := strconv.ParseInt(value, 10, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		MemStorage.AddCounter(name, num)
		w.WriteHeader(http.StatusOK)
		return
	}
}

func GaugePage(w http.ResponseWriter, r *http.Request) {
	if MemStorage == nil {
		panic("storage not initialized")
	}

	w.Header().Set("Content-Type", "text/plain")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	MemStorage.AddGauge(name, num)
	w.WriteHeader(http.StatusOK)
}

func GetValue(w http.ResponseWriter, r *http.Request) {
	typeV := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	switch typeV {
	case "gauge":
		if gauge, ok := MemStorage.GetGauge(name); !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			_, err := io.WriteString(w, fmt.Sprintf("%v", gauge))
			if err != nil {
				logger.Log.Error("[GetValue gauge]", zap.Error(err))
			}
		}
	case "counter":
		if counter, ok := MemStorage.GetCounter(name); !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			_, err := io.WriteString(w, fmt.Sprintf("%v", counter))
			if err != nil {
				logger.Log.Error("[GetValue counter]", zap.Error(err))
			}
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func AllMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := MemStorage.GetAllMetrics()
	marshal, _ := json.Marshal(metrics)
	_, err := w.Write(marshal)
	if err != nil {
		logger.Log.Error("[AllMetrics]", zap.Error(err))
	}
}
