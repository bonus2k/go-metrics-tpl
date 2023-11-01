package main

import (
	"flag"
	"github.com/pkg/errors"
	"os"
	"strconv"
)

var runAddr string
var runLog string
var storeInterval int
var fileStore string
var runRestoreMetrics bool
var dbConn string
var signPass string

func parseFlags() error {
	flag.StringVar(&runAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&runLog, "l", "info", "log level")
	flag.StringVar(&dbConn, "d", "", "database name and connection information")
	flag.StringVar(&signPass, "k", "", "signature for HashSHA256")
	flag.IntVar(&storeInterval, "i", 300, "metrics saving interval")
	flag.StringVar(&fileStore, "f", "/tmp/metrics-db.json", "file path for saving metrics")
	flag.BoolVar(&runRestoreMetrics, "r", true, "restore metrics")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		runAddr = envRunAddr
	}
	if envSignPass := os.Getenv("KEY"); envSignPass != "" {
		signPass = envSignPass
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		var err error
		storeInterval, err = strconv.Atoi(envStoreInterval)
		if err != nil {
			return errors.Wrap(err, "STORE_INTERVAL is not correct")
		}
	}
	if envFileStore := os.Getenv("FILE_STORAGE_PATH"); envFileStore != "" {
		fileStore = envFileStore
	}
	if envRunRestoreMetrics := os.Getenv("RESTORE"); envRunRestoreMetrics != "" {
		b, err := strconv.ParseBool(envRunRestoreMetrics)
		if err != nil {
			return errors.Wrap(err, "RESTORE is not correct")
		}
		runRestoreMetrics = b
	}
	if envDBConn := os.Getenv("DATABASE_DSN"); envDBConn != "" {
		dbConn = envDBConn
	}

	return nil
}
