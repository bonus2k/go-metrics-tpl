package clients

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	m "github.com/bonus2k/go-metrics-tpl/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strconv"
)

type Connect struct {
	Server   string
	Protocol string
	Client   resty.Client
}

func (con *Connect) SendToGauge(mm map[string]string) error {
	address := fmt.Sprintf("%s://%s/update/", con.Protocol, con.Server)
	for k, v := range mm {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			logger.Log.Error("can't parse float")
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
		logger.Log.Error("[SendToCounter]", zap.Error(err))
	}
}
