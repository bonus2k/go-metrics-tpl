// Agent осуществляет сбор метрик и их отправку на сервер
package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/agent/clients"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/interface/rest"
	"github.com/go-resty/resty/v2"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/bonus2k/go-metrics-tpl/internal/agent/services"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	_, err := fmt.Fprintf(os.Stdout, "Build version: %s \n", buildVersion)
	if err != nil {
		logger.Exit(err, 1)
	}
	_, err = fmt.Fprintf(os.Stdout, "Build date: %s \n", buildDate)
	if err != nil {
		logger.Exit(err, 1)
	}
	_, err = fmt.Fprintf(os.Stdout, "Build commit: %s \n", buildCommit)
	if err != nil {
		logger.Exit(err, 1)
	}

	if err = parseFlags(); err != nil {
		logger.Exit(err, 1)
	}

	err = logger.Initialize(runLog)
	if err != nil {
		logger.Exit(err, 1)
	}

	pollTicker := time.NewTicker(pollInterval)
	buf := getSizeBuf(reportInterval, pollInterval)
	resultErr := make(chan error)
	chanelMetrics := services.NewChanelMetrics(buf, resultErr)
	chanelMetrics.GetMetrics(pollTicker)
	chanelMetrics.GetPSUtilMetrics(pollTicker)
	chanelResult := chanelMetrics.GetChanelResult()

	sha256 := rest.NewSignSHA256(signPass)
	crypto, err := rest.NewEncrypt(cryptoKey)
	if err != nil {
		logger.Exit(err, 1)
	}
	res := resty.New().
		SetPreRequestHook(func(client *resty.Client, request *http.Request) error {
			err = crypto.EncryptRequest(client, request)
			if err != nil {
				return err
			}
			err = sha256.AddSignToReq(client, request)
			if err != nil {
				return err
			}
			err = rest.GzipReqCompression(client, request)
			if err != nil {
				return err
			}

			return nil
		}).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(9 * time.Second).
		SetCloseConnection(true)
	client := clients.Connect{Server: connectAddr, Protocol: "http", Client: res}

	pool := services.NewPool(client)
	for i := 0; i < rateLimitRoutines; i++ {
		num := i
		reportTicker := time.NewTicker(reportInterval)
		go func() {
			pool.BatchReport(chanelResult, resultErr, reportTicker, num)
		}()
	}

	if len(pprofAddr) != 0 {
		go func() {
			logger.Log.Infof("Run pprof on address %s", pprofAddr)
			err := http.ListenAndServe(pprofAddr, nil)
			logger.Log.Error("pprof server", err)
		}()
	}

	logger.Log.Infof("Connect to server %s, report interval=%d, poll interval=%d, rate limi=%d",
		connectAddr, reportInterval, pollInterval, rateLimitRoutines)
	for err := range resultErr {
		logger.Log.Error("Error poll worker send batch metrics", err)
	}
}

func getSizeBuf(reportInterval time.Duration, pollInterval time.Duration) int {
	i := reportInterval.Seconds() / pollInterval.Seconds()
	if i < 1 {
		i = 1
	}
	return int(i * 2)
}
