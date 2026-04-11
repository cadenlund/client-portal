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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate_User(t *testing.T) {
	// Generic context
	ctx := context.Background()

	// Create our test models to use
	user := CreateUserParams{
		Email:        "test@example.com",
		PasswordHash: "hashedPassword",
		Name:         util.Ptr("John Doe"), // Needed because we return address of for a pointer type
		AvatarUrl:    nil,
	}

	t.Run("Create user: success", func(t *testing.T) {
		//1. First create the queries dependency struct (Generated from sqlc)
		// Use a transaction here that rolls back on every subtest
		q := New(testutil.WithTx(t, testPool))

		//2. Insert user
		got, err := q.CreateUser(ctx, user)

		//3. Check error
		require.NoError(t, err)

		//4. Check inserted
		assert.Equal(t, user.Email, got.Email)
		assert.Equal(t, user.Name, got.Name)
		assert.Equal(t, user.PasswordHash, got.PasswordHash)
		assert.Nil(t, got.AvatarUrl)

		//5. Check defaults
		assert.NotEqual(t, uuid.Nil, got.ID)    // check if uuid is not the zero uuid
		assert.Equal(t, UserRoleUser, got.Role) // Use UserRoleUser from generated models
		assert.WithinDuration(t, time.Now(), got.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), got.UpdatedAt, time.Second)
	})

	t.Run("Create user: duplicate email", func(t *testing.T) {
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
