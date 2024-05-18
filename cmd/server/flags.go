package main

import (
	"encoding/json"
	"flag"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var runAddr string
var pprofAddr string
var runLog string
var storeInterval time.Duration
var fileStore string
var runRestoreMetrics bool
var dbConn string
var signPass string
var cryptoKey string
var configFile string

type config struct {
	Address        string `json:"address", omitempty`
	Restore        bool   `json:"restore", omitempty`
	Store_interval string `json:"store_interval", omitempty`
	Store_file     string `json:"store_file", omitempty`
	Database_dsn   string `json:"database_dsn", omitempty`
	Crypto_key     string `json:"crypto_key", omitempty`
}

var defaultStoreInterval = time.Minute * 5

func parseFlags() error {
	flag.StringVar(&runAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&configFile, "c", "", "path to config file")
	flag.StringVar(&cryptoKey, "crypto-key", "", "file with private key")
	flag.StringVar(&pprofAddr, "prof", "", "run pprof")
	flag.StringVar(&runLog, "l", "info", "log level")
	flag.StringVar(&dbConn, "d", "", "database name and connection information")
	flag.StringVar(&signPass, "k", "", "signature for HashSHA256")
	flag.DurationVar(&storeInterval, "i", defaultStoreInterval, "metrics saving interval")
	flag.StringVar(&fileStore, "f", "/tmp/metrics-db.json", "file path for saving metrics")
	flag.BoolVar(&runRestoreMetrics, "r", true, "restore metrics")
	flag.Parse()

	if err := parseEnv(); err != nil {
		return err
	}

	return nil
}

func parseEnv() error {

	if envConfigFile, ok := os.LookupEnv("CONFIG"); ok {
		configFile = envConfigFile
	}

	var conf config
	if configFile != "" {
		err := parseConfig(&conf)
		if err != nil {
			return err
		}
	}

	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		runAddr = envRunAddr
	} else if runAddr == "" {
		runAddr = conf.Address
	}

	if envCryptoKey, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		cryptoKey = envCryptoKey
	} else if cryptoKey == "" {
		cryptoKey = conf.Crypto_key
	}

	if envStoreInterval, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		s, err := time.ParseDuration(envStoreInterval)
		if err != nil {
			return errors.Wrap(err, "STORE_INTERVAL is not correct")
		}
		storeInterval = s
	} else if storeInterval == defaultStoreInterval {
		s, err := time.ParseDuration(conf.Store_interval)
		if err != nil {
			return errors.Wrap(err, "store_interval is not correct")
		}
		storeInterval = s
	}

	if envFileStore, ok := os.LookupEnv("STORE_FILE"); ok {
		fileStore = envFileStore
	} else if fileStore == "" {
		fileStore = conf.Store_file
	}

	if envRunRestoreMetrics, ok := os.LookupEnv("RESTORE"); ok {
		b, err := strconv.ParseBool(envRunRestoreMetrics)
		if err != nil {
			return errors.Wrap(err, "RESTORE is not correct")
		}
		runRestoreMetrics = b
	} else if !runRestoreMetrics {
		runRestoreMetrics = conf.Restore
	}

	if envDBConn, ok := os.LookupEnv("DATABASE_DSN"); ok {
		dbConn = envDBConn
	} else if dbConn == "" {
		dbConn = conf.Database_dsn
	}

	if envPprofAdr, ok := os.LookupEnv("PPROF_ADDRESS"); ok {
		pprofAddr = envPprofAdr
	}

	if envSignPass, ok := os.LookupEnv("KEY"); ok {
		signPass = envSignPass
	}

	return nil
}

func parseConfig(conf *config) error {
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.Log.Error("close file failed", err)
		}
	}()

	buf, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, conf)
	if err != nil {
		return err
	}

	return nil
}
