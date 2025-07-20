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
		Logger: l.With().Str("component", "fx_lifecycle").Logger(),
	}
}

type logger struct {
	Logger zerolog.Logger
}

func (l *logger) Printf(format string, args ...any) {
	l.Logger.Info().Msgf(format, args...)
}

func (l *logger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Logger.Info().
			Str("func", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStart Hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Logger.Warn().Err(e.Err).
				Str("func", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStart hook failed")
		} else {
			l.Logger.Info().
				Str("func", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStart hook executed")
		}
	case *fxevent.OnStopExecuting:
		l.Logger.Info().Str("func", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Logger.Warn().Err(e.Err).
				Str("func", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStop hook failed")
		} else {
			l.Logger.Info().Str("func", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStop hook executed")
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.Logger.Warn().Err(e.Err).
				Str("type", e.TypeName).
				Msg("supplied")
		} else {
			l.Logger.Info().
				Str("type", e.TypeName).
				Msg("supplied")
		}
	case *fxevent.Provided:
		for _, rangeType := range e.OutputTypeNames {
			l.Logger.Info().Str("type", rangeType).
				Str("constructor", e.ConstructorName).
				Msg("provided")
		}
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).
				Msg("error encountered while applying options")
		}
	case *fxevent.Invoking:
		// Do nothing. Will log on Invoked.
	case *fxevent.Invoked:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Str("stacktrace", e.Trace).
				Str("function", e.FunctionName).Msg("invoke failed")
		} else {
			l.Logger.Info().Str("function", e.FunctionName).Msg("invoked")
		}
	case *fxevent.Stopping:
		l.Logger.Info().Str("signal", strings.ToUpper(e.Signal.String())).Msg("received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("stop failed")
		}
	case *fxevent.RollingBack:
		l.Logger.Error().Err(e.StartErr).Msg("start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("start failed")
		} else {
			// TODO: add more context, or let it be empty
			// l.Logger.Info().Msg("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("custom logger initialization failed")
		} else {
			l.Logger.Info().Str("function", e.ConstructorName).Msg("initialized custom fxevent.Logger")
		}
	}
}
