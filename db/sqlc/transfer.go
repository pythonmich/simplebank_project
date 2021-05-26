package db

import (
	"context"
	"database/sql"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (
       from_account_id, to_account_id, amount
) VALUES (
    $1, $2, $3
) RETURNING id,from_account_id,to_account_id,amount,created_at 
`

type CreateTransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount float64 `json:"account"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var transfer Transfer
	err := row.Scan(
		&transfer.Id,
		&transfer.FromAccountID,
		&transfer.ToAccountID,
		&transfer.Amount,
		&transfer.CreatedAt,
		)
	return transfer, err
}

const getTransfer = `-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id=$1
LIMIT 1`

func (q Queries) GetTransfer(ctx context.Context, id int64)(Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var transfer Transfer
	err := row.Scan(
		&transfer.Id,
		&transfer.FromAccountID,
		&transfer.ToAccountID,
		&transfer.Amount,
		&transfer.CreatedAt,
	)
	return transfer, err
}

const listTransfers = `-- name: ListTransfers :many
SELECT * FROM transfers
WHERE 
	from_account_id=$1 
			OR 
	to_account_id=$2
ORDER BY id
LIMIT $3
OFFSET $4;
`

type ListTransfersParams struct {
	FromAccountID int64
	ToAccountID int64
	Limit int32
	Offset int32
}
// ListTransfers gets an Account.Owner transfers from his account to other Account or transfers are sent to the Account.Owner.
// Eg. If Account.ID is 7 it will return a slice Transfer if the records are available where the Transfer.FromAccountID or Transfer.ToAccountID == 7.
func (q Queries) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)  {
	rows, err := q.db.QueryContext(ctx, listTransfers, arg.FromAccountID, arg.ToAccountID, arg.Limit, arg.Offset)
	if err != nil{
		return nil, err
	}

	defer func(r *sql.Rows) {
		err = r.Close()
	}(rows)
	var transfers []Transfer
	for rows.Next() {
		var transfer Transfer
		err := rows.Scan(
				&transfer.Id,
				&transfer.FromAccountID,
				&transfer.ToAccountID,
				&transfer.Amount,
				&transfer.CreatedAt,
			)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}
	if err := rows.Close(); err != nil{
		return nil, err
	}
	if err := rows.Err(); err != nil{
		return nil, err
	}
	return transfers, err
}