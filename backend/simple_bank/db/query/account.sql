-- name: CreateAccount :one
INSERT INTO accounts (balance, owner, currency)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccount :one
SELECT *
from accounts
where id = $1
LIMIT 1;

-- name: ListAccounts :many
SELECT *
from accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccounts :one
UPDATE accounts 
SET balance = $1
WHERE id = $2
RETURNING *;

-- name: DeleteAccount :exec
DELETE from accounts
WHERE id = $1;