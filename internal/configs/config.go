package configs

import "fmt"

type Provider interface {
	Provide() error
}

type Config struct {
	App      App
	Server   Server
	Database Database
	JWT      JWT
}

func ProvideConfig() (*Config, error) {
	config := &Config{}

	providers := []Provider{
		&config.App,
		&config.Server,
		&config.Database,
		&config.JWT,
	}

	for _, p := range providers {
		if err := p.Provide(); err != nil {
			return nil, fmt.Errorf("failed to provide config: %w", err)
		}
	}

	return config, nil
}
