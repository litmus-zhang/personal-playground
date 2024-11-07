-- name: CreateUser :one
INSERT INTO users (username, hashed_password, full_name, email)
VALUES ($1, $2, $3, $4)
RETURNING username, email, full_name, created_at, password_changed_at;

-- name: GetUser :one
SELECT username, email, full_name, created_at, password_changed_at
from users
where username = $1
LIMIT 1;