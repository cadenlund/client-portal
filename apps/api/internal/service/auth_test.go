// Author: Caden Lund
// Created: 4/18/2026
// Last updated: 4/18/2026
// Notes:

package service

import (
	"testing"

	"github.com/cadenlund/client-portal/apps/api/internal/auth"
	"github.com/cadenlund/client-portal/apps/api/internal/domain"
	"github.com/cadenlund/client-portal/apps/api/internal/repository"
	"github.com/cadenlund/client-portal/apps/api/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		//1. First Define our params
		email := "test@test.com"
		password := "TestPass123"
		name := "John Doe"

		//2. Setup deps
		// Use a transaction here that rolls back on every subtest
		q := repository.New(testutil.WithTx(t, testPool)) // pass t to register rollback on cleanup

		h := auth.NewHasher(auth.DefaultConfig)

		s := NewAuthService(q, h)

		//3. Register user
		actual, err := s.Register(ctx, email, password, name)
		require.NoError(t, err)

		//4. Assert return
		assert.Equal(t, email, actual.Email)
		assert.NotEmpty(t, actual.PasswordHash)
		assert.Equal(t, name, actual.Name)
	})

	t.Run("Duplicate email", func(t *testing.T) {
		//1. Define params
		email := "test@test.com"
		password := "TestPass123"
		name := "John Doe"

		//2. Setup deps
		q := repository.New(testutil.WithTx(t, testPool))

		h := auth.NewHasher(auth.DefaultConfig)

		s := NewAuthService(q, h)

		//3. Register user
		_, err := s.Register(ctx, email, password, name)
		require.NoError(t, err)

		//4. Duplicate register
		_, err = s.Register(ctx, email, password, name)
		require.ErrorIs(t, err, domain.ErrEmailTaken)
	})
}

func TestLogin(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		//1. Define params
		email := "test@test.com"
		password := "TestPass123"
		name := "John Doe"

		//2. Setup deps
		q := repository.New(testutil.WithTx(t, testPool))

		h := auth.NewHasher(auth.DefaultConfig)

		s := NewAuthService(q, h)

		//3. Register user
		expected, err := s.Register(ctx, email, password, name)
		require.NoError(t, err)

		//4. Try to login with same user
		actual, err := s.Login(ctx, email, password)
		require.NoError(t, err)

		//5. Assertions
		assert.Equal(t, expected.Email, actual.Email)
		assert.Equal(t, expected.PasswordHash, actual.PasswordHash)
		assert.Equal(t, expected.Name, actual.Name)

	})

	t.Run("Wrong password", func(t *testing.T) {
		//1. Define params
		email := "test@test.com"
		password := "TestPass123"
		wrongPass := "WrongPass123"
		name := "John Doe"

		//2. Setup deps
		q := repository.New(testutil.WithTx(t, testPool))

		h := auth.NewHasher(auth.DefaultConfig)

		s := NewAuthService(q, h)

		//3. Register user
		_, err := s.Register(ctx, email, password, name)
		require.NoError(t, err)

		//4. Try to login with wrong password
		_, err = s.Login(ctx, email, wrongPass)
		require.ErrorIs(t, err, domain.ErrInvalidCredentials)

	})

	t.Run("Email not found", func(t *testing.T) {
		//1. Define params
		email := "test@test.com"
		wrongEmail := "wrong@wrong.com"
		password := "TestPass123"
		name := "John Doe"

		//2. Setup deps
		q := repository.New(testutil.WithTx(t, testPool))

		h := auth.NewHasher(auth.DefaultConfig)

		s := NewAuthService(q, h)

		//3. Register user
		_, err := s.Register(ctx, email, password, name)
		require.NoError(t, err)

		//4. Try to login with wrong email
		_, err = s.Login(ctx, wrongEmail, password)
		require.ErrorIs(t, err, domain.ErrInvalidCredentials)
	})
}
