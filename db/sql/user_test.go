package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/asshiddiq1306/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	password := util.RandomString(6)
	hashPassword, err := util.HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)

	user, err := testQuery.CreateNewUser(context.Background(), CreateNewUserArgs{
		Username:       util.RandomName(),
		HashedPassword: hashPassword,
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, user)
	return user
}

func TestCreateNewUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByID(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQuery.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestGetListUser(t *testing.T) {
	n := 10
	users := make([]User, n)
	for i := 0; i < n; i++ {
		users[i] = createRandomUser(t)
	}

	arg := GetListUserArgs{
		Limit:  5,
		Offset: 5,
	}

	userList, err := testQuery.GetListUser(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, userList, 5)

	for _, user := range userList {
		require.NotEmpty(t, user)
	}
}

func TestDeleteUserByID(t *testing.T) {
	user1 := createRandomUser(t)

	err := testQuery.DeleteUserByID(context.Background(), user1.Username)
	require.NoError(t, err)

	user2, err := testQuery.GetUserByUsername(context.Background(), user1.Username)
	require.Error(t, err)
	require.Empty(t, user2)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
