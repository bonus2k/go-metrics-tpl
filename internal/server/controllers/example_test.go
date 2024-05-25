package controllers

import (
	"context"
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/server/repositories"
	"io"
	"net/http"
	"net/http/httptest"
)

// демонстрирует пример отправки мертик типа count
func Example_counterPage() {
	var nameMetric = "aCount"
	var valueMetric = "100"
	url := fmt.Sprintf("/update/counter/%s/%s", nameMetric, valueMetric)

	code := saveMetric(url)
	fmt.Println(code)
	// Output: 200
}

// демонстрирует пример отправки мертик типа gauge
func Example_gaugePage() {
	var nameMetric = "aGauge"
	var valueMetric = "99"

	url := fmt.Sprintf("/update/gauge/%s/%s", nameMetric, valueMetric)

	code := saveMetric(url)
	fmt.Println(code)
	// Output: 200
}

// демонстрирует пример получения мертик
func Example_getValue() {
	var nameMetricGauge = "aGauge"
	urlGauge := fmt.Sprintf("/value/gauge/%s", nameMetricGauge)
	codeGauge, bodyGauge := sendRequest(urlGauge)
	fmt.Println(codeGauge)
	fmt.Println(string(bodyGauge))

	var nameMetricCounter = "aCount"
	urlCounter := fmt.Sprintf("/value/counter/%s", nameMetricCounter)
	codeCounter, bodyCounter := sendRequest(urlCounter)
	fmt.Println(codeCounter)
	fmt.Println(string(bodyCounter))

	// Output: 200
	//100
	//200
	//99
}

// демонстрирует пример получения всех мертик
func Example_allMetrics() {
	urlAllMetrics := "/"
	code, body := sendRequest(urlAllMetrics)
	fmt.Println(code)
	fmt.Println(string(body))
	// Output: 200
	//[{"Name":"aGauge","Value":"100"},{"Name":"aCount","Value":"99"}]
}

var sendRequest = getMetric()

func saveMetric(url string) int {
	storage := *repositories.NewMemStorage(false)
	server := httptest.NewServer(MetricsRouter(&storage, "", ""))
	defer server.Close()
	req, err := http.NewRequest(http.MethodPost, server.URL+url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := server.Client().Do(req)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()
	code := resp.StatusCode
	return code
}

func getMetric() func(url string) (int, []byte) {
	storage := *repositories.NewMemStorage(false)
	defer repositories.ResetMemStorage()
	storage.AddGauge(context.TODO(), "aGauge", 100)
	storage.AddCounter(context.TODO(), "aCount", 99)
	return func(url string) (int, []byte) {
		server := httptest.NewServer(MetricsRouter(&storage, "", ""))
		defer server.Close()
		req, err := http.NewRequest(http.MethodGet, server.URL+url, nil)
		if err != nil {
			panic(err)
		}
		resp, err := server.Client().Do(req)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				panic(err)
			}
		}()

		code := resp.StatusCode
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		return code, respBody
	}
}
