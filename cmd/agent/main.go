package main

import (
	"github.com/bonus2k/go-metrics-tpl/cmd/agent/clients"
	"github.com/bonus2k/go-metrics-tpl/cmd/agent/services"
	"time"
)

var mapMetrics map[string]string

func main() {
	count := services.GetPollCount()
	go func() {
		for {
			mapMetrics = services.GetMapMetrics()
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			services.AddRandomValue(mapMetrics)
			clients.SendToGauge(mapMetrics)
			clients.SendToCounter("PollCount", count())
		}
	}()
	time.Sleep(100 * time.Hour)
}
