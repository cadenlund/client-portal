// Author: Caden Lund
// Created: 4/12/2026
// Last updated: 4/12/2026
// Notes: - AuthService owns credential flows: Register, Login, and later
// password reset / change. User CRUD (profile, avatar, etc.) lives in
// UserService when that exists.

package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cadenlund/client-portal/apps/api/internal/auth"
	"github.com/cadenlund/client-portal/apps/api/internal/domain"
	"github.com/cadenlund/client-portal/apps/api/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// AuthService handles authentication and credential management.
type AuthService struct {
	repo   *repository.Queries
	hasher *auth.Hasher
}

// NewAuthService wires the repo and password hasher into an AuthService.
func NewAuthService(repo *repository.Queries, hasher *auth.Hasher) *AuthService {
	return &AuthService{repo: repo, hasher: hasher}
}

// Register creates a new user with a hashed password.
func (s *AuthService) Register(ctx context.Context, email, password, name string) (*repository.User, error) {
	//1. Hash password
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("register - failed to hash: %w", err)
	}

	//2. Create user
	user, err := s.repo.CreateUser(ctx, repository.CreateUserParams{
		Email:        email,
		PasswordHash: hash,
		Name:         &name,
		AvatarUrl:    nil,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrEmailTaken
		}

		// Default error
		return nil, fmt.Errorf("register - error inserting user: %w", err)
	}

	return &user, nil
}

// Login verifies credentials and returns the matching user.
func (s *AuthService) Login(ctx context.Context, email, password string) (*repository.User, error) {
	//1. Fetch user
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInvalidCredentials
		}

		// Default
		return nil, fmt.Errorf("login: failed to get user: %w", err)
	}

	//2. Compare and hash passwords
	match, err := s.hasher.Verify(password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("login: failed to compare & hash: %w", err)
	}

	//3. Check match
	if !match {
		return nil, domain.ErrInvalidCredentials
	}

	return &user, nil
}
