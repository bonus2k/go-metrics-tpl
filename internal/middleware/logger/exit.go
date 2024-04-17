package logger

import (
	"github.com/rs/zerolog"
	"os"
	sync2 "sync"
)

var log zerolog.Logger
var once sync2.Once

func Exit(err error, code int) {
	once.Do(func() {
		log = zerolog.
			New(os.Stdout).
			Level(zerolog.ErrorLevel).
			With().Timestamp().
			Logger()
	})

	log.Err(err).Msg("")
	os.Exit(1)
}
