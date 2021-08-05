package token

import (
	"testing"
	"time"

	"github.com/asshiddiq1306/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomName()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)
	accessToken, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)

	payload, err := maker.VerifyToken(accessToken)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	accessToken, err := maker.CreateToken(util.RandomName(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)

	payload, err := maker.VerifyToken(accessToken)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrExpiredToken.Error())
}
