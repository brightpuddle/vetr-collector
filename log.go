package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
)

const logFile = "collector.log"

var logMux sync.Mutex

// Logger aliases the zerolog.Logger
type Logger = zerolog.Logger

// MultiLevelWriter writes logs to file and console
type MultiLevelWriter struct {
	file    io.Writer
	console io.Writer
}

func (w MultiLevelWriter) Write(p []byte) (int, error) {
	logMux.Lock()
	count, err := w.file.Write(p)
	logMux.Unlock()
	return count, err
}

// WriteLevel writes log data for a given log level
func (w MultiLevelWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level >= zerolog.InfoLevel {
		n, err := w.console.Write(p)
		if err != nil {
			return n, err
		}
	}
	return w.file.Write(p)
}

func newLogger() Logger {
	file, err := os.Create(logFile)
	if err != nil {
		panic(fmt.Sprintf("cannot create log file %s", logFile))
	}

	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.DurationFieldInteger = true

	writer := MultiLevelWriter{
		file:    file,
		console: zerolog.ConsoleWriter{Out: colorable.NewColorableStdout()},
	}
	return zerolog.New(writer).With().Timestamp().Logger()
}
