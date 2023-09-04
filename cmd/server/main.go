package main

import (
	"github.com/bonus2k/go-metrics-tpl/cmd/server/controllers"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
	})
	mux.HandleFunc(`/update/gauge/`, controllers.GaugePage)
	mux.HandleFunc(`/update/counter/`, controllers.CounterPage)
	return http.ListenAndServe(`:8080`, mux)
}
