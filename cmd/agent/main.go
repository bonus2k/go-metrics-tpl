package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/agent/clients"
	"github.com/bonus2k/go-metrics-tpl/internal/agent/services"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/interface/rest"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func main() {
	mapMetrics := make(map[string]string)
	parseFlags()
	logger.Initialize(runLog)
	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	sendReport := batchReport(&mapMetrics)
	for {
		select {
		case <-reportTicker.C:
			sendReport()
		case <-pollTicker.C:
			mapMetrics = services.GetMapMetrics()
		}
	}

}

func batchReport(mapMetrics *map[string]string) func() {
	count := services.GetPollCount()
	resty := *resty.New().
		SetPreRequestHook(rest.GzipReqCompression).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(9 * time.Second)
	client := clients.Connect{Server: connectAddr, Protocol: "http", Client: resty}
	logger.Log.Info(fmt.Sprintf("Connect to server %s, report interval=%d, poll interval=%d", connectAddr, reportInterval, pollInterval))
	return func() {
		services.AddRandomValue(*mapMetrics)
		metrics, err := convertGaugeToMetrics(mapMetrics)
		if err != nil {
			logger.Log.Error("can't convert gauge", zap.Error(err))
		}

		i := count()
		metrics = append(metrics, models.Metrics{ID: "PollCount", Delta: &i, MType: "counter"})
		err = client.SendBatchMetrics(metrics)

		if err != nil {
			logger.Log.Error("can't send metrics to updates", zap.Error(err))
		}
	}
}

func convertGaugeToMetrics(metrics *map[string]string) ([]models.Metrics, error) {
	listM := make([]models.Metrics, 0)
	for k, v := range *metrics {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		m := models.Metrics{ID: k, Value: &value, MType: "gauge"}
		listM = append(listM, m)
	}
	return listM, nil
}

func report(mapMetrics *map[string]string) func() {
	m := mapMetrics
	count := services.GetPollCount()
	client := clients.Connect{Server: connectAddr, Protocol: "http", Client: *resty.New().SetPreRequestHook(rest.GzipReqCompression)}
	logger.Log.Info(fmt.Sprintf("Connect to server %s, report interval=%d, poll interval=%d", connectAddr, reportInterval, pollInterval))
	return func() {
		services.AddRandomValue(*m)
		err := client.SendToGauge(*m)
		if err != nil {
			logger.Log.Error("can't send gauge metric", zap.Error(err))
		}
		client.SendToCounter("PollCount", count())
	}
}
