package controllers

import (
	"net/http"
	"strconv"
	"strings"
)

func CounterPage(w http.ResponseWriter, r *http.Request) {
	if MemStorage == nil {
		panic("storage not initialized")
	}
	w.Header().Set("Content-Type", "text/plain")
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
		MemStorage.AddCounter(name, num)
		w.WriteHeader(http.StatusOK)
		return
	}
}
