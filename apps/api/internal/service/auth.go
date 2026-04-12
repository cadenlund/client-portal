// Author: Caden Lund
// Created: 4/12/2026
// Last updated: 4/12/2026
// Notes: - AuthService owns credential flows: Register, Login, and later
// password reset / change. User CRUD (profile, avatar, etc.) lives in
// UserService when that exists.

package service

import (
	"context"

	"github.com/cadenlund/client-portal/apps/api/internal/auth"
	"github.com/cadenlund/client-portal/apps/api/internal/repository"
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
	return nil, nil
}

// Login verifies credentials and returns the matching user.
func (s *AuthService) Login(ctx context.Context, email, password string) (*repository.User, error) {
	return nil, nil
}
