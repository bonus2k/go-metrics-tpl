package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/cmd/agent/clients"
	"github.com/bonus2k/go-metrics-tpl/cmd/agent/services"
	"os"
	"sync"
	"time"
)

var mapMetrics map[string]string

func main() {
	parseFlags()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			mapMetrics = services.GetMapMetrics()
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}()

	go func() {
		count := services.GetPollCount()
		client := clients.Connect{Server: connectAddr, Protocol: "http"}
		fmt.Fprintf(os.Stdout, "Connect to server %s, report interval=%d, poll interval=%d\n", connectAddr, reportInterval, pollInterval)
		defer wg.Done()
		for {
			time.Sleep(time.Duration(reportInterval) * time.Second)
			services.AddRandomValue(mapMetrics)
			if _, err := client.SendToGauge(mapMetrics); err != nil {
				fmt.Fprintf(os.Stderr, "[SendToGauge] Error %v\n", err)
			}
			if _, err := client.SendToCounter("PollCount", count()); err != nil {
				fmt.Fprintf(os.Stderr, "[SendToCounter] Error %v\n", err)
			}
		}
	}()
	wg.Wait()
}
