package main

import (
	"github.com/bonus2k/go-metrics-tpl/cmd/server/repository"
	"net/http"
	"strconv"
	"strings"
)

var mem repository.MemStorage

func main() {
	mem = repository.NewMemStorage()
	if err := run(); err != nil {
		panic(err)
	}
}

func gaugePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	segments := strings.Split(r.URL.Path, "/")
	name := strings.TrimSpace(segments[3])
	if len(name) == 0 || len(segments) < 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if float, err := strconv.ParseFloat(segments[4], 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		mem.AddGauge(name, float)
		w.WriteHeader(http.StatusOK)
		return
	}
}

func counterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	segments := strings.Split(r.URL.Path, "/")
	name := strings.TrimSpace(segments[3])
	if len(name) == 0 || len(segments) < 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if num, err := strconv.ParseInt(segments[4], 10, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		mem.AddCounter(name, num)
		w.WriteHeader(http.StatusOK)
		return
	}
}

func run() error {

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
	})
	mux.HandleFunc(`/update/gauge/`, gaugePage)
	mux.HandleFunc(`/update/counter/`, counterPage)
	return http.ListenAndServe(`:8080`, mux)
}
