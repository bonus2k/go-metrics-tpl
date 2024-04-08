package main

import (
	"flag"
	"github.com/pkg/errors"

	"os"
	"strconv"
)

var connectAddr string
var pprofAddr string
var reportInterval int
var pollInterval int
var runLog string
var signPass string
var rateLimitRoutines int

func parseFlags() error {
	flag.StringVar(&connectAddr, "a", "localhost:8080", "address and port for connecting to server")
	flag.StringVar(&pprofAddr, "prof", "", "run pprof")
	flag.IntVar(&reportInterval, "r", 10, "timer of report interval for send metrics")
	flag.IntVar(&rateLimitRoutines, "l", 1, "count of routines")
	flag.IntVar(&pollInterval, "p", 2, "timer of poll interval for metrics")
	flag.StringVar(&runLog, "log", "info", "log level")
	flag.StringVar(&signPass, "k", "", "signature for HashSHA256")
	flag.Parse()

	if envConnectAddr := os.Getenv("ADDRESS"); envConnectAddr != "" {
		connectAddr = envConnectAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		r, err := strconv.Atoi(envReportInterval)
		if err != nil {
			return errors.Wrap(err, "REPORT_INTERVAL is not correct")
		}
		reportInterval = r
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		p, err := strconv.Atoi(envPollInterval)
		if err != nil {
			return errors.Wrap(err, "POLL_INTERVAL is not correct")
		}
		pollInterval = p
	}

	if envSignPass := os.Getenv("KEY"); envSignPass != "" {
		signPass = envSignPass
	}
	if envRateLimitRoutines := os.Getenv("RATE_LIMIT"); envRateLimitRoutines != "" {
		rateLimitRoutines, _ = strconv.Atoi(envRateLimitRoutines)
	}
	if envPprofAdr := os.Getenv("PPROF_ADDRESS"); envPprofAdr != "" {
		pprofAddr = envPprofAdr
	}
	return nil
}
