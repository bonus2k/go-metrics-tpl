// Package clients реализует интеграцию с сервисом Server
package clients

import (
	"fmt"
	"strconv"

	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	m "github.com/bonus2k/go-metrics-tpl/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// Connect задает параметры для подключения к серверу
type Connect struct {
	Server   string
	Protocol string
	Client   *resty.Client
}

// SendBatchMetrics отправляет метрики на сервер пакетом
// за одну HTTP сессию может быть пераданно более одной метрики
func (con *Connect) SendBatchMetrics(listMetrics []m.Metrics) error {
	address := fmt.Sprintf("%s://%s/updates/", con.Protocol, con.Server)
	_, err := con.Client.R().
		SetHeader(m.KeyContentType, m.TypeJSONContent).
		SetBody(listMetrics).
		Post(address)
	if err != nil {
		return errors.Wrap(err, "[SendBatchMetrics]")
	}
	return nil
}

// SendToGauge отправляет на сервер по одной метрике Gauge
// за одну HTTP сессию может быть пераданно не более одной метрики
// Deprecated: используйте SendBatchMetrics если необходимо передать более одной метрики на сервер
func (con *Connect) SendToGauge(mm map[string]string) error {
	address := fmt.Sprintf("%s://%s/update/", con.Protocol, con.Server)
	for k, v := range mm {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			logger.Log.Error("can't parse float", err)
			continue
		}
		_, err = con.Client.R().
			SetHeader(m.KeyContentType, m.TypeJSONContent).
			SetBody(m.Metrics{
				ID:    k,
				MType: "gauge",
				Value: &value,
			}).
			Post(address)
		if err != nil {
			return errors.Wrap(err, "[SendToGauge]")
		}
	}
	return nil
}

// SendToCounter отправляет на сервер по одной метрике Counter
// за одну HTTP сессию может быть пераданно не более одной метрики
// Deprecated: используйте SendBatchMetrics если необходимо передать более одной метрики на сервер
func (con *Connect) SendToCounter(name string, value int64) {

	address := fmt.Sprintf("%s://%s/update/", con.Protocol, con.Server)
	_, err := con.Client.R().
		SetHeader(m.KeyContentType, m.TypeJSONContent).
		SetBody(m.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &value,
		}).
		Post(address)

	if err != nil {
		logger.Log.Error("[SendToCounter]", err)
	}
}
