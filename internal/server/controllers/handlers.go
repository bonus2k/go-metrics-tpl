package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

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
