// Author: Caden Lund
// Created: 4/12/2026
// Last updated: 4/12/2026
// Notes: - Wraps github.com/alexedwards/argon2id behind a small
// AuthConfig/Auth surface so callers in the api don't depend on it directly.
// - Hash format (PHC string): $argon2id$v=19$m=<KiB>,t=<iters>,p=<lanes>$<b64Salt>$<b64Key>

// Argon2id notes:
// - Increasing memory causes the algorithm to
// allocate more psuedo random data to r/w over
// - Increasing memory decreases brute force gpu attacks
// - Iterations causes the algorithm to go over that memory
// multiple times
// - Iterations increases CPU on our server
// - Parallelism changes output and our server can use more cores
// - Salt is random bytes mixed into hash so users with
// same password have a different hash
// - Would allow hackers to compute rainbow tables
// tables of precomputed common passwords
// 16 bytes = 128 bits of entropy
// - Key length is how long the output hash is
// - longer = harder to guess

package auth

import (
	"github.com/alexedwards/argon2id"
)

// AuthConfig holds Argon2id tuning parameters.
type AuthConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultAuthConfig is the OWASP minimum for Argon2id.
var DefaultAuthConfig = AuthConfig{
	Memory:      19 * 1024,
	Iterations:  2,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

// Auth hashes and verifies passwords using a fixed config.
type Auth struct {
	params *argon2id.Params
}

func NewAuth(cfg AuthConfig) *Auth {
	return &Auth{
		params: &argon2id.Params{
			Memory:      cfg.Memory,
			Iterations:  cfg.Iterations,
			Parallelism: cfg.Parallelism,
			SaltLength:  cfg.SaltLength,
			KeyLength:   cfg.KeyLength,
		},
	}
}

// Hash returns a PHC-formatted Argon2id hash of password.
func (a *Auth) Hash(password string) (string, error) {
	return argon2id.CreateHash(password, a.params)
}

// Verify reports whether password matches encodedHash.
func (a *Auth) Verify(password, encodedHash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, encodedHash)
}
