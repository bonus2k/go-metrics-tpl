package main

import (
	"flag"
	"time"
)

var connectAddr string
var reportInterval time.Duration
var pollInterval time.Duration

func parseFlags() {
	flag.StringVar(&connectAddr, "a", "localhost:8080", "address and port for connecting to server")
	flag.DurationVar(&reportInterval, "r", 10, "timer of report interval for send metrics")
	flag.DurationVar(&pollInterval, "p", 2, "timer of poll interval for metrics")
	flag.Parse()
}
