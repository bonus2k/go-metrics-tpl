// Server принимает метрики и сохраняет их в БД или в JSON файле в зависимости от настройки
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/server/controllers"
	"github.com/bonus2k/go-metrics-tpl/internal/server/migrations"
	"github.com/bonus2k/go-metrics-tpl/internal/server/repositories"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	idleConnsClosed := make(chan struct{})
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

	runtime.SetCPUProfileRate(0)
	if err = parseFlags(); err != nil {
		logger.Exit(err, 1)
	}

	err = logger.Initialize(runLog)
	if err != nil {
		logger.Exit(err, 1)
	}

	var storage *repositories.Storage

	if len(dbConn) != 0 {
		err := migrations.Start(dbConn)
		if err != nil {
			logger.Exit(err, 1)
		}
		storage, err = repositories.NewDBStorage(dbConn)
		if err != nil {
			logger.Log.Error("connect to db", err)
		}
	} else {
		storage = repositories.NewMemStorage(storeInterval == 0)
		memService, err := repositories.NewMemStorageService(storeInterval, fileStore, runRestoreMetrics, storage)
		if err != nil {
			logger.Exit(err, 1)
		}

		if storeInterval != 0 {
			saveMemTicker := time.NewTicker(time.Duration(storeInterval) * time.Second)
			go func() {
				for range saveMemTicker.C {
					err := memService.Save()
					if err != nil {
						logger.Log.Error("save metrics ", err)
					}
				}
			}()
		}
	}

	if len(pprofAddr) != 0 {
		go func() {
			logger.Log.Infof("Run pprof on address %s", pprofAddr)
			err := http.ListenAndServe(pprofAddr, nil)
			logger.Log.Error("pprof server", err)
		}()
	}

	var srv = http.Server{Addr: runAddr, Handler: controllers.MetricsRouter(storage, signPass, cryptoKey)}
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigint
		err = srv.Shutdown(context.Background())
		close(idleConnsClosed)
	}()

	logger.Log.Infof("Running server on %s log level %s", runAddr, runLog)
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Exit(err, 1)
	}

	<-idleConnsClosed
	logger.Log.Info("server shutdown gracefully")
}
