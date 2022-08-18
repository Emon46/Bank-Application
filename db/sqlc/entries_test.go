package db

import (
	"context"
	"github.com/emon46/bank-application/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.TODO(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	return entry
}

func Test_CreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func Test_GetEntry(t *testing.T) {
	acc := createRandomAccount(t)
	entry := createRandomEntry(t, acc)

	returnedEntry, err := testQueries.GetEntry(context.TODO(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, returnedEntry)
	require.Equal(t, entry.ID, returnedEntry.ID)
	require.Equal(t, entry.AccountID, returnedEntry.AccountID)
	require.Equal(t, entry.Amount, returnedEntry.Amount)
	require.Equal(t, entry.CreatedAt, returnedEntry.CreatedAt)
}

func Test_ListEntry(t *testing.T) {
	acc := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, acc)
	}
	arg := ListEntriesParams{
		AccountID: acc.ID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntries(context.TODO(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
