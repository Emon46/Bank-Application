-- name: CreateAccount :one
INSERT INTO account (
    owner, balance, currency
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM account
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM account
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccount :many
SELECT * FROM account
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE account
set balance = $2
WHERE id = $1
RETURNING *;

-- name: AddToAccountBalance :one
UPDATE account
set balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
    RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM account
WHERE id = $1;