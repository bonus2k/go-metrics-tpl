package main

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"log"
	"os"
	"strconv"
)

var connectAddr string
var reportInterval int
var pollInterval int
var runLog string
var signPass string
var rateLimitRoutines int

type config struct {
	ConnectAddr    string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func parseFlags() {
	var cfg config
	flag.StringVar(&connectAddr, "a", "localhost:8080", "address and port for connecting to server")
	flag.IntVar(&reportInterval, "r", 10, "timer of report interval for send metrics")
	flag.IntVar(&rateLimitRoutines, "l", 1, "count of routines")
	flag.IntVar(&pollInterval, "p", 2, "timer of poll interval for metrics")
	flag.StringVar(&runLog, "log", "info", "log level")
	flag.StringVar(&signPass, "k", "", "signature for HashSHA256")
	flag.Parse()

	opts := env.Options{
		OnSet: func(tag string, value interface{}, isDefault bool) {
			if value == "" {
				return
			}
			switch tag {
			case "ADDRESS":
				connectAddr = value.(string)
			case "REPORT_INTERVAL":
				reportInterval, _ = strconv.Atoi(value.(string))
			case "POLL_INTERVAL":
				pollInterval, _ = strconv.Atoi(value.(string))
			default:
				return
			}
		},
	}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		log.Fatal(err)
	}
	if envSignPass := os.Getenv("KEY"); envSignPass != "" {
		signPass = envSignPass
	}
	if envRateLimitRoutines := os.Getenv("RATE_LIMIT"); envRateLimitRoutines != "" {
		rateLimitRoutines, _ = strconv.Atoi(envRateLimitRoutines)
	}
}
