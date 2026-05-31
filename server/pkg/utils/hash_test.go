package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("hash a password successfully", func(t *testing.T) {
		password := "testpassword123"
		hash, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash) // hash should be different from plaintext
	})

	t.Run("hash starts with bcrypt prefix", func(t *testing.T) {
		hash, err := HashPassword("mypassword")
		assert.NoError(t, err)
		assert.Contains(t, hash, "$2a$")
	})
}

func TestCheckPassword(t *testing.T) {
	t.Run("correct password returns true", func(t *testing.T) {
		password := "testpassword123"
		hash, err := HashPassword(password)
		assert.NoError(t, err)

		result := CheckPassword(password, hash)
		assert.True(t, result)
	})

	t.Run("incorrect password returns false", func(t *testing.T) {
		hash, err := HashPassword("correctpassword")
		assert.NoError(t, err)

		result := CheckPassword("wrongpassword", hash)
		assert.False(t, result)
	})

	t.Run("empty password against valid hash returns false", func(t *testing.T) {
		hash, err := HashPassword("somepassword")
		assert.NoError(t, err)

		result := CheckPassword("", hash)
		assert.False(t, result)
	})
}
