package jwthelper

import (
	"errors"
	"fmt"
	"time"

	"github.com/evalabs-id/go-clean-architecture/internal/configs"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email,omitempty"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidTokenType = errors.New("invalid token type")
	ErrInvalidClaims    = errors.New("invalid token claims")
)

// GenerateTokenPair creates both access and refresh tokens
func GenerateTokenPair(config *configs.Config, userID, email string) (*TokenPair, error) {
	accessToken, accessExp, err := GenerateAccessToken(config, userID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, _, err := GenerateRefreshToken(config, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessExp.Unix(),
	}, nil
}

// GenerateAccessToken creates a new access token
func GenerateAccessToken(config *configs.Config, userID, email string) (string, time.Time, error) {
	duration, err := getAccessTokenDuration(config)
	if err != nil {
		return "", time.Time{}, err
	}

	expirationTime := time.Now().Add(duration)

	claims := TokenClaims{
		UserID:    userID,
		Email:     email,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.JWT.Issuer,
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWT.SecretKey))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expirationTime, nil
}

// GenerateRefreshToken creates a new refresh token
func GenerateRefreshToken(config *configs.Config, userID string) (string, time.Time, error) {
	duration, err := getRefreshTokenDuration(config)
	if err != nil {
		return "", time.Time{}, err
	}

	expirationTime := time.Now().Add(duration)

	claims := TokenClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.JWT.Issuer,
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWT.SecretKey))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expirationTime, nil
}

// ValidateToken validates and parses a JWT token
func ValidateToken(config *configs.Config, tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWT.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// ValidateAccessToken validates specifically access tokens
func ValidateAccessToken(config *configs.Config, tokenString string) (*TokenClaims, error) {
	claims, err := ValidateToken(config, tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

// ValidateRefreshToken validates specifically refresh tokens
func ValidateRefreshToken(config *configs.Config, tokenString string) (*TokenClaims, error) {
	claims, err := ValidateToken(config, tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

// RefreshAccessToken creates a new access token using a valid refresh token
func RefreshAccessToken(config *configs.Config, refreshTokenString string) (*TokenPair, error) {
	refreshClaims, err := ValidateRefreshToken(config, refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate new token pair
	return GenerateTokenPair(config, refreshClaims.UserID, refreshClaims.Email)
}

// ExtractUserID extracts user ID from token without full validation (for logging purposes)
func ExtractUserID(config *configs.Config, tokenString string) string {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.SecretKey), nil
	})

	if err != nil {
		return ""
	}

	if claims, ok := token.Claims.(*TokenClaims); ok {
		return claims.UserID
	}

	return ""
}

// IsTokenExpired checks if token is expired without validating signature
func IsTokenExpired(config *configs.Config, tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.SecretKey), nil
	})

	if err != nil {
		return true
	}

	if claims, ok := token.Claims.(*TokenClaims); ok {
		return claims.ExpiresAt.Time.Before(time.Now())
	}

	return true
}

// GetTokenRemainingTime returns how much time is left before token expires
func GetTokenRemainingTime(config *configs.Config, tokenString string) time.Duration {
	claims, err := ValidateToken(config, tokenString)
	if err != nil {
		return 0
	}

	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining < 0 {
		return 0
	}

	return remaining
}

// Helper functions for duration parsing
func getAccessTokenDuration(config *configs.Config) (time.Duration, error) {
	duration, err := time.ParseDuration(config.JWT.AccessTokenDuration)
	if err != nil {
		return 0, fmt.Errorf("invalid access token duration: %w", err)
	}
	return duration, nil
}

func getRefreshTokenDuration(config *configs.Config) (time.Duration, error) {
	duration, err := time.ParseDuration(config.JWT.RefreshTokenDuration)
	if err != nil {
		return 0, fmt.Errorf("invalid refresh token duration: %w", err)
	}
	return duration, nil
}
