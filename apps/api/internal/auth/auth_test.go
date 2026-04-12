// Author: Caden Lund
// Created: 4/12/2026
// Last updated: 4/12/2026

package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var auth = NewAuth(DefaultAuthConfig)

func TestVerify_Round_Trip(t *testing.T) {
	const password = "Password123123"

	//1. Hash
	hash, err := auth.Hash(password)
	require.NoError(t, err)

	t.Run("Correct password", func(t *testing.T) {
		match, err := auth.Verify(password, hash)
		require.NoError(t, err)
		assert.True(t, match, "Correct password should match")
	})

	t.Run("Incorrect password", func(t *testing.T) {
		match, err := auth.Verify("wrong password", hash)
		require.NoError(t, err)
		assert.False(t, match, "Wrong password should not match")
	})

}
