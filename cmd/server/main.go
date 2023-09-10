package main

import (
	"github.com/bonus2k/go-metrics-tpl/cmd/server/controllers"
	"net/http"
)

func main() {
	err := http.ListenAndServe(`:8080`, controllers.MetricsRouter())
	if err != nil {
		panic(err)
	}
}
