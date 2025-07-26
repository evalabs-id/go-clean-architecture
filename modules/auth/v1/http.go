package authv1

import (
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/evalabs-id/go-clean-architecture/internal/providers/logger"
	"github.com/evalabs-id/go-clean-architecture/pkg/httphelper"
)

var (
	sonicConfig = sonic.ConfigDefault
)

type Http struct {
	service Service
}

// ProvideHttp creates a new auth HTTP handler
func ProvideHttp(service Service) *Http {
	return &Http{
		service: service,
	}
}

// SignUp handles user registration
func (h *Http) SignUp(w http.ResponseWriter, r *http.Request) http.Handler {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var req SignUpRequest
	decoder := sonicConfig.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Error().Err(err).Msg("failed to decode sign up request")
		return httphelper.SimpleError(
			"Invalid request body",
			http.StatusBadRequest,
			"INVALID_REQUEST_BODY",
		)
	}

	response, err := h.service.SignUp(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to sign up user")

		// Check for specific error types
		errMsg := err.Error()
		if errMsg == "email already exists" {
			return httphelper.SimpleError(
				"Email already exists",
				http.StatusConflict,
				"EMAIL_ALREADY_EXISTS",
			)
		}

		if errMsg[:17] == "validation failed" {
			return httphelper.SimpleError(
				"Validation failed",
				http.StatusBadRequest,
				"VALIDATION_FAILED",
				err.Error(),
			)
		}

		return httphelper.SimpleError(
			"Failed to create user account",
			http.StatusInternalServerError,
			"SIGN_UP_FAILED",
		)
	}

	log.Info().
		Str("user_email", response.User.Email).
		Int("user_id", response.User.ID).
		Msg("user signed up successfully")

	return httphelper.Created(response)
}

// SignIn handles user authentication
func (h *Http) SignIn(w http.ResponseWriter, r *http.Request) http.Handler {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var req SignInRequest
	decoder := sonicConfig.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Error().Err(err).Msg("failed to decode sign in request")
		return httphelper.SimpleError(
			"Invalid request body",
			http.StatusBadRequest,
			"INVALID_REQUEST_BODY",
		)
	}

	response, err := h.service.SignIn(ctx, req)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("failed to sign in user")

		// Check for specific error types
		errMsg := err.Error()
		if errMsg == "invalid email or password" {
			return httphelper.SimpleError(
				"Invalid email or password",
				http.StatusUnauthorized,
				"INVALID_CREDENTIALS",
			)
		}

		if errMsg == "user account is inactive" {
			return httphelper.SimpleError(
				"User account is inactive",
				http.StatusForbidden,
				"ACCOUNT_INACTIVE",
			)
		}

		if errMsg[:17] == "validation failed" {
			return httphelper.SimpleError(
				"Validation failed",
				http.StatusBadRequest,
				"VALIDATION_FAILED",
				err.Error(),
			)
		}

		return httphelper.SimpleError(
			"Failed to sign in",
			http.StatusInternalServerError,
			"SIGN_IN_FAILED",
		)
	}

	log.Info().
		Str("user_email", response.User.Email).
		Int("user_id", response.User.ID).
		Msg("user signed in successfully")

	return httphelper.OK(response)
}
