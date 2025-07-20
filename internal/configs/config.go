package configs

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	App      App
	Server   Server
	Database Database
}

type App struct {
	Name        string
	Version     string
	Environment string
}

type Server struct {
	Port     int
	Debug    bool
	Timezone string
	ApiKey   string
}

type Database struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	SSLMode  string
}

func ProvideConfig() (*Config, error) {
	var config Config

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	config.App.Name = os.Getenv("APP_NAME")
	if config.App.Name == "" {
		return fmt.Errorf("app name cannot be empty")
	}
	config.App.Version = os.Getenv("APP_VERSION")
	if config.App.Version == "" {
		return fmt.Errorf("app version cannot be empty")
	}
	config.App.Environment = os.Getenv("APP_ENVIRONMENT")
	if config.App.Environment == "" {
		return fmt.Errorf("app environment cannot be empty")
	}
	config.Server.Port = 8080 // default port
	if port := os.Getenv("SERVER_PORT"); port != "" {
		var err error
		config.Server.Port, err = strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid server port: %w", err)
		}
	}
	config.Server.Debug = false // default debug mode
	if debug := os.Getenv("SERVER_DEBUG"); debug != "" {
		var err error
		config.Server.Debug, err = strconv.ParseBool(debug)
		if err != nil {
			return fmt.Errorf("invalid server debug value: %w", err)
		}
	}
	config.Server.Timezone = "UTC" // default timezone
	if timezone := os.Getenv("SERVER_TIMEZONE"); timezone != "" {
		config.Server.Timezone = timezone
	}
	config.Server.ApiKey = os.Getenv("API_KEY")
	if config.Server.ApiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	config.Database.Driver = os.Getenv("DB_DRIVER")
	if config.Database.Driver == "" {
		return fmt.Errorf("database driver cannot be empty")
	}
	config.Database.Host = os.Getenv("DB_HOST")
	if config.Database.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}
	config.Database.Port = 5432 // default port
	if port := os.Getenv("DB_PORT"); port != "" {
		var err error
		config.Database.Port, err = strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid database port: %w", err)
		}
	}
	config.Database.Username = os.Getenv("DB_USERNAME")
	if config.Database.Username == "" {
		return fmt.Errorf("database username cannot be empty")
	}
	config.Database.Password = os.Getenv("DB_PASSWORD")
	if config.Database.Password == "" {
		return fmt.Errorf("database password cannot be empty")
	}
	config.Database.Name = os.Getenv("DB_NAME")
	if config.Database.Name == "" {
		return fmt.Errorf("database name cannot be empty")
	}
	config.Database.SSLMode = "disable" // default SSL mode
	if sslMode := os.Getenv("DB_SSLMODE"); sslMode != "" {
		config.Database.SSLMode = sslMode
	}

	return nil
}
