package logger

import (
	"context"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

type ContextKey string

const LoggerContextKey ContextKey = "logger"

func InjectLoggerToContext(ctx context.Context, logger zerolog.Logger) context.Context {
	if requestID := middleware.GetReqID(ctx); requestID != "" {
		logger = logger.With().Str("request_id", requestID).Logger()
	}

	return context.WithValue(ctx, LoggerContextKey, logger)
}

func FromContext(ctx context.Context) zerolog.Logger {
	if logger, ok := ctx.Value(LoggerContextKey).(zerolog.Logger); ok {
		return logger
	}

	return zerolog.Ctx(ctx).With().Logger()
}

func Default(l zerolog.Logger) fxevent.Logger {
	return &logger{
		Logger: l.With().Logger(),
	}
}

type logger struct {
	Logger zerolog.Logger
}

func (l *logger) Printf(format string, args ...any) {
	l.Logger.Info().Msgf(format, args...)
}

func (l *logger) simplifyFunctionName(fullname string) string {
	parts := strings.Split(fullname, "/")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]

		if strings.Contains(lastPart, ".") {
			funcParts := strings.Split(lastPart, ".")
			if len(funcParts) >= 2 {
				return funcParts[len(funcParts)-2] + "." + funcParts[len(funcParts)-1]
			}
		}

		return lastPart
	}
	return fullname
}

func (l *logger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Logger.Info().Msgf("Starting: %s", l.simplifyFunctionName(e.FunctionName))
	case *fxevent.OnStartExecuted:
		funcName := l.simplifyFunctionName(e.FunctionName)
		if e.Err != nil {
			l.Logger.Error().Msgf("Start failed: %s - %v", funcName, e.Err)
		} else {
			l.Logger.Info().Msgf("Started: %s (%s)", funcName, e.Runtime.String())
		}
	case *fxevent.OnStopExecuting:
		l.Logger.Info().Msgf("Stopping: %s", l.simplifyFunctionName(e.FunctionName))
	case *fxevent.OnStopExecuted:
		funcName := l.simplifyFunctionName(e.FunctionName)
		if e.Err != nil {
			l.Logger.Error().Msgf("Stop failed: %s - %v", funcName, e.Err)
		} else {
			l.Logger.Info().Msgf("Stopped: %s (%s)", funcName, e.Runtime.String())
		}
	case *fxevent.Supplied:
		funcName := l.simplifyFunctionName(e.TypeName)
		if e.Err != nil {
			l.Logger.Error().Msgf("Supply failed: %s - %v", funcName, e.Err)
		} else {
			l.Logger.Debug().Msgf("Supplied: %s", funcName)
		}
	case *fxevent.Provided:
		funcName := l.simplifyFunctionName(e.ConstructorName)
		if e.Err != nil {
			l.Logger.Error().Msgf("Provide failed: %s - %v", funcName, e.Err)
		} else {
			types := strings.Join(e.OutputTypeNames, ", ")
			l.Logger.Debug().Msgf("Provided: %s -> [%s]", funcName, types)
		}
	case *fxevent.Invoked:
		funcName := l.simplifyFunctionName(e.FunctionName)
		if e.Err != nil {
			l.Logger.Error().Msgf("Invoke failed: %s - %v", funcName, e.Err)
		} else {
			l.Logger.Debug().Msgf("Invoked: %s", funcName)
		}
	case *fxevent.Stopping:
		l.Logger.Info().Msgf("Received signal: %s", strings.ToUpper(e.Signal.String()))
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Logger.Error().Msgf("Stop failed: %v", e.Err)
		} else {
			l.Logger.Info().Msg("Application stopped successfully")
		}
	case *fxevent.RollingBack:
		l.Logger.Error().Msgf("Rolling back due to start failure: %v", e.StartErr)
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Logger.Error().Msgf("Rollback failed: %v", e.Err)
		} else {
			l.Logger.Info().Msg("Rollback completed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Logger.Error().Msgf("Application start failed: %v", e.Err)
		} else {
			l.Logger.Info().Msg("Application started successfully!")
		}
	case *fxevent.LoggerInitialized:
		constructorName := l.simplifyFunctionName(e.ConstructorName)
		if e.Err != nil {
			l.Logger.Error().Msgf("Logger init failed: %s - %v", constructorName, e.Err)
		} else {
			l.Logger.Info().Msgf("Logger initialized: %s", constructorName)
		}
	}
}
