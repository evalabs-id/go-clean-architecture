package configs

import (
	"fmt"
	"os"
	"strconv"
)

type Database struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	SSLMode  string
}

func (d *Database) Provide() error {
	d.Driver = os.Getenv("DB_DRIVER")
	if d.Driver == "" {
		return fmt.Errorf("database driver cannot be empty")
	}
	d.Host = os.Getenv("DB_HOST")
	if d.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}
	d.Port = 5432 // default port
	if port := os.Getenv("DB_PORT"); port != "" {
		var err error
		d.Port, err = strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid database port: %w", err)
		}
	}
	d.Username = os.Getenv("DB_USERNAME")
	if d.Username == "" {
		return fmt.Errorf("database username cannot be empty")
	}
	d.Password = os.Getenv("DB_PASSWORD")
	if d.Password == "" {
		return fmt.Errorf("database password cannot be empty")
	}
	d.Name = os.Getenv("DB_NAME")
	if d.Name == "" {
		return fmt.Errorf("database name cannot be empty")
	}
	d.SSLMode = os.Getenv("DB_SSLMODE")
	if d.SSLMode == "" {
		d.SSLMode = "disable" // default SSL mode
	}

	return nil
}
