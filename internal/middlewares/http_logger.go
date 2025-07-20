package middlewares

import (
	"net/http"
	"time"

	"github.com/evalabs-id/go-clean-architecture/pkg/helper/httphelper"
	"github.com/evalabs-id/go-clean-architecture/pkg/log"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type HttpMiddleware struct {
	Logger zerolog.Logger
}

func ProvideHttpMiddleware(logger zerolog.Logger) *HttpMiddleware {
	return &HttpMiddleware{
		Logger: logger.With().Str("component", "http_server").Logger(),
	}
}

func (m *HttpMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		requestID := middleware.GetReqID(r.Context())

		reqLogger := m.Logger.With().
			Str("request_id", requestID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_ip", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Str("referer", r.Referer()).
			Str("proto", r.Proto).
			Logger()

		defer func() {
			reqLogger.Info().
				Int("status", ww.Status()).
				Int("bytes_written", ww.BytesWritten()).
				Dur("duration", time.Since(start)).
				Msg("Request completed")
		}()

		next.ServeHTTP(ww, r)
	})
}

func (m *HttpMiddleware) Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.Logger.Error().
					Interface("recover_info", err).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Msg("panic recovered in http server")

				httphelper.SimpleError(
					"Internal Server Error",
					http.StatusInternalServerError,
					"PANIC_RECOVERED",
				).ServeHTTP(w, r)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (m *HttpMiddleware) LoggerContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := log.InjectLoggerToContext(r.Context(), m.Logger)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
