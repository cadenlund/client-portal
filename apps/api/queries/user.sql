-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, avatar_url) 
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET 
    name = COALESCE($2, name),
    avatar_url = COALESCE($3, avatar_url)
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET password_hash = $2
WHERE id = $1
RETURNING *;

-- name: ClearUserAvatar :one
UPDATE users
SET avatar_url = NULL
WHERE id = $1
RETURNING *;
