package telemetry

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
				Filename:   "wikipedia-demo.log",
				MaxSize:    5, //
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   true,
			}

			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		}

		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

		infoSampler := &zerolog.BurstSampler{
			Burst:  3,
			Period: 1 * time.Second,
		}

		warnSampler := &zerolog.BurstSampler{
			Burst:  3,
			Period: 1 * time.Second,
			// Log every 5th message after exceeding the burst rate of 3 messages per
			// second
			NextSampler: &zerolog.BasicSampler{N: 5},
		}

		errorSampler := &zerolog.BasicSampler{N: 2}

		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("[%s]", i))
			},
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("| %s |", i)
			},
			FormatCaller: func(i interface{}) string {
				return filepath.Base(fmt.Sprintf("%s", i))
			},
			PartsExclude: []string{
				zerolog.TimestampFieldName,
			},
		}).
			Level(zerolog.TraceLevel).
			With().
			Timestamp().
			Caller().
			Int("pid", os.Getpid()).
			Logger().
			Sample(zerolog.LevelSampler{
				WarnSampler:  warnSampler,
				InfoSampler:  infoSampler,
				ErrorSampler: errorSampler,
			})
	})

	return logger
}
