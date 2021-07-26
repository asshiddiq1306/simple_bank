package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/asshiddiq1306/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateNewAccountArgs{
		Owner:    util.RandomName(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQuery.CreateNewAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateNewAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccountByID(t *testing.T) {
	account1 := createRandomAccount(t)

	account2, err := testQuery.GetAccountByID(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestGetListAccounts(t *testing.T) {
	n := 10
	accounts := make([]Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = createRandomAccount(t)
	}

	arg := GetListAccountsArgs{
		Limit:  5,
		Offset: 5,
	}

	accountList, err := testQuery.GetListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accountList, 5)

	for _, account := range accountList {
		require.NotEmpty(t, account)
	}
}

func TestUpdateAccountByID(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountByIDArgs{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQuery.UpdateAccountByID(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

func TestDeleteAccountByID(t *testing.T) {
	account := createRandomAccount(t)

	err := testQuery.DeleteAccountBuID(context.Background(), account.ID)
	require.NoError(t, err)

	account1, err := testQuery.GetAccountByID(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account1)

}
