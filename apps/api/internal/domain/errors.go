// Author: Caden Lund
// Created: 4/12/2026
// Last updated: 4/12/2026
// Notes: - Sentinel errors the service layer emits for known business failures.
// Handlers branch on these via errors.Is and map to HTTP status codes.

package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already registered")
)
