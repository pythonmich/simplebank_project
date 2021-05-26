package db

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{})(sql.Result, error)
	PrepareContext(ctx context.Context, query string)(*sql.Stmt,error)
	QueryContext(ctx context.Context, query string, args ...interface{})(*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db:db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{db: tx}
}