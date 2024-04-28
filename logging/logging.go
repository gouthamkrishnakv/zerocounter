package logging

import (
	// stdlib imports
	"os"
	"time"

	// third-party imports
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// -- variables --

var l *zerolog.Logger = nil

// -- functions --

func L() *zerolog.Logger {
	return l
}

func Initialize() error {
	// set up time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// set up zerologger stack-trace
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// process pid
	pid := os.Getpid()

	// set up human-friendly logging
	// TODO: add a mutltiwriter logger so that both console and a log-file
	// would be created
	logger := zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).Level(zerolog.TraceLevel).
		With().Timestamp().Caller().Int("pid", pid).Logger()
	l = &logger

	return nil
}
