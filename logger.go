package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(config *Config) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.IsDevelopment() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if config.IsDevelopment() {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	log.Info().
		Str("environment", config.Environment).
		Str("port", config.Port).
		Msg("Logger initialized")
}
