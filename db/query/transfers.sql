-- CRUD

-- name: CreateTransfer :one
INSERT INTO transfers (
       from_account_id, to_account_id, amount
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id=$1
LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE from_account_id=$1 OR to_account_id=$2
ORDER BY id
LIMIT $3
OFFSET $4;

--TODO: implement update from_account amount and to_account amount, later on implement a lifetime in which are record must be deleted

-- name: UpdateToTransfer :one
UPDATE transfers SET amount=$2
WHERE to_account_id=$1;

-- name: UpdateFromTransfer :one
UPDATE transfers SET amount=$2
WHERE to_account_id=$1;

