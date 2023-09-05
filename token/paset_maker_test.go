package token

import (
	"testing"
	"time"

	"github.com/mhg14/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPASETOMaker(t *testing.T) {
	maker, err := NewPASETOMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt.Time, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt.Time, time.Second)
}

func TestExpiredPASETOToken(t *testing.T) {
	maker, err := NewPASETOMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := maker.CreateToken(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
