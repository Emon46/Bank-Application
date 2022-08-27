package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_TransferTx(t *testing.T) {
	store := NewStore(testDB)
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)
	n := 50
	amount := int64(5)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.TODO(), TransferTxParams{
				FromAccountID: fromAcc.ID,
				ToAccountID:   toAcc.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	visited := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results
		transfer := result.Transfer
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAcc.ID, transfer.FromAccountID)
		require.Equal(t, toAcc.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.TODO(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAcc.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.TODO(), fromEntry.ID)
		require.NoError(t, err)
		//check entries
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toAcc.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.TODO(), toEntry.ID)
		require.NoError(t, err)

		// check accounts'

		returnedFromAcc := result.FromAccount
		require.NotEmpty(t, returnedFromAcc)
		require.Equal(t, fromAcc.ID, returnedFromAcc.ID)

		returnedToAcc := result.ToAccount
		require.NotEmpty(t, returnedToAcc)
		require.Equal(t, toAcc.ID, returnedToAcc.ID)
		// check accounts balance after transactions
		diff1 := fromAcc.Balance - returnedFromAcc.Balance
		diff2 := returnedToAcc.Balance - toAcc.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1%amount == 0 && diff1 > 0)

		require.NotContains(t, visited, diff1/amount)
		visited[int(diff1/amount)] = true
	}
	updatedFromAcc, err := testQueries.GetAccount(context.TODO(), fromAcc.ID)
	require.NoError(t, err)
	updatedToAcc, err := testQueries.GetAccount(context.TODO(), toAcc.ID)
	require.NoError(t, err)
	require.Equal(t, fromAcc.Balance, updatedFromAcc.Balance+int64(n)*amount)
	require.Equal(t, toAcc.Balance, updatedToAcc.Balance-int64(n)*amount)
}

func Test_TransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)
	n := 10
	amount := int64(5)
	errs := make(chan error)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			go func() {
				_, err := store.TransferTx(context.TODO(), TransferTxParams{
					FromAccountID: fromAcc.ID,
					ToAccountID:   toAcc.ID,
					Amount:        amount,
				})
				errs <- err
			}()
		} else {
			go func() {
				_, err := store.TransferTx(context.TODO(), TransferTxParams{
					FromAccountID: toAcc.ID,
					ToAccountID:   fromAcc.ID,
					Amount:        amount,
				})
				errs <- err
			}()
		}
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	updatedFromAcc, err := testQueries.GetAccount(context.TODO(), fromAcc.ID)
	require.NoError(t, err)
	updatedToAcc, err := testQueries.GetAccount(context.TODO(), toAcc.ID)
	require.NoError(t, err)

	require.Equal(t, fromAcc.Balance, updatedFromAcc.Balance)
	require.Equal(t, toAcc.Balance, updatedToAcc.Balance)
}
