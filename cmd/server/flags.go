package main

import (
	"flag"
	"os"
)

var runAddr string
var runLog string

func parseFlags() {
	flag.StringVar(&runAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&runLog, "l", "info", "log level")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		runAddr = envRunAddr
	}
}
