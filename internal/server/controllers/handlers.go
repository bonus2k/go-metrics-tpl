package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	m "github.com/bonus2k/go-metrics-tpl/internal/models"
	"github.com/bonus2k/go-metrics-tpl/internal/server/repositories"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type controller struct {
	mem repositories.Storage
}

func NewController(mem *repositories.Storage) *controller {
	return &controller{mem: *mem}
}

func (c *controller) SaveMetrics(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("decoding request")
	metrics := make([]m.Metrics, 0)
	dec := json.NewDecoder(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			logger.Log.Error("SaveMetrics", err)
		}
	}()
	if err := dec.Decode(&metrics); err != nil {
		logger.Log.Error("cannot decode request JSON body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := c.mem.AddMetrics(r.Context(), metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *controller) SaveMetric(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("decoding request")
	var metric m.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metric); err != nil {
		logger.Log.Error("cannot decode request JSON body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(m.KeyContentType, m.TypeJSONContent)
	switch strings.ToLower(metric.MType) {
	case "gauge":
		logger.Log.Debugf("save gauge metric %v", metric)
		err := c.mem.AddGauge(r.Context(), metric.ID, *metric.Value)
		if err != nil {
			logger.Log.Error("can't save gauge", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "counter":
		logger.Log.Debugf("save counter metric %v", metric)
		err := c.mem.AddCounter(r.Context(), metric.ID, *metric.Delta)
		if err != nil {
			logger.Log.Error("can't save counter", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		logger.Log.Debugf("type of metrics not reconcile %v", metric)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		logger.Log.Error("error encoding response", err)
		return
	}
}

func (c *controller) GetMetric(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("decoding request")
	var metric m.Metrics
	dec := json.NewDecoder(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			logger.Log.Error("GetMetric", err)
		}
	}()
	if err := dec.Decode(&metric); err != nil {
		logger.Log.Error("cannot decode request JSON body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Log.Infof("got metric %v", metric)
	w.Header().Set(m.KeyContentType, m.TypeJSONContent)
	switch strings.ToLower(metric.MType) {
	case "gauge":
		gauge, err := c.mem.GetGauge(r.Context(), metric.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric.Value = &gauge
	case "counter":
		counter, err := c.mem.GetCounter(r.Context(), metric.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric.Delta = &counter
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		logger.Log.Error("error encoding response", err)
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
}

func (c *controller) CounterPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(m.KeyContentType, m.TypeHTMLContent)
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if num, err := strconv.ParseInt(value, 10, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		err := c.mem.AddCounter(r.Context(), name, num)
		if err != nil {
			logger.Log.Error("can't save counter", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (c *controller) GaugePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(m.KeyContentType, m.TypeHTMLContent)
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.mem.AddGauge(r.Context(), name, num)
	if err != nil {
		logger.Log.Error("can't save gauge", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *controller) GetValue(w http.ResponseWriter, r *http.Request) {
	typeV := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	switch typeV {
	case "gauge":
		if gauge, err := c.mem.GetGauge(r.Context(), name); err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Header().Set(m.KeyContentType, m.TypeHTMLContent)
			_, err := io.WriteString(w, fmt.Sprintf("%v", gauge))
			if err != nil {
				logger.Log.Error("[GetValue gauge]", err)
			}
		}
	case "counter":
		if counter, err := c.mem.GetCounter(r.Context(), name); err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Header().Set(m.KeyContentType, m.TypeHTMLContent)
			_, err := io.WriteString(w, fmt.Sprintf("%v", counter))
			if err != nil {
				logger.Log.Error("[GetValue counter]", err)
			}
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (c *controller) AllMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := c.mem.GetAllMetrics(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	marshal, _ := json.Marshal(metrics)
	w.Header().Set(m.KeyContentType, m.TypeHTMLContent)
	_, err = w.Write(marshal)
	if err != nil {
		logger.Log.Error("[AllMetrics]", err)
	}
}

func (c *controller) Ping(w http.ResponseWriter, r *http.Request) {
	if err := c.mem.CheckConnection(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
