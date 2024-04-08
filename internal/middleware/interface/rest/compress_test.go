package rest

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkGzipReqDecompression(b *testing.B) {
	testStr := "test, test, test, test, test, test, test, test, test, test, test, test"

	b.Run("compress with file", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			body := testCompressBody(testStr)
			if resp, resBody := testRequest("POST", body, GzipReqDecompression); testStr != resBody && resp.StatusCode != 200 {
				b.Error(resp, resBody)
			}
		}

	})

	b.Run("compress with reader", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			body := testCompressBody(testStr)
			if resp, resBody := testRequest("POST", body, GzipReqDecompression1); testStr != resBody && resp.StatusCode != 200 {
				b.Error(resp, resBody)
			}
		}
	})

}

func testRequest(method string, body io.Reader, fDecompression func(h http.Handler) http.Handler) (*http.Response, string) {
	req, err := http.NewRequest(method, "/", body)
	if err != nil {
		return nil, ""
	}

	req.Close = true
	req.Header.Add("Content-Encoding", "gzip")
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rw := httptest.NewRecorder()
	handler := fDecompression(testHandler)
	handler.ServeHTTP(rw, req)

	return rw.Result(), string(rw.Body.Bytes())
}

func testCompressBody(body string) io.Reader {
	b := []byte(body)
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	gzipw, _ := gzip.NewWriterLevel(buffer, gzip.BestCompression)
	defer gzipw.Close()
	gzipw.Write(b)
	return buffer
}
