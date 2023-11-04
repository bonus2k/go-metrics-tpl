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
	"net/http"
	"time"
)

func main() {
	parseFlags()
	logger.Initialize(runLog)
	//reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	//По заданию немного не понятно как должен вести себя каждый воркер относительно
	//таймаута reportTicker:
	//Для каждого воркера существует свой тикер или для всех один.
	//Если для всех один то отподает необходимость в буферезированной очереди,
	//так как сборщик статистики упирается в блокируемую очередь
	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	buf := getSizeBuf(reportInterval, pollInterval)
	resultErr := make(chan error)
	chanelMetrics := services.NewChanelMetrics(buf, resultErr)
	chanelMetrics.GetMetrics(pollTicker)
	chanelMetrics.GetPSUtilMetrics(pollTicker)
	chanelResult := chanelMetrics.GetChanelResult()

	pool := services.NewPool(signPass, connectAddr)
	for i := 0; i < rateLimitRoutines; i++ {
		num := i
		reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
		go func() {
			pool.BatchReport(chanelResult, resultErr, reportTicker, num)
		}()
	}
	logger.Log.Info(fmt.Sprintf("Connect to server %s, report interval=%d, poll interval=%d, rate limi=%d",
		connectAddr, reportInterval, pollInterval, rateLimitRoutines))
	for err := range resultErr {
		logger.Log.Error("Error poll worker send batch metrics", zap.Error(err))
	}
}

func getSizeBuf(reportInterval int, pollInterval int) int {
	i := reportInterval / pollInterval
	if i < 1 {
		i = 1
	}
	return i * 2
}

func batchReport(mapMetrics *map[string]string, pass string) func() {
	count := services.GetPollCount()
	sha256 := rest.NewSignSHA256(pass)
	res := resty.New().
		SetPreRequestHook(func(client *resty.Client, request *http.Request) error {
			err := sha256.AddSignToReq(client, request)
			if err != nil {
				return err
			}
			return rest.GzipReqCompression(client, request)
		}).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(9 * time.Second)
	client := clients.Connect{Server: connectAddr, Protocol: "http", Client: res}
	logger.Log.Info(fmt.Sprintf("Connect to server %s, report interval=%d, poll interval=%d", connectAddr, reportInterval, pollInterval))
	return func() {
		metrics, err := models.ConvertGaugeToMetrics(mapMetrics)
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

func report(mapMetrics *map[string]string) func() {
	m := mapMetrics
	count := services.GetPollCount()
	client := clients.Connect{Server: connectAddr, Protocol: "http", Client: resty.New().SetPreRequestHook(rest.GzipReqCompression)}
	logger.Log.Info(fmt.Sprintf("Connect to server %s, report interval=%d, poll interval=%d", connectAddr, reportInterval, pollInterval))
	return func() {
		err := client.SendToGauge(*m)
		if err != nil {
			logger.Log.Error("can't send gauge metric", zap.Error(err))
		}
		client.SendToCounter("PollCount", count())
	}
}
