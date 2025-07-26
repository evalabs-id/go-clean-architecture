package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/evalabs-id/go-clean-architecture/internal/configs"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func HttpServer(lc fx.Lifecycle, router *chi.Mux, cfg *configs.Config) *http.Server {
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Addr:              addr,
		Handler:           router,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				log.Info().Str("port", fmt.Sprintf("%d", cfg.Server.Port)).Msg("Starting HTTP server")
				err := srv.ListenAndServe()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error().Err(err).Msg("error starting http server")
					os.Exit(1)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			toCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := srv.Shutdown(toCtx); err != nil {
				log.Error().Err(err).Msg("can't shutdown http server")
				return err
			}
			log.Info().Msg("http server shutdown")
			return nil
		},
	})

	return srv
}
