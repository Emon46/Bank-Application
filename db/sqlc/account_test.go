package db

import (
	"context"
	"github.com/emon46/bank-application/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomName(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.TODO(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func Test_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func Test_GetAccount(t *testing.T) {
	acc := createRandomAccount(t)

	returnedAcc, err := testQueries.GetAccount(context.TODO(), acc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, returnedAcc)
	require.Equal(t, acc.ID, returnedAcc.ID)
	require.Equal(t, acc.Owner, returnedAcc.Owner)
	require.Equal(t, acc.Balance, returnedAcc.Balance)
	require.Equal(t, acc.Currency, returnedAcc.Currency)
	require.Equal(t, acc.CreatedAt, returnedAcc.CreatedAt)
}

func Test_UpdateAccount(t *testing.T) {
	acc := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      acc.ID,
		Balance: util.RandomMoney(),
	}
	returnedAcc, err := testQueries.UpdateAccount(context.TODO(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, returnedAcc)
	require.Equal(t, acc.ID, returnedAcc.ID)
	require.Equal(t, acc.Owner, returnedAcc.Owner)
	require.Equal(t, arg.Balance, returnedAcc.Balance)
	require.Equal(t, acc.Currency, returnedAcc.Currency)
	require.Equal(t, acc.CreatedAt, returnedAcc.CreatedAt)
}

func Test_ListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccount(context.TODO(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
