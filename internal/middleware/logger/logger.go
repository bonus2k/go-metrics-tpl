// Package logger реализует интерфейс логгера
package logger

import (
	"bytes"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger interface {
	Initialize(level string) Logger
	Info(string string)
	Infof(string string, any ...any)
	Debug(string string)
	Debugf(string string, any ...any)
	Error(string string, err error)
	Trace(string string)
	Warn(string string)
	Panic(string string)
}

type loggerImpl struct {
	logger zerolog.Logger
}

func (l *loggerImpl) Info(string string) {
	l.logger.Info().Msg(string)
}

func (l *loggerImpl) Infof(string string, any ...any) {
	l.logger.Info().Msgf(string, any...)
}

func (l *loggerImpl) Debug(string string) {
	l.logger.Debug().Msg(string)
}

func (l *loggerImpl) Debugf(string string, any ...any) {
	l.logger.Debug().Msgf(string, any...)
}

func (l *loggerImpl) Error(string string, err error) {
	l.logger.Error().Err(err).Msg(string)
}

func (l *loggerImpl) Trace(string string) {
	l.logger.Trace().Msg(string)
}

func (l *loggerImpl) Warn(string string) {
	l.logger.Warn().Msg(string)
}

func (l *loggerImpl) Panic(string string) {
	l.logger.Panic().Msg(string)
}

var Log loggerImpl

// Initialize Инизиализирует логгер с заданным уровнем логирования
func Initialize(level string) error {
	parseLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}
	Log = loggerImpl{zerolog.
		New(os.Stdout).
		Level(parseLevel).
		With().Timestamp().
		Logger()}
	return nil
}

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

// MiddlewareLog логирует HTTP сессии
func MiddlewareLog(h http.Handler) http.Handler {
	var lw *loggingResponseWriter
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var buf bytes.Buffer
		lw = &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   &responseData{status: 0, size: 0},
			body:           buf,
		}
		defer func() {
			err := r.Body.Close()
			if err != nil {
				Log.logger.Error().Err(err).Msg("")
			}
		}()
		h.ServeHTTP(lw, r)
		duration := time.Since(start)
		Log.logger.Info().
			Str("uri", r.RequestURI).
			Str("method", r.Method).
			Int("status", lw.responseData.status).
			Dur("duration", duration).
			Int("size", lw.responseData.size).
			Str("body", lw.body.String()).
			Msg("")

	})
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size = size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
