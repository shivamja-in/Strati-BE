package telemetry

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger zerolog.Logger
var once sync.Once

func Telemetry() zerolog.Logger {
	once.Do(func() {
		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		if os.Getenv("APP_ENV") != "development" {
			fileLogger := &lumberjack.Logger{
				Filename:   "demo.log",
				MaxSize:    5,
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   true,
			}

			output = zerolog.MultiLevelWriter(fileLogger, os.Stderr)
		}

		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

		logLevel := zerolog.InfoLevel
		if level := os.Getenv("LOG_LEVEL"); level != "" {
			if parsedLevel, err := zerolog.ParseLevel(level); err == nil {
				logLevel = parsedLevel
			}
		}

		logger = zerolog.New(output).
			Level(logLevel).
			With().
			Timestamp().
			Caller().
			Int("pid", os.Getpid()).
			Logger()
	})

	return logger
}
