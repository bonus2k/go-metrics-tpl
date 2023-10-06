package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/agent/clients"
	"github.com/bonus2k/go-metrics-tpl/internal/agent/services"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/interface/rest"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/go-resty/resty/v2"
	"time"
)

func main() {
	mapMetrics := make(map[string]string)
	parseFlags()
	logger.Initialize(runLog)
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
	client := clients.Connect{Server: connectAddr, Protocol: "http", Client: *resty.New().SetPreRequestHook(rest.GzipReqCompression)}
	logger.Log.Info(fmt.Sprintf("Connect to server %s, report interval=%d, poll interval=%d", connectAddr, reportInterval, pollInterval))
	return func() {
		services.AddRandomValue(*m)
		client.SendToGauge(*m)
		client.SendToCounter("PollCount", count())
	}
}
