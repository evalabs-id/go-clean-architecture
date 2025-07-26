package configs

import (
	"fmt"
	"os"
)

type JWT struct {
	SecretKey            string
	AccessTokenDuration  string // in minutes
	RefreshTokenDuration string // in hours
	Issuer               string
}

func (j *JWT) Provide() error {
	j.SecretKey = os.Getenv("JWT_SECRET_KEY")
	if j.SecretKey == "" {
		return fmt.Errorf("JWT secret key cannot be empty")
	}

	j.AccessTokenDuration = "15m" // default 15 minutes
	if accessDuration := os.Getenv("JWT_ACCESS_TOKEN_DURATION"); accessDuration != "" {
		j.AccessTokenDuration = accessDuration
	}
	j.RefreshTokenDuration = "24h" // default 24 hours
	if refreshDuration := os.Getenv("JWT_REFRESH_TOKEN_DURATION"); refreshDuration != "" {
		j.RefreshTokenDuration = refreshDuration
	}
	j.Issuer = os.Getenv("JWT_ISSUER")
	if j.Issuer == "" {
		j.Issuer = "default-issuer" // default issuer
	}
	return nil
}
