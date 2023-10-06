package rest

import (
	"bytes"
	"compress/gzip"
	"github.com/bonus2k/go-metrics-tpl/internal/models"
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

func GzipCompression(c *resty.Client, r *http.Request) error {
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
		c.Header.Add(models.KeyContentEncoding, models.TypeEncodingContent)
	}
	return nil
}

func GzipDecompression(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get(models.KeyContentEncoding), models.TypeEncodingContent) {
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
			w.Write(body)
		}
	})
}

func isPayloadSupported(h http.Header) bool {
	for _, c := range contentIsCompressed {
		if strings.Contains(h.Get(models.KeyContentType), c) {
			return true
		}
	}
	return false
}
