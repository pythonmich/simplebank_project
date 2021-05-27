package db

import (
	"context"
	"database/sql"
)

const createAccount = `-- name: CreateAccount :one
	INSERT INTO accounts (
		owner,
		balance,
		currency
	) VALUES (
		$1, $2, $3
	) RETURNING id, owner, balance, currency, created_at
`

type CreateAccountParams struct {
	Owner string `json:"owner"`
	Balance float64 `json:"balance"`
	Currency string `json:"currency"`
}


func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams)(Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.Owner,arg.Balance, arg.Currency)
	var account Account
	err := row.Scan(
		&account.ID,
		&account.Owner,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
	)
	return account, err
}

const getAccount = `-- name: GetAccount :one	
	SELECT id, owner, balance, currency, created_at FROM accounts
	WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var account Account
	err := row.Scan(
		&account.ID,
		&account.Owner,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		)
	return account, err
}

const getAccountForUpdate = `-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountForUpdate, id)
	var account Account
	err := row.Scan(
		&account.ID,
		&account.Owner,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
	)
	return account, err
}

type ListAccountParams struct {
	Limit int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

const listAccounts = `-- name: ListAccounts :many
	SELECT id, owner, balance, currency, created_at FROM accounts
	ORDER BY id
	LIMIT $1
	OFFSET $2;
`
func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountParams) ([]Account, error)  {
	rows, err := q.db.QueryContext(ctx, listAccounts, arg.Limit, arg.Offset)
	if err != nil{
		return nil, err
	}
	defer func(r *sql.Rows){
		err = r.Close()
	}(rows)

	var listAccounts = []Account{}
	for rows.Next(){
		var account Account
		if err := rows.Scan(
			&account.ID,
			&account.Owner,
			&account.Balance,
			&account.Currency,
			&account.CreatedAt,
		); err != nil{
			return nil, err
		}
		listAccounts = append(listAccounts, account)
	}
	if err := rows.Close(); err != nil{
		return nil, err
	}
	return listAccounts, err
}

const updateAccount = `-- name: UpdateAccount :one
UPDATE accounts SET balance = $2
WHERE id = $1 RETURNING id, owner, balance, currency, created_at
`

type UpdateAccountParams struct {
	Id int64 `json:"id"`
	Balance float64 `json:"balance"`
}

func (q Queries) UpdateAccount(ctx context.Context, args UpdateAccountParams) (Account,error) {
	row := q.db.QueryRowContext(ctx,updateAccount, args.Id, args.Balance)
	var account Account
	err := row.Scan(
		&account.ID,
		&account.Owner,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
	)
	return account, err
}
const addAccountBalance = `-- name: AddAccountBalance :one
UPDATE accounts SET balance = balance + $2
WHERE id = $1
RETURNING id, owner, balance, currency, created_at`

type AddAccountBalanceParams struct {
	ID int64 `json:"id"`
	Amount float64 `json:"amount"`
}

func (q Queries) AddAccountBalance(ctx context.Context, args AddAccountBalanceParams) (Account,error) {
	row := q.db.QueryRowContext(ctx,addAccountBalance, args.ID, args.Amount)
	var account Account
	err := row.Scan(
		&account.ID,
		&account.Owner,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
	)
	return account, err
}


const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1
`


func (q Queries) DeleteAccount(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}
//func (q Queries) UpdateAccount(ctx context.Context, args UpdateAccountParams) error {
//	_, err := q.db.ExecContext(ctx,updateAccount, args)
//	if err != nil{
//		return err
//	}
//	return err
//}