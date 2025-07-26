package configs

import (
	"fmt"
	"os"
)

type App struct {
	Name        string
	Version     string
	Environment string
}

func (a *App) Provide() error {
	a.Name = os.Getenv("APP_NAME")
	if a.Name == "" {
		return fmt.Errorf("app name is required")
	}
	a.Version = os.Getenv("APP_VERSION")
	if a.Version == "" {
		return fmt.Errorf("app version is required")
	}
	a.Environment = os.Getenv("APP_ENVIRONMENT")
	if a.Environment == "" {
		return fmt.Errorf("app environment is required")
	}

	return nil
}
