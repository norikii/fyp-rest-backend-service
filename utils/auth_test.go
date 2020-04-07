package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	userId  = "password123456"
	userName  = "password123456"
	userEmail  = "password123456"
)


func TestCreateJWTToken_Success(t *testing.T) {
	tkn, err := CreateJWTToken(userId, userName, userEmail)

	assert.NoError(t, err)
	assert.NotEmpty(t, tkn)
}

func TestIsValidJWTToken_ValidToken(t *testing.T) {
	tkn, err := CreateJWTToken(userId, userName, userEmail)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	isValid, err := IsValidJWTToken(tkn)
	assert.NoError(t, err)
	assert.True(t, isValid)
}

