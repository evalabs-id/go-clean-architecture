package configs

import (
	"fmt"
	"os"
	"strconv"
)

type Server struct {
	Port     int
	Debug    bool
	Timezone string
	ApiKey   string
}

func (s *Server) Provide() error {
	s.Port = 5000 // default port
	if port := os.Getenv("SERVER_PORT"); port != "" {
		var err error
		s.Port, err = strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid server port: %w", err)
		}
	}

	s.Debug = false // default debug mode
	if debug := os.Getenv("SERVER_DEBUG"); debug != "" {
		var err error
		s.Debug, err = strconv.ParseBool(debug)
		if err != nil {
			return fmt.Errorf("invalid server debug value: %w", err)
		}
	}

	s.Timezone = "UTC" // default timezone
	if tz := os.Getenv("SERVER_TIMEZONE"); tz != "" {
		s.Timezone = tz
	}

	s.ApiKey = os.Getenv("SERVER_API_KEY")
	if s.ApiKey == "" {
		return fmt.Errorf("server API key is required")
	}

	return nil
}
