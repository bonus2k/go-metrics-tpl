package services

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"strconv"
	"testing"
)

func TestAddRandomValue(t *testing.T) {
	type args struct {
		m map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "add random value",
			args: args{m: make(map[string]string)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddRandomValue(tt.args.m)
			rnd := tt.args.m["RandomValue"]
			_, err := strconv.ParseFloat(rnd, 32)
			require.NoError(t, err)
		})
	}
}

func TestGetMapMetrics(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "check metrics",
			want: memStatConst,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMapMetrics()
			require.Equal(t, len(got), len(tt.want))
			for _, s := range tt.want {
				_, found := got[s]
				assert.True(t, found)
			}
		})
	}
}

func TestGetPollCount(t *testing.T) {
	tests := []struct {
		name    string
		iterate int
		want    int64
	}{
		{
			name:    "first iterate",
			iterate: 1,
			want:    1,
		},
		{
			name:    "second iterate",
			iterate: 2,
			want:    2,
		},
		{
			name:    "tenth iterate",
			iterate: 10,
			want:    10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPollCount(); !reflect.DeepEqual(got, tt.want) {
				var count int64
				for i := 0; i < tt.iterate; i++ {
					count = got()
				}
				assert.Equal(t, tt.want, count)
			}
		})
	}
}
