package controllers

import (
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/interface/rest"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/server/repositories"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func MetricsRouter(mem *repositories.Storage, pass string) chi.Router {
	var sha256 *rest.SignSHA256
	ctrl := NewController(mem)
	router := chi.NewRouter()
	if pass != "" {
		sha256 = rest.NewSignSHA256(pass)
	}
	router.Use(
		rest.GzipReqDecompression,
		rest.GzipResCompression,
		logger.MiddlewareLog,
		sha256.AddSignToRes,
		sha256.CheckSignReq,
	)
	router.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{name}/{value}", ctrl.GaugePage)
		r.Post("/counter/{name}/{value}", ctrl.CounterPage)
		r.Post("/*", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
		})
		r.Post("/gauge/*", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		})
		r.Post("/counter/*", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		})

	})
	router.Post("/update/", ctrl.SaveMetric)
	router.Post("/updates/", ctrl.SaveMetrics)
	router.Post("/value/", ctrl.GetMetric)
	router.Get("/", ctrl.AllMetrics)
	router.Get("/ping", ctrl.Ping)
	router.Get("/value/{type}/{name}", ctrl.GetValue)
	return router
}
