package db

import (
	"context"
	"database/sql"
)

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries(
        account_id, amount 
)VALUES (
    $1, $2
)RETURNING id, account_id, amount, created_at
`

type CreateEntryParams struct {
	AccountID int64 `json:"account_id"`
	Amount float64 `json:"amount"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error){
	row := q.db.QueryRowContext(ctx, createEntry, arg.AccountID, arg.Amount)
	var entry Entry
	err := row.Scan(
		&entry.ID,
		&entry.AccountID,
		&entry.Amount,
		&entry.CreatedAt,
	)
	return entry,err
}

const getEntry = `-- name: GetEntry :one
SELECT id, account_id, amount, created_at FROM entries
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetEntry(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntry, id)
	var entry Entry
	err := row.Scan(
		&entry.ID,
		&entry.AccountID,
		&entry.Amount,
		&entry.CreatedAt,
	)
	return entry,err
}

const getAllEntries = `-- name: GetAllEntries :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type GetAllEntriesParams struct {
	AccountID int64 `json:"account_id"`
	Limit int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllEntries(ctx context.Context, arg GetAllEntriesParams) ([]Entry, error) {
	row, err := q.db.QueryContext(ctx, getAllEntries, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil{
		return nil, err
	}
	defer func(r *sql.Rows) {
		err = r.Close()
	}(row)

	var entries []Entry

	for row.Next() {
		var entry Entry
		if err := row.Scan(
			&entry.ID,
			&entry.AccountID,
			&entry.Amount,
			&entry.CreatedAt,
		); err != nil{
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := row.Close(); err != nil {
		return nil, err
	}
	return entries, err
}