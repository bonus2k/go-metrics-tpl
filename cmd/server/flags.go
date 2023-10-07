package main

import (
	"flag"
	"os"
	"strconv"
)

var runAddr string
var runLog string
var storeInterval int
var fileStore string
var runRestoreMetrics bool

func parseFlags() {
	flag.StringVar(&runAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&runLog, "l", "info", "log level")
	flag.IntVar(&storeInterval, "i", 300, "metrics saving interval")
	flag.StringVar(&fileStore, "f", "/tmp/metrics-db.json", "file path for saving metrics")
	flag.BoolVar(&runRestoreMetrics, "r", true, "restore metrics")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		runAddr = envRunAddr
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		storeInterval, _ = strconv.Atoi(envStoreInterval)
	}
	if envFileStore := os.Getenv("FILE_STORAGE_PATH"); envFileStore != "" {
		fileStore = envFileStore
	}
	if envRunRestoreMetrics := os.Getenv("RESTORE"); envRunRestoreMetrics != "" {
		b, err := strconv.ParseBool(envRunRestoreMetrics)
		if err == nil {
			runRestoreMetrics = b
		}
	}

}
