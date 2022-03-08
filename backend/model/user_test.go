package model

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	passwordTokens = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	userIdTokens   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

func randomeString(tokens string, n int) string {
	var sb strings.Builder
	k := len(tokens)

	for i := 0; i < n; i++ {
		c := tokens[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func TestValidatePlainPassword(t *testing.T) {
	for i := 0; i < 8; i++ {
		plain := randomeString(passwordTokens, i)
		res := ValidateUserPassword(plain)
		require.Error(t, res)
	}

	for i := 8; i <= 127; i++ {
		plain := randomeString(passwordTokens, i)
		res := ValidateUserPassword(plain)
		require.NoError(t, res)
	}

	plain := randomeString(passwordTokens, 128)
	res := ValidateUserPassword(plain)
	require.Error(t, res)
}

func TestValidateUserId(t *testing.T) {
	plain := randomeString(userIdTokens, 0)
	res := ValidateUserId(plain)
	require.Error(t, res)

	plain = randomeString(userIdTokens, 33)
	res = ValidateUserId(plain)
	require.Error(t, res)

	for i := 1; i <= 32; i++ {
		plain := randomeString(userIdTokens, i)
		res := ValidateUserId(plain)
		require.NoError(t, res)
	}
}

func TestValidateUserName(t *testing.T) {
	name := randomeString(userIdTokens, 0)
	res := ValidateUserName(name)
	require.Error(t, res)

	name = randomeString(userIdTokens, 33)
	res = ValidateUserName(name)
	require.Error(t, res)

	for i := 1; i <= 32; i++ {
		name := randomeString(userIdTokens, i)
		res = ValidateUserName(name)
		require.NoError(t, res)
	}
}
