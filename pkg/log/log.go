// Package logger provides a simple wrapper around zerolog.
package log

import (
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const logFile = "aci-vetr.log"

// Logger aliases the zerolog.Logger
type Logger = zerolog.Logger

var (
	logger = New()

	// Convenience shortcuts for logging levels
	Debug = logger.Debug
	Info  = logger.Info
	Warn  = logger.Warn
	Error = logger.Error
	Fatal = logger.Fatal
	Panic = logger.Panic
	With  = logger.With
)

// New creates a new multi-level logger
func New() Logger {
	if testing.Testing() {
		return zerolog.Nop()
	}
	// If filename is specified, open file and assume file logging
	file, err := os.Create(logFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create log file")
	}

	// Levels default to zero, i.e. debug
	return zerolog.New(

		zerolog.ConsoleWriter{
			Out:     io.MultiWriter(os.Stderr, file),
			NoColor: runtime.GOOS == "windows",
		}).With().Timestamp().Logger()
}

func init() {
	// defaults
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.DurationFieldInteger = true
}
