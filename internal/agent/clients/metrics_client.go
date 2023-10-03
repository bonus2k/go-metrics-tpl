package clients

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"strconv"
)

type Connect struct {
	Server   string
	Protocol string
	Client   resty.Client
}

func (con *Connect) SendToGauge(m map[string]string) {
	address := fmt.Sprintf("%s://%s/update/", con.Protocol, con.Server)
	for k, v := range m {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			logger.Log.Error("can't parse float")
			continue
		}
		resp, err := con.Client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(models.Metrics{
				ID:    k,
				MType: "gauge",
				Value: &value,
			}).
			Post(address)

		if err != nil {
			logger.Log.Error("error response", zap.String("code", resp.Status()))
		}
	}
}

func (con *Connect) SendToCounter(name string, value int64) {

	address := fmt.Sprintf("%s://%s/update/", con.Protocol, con.Server)
	resp, err := con.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(models.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &value,
		}).
		Post(address)

	if err != nil {
		logger.Log.Error("error response", zap.String("code", resp.Status()))

	}
}
