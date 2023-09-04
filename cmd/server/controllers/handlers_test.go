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
			want:    want{contentType: "text/plain", statusCode: 405},
		},
		{
			name:    "test counter page 404",
			request: "/update/counter/",
			method:  http.MethodPost,
			want:    want{contentType: "text/plain", statusCode: 404},
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
			request := httptest.NewRequest(tt.method, tt.request, nil)
			resp := httptest.NewRecorder()
			CounterPage(resp, request)
			result := resp.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestGaugePage(t *testing.T) {
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
			want:    want{contentType: "text/plain", statusCode: 405},
		},
		{
			name:    "test gauge page 404",
			request: "/update/gauge/",
			method:  http.MethodPost,
			want:    want{contentType: "text/plain", statusCode: 404},
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
			request := httptest.NewRequest(tt.method, tt.request, nil)
			resp := httptest.NewRecorder()
			GaugePage(resp, request)
			result := resp.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}
