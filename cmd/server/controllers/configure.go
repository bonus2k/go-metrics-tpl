package controllers

import (
	"github.com/bonus2k/go-metrics-tpl/cmd/server/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var MemStorage repositories.MemStorage

func init() {
	MemStorage = repositories.NewMemStorage()
}

func MetricsRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{name}/{value}", GaugePage)
		r.Post("/counter/{name}/{value}", CounterPage)
	})
	router.Get("/", AllMetrics)
	router.Get("/value/{type}/{name}", GetValue)
	return router
}
