package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTGenerateAndValidateTokenPair(t *testing.T) {
	manager := NewJWTManager("test-secret-key", time.Hour, 7*24*time.Hour)

	t.Run("successfully generate and validate token pair", func(t *testing.T) {
		userID := uint(1)
		username := "testuser"

		tokens, err := manager.GenerateTokenPair(userID, username)
		require.NoError(t, err)
		require.NotEmpty(t, tokens.AccessToken)
		require.NotEmpty(t, tokens.RefreshToken)

		// Validate access token
		claims, err := manager.ValidateToken(tokens.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, username, claims.Username)
	})

	t.Run("access token expires correctly", func(t *testing.T) {
		shortManager := NewJWTManager("test-secret", 1*time.Second, 1*time.Second)
		tokens, err := shortManager.GenerateTokenPair(1, "test")
		require.NoError(t, err)

		// Token should be valid immediately
		_, err = shortManager.ValidateToken(tokens.AccessToken)
		require.NoError(t, err)

		// Wait for token to expire
		time.Sleep(1100 * time.Millisecond)

		// Token should now be invalid
		_, err = shortManager.ValidateToken(tokens.AccessToken)
		require.Error(t, err)
	})

	t.Run("reject invalid token", func(t *testing.T) {
		_, err := manager.ValidateToken("invalid-token-string")
		require.Error(t, err)
	})

	t.Run("reject token signed with different secret", func(t *testing.T) {
		manager1 := NewJWTManager("secret1", time.Hour, time.Hour)
		manager2 := NewJWTManager("secret2", time.Hour, time.Hour)

		tokens, err := manager1.GenerateTokenPair(1, "test")
		require.NoError(t, err)

		_, err = manager2.ValidateToken(tokens.AccessToken)
		require.Error(t, err)
	})
}

func TestJWTExpirationSettings(t *testing.T) {
	t.Run("access token has shorter TTL than refresh token", func(t *testing.T) {
		manager := NewJWTManager("test-secret", 1*time.Hour, 7*24*time.Hour)
		tokens, err := manager.GenerateTokenPair(1, "test")
		require.NoError(t, err)

		accessClaims, _ := manager.ValidateToken(tokens.AccessToken)
		refreshClaims, _ := manager.ValidateToken(tokens.RefreshToken)

		accessExp := accessClaims.ExpiresAt.Time
		refreshExp := refreshClaims.ExpiresAt.Time

		assert.True(t, refreshExp.After(accessExp), "refresh token should expire later than access token")
	})

	t.Run("token pair contains correct expires_in", func(t *testing.T) {
		manager := NewJWTManager("test-secret", 30*time.Minute, 7*24*time.Hour)
		tokens, err := manager.GenerateTokenPair(1, "test")
		require.NoError(t, err)

		assert.Equal(t, int64(1800), tokens.ExpiresIn)
	})
}
