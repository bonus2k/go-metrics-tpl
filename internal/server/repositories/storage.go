package repositories

type Storage interface {
	AddGauge(string, float64)
	GetGauge(string) (float64, bool)
	AddCounter(string, int64)
	GetCounter(string) (int64, bool)
	GetAllMetrics() []Metric
	CheckConnection() error
}
