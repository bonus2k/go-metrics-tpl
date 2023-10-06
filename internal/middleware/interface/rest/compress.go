package rest

import (
	"bytes"
	"compress/gzip"
	m "github.com/bonus2k/go-metrics-tpl/internal/models"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"strings"
)

var contentIsCompressed = []string{
	"application/javascript",
	"application/json",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}

func GzipReqCompression(c *resty.Client, r *http.Request) error {
	if r.Body != nil && isPayloadSupported(r.Header) {
		var buf bytes.Buffer
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}

		gzipw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
		if err != nil {
			return err
		}

		_, err = gzipw.Write(body)
		if err != nil {
			return err
		}

		gzipw.Close()
		r.ContentLength = int64(buf.Len())
		r.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
		c.Header.Add(m.KeyContentEncoding, m.TypeEncodingContent)
	}
	return nil
}

func GzipReqDecompression(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get(m.KeyContentEncoding), m.TypeEncodingContent) {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()
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

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

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
		defer gz.Close()
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
