package services

import (
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/models"
	"sync"
	"time"

	"github.com/bonus2k/go-metrics-tpl/internal/agent/clients"
)

type PoolWorcker struct {
	client      clients.Connect
	count       func() int64
	workerGroup sync.WaitGroup
	shutdown    bool
}

var pool *PoolWorcker
var once sync.Once

// NewPool создает пул worker которые собирают метрики и отправляют их на сервер
func NewPool(client clients.Connect) *PoolWorcker {
	once.Do(func() {
		pool = &PoolWorcker{client: client, count: GetPollCount(), shutdown: false}
	})
	return pool
}

// BatchReport осуществляет сбор метрик и их отправку в сервис Server,
// в вслучае HTTP ошибки отправялет в канал error сообщение об ошибке
func (p *PoolWorcker) BatchReport(jobs <-chan map[string]string, errors chan<- error, ticker *time.Ticker, goRoutine int) {
	p.workerGroup.Add(1)
out:
	for range ticker.C {
		if !pool.shutdown {
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
				errors <- fmt.Errorf("worker %d can't send metrics to updates %w", goRoutine, err)
				continue
			}
			logger.Log.Info(fmt.Sprintf("worker %d send count %d", goRoutine, count))
		} else {
			logger.Log.Info(fmt.Sprintf("worker %d is shutdown", goRoutine))
			p.workerGroup.Done()
			ticker.Reset(time.Second)
			break out
		}

	}
}

func (p *PoolWorcker) Shutdown() {
	p.shutdown = true
	p.workerGroup.Wait()
}
