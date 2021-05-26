package db

import (
	"GoBankProject/util"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T, acc Account) Entry {
	arg := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestQueries_CreateEntry(t *testing.T) {
	acc1 := createRandomAccount(t)
	createRandomEntry(t, acc1)
}

func TestQueries_GetEntry(t *testing.T) {
	acc1 := createRandomAccount(t)
	entry1 := createRandomEntry(t, acc1)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt,time.Second)
}

func TestQueries_ListTransfers(t *testing.T) {
	acc1 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, acc1)
	}
	arg := GetAllEntriesParams{
		AccountID: acc1.ID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.GetAllEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}

}