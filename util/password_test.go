package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)
	hashedPassword, err := GetHashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)
	//check for mismatched password

	wrongPass := RandomString(6)
	err = CheckPassword(wrongPass, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	password1 := RandomString(6)
	hashedPassword1, err := GetHashPassword(password1)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)
	require.NotEqual(t, hashedPassword, hashedPassword1)

}
