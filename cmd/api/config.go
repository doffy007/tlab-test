package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type config struct {
	LogLevel zerolog.Level `envconfig:"log_level" default:"0"`
	Port     int           `envconfig:"port" default:"8080"`
}

var cfg config

var _ = func() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var err error

	if err = envconfig.Process("", &cfg); err != nil {
		log.Fatal().Err(err).Msg("")
	}
	zerolog.SetGlobalLevel(cfg.LogLevel)

	t := time.Now()

	log.Info().
		Interface("config", cfg).
		Str("location", fmt.Sprintf("%s %s", t.Location().String(), t.Format("15:04"))).
		Msg("initialize api")

	return nil
}()
