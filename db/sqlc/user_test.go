package db

import (
	"context"
	"github.com/emon46/bank-application/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	passwd := util.RandomString(6)
	hashedPasswd, err := util.GetHashPassword(passwd)
	arg := CreateUserParams{
		Username:       util.RandomName(),
		HashedPassword: hashedPasswd,
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.TODO(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func Test_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func Test_GetUser(t *testing.T) {
	user := createRandomUser(t)

	returnedUser, err := testQueries.GetUser(context.TODO(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, returnedUser)
	require.Equal(t, user.Username, returnedUser.Username)
	require.Equal(t, user.HashedPassword, returnedUser.HashedPassword)
	require.Equal(t, user.FullName, returnedUser.FullName)
	require.Equal(t, user.Email, returnedUser.Email)

	require.WithinDuration(t, user.PasswordChangedAt, returnedUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user.CreatedAt, returnedUser.CreatedAt, time.Second)
}
