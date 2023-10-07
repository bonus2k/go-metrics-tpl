package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/server/controllers"
	"github.com/bonus2k/go-metrics-tpl/internal/server/repositories"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	parseFlags()
	memService := repositories.NewMemStorageService(storeInterval, fileStore, runRestoreMetrics)
	saveMemTicker := time.NewTicker(time.Duration(storeInterval) * time.Second)
	go saveMetrics(saveMemTicker, memService)
	var err error
	err = multierr.Append(err, logger.Initialize(runLog))
	logger.Log.Info(fmt.Sprintf("Running server on %s log level %s", runAddr, runLog))
	err = multierr.Append(err, http.ListenAndServe(runAddr, controllers.MetricsRouter()))
	if err != nil {
		panic(err)
	}

}

func saveMetrics(saveMemTicker *time.Ticker, memService *repositories.MemStorageService) {
	for {
		select {
		case <-saveMemTicker.C:
			err := memService.Save()
			logger.Log.Error("save metrics", zap.Error(err))
		}
	}
}
