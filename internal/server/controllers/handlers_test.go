package controllers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/server/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounterPage(t *testing.T) {
	storage := repositories.NewMemStorage(false)
	server := httptest.NewServer(MetricsRouter(storage, ""))
	defer server.Close()
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "test counter page 200",
			request: "/update/counter/aCount/100",
			method:  http.MethodPost,
			want:    want{contentType: "text/html", statusCode: 200},
		},
		{
			name:    "test counter page 405",
			request: "/update/counter/aCount/100",
			method:  http.MethodPut,
			want:    want{contentType: "", statusCode: 405},
		},
		{
			name:    "test counter page 404",
			request: "/update/counter/",
			method:  http.MethodPost,
			want:    want{contentType: "", statusCode: 404},
		},
		{
			name:    "test counter page 400",
			request: "/update/counter/aCount/ttt",
			method:  http.MethodPost,
			want:    want{contentType: "text/html", statusCode: 400},
		},
		{
			name:    "test update page 400",
			request: "/update/counter/aCount/ttt",
			method:  http.MethodPost,
			want:    want{contentType: "text/html", statusCode: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, server, tt.method, tt.request)
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					logger.Log.Error("Test", err)
				}
			}()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func TestGaugePage(t *testing.T) {
	storage := repositories.NewMemStorage(false)
	server := httptest.NewServer(MetricsRouter(storage, ""))
	defer server.Close()
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "test gauge page 200",
			request: "/update/gauge/aGauge/100",
			method:  http.MethodPost,
			want:    want{contentType: "text/html", statusCode: 200},
		},
		{
			name:    "test gauge page 405",
			request: "/update/gauge/aGauge/100",
			method:  http.MethodPut,
			want:    want{contentType: "", statusCode: 405},
		},
		{
			name:    "test gauge page 404",
			request: "/update/gauge/",
			method:  http.MethodPost,
			want:    want{contentType: "", statusCode: 404},
		},
		{
			name:    "test gauge page 400",
			request: "/update/gauge/aGauge/ttt",
			method:  http.MethodPost,
			want:    want{contentType: "text/html", statusCode: 400},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, server, tt.method, tt.request)
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					logger.Log.Error("Test", err)
				}
			}()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func TestGetValue(t *testing.T) {
	storage := *repositories.NewMemStorage(false)
	server := httptest.NewServer(MetricsRouter(&storage, ""))
	defer server.Close()
	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "test value aGauge 200",
			request: "/value/gauge/aGauge",
			method:  http.MethodGet,
			want:    want{contentType: "text/html", statusCode: 200, body: "100"},
		},
		{
			name:    "test value aCount 200",
			request: "/value/counter/aCount",
			method:  http.MethodGet,
			want:    want{contentType: "text/html", statusCode: 200, body: "1099"},
		},
		{
			name:    "test value gauge page 404",
			request: "/value/gauge/unknown",
			method:  http.MethodGet,
			want:    want{contentType: "", statusCode: 404},
		},
		{
			name:    "test value counter page 404",
			request: "/value/counter/unknown",
			method:  http.MethodGet,
			want:    want{contentType: "", statusCode: 404},
		},
		{
			name:    "test value page 400",
			request: "/value/gauge/unknown",
			method:  http.MethodGet,
			want:    want{contentType: "", statusCode: 404},
		},
	}

	storage.AddGauge(context.TODO(), "aGauge", 100)
	storage.AddCounter(context.TODO(), "aCount", 999)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, server, tt.method, tt.request)
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					logger.Log.Error("Test", err)
				}
			}()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.body, body)
		})
	}
}

func TestAllMetrics(t *testing.T) {
	storage := *repositories.NewMemStorage(false)
	server := httptest.NewServer(MetricsRouter(&storage, ""))
	defer server.Close()
	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "test all metrics 200",
			request: "/",
			method:  http.MethodGet,
			want:    want{contentType: "text/html", statusCode: 200, body: "[{\"Name\":\"aGauge\",\"Value\":\"100\"},{\"Name\":\"aCount\",\"Value\":\"2098\"}]"},
		},
		{
			name:    "test all metrics 405",
			request: "/",
			method:  http.MethodPost,
			want:    want{contentType: "", statusCode: 405},
		},
	}

	storage.AddGauge(context.TODO(), "aGauge", 100)
	storage.AddCounter(context.TODO(), "aCount", 999)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, server, tt.method, tt.request)
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					logger.Log.Error("Test", err)
				}
			}()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.body, body)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			logger.Log.Error("Test", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
