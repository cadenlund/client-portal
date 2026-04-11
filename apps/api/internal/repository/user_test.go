// Author: Caden Lund
// Created: 4/11/2026
// Last updated: 4/11/2026
// Notes: - the Ptr function stores anything in a variable and returns its address
// 		  - needed for *string input because *string type only takes &of
//        - For enums, use what model generates for you

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/cadenlund/client-portal/internal/testutil"
	"github.com/cadenlund/client-portal/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Use generated params struct for each method
// cant use := in global scope
var user = CreateUserParams{
	Email:        "test@example.com",
	PasswordHash: "hashedPassword",
	Name:         util.Ptr("John Doe"), // Needed because we return address of for a pointer type
	AvatarUrl:    nil,
}

var ctx = context.Background() // generic, never canceled

func TestCreate_User(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		//1. First create the queries dependency struct (Generated from sqlc)
		// Use a transaction here that rolls back on every subtest
		q := New(testutil.WithTx(t, testPool)) // pass t to register rollback on cleanup

		//2. Insert user
		actual, err := q.CreateUser(ctx, user)

		//3. Check error
		require.NoError(t, err)

		//4. Check inserted
		assert.Equal(t, user.Email, actual.Email)
		assert.Equal(t, user.Name, actual.Name)
		assert.Equal(t, user.PasswordHash, actual.PasswordHash)
		assert.Nil(t, actual.AvatarUrl)

		//5. Check defaults
		assert.NotEqual(t, uuid.Nil, actual.ID)    // check if uuid is not the zero uuid
		assert.Equal(t, UserRoleUser, actual.Role) // Use UserRoleUser from generated models
		assert.WithinDuration(t, time.Now(), actual.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), actual.UpdatedAt, time.Second)
	})

	t.Run("Duplicate email", func(t *testing.T) {
		q := New(testutil.WithTx(t, testPool))
		//1. First insert succeeds
		_, err := q.CreateUser(ctx, user)
		require.NoError(t, err)

		//2. Second insert with same email fails
		_, err = q.CreateUser(ctx, user)

		//3. Create empty pgconn error, check if they are the same type
		var pgErr *pgconn.PgError
		require.ErrorAs(t, err, &pgErr)
		assert.Equal(t, "23505", pgErr.Code) // 23505  means unique contstraint violation
	})

}

func TestGet_user_by_email(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		q := New(testutil.WithTx(t, testPool))

		//1. Insert user
		expected, err := q.CreateUser(ctx, user)
		require.NoError(t, err)

		//2. Get by email
		actual, err := q.GetUserByEmail(ctx, user.Email)
		require.NoError(t, err)

		//3. assert
		assert.Equal(t, expected, actual)
	})

	t.Run("Not found", func(t *testing.T) {
		q := New(testutil.WithTx(t, testPool))

		//1. Get by email
		_, err := q.GetUserByEmail(ctx, user.Email)

		//2. Require error
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestGet_user_by_ID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		q := New(testutil.WithTx(t, testPool))

		//1. Insert user
		expected, err := q.CreateUser(ctx, user)
		require.NoError(t, err)

		//2. Get by ID
		actual, err := q.GetUserByID(ctx, expected.ID)
		require.NoError(t, err)

		//3. Assert
		assert.Equal(t, expected, actual)
	})

	t.Run("Not found", func(t *testing.T) {
		q := New(testutil.WithTx(t, testPool))

		//1. Get by id
		_, err := q.GetUserByID(ctx, uuid.New()) // Use uuid.New() to create a random uuid v4

		//2. Require error
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})

}
