package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	log.Println("-> Before: ", acc1.Balance, acc2.Balance)
//	 run a concurrent transfer transaction
	n := 5
	amount := float64(10)
	errs := make(chan  error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		//txName := fmt.Sprintf("tx -> %d", i + 1)
		go func() {
			//ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(context.Background(),TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	// get results from channel
	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)

		result := <- results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.Id)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.Id)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc2.ID, toAccount.ID)

		// check account's balances
		log.Println("-> tx: ", fromAccount.Balance, toAccount.Balance)
		dif1 := acc1.Balance - fromAccount.Balance
		dif2 := toAccount.Balance - acc2.Balance
		require.Equal(t, dif1, dif2)
		require.True(t, dif1 > 0)
		require.True(t, int64(dif1) % int64(amount) == 0)

		k := int(dif1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
//	check the final updated balance of the two accounts
	updateAcc1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updateAcc2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, int64(acc1.Balance) - int64(n) * int64(amount), int64(updateAcc1.Balance))

	require.Equal(t, int64(acc2.Balance) + int64(n) * int64(amount), int64(updateAcc2.Balance))

}

func TestStore_TransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	log.Println("-> Before: ", acc1.Balance, acc2.Balance)
	//	 run a concurrent transfer transaction
	n := 10
	amount := float64(10)
	errs := make(chan  error)
	for i := 0; i < n; i++ {
		fromAccountId := acc1.ID
		toAccountId := acc2.ID

		// Algorithm of half of the transactions to send money from account 2 to account 1
		if i % 2 == 1{
			fromAccountId = acc2.ID
			toAccountId =  acc1.ID
		}
		//txName := fmt.Sprintf("tx -> %d", i + 1)
		go func() {
			//ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(context.Background(),TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	// get results from channel
	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)
	}
	//	check the final updated balance of the two accounts
	updateAcc1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updateAcc2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, int64(acc1.Balance), int64(updateAcc1.Balance))

	require.Equal(t, int64(acc2.Balance), int64(updateAcc2.Balance))

}
