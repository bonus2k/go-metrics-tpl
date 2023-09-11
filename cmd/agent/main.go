package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/cmd/agent/clients"
	"github.com/bonus2k/go-metrics-tpl/cmd/agent/services"
	"time"
)

var mapMetrics map[string]string

func main() {
	parseFlags()
	count := services.GetPollCount()
	client := clients.Connect{Server: connectAddr, Protocol: "http"}
	fmt.Println("Connect to server", connectAddr)
	go func() {
		for {
			mapMetrics = services.GetMapMetrics()
			time.Sleep(pollInterval * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(reportInterval * time.Second)
			services.AddRandomValue(mapMetrics)
			client.SendToGauge(mapMetrics)
			client.SendToCounter("PollCount", count())
		}
	}()
	time.Sleep(100 * time.Hour)
}
