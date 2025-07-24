package authv1

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/evalabs-id/go-clean-architecture/internal/queries"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}

type repository struct {
	db *sqlx.DB
}

// ProvideRepository creates a new auth repository
func ProvideRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

// CreateUser creates a new user in the database
func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var newUser User

	err := r.db.QueryRowxContext(
		ctx,
		queries.InsertUser,
		user.Email,
		user.Name,
		user.Password,
		user.IsActive,
	).StructScan(&newUser)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}

// GetUserByEmail retrieves a user by email
func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	err := r.db.QueryRowxContext(
		ctx,
		queries.GetUserByEmail,
		email,
	).StructScan(&user)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *repository) GetUserByID(ctx context.Context, id int) (*User, error) {
	var user User

	err := r.db.QueryRowxContext(
		ctx,
		queries.GetUserByID,
		id,
	).StructScan(&user)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// EmailExists checks if an email already exists in the database
func (r *repository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int

	err := r.db.QueryRowxContext(
		ctx,
		queries.CountCheckEmail,
		email,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}
