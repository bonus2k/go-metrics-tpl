package main

import (
	"flag"
)

var connectAddr string
var reportInterval int
var pollInterval int

func parseFlags() {
	flag.StringVar(&connectAddr, "a", "localhost:8080", "address and port for connecting to server")
	flag.IntVar(&reportInterval, "r", 10, "timer of report interval for send metrics")
	flag.IntVar(&pollInterval, "p", 2, "timer of poll interval for metrics")
	flag.Parse()
}
