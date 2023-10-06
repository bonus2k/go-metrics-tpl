package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/server/controllers"
	"go.uber.org/multierr"
	"net/http"
)

func main() {
	parseFlags()
	var err error
	err = multierr.Append(err, logger.Initialize(runLog))
	logger.Log.Info(fmt.Sprintf("Running server on %s log level %s", runAddr, runLog))
	err = multierr.Append(err, http.ListenAndServe(runAddr, controllers.MetricsRouter()))
	if err != nil {
		panic(err)
	}
}
