package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/asshiddiq1306/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := CreateNewTransferArgs{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQuery.CreateNewTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	return transfer
}

func TestCreateNewTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransferByID(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	transfer2, err := testQuery.GetTransferByID(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
}

func TestGetTransfersList(t *testing.T) {
	n := 10
	transfers := make([]Transfer, n)
	for i := 0; i < n; i++ {
		transfers[i] = createRandomTransfer(t)
	}

	arg := GetTransfersListArgs{
		Limit:  5,
		Offset: 5,
	}

	transferList, err := testQuery.GetTransfersList(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transferList, 5)

	for _, transfer := range transferList {
		require.NotEmpty(t, transfer)
	}
}

func TestUpdateTransferByID(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	arg := UpdateTransferByIDArgs{
		ID:     transfer1.ID,
		Amount: util.RandomMoney(),
	}
	transfer2, err := testQuery.UpdateTransferByID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, arg.Amount, transfer2.Amount)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
}

func TestDeleteTransferByID(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	err := testQuery.DeleteTransferByID(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQuery.GetTransferByID(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, transfer2)
}
