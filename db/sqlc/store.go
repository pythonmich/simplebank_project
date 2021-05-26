package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store will provide all functions to run db queries individually and as well as their combination within a transaction
type Store struct {
	*Queries
	// db is required to create a new db transaction
	db *sql.DB
}
// NewStore create a new store object
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a function within a database transaction
func (s *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil{
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil{
		if rbErr := tx.Rollback(); rbErr != nil{
			return fmt.Errorf("tx error: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains all necessary input parameters to transfer money btwn 2 accounts
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount float64 `json:"amount"`
}

// TransferTxResult it contains the results of the transfer transaction
type TransferTxResult struct {
	Transfer      Transfer `json:"transfer"`
	// FromAccount use from Account.ID
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry     Entry    `json:"from_entry"`
	ToEntry       Entry    `json:"to_entry"`
}

//var txKey = struct {}{}

// TransferTx performs a money transfer from one account to another
// it creates a new transfer record, add account entries and update account balance within a single db transaction
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		//txName := ctx.Value(txKey)
		//log.Println(txName, "Create Transfer")
		// this records a Transfer Record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		// From Account Entry
		//log.Println(txName, "Create Entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    - arg.Amount,
		})
		if err != nil {
			return err
		}
		//log.Println(txName, "Create Entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:     arg.Amount,
		})
		if err != nil {
			return err
		}

		//log.Println(txName, "Update Account 1")
		if arg.FromAccountID < arg.ToAccountID{
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)

		}else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

		}
		return nil
	})
	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountID1 int64,
	amount1 float64, accountID2 int64, amount2 float64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}