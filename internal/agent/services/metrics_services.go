package services

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"reflect"
	"runtime"
	"sync/atomic"
	"time"
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

type ChanelMetrics struct {
	outResult chan map[string]string
	outError  chan error
}

func NewChanelMetrics(sizeBuf int, outError chan error) *ChanelMetrics {
	out := make(chan map[string]string, sizeBuf)
	return &ChanelMetrics{outResult: out, outError: outError}
}

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

func (ch *ChanelMetrics) GetMetrics(ticker *time.Ticker) {
	go func() {
		for range ticker.C {
			ch.outResult <- GetMapMetrics()
			logger.Log.Info("send metrics to out chan")
		}
	}()
}

func (ch *ChanelMetrics) GetPSUtilMetrics(ticker *time.Ticker) {
	go func() {
		for range ticker.C {
			metrics, err := GetGoPSUtilMapMetrics()
			if err != nil {
				ch.outError <- fmt.Errorf("can't get PSUtil metrics, %w", err)
				continue
			}
			ch.outResult <- metrics
			logger.Log.Info("send PSUtil metrics to out chan")
		}
	}()
}

func (ch *ChanelMetrics) GetChanelResult() chan map[string]string {
	return ch.outResult
}

func GetPollCount() func() int64 {
	var count atomic.Int64
	return func() int64 {
		count.Add(1)
		return count.Load()
	}
}
