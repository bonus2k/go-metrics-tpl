package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/agent/clients"
	"github.com/bonus2k/go-metrics-tpl/internal/agent/services"
	"net/http"
	"os"
	"time"
)

func main() {
	mapMetrics := make(map[string]string)
	parseFlags()
	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	sendReport := report(&mapMetrics)
	for {
		select {
		case <-reportTicker.C:
			sendReport()
		case <-pollTicker.C:
			mapMetrics = services.GetMapMetrics()
		}
	}

}

func report(mapMetrics *map[string]string) func() {
	m := mapMetrics
	count := services.GetPollCount()
	client := clients.Connect{Server: connectAddr, Protocol: "http", Client: *http.DefaultClient}
	fmt.Fprintf(os.Stdout, "Connect to server %s, report interval=%d, poll interval=%d\n", connectAddr, reportInterval, pollInterval)
	return func() {
		services.AddRandomValue(*m)
		if _, err := client.SendToGauge(*m); err != nil {
			fmt.Fprintf(os.Stderr, "[SendToGauge] Error %v\n", err)
		}
		if _, err := client.SendToCounter("PollCount", count()); err != nil {
			fmt.Fprintf(os.Stderr, "[SendToCounter] Error %v\n", err)
		}
	}
}
