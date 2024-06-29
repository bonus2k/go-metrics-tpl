package controllers

import (
	"net/http"

	"github.com/bonus2k/go-metrics-tpl/internal/middleware/interface/rest"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/server/repositories"
	"github.com/go-chi/chi/v5"
)

// MetricsRouter создает роутер для HTTP сервера
func MetricsRouter(mem *repositories.Storage, pass string, file string, trustNet string) chi.Router {
	ctrl := NewController(mem)
	router := chi.NewRouter()
	sha256 := rest.NewSignSHA256(pass)
	decrypt, err := rest.NewDecrypt(file)
	subnet := rest.NewTrustSubnet(trustNet)
	if err != nil {
		logger.Exit(err, 1)
	}
	router.Use(
		subnet.CheckRealIp,
		rest.GzipReqDecompression1,
		rest.GzipResCompression,
		logger.MiddlewareLog,
		sha256.AddSignToRes,
		sha256.CheckSignReq,
		decrypt.DecryptRequest,
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
