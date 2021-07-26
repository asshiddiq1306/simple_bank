package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/asshiddiq1306/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account1 := createRandomAccount(t)
	arg := CreateNewEntryArgs{
		AccountID: account1.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQuery.CreateNewEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, account1.ID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	return entry
}

func TestCreateNewEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntryByID(t *testing.T) {
	entry1 := createRandomEntry(t)

	entry2, err := testQuery.GetEntryByID(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestGetEntryList(t *testing.T) {
	n := 10
	entries := make([]Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = createRandomEntry(t)
	}

	arg := GetEntryListArgs{
		Limit:  5,
		Offset: 5,
	}

	entryList, err := testQuery.GetEntryList(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entryList, 5)

	for _, entry := range entryList {
		require.NotEmpty(t, entry)
	}
}

func TestUpdateEntryByID(t *testing.T) {
	entry1 := createRandomEntry(t)

	arg := UpdateEntryByIDArgs{
		ID:     entry1.ID,
		Amount: util.RandomMoney(),
	}

	entry2, err := testQuery.UpdateEntryByID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, arg.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
}

func TestDeleteEntryByID(t *testing.T) {
	entry1 := createRandomEntry(t)

	err := testQuery.DeleteEntryByID(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQuery.GetEntryByID(context.Background(), entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}
