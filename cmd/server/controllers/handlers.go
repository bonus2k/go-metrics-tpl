package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
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

	if num, err := strconv.ParseFloat(value, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		MemStorage.AddGauge(name, num)
		w.WriteHeader(http.StatusOK)
		return
	}
}

func GetValue(w http.ResponseWriter, r *http.Request) {
	typeV := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	switch typeV {
	case "gauge":
		if gauge, ok := MemStorage.GetGauge(name); !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			io.WriteString(w, fmt.Sprintf("%v", gauge))
		}
	case "counter":
		if counter, ok := MemStorage.GetCounter(name); !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			io.WriteString(w, fmt.Sprintf("%v", counter))
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func AllMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := MemStorage.GetAllMetrics()
	marshal, _ := json.Marshal(metrics)
	w.Write(marshal)
}
