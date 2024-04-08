package main

import (
	"github.com/bonus2k/go-metrics-tpl/internal/agent/services"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

func main() {
	if err := parseFlags(); err != nil {
		os.Exit(1)
	}

	err := logger.Initialize(runLog)
	if err != nil {
		os.Exit(1)
	}

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

func getSizeBuf(reportInterval int, pollInterval int) int {
	i := reportInterval / pollInterval
	if i < 1 {
		i = 1
	}
	return i * 2
}
