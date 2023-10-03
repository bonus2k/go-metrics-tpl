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
	router.Route("/", func(r chi.Router) {
		r.Post("/update/", SaveMetric)
		r.Post("/value/", GetMetric)
		r.Post("/update/gauge/{name}/{value}", GaugePage)
		r.Post("/update/counter/{name}/{value}", CounterPage)
		r.Post("/update/*", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
		})
		r.Post("/update/gauge/*", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		})
		r.Post("/update/counter/*", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
		})

	})
	router.Get("/", AllMetrics)
	router.Get("/value/{type}/{name}", GetValue)
	return router
}
