package controllers

import (
	"github.com/bonus2k/go-metrics-tpl/internal/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/server/repositories"
	"github.com/go-chi/chi/v5"
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
	})
	router.Get("/", AllMetrics)
	return router
}
