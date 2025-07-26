package database

import (
	"context"
	"fmt"
	"time"

	"github.com/evalabs-id/go-clean-architecture/internal/configs"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func ProvideDB(lc fx.Lifecycle, cfg *configs.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.Driver, cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.SSLMode)

	db, err := sqlx.ConnectContext(context.Background(), cfg.Database.Driver, dsn)

	if err != nil {
		log.Error().Err(err).Msg("cannot connect to database")
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(95)

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			log.Info().Msg("closing database connection")
			_ = db.Close()
			return nil
		},
	})

	return db, nil
}
