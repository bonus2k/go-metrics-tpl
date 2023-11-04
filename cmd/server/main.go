package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/server/controllers"
	"github.com/bonus2k/go-metrics-tpl/internal/server/migrations"
	"github.com/bonus2k/go-metrics-tpl/internal/server/repositories"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	if err := parseFlags(); err != nil {
		panic(err)
	}
	var storage *repositories.Storage

	if len(dbConn) != 0 {
		err := migrations.Start(dbConn)
		if err != nil {
			logger.Log.Error("error migration db", zap.Error(err))
			panic(err)
		}
		storage, err = repositories.NewDBStorage(dbConn)
		if err != nil {
			logger.Log.Error("connect to db", zap.Error(err))
		}
	} else {
		storage = repositories.NewMemStorage(storeInterval == 0)
		memService, err := repositories.NewMemStorageService(storeInterval, fileStore, runRestoreMetrics, storage)
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
	}

	err := logger.Initialize(runLog)
	if err != nil {
		panic(err)
	}
	logger.Log.Info(fmt.Sprintf("Running server on %s log level %s", runAddr, runLog))

	err = http.ListenAndServe(runAddr, controllers.MetricsRouter(storage, signPass))
	if err != nil {
		panic(err)
	}
}
