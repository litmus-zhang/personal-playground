-- name: CreateAccount :one
INSERT INTO accounts (balance, owner, currency)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccount :one
SELECT *
from accounts
where id = $1
LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT *
from accounts
where id = $1
LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT *
from accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts 
SET balance = $1
WHERE id = $2
RETURNING *;

-- name: UpdateAccountBalance :one
UPDATE accounts 
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE from accounts
WHERE id = $1;