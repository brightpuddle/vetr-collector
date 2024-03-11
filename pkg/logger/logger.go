// Package logger provides a simple wrapper around zerolog.
package logger

import (
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// Logger aliases the zerolog.Logger
type Logger = zerolog.Logger

var (
	log     *zerolog.Logger
	logMux  sync.Mutex
	windows = runtime.GOOS == "windows"
)

// Get returns the current logger instance.
func Get() *Logger {
	if log == nil {
		l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: windows}).
			Level(zerolog.WarnLevel).
			With().
			Timestamp().
			Logger()
		log = &l
	}
	return log
}

// Set sets the logger instance.
func Set(l *Logger) {
	log = l
}

// SetLevel sets the logging level.
func SetLevel(level zerolog.Level) {
	log := Get().Level(level)
	Set(&log)
}

// MultiLevelWriter writes logs to file and console
type MultiLevelWriter struct {
	file         io.Writer
	console      io.Writer
	consoleLevel zerolog.Level
}

func (w MultiLevelWriter) Write(p []byte) (int, error) {
	logMux.Lock()
	count, err := w.file.Write(p)
	logMux.Unlock()
	return count, err
}

// WriteLevel writes log data for a given log level
func (w MultiLevelWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level >= w.consoleLevel {
		n, err := w.console.Write(p)
		if err != nil {
			return n, err
		}
	}
	return w.file.Write(p)
}

// Level is a logger level
type Level = zerolog.Level

const (
	// DebugLevel is debug level logging
	DebugLevel = iota + zerolog.DebugLevel
	// InfoLevel is info level logging
	InfoLevel
	// WarnLevel is warning level logging
	WarnLevel
	// ErrorLevel is error level logging
	ErrorLevel
	// FatalLevel is fatal level logging
	FatalLevel
	// PanicLevel is panic level logging
	PanicLevel
)

// Config is a logger configuration
type Config struct {
	ConsoleLevel Level
	FileLevel    Level
	Filename     string
	FileOut      io.Writer
	ConsoleOut   io.Writer
}

type fileWriter struct {
	io.Writer
}

// New creates a new multi-level logger
func New(cfg Config) (*Logger, error) {
	// If filename is specified, open file and assume file logging
	if cfg.Filename != "" {
		file, err := os.Create(cfg.Filename)
		if err != nil {
			return nil, err
		}
		cfg.FileOut = file
	}

	// Only log to file if specified
	if cfg.FileOut == nil {
		cfg.FileOut = io.Discard
	}
	// Log to stderr unless otherwise specified
	if cfg.ConsoleOut == nil {
		cfg.ConsoleOut = os.Stderr
	}

	// Levels default to zero, i.e. debug
	consoleWriter := zerolog.ConsoleWriter{
		Out:     cfg.ConsoleOut,
		NoColor: windows,
	}
	writer := MultiLevelWriter{
		file:         cfg.FileOut,
		console:      consoleWriter,
		consoleLevel: cfg.ConsoleLevel,
	}
	log := zerolog.New(writer).With().Timestamp().Logger()
	Set(&log)
	return &log, nil
}

func init() {
	// defaults
	zerolog.SetGlobalLevel(DebugLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.DurationFieldInteger = true
}
