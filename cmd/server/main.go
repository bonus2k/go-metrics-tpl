package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/server/controllers"
	"github.com/bonus2k/go-metrics-tpl/internal/server/repositories"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	if err := parseFlags(); err != nil {
		panic(err)
	}

	storage := *repositories.NewMemStorage(storeInterval == 0)
	if err := storage.CheckConnection(); err != nil {
		panic(err)
	}
	dbStorage, err := repositories.NewDbStorage(dbConn)
	if err != nil {
		logger.Log.Error("connect to db", zap.Error(err))
	}

	memService, err := repositories.NewMemStorageService(storeInterval, fileStore, runRestoreMetrics, &storage)
	if err != nil {
		panic(err)
	}

	if storeInterval != 0 {
		saveMemTicker := time.NewTicker(time.Duration(storeInterval) * time.Second)
		go func() {
			for range saveMemTicker.C {
				err := memService.Save()
				if err != nil {
					logger.Log.Error("save metrics ", zap.Error(err))
				}
			}
		}()
	}

	err = logger.Initialize(runLog)
	if err != nil {
		panic(err)
	}
	logger.Log.Info(fmt.Sprintf("Running server on %s log level %s", runAddr, runLog))

	err = http.ListenAndServe(runAddr, controllers.MetricsRouter(&storage, dbStorage))
	if err != nil {
		panic(err)
	}
}
