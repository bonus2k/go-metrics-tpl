package logger

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"go.uber.org/zap"
)

func BenchmarkLogging(b *testing.B) {
	str := "test string"
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	logger, _ := cfg.Build()
	log := zlog.Level(zerolog.InfoLevel)

	b.Run("zap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Debug(str)
		}
	})

	b.Run("zap with format Sprintf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Debug(fmt.Sprintf("str %s \n", str))
		}
	})

	b.Run("zap with format zap.String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Debug(str, zap.String("str", str))
		}
	})

	//b.Run("slog", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		slog.Debug(str)
	//	}
	//})

	b.Run("zlog", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.Debug().Msg(str)
		}
	})

	b.Run("zlog with format Sprintf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.Debug().Msg(fmt.Sprintf("str %s", str))
		}
	})

	b.Run("zlog with format Msgf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			log.Debug().Msgf("str %s", str)
		}
	})
}
