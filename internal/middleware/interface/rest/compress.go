// Package rest реализует middleware для сжатия и проверки cookie
package rest

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	m "github.com/bonus2k/go-metrics-tpl/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

var contentIsCompressed = []string{
	"application/javascript",
	"application/json",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}

// GzipReqCompression сжимает request по средствам gzip flate.BestCompression
func GzipReqCompression(c *resty.Client, r *http.Request) error {
	if r.Body == nil || !isPayloadSupported(r.Header) {
		return nil
	}

	var buf bytes.Buffer
	body, err := io.ReadAll(r.Body)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			logger.Log.Error("GzipReqCompression", err)
		}
	}()
	if err != nil {
		return errors.Wrap(err, "create gzip writer")
	}

	gzipw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	defer func() {
		err = gzipw.Close()
		if err != nil {
			logger.Log.Error("GzipReqCompression", err)
		}
	}()
	if err != nil {
		return errors.Wrap(err, "level compression is invalid")
	}

	_, err = gzipw.Write(body)
	if err != nil {
		return errors.Wrap(err, "can't write in body")
	}

	r.ContentLength = int64(buf.Len())
	r.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	r.Header.Add(m.KeyContentEncoding, m.TypeEncodingContent)

	return nil
}

// GzipReqDecompression разархивирует request по средствам gzip flate.BestCompression
func GzipReqDecompression(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get(m.KeyContentEncoding), m.TypeEncodingContent) {
			gz, err := gzip.NewReader(r.Body)
			defer func() {
				err = r.Body.Close()
				if err != nil {
					logger.Log.Error("GzipReqDecompression", err)
				}
			}()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() {
				err = gz.Close()
				if err != nil {
					logger.Log.Error("GzipReqDecompression", err)
				}
			}()
			body, err := io.ReadAll(gz)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(body))
		}
		h.ServeHTTP(w, r)
	})
}

func GzipReqDecompression1(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get(m.KeyContentEncoding), m.TypeEncodingContent) {
			gz, err := gzip.NewReader(r.Body)
			defer func() {
				err = r.Body.Close()
				if err != nil {
					logger.Log.Error("GzipReqDecompression", err)
				}
			}()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() {
				err = gz.Close()
				if err != nil {
					logger.Log.Error("GzipReqDecompression", err)
				}
			}()
			r.Body = gz
		}
		h.ServeHTTP(w, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipResCompression сжимает response по средствам gzip flate.BestCompression
func GzipResCompression(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get(m.KeyAcceptEncoding), m.TypeEncodingContent) {
			h.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			err = gz.Close()
			if err != nil {
				logger.Log.Error("GzipResCompression", err)
			}
		}()
		w.Header().Set(m.KeyContentEncoding, m.TypeEncodingContent)
		h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func isPayloadSupported(h http.Header) bool {
	for _, c := range contentIsCompressed {
		if strings.Contains(h.Get(m.KeyContentType), c) {
			return true
		}
	}
	return false
}
