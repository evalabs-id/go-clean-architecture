package providers

import (
	"io"
	"os"
	"time"

	"github.com/evalabs-id/go-clean-architecture/internal/configs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func ProvideZerolog(cfg *configs.Config) zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimestampFieldName = "timestamp"
	zerolog.CallerFieldName = "caller"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	level := getLogLevel(cfg.Server.Debug)

	output := getLogOutput()

	logger := getLogger(cfg.App.Environment, output, level)

	return logger
}

func getLogLevel(debug bool) zerolog.Level {
	if debug {
		return zerolog.DebugLevel
	}
	return zerolog.InfoLevel
}

func getLogOutput() io.Writer {
	// for now, always log to stdout
	// for future, make this configurable to io.Writer
	return os.Stdout
}

func getLogger(env string, output io.Writer, level zerolog.Level) zerolog.Logger {
	switch env {
	case "development", "dev":
		return developmentLogger(output, level)
	case "production", "prod":
		return productionLogger(output, level)
	default:
		return developmentLogger(output, level)
	}
}

func developmentLogger(output io.Writer, level zerolog.Level) zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        output,
		TimeFormat: time.RFC3339,
	}

	return zerolog.New(consoleWriter).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()
}

func productionLogger(output io.Writer, level zerolog.Level) zerolog.Logger {
	return zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()
}
