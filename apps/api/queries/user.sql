-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, avatar_url) 
VALUES ($1, $2, $3, $4)
RETURNING *;