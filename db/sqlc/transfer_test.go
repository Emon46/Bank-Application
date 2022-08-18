package db

import (
	"context"
	"github.com/emon46/bank-application/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomTransfer(t *testing.T, fromAcc, toAcc Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.TODO(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	return transfer
}

func Test_CreateTransfer(t *testing.T) {
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)
	createRandomTransfer(t, fromAcc, toAcc)
}

func Test_GetTransfer(t *testing.T) {
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)
	transfer := createRandomTransfer(t, fromAcc, toAcc)

	returnedTransfer, err := testQueries.GetTransfer(context.TODO(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, returnedTransfer)
	require.Equal(t, transfer.ID, returnedTransfer.ID)
	require.Equal(t, transfer.FromAccountID, returnedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, returnedTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, returnedTransfer.Amount)
	require.Equal(t, transfer.CreatedAt, returnedTransfer.CreatedAt)
}

func Test_ListTransfer(t *testing.T) {
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)
	for i := 0; i < 5; i++ {
		createRandomTransfer(t, fromAcc, toAcc)
		createRandomTransfer(t, toAcc, fromAcc)
	}
	arg := ListTransfersParams{
		FromAccountID: fromAcc.ID,
		ToAccountID:   fromAcc.ID,
		Limit:         5,
		Offset:        5,
	}
	transfers, err := testQueries.ListTransfers(context.TODO(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.NotEmpty(t, transfer)
		require.True(t, (transfer.FromAccountID == toAcc.ID || transfer.ToAccountID == toAcc.ID) && (transfer.FromAccountID == fromAcc.ID || transfer.ToAccountID == fromAcc.ID))
	}
}
