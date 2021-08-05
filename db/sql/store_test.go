package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	amount := int64(10)

	store := NewStore(testDB)

	n := 5
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxArg{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		require.NoError(t, err)
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.ID)

		_, err = store.GetTransferByID(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		_, err = store.GetEntryByID(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		_, err = store.GetEntryByID(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.Equal(t, account1.Owner, fromAccount.Owner)
		require.Equal(t, account1.Currency, fromAccount.Currency)

		_, err = store.GetAccountByID(context.Background(), fromAccount.ID)
		require.NoError(t, err)

		toAccount := result.ToAccount
		require.Equal(t, account2.Owner, toAccount.Owner)
		require.Equal(t, account2.Currency, toAccount.Currency)

		_, err = store.GetAccountByID(context.Background(), toAccount.ID)
		require.NoError(t, err)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	updateAccount1, err := store.GetAccountByID(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := store.GetAccountByID(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	amount := int64(10)

	store := NewStore(testDB)

	n := 10
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxArg{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs

		require.NoError(t, err)
	}

	updateAccount1, err := store.GetAccountByID(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := store.GetAccountByID(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
