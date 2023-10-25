package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	m "github.com/bonus2k/go-metrics-tpl/internal/models"
	"github.com/bonus2k/go-metrics-tpl/internal/server/repositories"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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

func (c *controller) SaveMetric(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("decoding request")
	var metric m.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metric); err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(m.KeyContentType, m.TypeJSONContent)
	switch strings.ToLower(metric.MType) {
	case "gauge":
		logger.Log.Debug("save", zap.Any("gauge", metric))
		c.mem.AddGauge(metric.ID, *metric.Value)
	case "counter":
		logger.Log.Debug("save", zap.Any("counter", metric))
		c.mem.AddCounter(metric.ID, *metric.Delta)
	default:
		logger.Log.Debug("default", zap.Any("metric", metric))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
}

func (c *controller) GetMetric(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("decoding request")
	var metric m.Metrics
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := dec.Decode(&metric); err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Log.Info("get metric", zap.Any("metric", metric))
	w.Header().Set(m.KeyContentType, m.TypeJSONContent)
	switch strings.ToLower(metric.MType) {
	case "gauge":
		gauge, ok := c.mem.GetGauge(metric.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		}
		metric.Value = &gauge
	case "counter":
		counter, ok := c.mem.GetCounter(metric.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		}
		metric.Delta = &counter
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
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
		c.mem.AddCounter(name, num)
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
	c.mem.AddGauge(name, num)
	w.WriteHeader(http.StatusOK)
}

func (c *controller) GetValue(w http.ResponseWriter, r *http.Request) {
	typeV := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	switch typeV {
	case "gauge":
		if gauge, ok := c.mem.GetGauge(name); !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Header().Set(m.KeyContentType, m.TypeHTMLContent)
			_, err := io.WriteString(w, fmt.Sprintf("%v", gauge))
			if err != nil {
				logger.Log.Error("[GetValue gauge]", zap.Error(err))
			}
		}
	case "counter":
		if counter, ok := c.mem.GetCounter(name); !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Header().Set(m.KeyContentType, m.TypeHTMLContent)
			_, err := io.WriteString(w, fmt.Sprintf("%v", counter))
			if err != nil {
				logger.Log.Error("[GetValue counter]", zap.Error(err))
			}
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (c *controller) AllMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := c.mem.GetAllMetrics()
	marshal, _ := json.Marshal(metrics)
	w.Header().Set(m.KeyContentType, m.TypeHTMLContent)
	_, err := w.Write(marshal)
	if err != nil {
		logger.Log.Error("[AllMetrics]", zap.Error(err))
	}
}

func (c *controller) Ping(w http.ResponseWriter, r *http.Request) {
	if err := c.mem.CheckConnection(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
