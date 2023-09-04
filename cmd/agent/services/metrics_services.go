package services

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
)

var mem runtime.MemStats
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

func GetMapMetrics() map[string]string {
	metrics := make(map[string]string)
	runtime.ReadMemStats(&mem)
	valuesMemStats := reflect.ValueOf(mem)
	for _, s := range memStatConst {
		value := valuesMemStats.FieldByName(s)
		metrics[s] = fmt.Sprintf("%v", value.Interface())
	}
	return metrics
}

func AddRandomValue(m map[string]string) {
	m["RandomValue"] = fmt.Sprintf("%.4f", rand.Float64()*100)
}

func GetPollCount() func() int64 {
	var count int64 = 0
	return func() int64 {
		count++
		return count
	}
}
