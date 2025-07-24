package authv1

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/evalabs-id/go-clean-architecture/internal/configs"
	"github.com/evalabs-id/go-clean-architecture/pkg/jwthelper"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignUp(ctx context.Context, req SignUpRequest) (*AuthResponse, error)
	SignIn(ctx context.Context, req SignInRequest) (*AuthResponse, error)
}

type service struct {
	repo   Repository
	config *configs.Config
}

// ProvideService creates a new auth service
func ProvideService(repo Repository, config *configs.Config) Service {
	return &service{
		repo:   repo,
		config: config,
	}
}

// SignUp creates a new user account
func (s *service) SignUp(ctx context.Context, req SignUpRequest) (*AuthResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if email already exists
	emailExists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}

	if emailExists {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &User{
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Name:     strings.TrimSpace(req.Name),
		Password: hashedPassword,
		IsActive: true,
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokenPair, err := jwthelper.GenerateTokenPair(s.config, strconv.Itoa(createdUser.ID), createdUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AuthResponse{
		User:      createdUser.ToUserResponse(),
		TokenPair: tokenPair,
	}, nil
}

// SignIn authenticates a user
func (s *service) SignIn(ctx context.Context, req SignInRequest) (*AuthResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	if err := s.verifyPassword(user.Password, req.Password); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Generate tokens
	tokenPair, err := jwthelper.GenerateTokenPair(s.config, strconv.Itoa(user.ID), user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AuthResponse{
		User:      user.ToUserResponse(),
		TokenPair: tokenPair,
	}, nil
}

// hashPassword hashes a password using bcrypt
func (s *service) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// verifyPassword verifies a password against its hash
func (s *service) verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
