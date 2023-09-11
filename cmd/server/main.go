package main

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/cmd/server/controllers"
	"net/http"
	"os"
)

func main() {
	parseFlags()
	fmt.Fprintf(os.Stdout, "Running server on %s\n", runAddr)
	err := http.ListenAndServe(runAddr, controllers.MetricsRouter())
	if err != nil {
		panic(err)
	}

}
