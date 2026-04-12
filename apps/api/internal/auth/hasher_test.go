// Author: Caden Lund
// Created: 4/12/2026
// Last updated: 4/12/2026

package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var hasher = NewHasher(DefaultConfig)

func TestVerify_Round_Trip(t *testing.T) {
	const password = "Password123123"
	const wrongPassowrd = "wrong password"

	hash, err := hasher.Hash(password)
	require.NoError(t, err)

	t.Run("Correct password", func(t *testing.T) {
		match, err := hasher.Verify(password, hash)
		require.NoError(t, err)
		assert.True(t, match, "Correct password should match")
	})

	t.Run("Incorrect password", func(t *testing.T) {
		match, err := hasher.Verify(wrongPassowrd, hash)
		require.NoError(t, err)
		assert.False(t, match, "Wrong password should not match")
	})

}
