package controllers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCounterPage(t *testing.T) {
	server := httptest.NewServer(MetricsRouter())
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
			want:    want{contentType: "text/plain", statusCode: 200},
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
			want:    want{contentType: "text/plain; charset=utf-8", statusCode: 404},
		},
		{
			name:    "test counter page 400",
			request: "/update/counter/aCount/ttt",
			method:  http.MethodPost,
			want:    want{contentType: "text/plain", statusCode: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, server, tt.method, tt.request)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func TestGaugePage(t *testing.T) {
	server := httptest.NewServer(MetricsRouter())
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
			want:    want{contentType: "text/plain", statusCode: 200},
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
			want:    want{contentType: "text/plain; charset=utf-8", statusCode: 404},
		},
		{
			name:    "test gauge page 400",
			request: "/update/gauge/aGauge/ttt",
			method:  http.MethodPost,
			want:    want{contentType: "text/plain", statusCode: 400},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, server, tt.method, tt.request)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
