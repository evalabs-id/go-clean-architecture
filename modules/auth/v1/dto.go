package authv1

import (
	"time"

	"github.com/evalabs-id/go-clean-architecture/pkg/jwthelper"
	validation "github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// User represents the user model
type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"-" db:"password"` // Hidden from JSON
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SignUpRequest represents the sign up request payload
type SignUpRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Validate validates the sign up request
func (s SignUpRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Email, validation.Required, is.Email),
		validation.Field(&s.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&s.Password, validation.Required, validation.Length(8, 255)),
	)
}

// SignInRequest represents the sign in request payload
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the sign in request
func (s SignInRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Email, validation.Required, is.Email),
		validation.Field(&s.Password, validation.Required),
	)
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User      UserResponse         `json:"user"`
	TokenPair *jwthelper.TokenPair `json:"tokens"`
}

// UserResponse represents the user data in responses
type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse converts User to UserResponse
func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
