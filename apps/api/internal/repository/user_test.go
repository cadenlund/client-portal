// Author: Caden Lund
// Created: 4/11/2026
// Last updated: 4/11/2026
// Notes: - the Ptr function stores anything in a variable and returns its address
// 		  - needed for *string input because *string type only takes &of
//        - For enums, use what model generates for you

package repository

import (
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

func TestCreate_User(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		user := CreateUserParams{
			Email:        "test@example.com",
			PasswordHash: "hashedPassword",
			Name:         util.Ptr("John Doe"), // Needed because we return address of for a pointer type
			AvatarUrl:    nil,
		}

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
		user := CreateUserParams{
			Email:        "test@example.com",
			PasswordHash: "hashedPassword",
			Name:         util.Ptr("John Doe"),
			AvatarUrl:    nil,
		}

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

func TestGet_User_By_Email(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		user := CreateUserParams{
			Email:        "test@example.com",
			PasswordHash: "hashedPassword",
			Name:         util.Ptr("John Doe"),
			AvatarUrl:    nil,
		}
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
		_, err := q.GetUserByEmail(ctx, "nonexistent@example.com")

		//2. Require error
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestGet_User_By_ID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		user := CreateUserParams{
			Email:        "test@example.com",
			PasswordHash: "hashedPassword",
			Name:         util.Ptr("John Doe"),
			AvatarUrl:    nil,
		}
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

func TestUpdate_User(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		user := CreateUserParams{
			Email:        "test@example.com",
			PasswordHash: "hashedPassword",
			Name:         util.Ptr("John Doe"),
			AvatarUrl:    nil,
		}
		update := UpdateUserParams{
			Name:      util.Ptr("Bobby Smith"),
			AvatarUrl: util.Ptr("Image_bucket.com"),
		}
		q := New(testutil.WithTx(t, testPool))

		//1. Insert user
		created, err := q.CreateUser(ctx, user)
		require.NoError(t, err)

		//2. Set ID to update user
		update.ID = created.ID

		//3. Update user
		actual, err := q.UpdateUser(ctx, update)
		require.NoError(t, err)

		//4. Assert
		assert.Equal(t, "Bobby Smith", *actual.Name)
		assert.Equal(t, "Image_bucket.com", *actual.AvatarUrl)
	})

	t.Run("Partial update", func(t *testing.T) {
		user := CreateUserParams{
			Email:        "test@example.com",
			PasswordHash: "hashedPassword",
			Name:         util.Ptr("John Doe"),
			AvatarUrl:    util.Ptr("https://old-avatar.com/img.png"),
		}
		q := New(testutil.WithTx(t, testPool))

		//1. Insert user
		created, err := q.CreateUser(ctx, user)
		require.NoError(t, err)

		//2. Create
		partialUpdate := UpdateUserParams{
			ID:        created.ID,
			Name:      nil,
			AvatarUrl: util.Ptr("Image_bucket.com"),
		}

		//3. Update user
		actual, err := q.UpdateUser(ctx, partialUpdate)
		require.NoError(t, err)

		//4. Assert - name unchanged, avatar updated
		assert.Equal(t, "John Doe", *actual.Name)
		assert.Equal(t, "Image_bucket.com", *actual.AvatarUrl)
	})

	t.Run("Not found", func(t *testing.T) {
		update := UpdateUserParams{
			ID:        uuid.New(),
			Name:      util.Ptr("New Name"),
			AvatarUrl: nil,
		}
		q := New(testutil.WithTx(t, testPool))

		//1. Try to update non-existent user
		_, err := q.UpdateUser(ctx, update)

		//2. Require error
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestUpdate_User_Password(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		user := CreateUserParams{
			Email:        "test@example.com",
			PasswordHash: "hashedPassword",
			Name:         util.Ptr("John Doe"),
			AvatarUrl:    util.Ptr("https://old-avatar.com/img.png"),
		}

		q := New(testutil.WithTx(t, testPool))

		//1. Insert user
		created, err := q.CreateUser(ctx, user)
		require.NoError(t, err)

		//2. Make update params (need id)
		update := UpdateUserPasswordParams{
			ID:           created.ID,
			PasswordHash: "NewHash",
		}

		//3. Update password
		updated, err := q.UpdateUserPassword(ctx, update)
		require.NoError(t, err)

		//4. Assert
		assert.Equal(t, "NewHash", updated.PasswordHash)
	})

	t.Run("Not found", func(t *testing.T) {
		update := UpdateUserPasswordParams{
			ID:           uuid.New(),
			PasswordHash: "NewHash",
		}
		q := New(testutil.WithTx(t, testPool))

		//1. Try to update non-existent user
		_, err := q.UpdateUserPassword(ctx, update)

		//2. Require error
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestClear_User_Avatar(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		user := CreateUserParams{
			Email:        "test@example.com",
			PasswordHash: "hashedPassword",
			Name:         util.Ptr("John Doe"),
			AvatarUrl:    util.Ptr("https://example.com/avatar.png"),
		}
		q := New(testutil.WithTx(t, testPool))

		//1. Insert user with avatar
		created, err := q.CreateUser(ctx, user)
		require.NoError(t, err)
		require.NotNil(t, created.AvatarUrl)

		//2. Clear avatar
		updated, err := q.ClearUserAvatar(ctx, created.ID)
		require.NoError(t, err)

		//3. Assert avatar is nil
		assert.Nil(t, updated.AvatarUrl)
	})

	t.Run("Not found", func(t *testing.T) {
		q := New(testutil.WithTx(t, testPool))

		//1. Try to clear avatar on non-existent user
		_, err := q.ClearUserAvatar(ctx, uuid.New())

		//2. Require error
		require.ErrorIs(t, err, pgx.ErrNoRows)
	})
}
