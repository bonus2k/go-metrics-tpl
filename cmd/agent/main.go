package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/cmd/agent/clients"
	"github.com/bonus2k/go-metrics-tpl/cmd/agent/services"
	"os"
	"time"
)

var mapMetrics map[string]string

func main() {
	parseFlags()
	count := services.GetPollCount()
	client := clients.Connect{Server: connectAddr, Protocol: "http"}
	fmt.Fprintf(os.Stdout, "Connect to server %s\n", connectAddr)
	go func() {
		for {
			mapMetrics = services.GetMapMetrics()
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(reportInterval) * time.Second)
			services.AddRandomValue(mapMetrics)
			client.SendToGauge(mapMetrics)
			client.SendToCounter("PollCount", count())
		}
	}()
	time.Sleep(100 * time.Hour)
}
