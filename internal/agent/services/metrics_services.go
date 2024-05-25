// Package services реализует работу worker и сбор метрик
package services

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

var memStats runtime.MemStats
var memStatConst = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

// ChanelMetrics канал для получения метрик worker
type ChanelMetrics struct {
	outResult    chan map[string]string
	outError     chan error
	metricsGroup sync.WaitGroup
	shutDown     bool
}

// NewChanelMetrics создает канал для получения метрик с задданым буфером и каналом для отправки ошибок
func NewChanelMetrics(sizeBuf int, outError chan error) *ChanelMetrics {
	out := make(chan map[string]string, sizeBuf)
	return &ChanelMetrics{outResult: out, outError: outError}
}

// GetMapMetrics считывает метрики согласно memStatConst и возвращает Map где ключ - имя метрики
func GetMapMetrics() map[string]string {
	metrics := make(map[string]string)
	runtime.ReadMemStats(&memStats)
	valuesMemStats := reflect.ValueOf(memStats)
	for _, s := range memStatConst {
		value := valuesMemStats.FieldByName(s)
		metrics[s] = fmt.Sprintf("%v", value.Interface())
	}

	metrics["RandomValue"] = fmt.Sprintf("%.4f", rand.Float64()*100)
	return metrics
}

// GetGoPSUtilMapMetrics предоставляет метрики по количеству/занятой RAM и утилизации CPU
func GetGoPSUtilMapMetrics() (map[string]string, error) {
	metrics := make(map[string]string)
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	cpuStat, err := cpu.Times(true)
	if err != nil {
		return nil, err
	}
	metrics["TotalMemory"] = fmt.Sprintf("%d", vmStat.Total)
	metrics["FreeMemory"] = fmt.Sprintf("%d", vmStat.Free)
	for i, stat := range cpuStat {
		cpuName := fmt.Sprintf("CPUutilization%d", i)
		metrics[cpuName] = fmt.Sprintf("%f", stat.Idle)
	}
	return metrics, nil
}

// GetMetrics собирает метрики используя GetMapMetrics и отправляет их worker
func (ch *ChanelMetrics) GetMetrics(ticker *time.Ticker) {
	ch.metricsGroup.Add(1)
	go func() {
	out:
		for range ticker.C {
			if !ch.shutDown {
				ch.outResult <- GetMapMetrics()
				logger.Log.Info("send metrics to out chan")
			} else {
				ch.outResult <- nil
				ch.metricsGroup.Done()
				break out
			}
		}
		logger.Log.Info("Metrics shutdown")
	}()
}

// GetPSUtilMetrics собирает метрики используя GetGoPSUtilMapMetric и отправляет их worker
func (ch *ChanelMetrics) GetPSUtilMetrics(ticker *time.Ticker) {
	ch.metricsGroup.Add(1)
	go func() {
	out:
		for range ticker.C {
			if !ch.shutDown {
				metrics, err := GetGoPSUtilMapMetrics()
				if err != nil {
					ch.outError <- fmt.Errorf("can't get PSUtil metrics, %w", err)
					continue
				}
				ch.outResult <- metrics
				logger.Log.Info("send PSUtil metrics to out chan")
			} else {
				ch.outResult <- nil
				ch.metricsGroup.Done()
				break out
			}
		}
		logger.Log.Info("PSUtilMetrics shutdown")
	}()
}

// GetChanelResult возвращает канал с метриками
func (ch *ChanelMetrics) GetChanelResult() chan map[string]string {
	return ch.outResult
}

func (ch *ChanelMetrics) Shutdown() {
	ch.shutDown = true
	ch.metricsGroup.Wait()
}

// GetPollCount инкрементирует метрику PollCount
func GetPollCount() func() int64 {
	var count atomic.Int64
	return func() int64 {
		count.Add(1)
		return count.Load()
	}
}
