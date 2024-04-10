package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorageImpl_AddCounter(t *testing.T) {

	type args struct {
		name  string
		value int64
		want  int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test#1 add Counter",
			args: args{name: "aCount", value: 5, want: 5},
		},
		{
			name: "test#2 add Counter",
			args: args{name: "aCount", value: 6, want: 11},
		},
		{
			name: "test#3 add Counter",
			args: args{name: "aCount", value: 100, want: 111},
		},
	}
	ms := &MemStorageImpl{
		Gauge:   nil,
		Counter: make(map[string]int64),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.AddCounter(context.TODO(), tt.args.name, tt.args.value)
			got := ms.Counter[tt.args.name]
			assert.Equal(t, tt.args.want, got)
		})
	}
}

func TestMemStorageImpl_AddGauge(t *testing.T) {

	type args struct {
		name  string
		value float64
		want  float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test#1 add Gauge",
			args: args{name: "aGauge", value: 5.005, want: 5.005},
		},
		{
			name: "test#2 add Gauge",
			args: args{name: "aGauge", value: -0.99, want: -0.99},
		},
		{
			name: "test#3 add Gauge",
			args: args{name: "aGauge", value: 0, want: 0},
		},
		{
			name: "test#4 add Gauge",
			args: args{name: "zGauge", value: 0.88, want: 0.88},
		},
	}
	ms := &MemStorageImpl{
		Gauge: make(map[string]float64),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.AddGauge(context.TODO(), tt.args.name, tt.args.value)
			got := ms.Gauge[tt.args.name]
			assert.Equal(t, tt.args.want, got)
		})
	}
}

func TestMemStorageImpl_GetCounter(t *testing.T) {

	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "test get count",
			args: args{name: "aCount", value: 99},
			want: 99,
		},
	}
	ms := &MemStorageImpl{
		Counter: make(map[string]int64),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.Counter[tt.args.name] = tt.args.value
			got, _ := ms.GetCounter(context.TODO(), tt.args.name)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMemStorageImpl_GetGauge(t *testing.T) {

	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "test get Gauge",
			args: args{name: "zGauge", value: 0.99},
			want: 0.99,
		},
	}
	ms := &MemStorageImpl{
		Gauge: make(map[string]float64),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.Gauge[tt.args.name] = tt.args.value
			got, _ := ms.GetGauge(context.TODO(), tt.args.name)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want Storage
	}{
		{
			name: "test mem storage",
			want: &MemStorageImpl{
				Gauge:   make(map[string]float64),
				Counter: make(map[string]int64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := *NewMemStorage(false)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMemStorageImpl_GetAllMetrics(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	tests := []struct {
		name   string
		fields fields
		want   []Metric
	}{
		{
			name: "all metrics",
			fields: fields{
				gauge:   map[string]float64{"aGauge": 999.0},
				counter: map[string]int64{"aCounter": 100},
			},
			want: []Metric{{Name: "aGauge", Value: "999"}, {Name: "aCounter", Value: "100"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorageImpl{
				Gauge:   tt.fields.gauge,
				Counter: tt.fields.counter,
			}
			metrics, _ := ms.GetAllMetrics(context.TODO())
			assert.Equalf(t, tt.want, metrics, "GetAllMetrics()")
		})
	}
}
