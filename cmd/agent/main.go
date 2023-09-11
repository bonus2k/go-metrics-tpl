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
	fmt.Fprintf(os.Stdout, "Connect to server %s, report interval=%d, poll interval=%d\n", connectAddr, reportInterval, pollInterval)
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
			if _, err := client.SendToGauge(mapMetrics); err != nil {
				fmt.Fprintf(os.Stdout, "[SendToGauge] Error %v\n", err)
			}
			if _, err := client.SendToCounter("PollCount", count()); err != nil {
				fmt.Fprintf(os.Stdout, "[SendToCounter] Error %v\n", err)
			}
		}
	}()
	time.Sleep(100 * time.Hour)
}
