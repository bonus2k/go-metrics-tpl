package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/cmd/server/controllers"
	"net/http"
)

func main() {
	parseFlags()
	fmt.Println("Running server on", runAddr)
	err := http.ListenAndServe(runAddr, controllers.MetricsRouter())
	if err != nil {
		panic(err)
	}

}
