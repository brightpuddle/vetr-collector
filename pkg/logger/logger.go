package logger

import (
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/gofiber/fiber/v2"
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
	log := Get()
	l := log.Level(level)
	Set(&l)
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
	if level >= zerolog.InfoLevel || level >= w.consoleLevel {
		n, err := w.console.Write(p)
		if err != nil {
			return n, err
		}
	}
	return w.file.Write(p)
}

// Config is a logger configuration
type Config struct {
	Debug bool
}

// New creates a new multi-level logger
func New(logFile string, cfgs ...Config) (*Logger, error) {
	var cfg Config
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}

	file, err := os.Create(logFile)
	if err != nil {
		return nil, err
	}

	consoleLevel := zerolog.WarnLevel
	if cfg.Debug {
		consoleLevel = zerolog.DebugLevel
	}

	writer := MultiLevelWriter{
		file:         zerolog.ConsoleWriter{Out: file, NoColor: true},
		console:      zerolog.ConsoleWriter{Out: os.Stderr, NoColor: windows},
		consoleLevel: consoleLevel,
	}
	log := zerolog.New(writer).With().Timestamp().Logger()
	Set(&log)
	return &log, err
}

// Middleware returns a fiber handler using the global log instance.
func Middleware(verbose bool) fiber.Handler {
	sublog := Get()

	return func(c *fiber.Ctx) error {
		chainErr := c.Next()

		msg := "Request"
		if chainErr != nil {
			msg = chainErr.Error()
		}

		code := c.Response().StatusCode()

		dumplogger := sublog.With().
			Int("status", code).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Str("user-agent", c.Get(fiber.HeaderUserAgent)).
			Logger()

		switch {
		case code >= 200 && code < 300:
			if verbose {
				dumplogger.Info().Msg(msg)
			} else {
				dumplogger.Debug().Msg(msg)
			}
		case code >= 300 && code < 400:
			dumplogger.Warn().Msg(msg)
		default:
			dumplogger.Error().Msg(msg)
		}
		return chainErr
	}
}

func init() {
	// defaults
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.DurationFieldInteger = true
}
