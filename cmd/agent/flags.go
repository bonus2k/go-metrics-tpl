package main

import (
	"encoding/json"
	"flag"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"io"
	"time"

	"github.com/pkg/errors"

	"os"
	"strconv"
)

var connectAddr string
var pprofAddr string
var reportInterval time.Duration
var pollInterval time.Duration
var runLog string
var signPass string
var rateLimitRoutines int
var cryptoKey string
var configFile string

type config struct {
	Address        string `json:"address,omitempty"`
	ReportInterval string `json:"report_interval,omitempty"`
	PollInterval   string `json:"poll_interval,omitempty"`
	CryptoKey      string `json:"crypto_key,omitempty"`
}

var defaultReportInterval = time.Second * 10
var defaultPollInterval = time.Second * 1

func parseFlags() error {
	flag.StringVar(&connectAddr, "a", "localhost:8080", "address and port for connecting to server")
	flag.StringVar(&cryptoKey, "crypto-key", "", "file with public key")
	flag.StringVar(&pprofAddr, "prof", "", "run pprof")
	flag.StringVar(&configFile, "c", "", "path to config file")
	flag.DurationVar(&reportInterval, "r", defaultReportInterval, "timer of report interval for send metrics")
	flag.IntVar(&rateLimitRoutines, "l", 1, "count of routines")
	flag.DurationVar(&pollInterval, "p", defaultPollInterval, "timer of poll interval for metrics")
	flag.StringVar(&runLog, "log", "info", "log level")
	flag.StringVar(&signPass, "k", "", "signature for HashSHA256")
	flag.Parse()

	err := parseEnv()
	if err != nil {
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

	if envConnectAddr, ok := os.LookupEnv("ADDRESS"); ok {
		connectAddr = envConnectAddr
	} else if connectAddr == "" {
		connectAddr = conf.Address
	}

	if envCryptoKey, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		cryptoKey = envCryptoKey
	} else if cryptoKey == "" {
		cryptoKey = conf.CryptoKey
	}

	if envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		r, err := time.ParseDuration(envReportInterval)
		if err != nil {
			return errors.Wrap(err, "REPORT_INTERVAL is not correct")
		}
		reportInterval = r
	} else if reportInterval == defaultReportInterval {
		r, err := time.ParseDuration(conf.ReportInterval)
		if err != nil {
			return err
		}
		reportInterval = r
	}

	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		p, err := time.ParseDuration(envPollInterval)
		if err != nil {
			return errors.Wrap(err, "POLL_INTERVAL is not correct")
		}
		pollInterval = p
	} else if pollInterval == defaultPollInterval {
		p, err := time.ParseDuration(conf.PollInterval)
		if err != nil {
			return err
		}
		pollInterval = p
	}

	if envSignPass, ok := os.LookupEnv("KEY"); ok {
		signPass = envSignPass
	}
	if envRateLimitRoutines, ok := os.LookupEnv("RATE_LIMIT"); ok {
		r, err := strconv.Atoi(envRateLimitRoutines)
		if err != nil {
			return err
		}
		rateLimitRoutines = r
	}
	if envPprofAdr, ok := os.LookupEnv("PPROF_ADDRESS"); ok {
		pprofAddr = envPprofAdr
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
