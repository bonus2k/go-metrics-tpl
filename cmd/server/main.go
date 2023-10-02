package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/server/controllers"
	"net/http"
)

func main() {
	parseFlags()
	err := logger.Initialize(runLog)
	if err != nil {
		panic(err)
	}
	logger.Log.Info(fmt.Sprintf("Running server on %s log level %s", runAddr, runLog))
	err = http.ListenAndServe(runAddr, controllers.MetricsRouter())
	if err != nil {
		panic(err)
	}

}
