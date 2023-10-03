package controllers

import (
	"github.com/bonus2k/go-metrics-tpl/internal/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/server/repositories"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var MemStorage repositories.MemStorage

func init() {
	MemStorage = repositories.NewMemStorage()
}

func MetricsRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(logger.MiddlewareLog)
	router.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{name}/{value}", GaugePage)
		r.Post("/counter/{name}/{value}", CounterPage)
		r.Post("/*", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
		})
		r.Post("/update/gauge/*", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		})
		r.Post("/update/counter/*", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		})

	})
	router.Post("/update/", SaveMetric)
	router.Post("/value/", GetMetric)
	router.Get("/", AllMetrics)
	router.Get("/value/{type}/{name}", GetValue)
	return router
}
