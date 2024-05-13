package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/bonus2k/go-metrics-tpl/internal/agent/clients"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/models"
)

type PoolWorcker struct {
	client clients.Connect
	count  func() int64
}

var pool *PoolWorcker
var once sync.Once

// NewPool создает пул worker которые собирают метрики и отправляют их на сервер
func NewPool(client clients.Connect) *PoolWorcker {
	once.Do(func() {
		pool = &PoolWorcker{client: client, count: GetPollCount()}
	})
	return pool
}

// BatchReport осуществляет сбор метрик и их отправку в сервис Server,
// в вслучае HTTP ошибки отправялет в канал error сообщение об ошибке
func (p *PoolWorcker) BatchReport(jobs <-chan map[string]string, errors chan<- error, ticker *time.Ticker, goRoutine int) {

	for range ticker.C {
		m := <-jobs
		metrics, err := models.ConvertGaugeToMetrics(&m)
		if err != nil {
			errors <- fmt.Errorf("can't convert gauge to metrics %w", err)
			continue
		}
		count := p.count()
		metrics = append(metrics, models.Metrics{ID: "PollCount", Delta: &count, MType: "counter"})
		err = p.client.SendBatchMetrics(metrics)
		if err != nil {
			//интересно какую логику подразумевает проект в отношении PollCount?
			//если не удалось отправить сообщение, стоит ли нам делать откат счетчика на 1
			errors <- fmt.Errorf("goRoutine %d can't send metrics to updates %w", goRoutine, err)
			continue
		}
		logger.Log.Info(fmt.Sprintf("goRoutine %d send count %d", goRoutine, count))
	}
}
