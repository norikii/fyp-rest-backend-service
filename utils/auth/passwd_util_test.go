package auth_test

import (
	"github.com/tatrasoft/fyp-rest-backend-service/utils/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validPass  = "password123456"
const invalidPass = "12345"

func TestHashAndSaltPwd_Success(t *testing.T) {
	hash, err := auth.HashAndSaltPwd(validPass)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestIsValidPassword_PassAndHashMatch(t *testing.T) {
	hash, err := auth.HashAndSaltPwd(validPass)
	require.NoError(t, err)

	isValid, err := auth.IsValidPassword(hash, validPass)
	require.NoError(t, err)
	assert.True(t, isValid)
}

func TestIsValidPassword_PassAndHashDoNotMatch(t *testing.T) {
	hash, err := auth.HashAndSaltPwd(validPass)
	require.NoError(t, err)

	isValid, err := auth.IsValidPassword(hash, invalidPass)
	require.Error(t, err)
	assert.False(t, isValid)
}


