package logger

import (
	"bytes"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		body         bytes.Buffer
		responseData *responseData
	}
)

var Log = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	logger, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = logger
	return nil
}

func MiddlewareLog(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var buf bytes.Buffer
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
			body:           buf,
		}

		h.ServeHTTP(&lw, r)
		duration := time.Since(start)
		sugar := Log.Sugar()

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
			"body", string(lw.body.Bytes()),
		)

	}
	return http.HandlerFunc(logFn)
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
