package main

import (
	"github.com/bonus2k/go-metrics-tpl/cmd/agent/clients"
	"github.com/bonus2k/go-metrics-tpl/cmd/agent/services"
	"time"
)

var mapMetrics map[string]string

func main() {
	count := services.GetPollCount()
	client := clients.Connect{Server: "localhost", Port: "8080", Protocol: "http"}
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
			client.SendToGauge(mapMetrics)
			client.SendToCounter("PollCount", count())
		}
	}()
	time.Sleep(100 * time.Hour)
}
