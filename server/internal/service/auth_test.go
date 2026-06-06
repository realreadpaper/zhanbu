package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"zhanbu/config"
	"zhanbu/internal/model"
	"zhanbu/internal/repository"
	apperrors "zhanbu/pkg/errors"
	"zhanbu/pkg/utils"
)

// setupTestDB creates an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.EmailVerification{}, &model.SendLog{}))
	return db
}

// setupAuthService creates an AuthService with an in-memory DB for testing.
func setupAuthService(t *testing.T) (*AuthService, *gorm.DB) {
	t.Helper()
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	jwtConfig := config.JWTConfig{
		Secret:     "test-secret-key",
		AccessTTL:  time.Hour,
		RefreshTTL: 7 * 24 * time.Hour,
	}
	jwtManager := utils.NewJWTManager(jwtConfig.Secret, jwtConfig.AccessTTL, jwtConfig.RefreshTTL)
	verRepo := repository.NewVerificationRepository(db)
	smtpConfig := config.SMTPConfig{Enabled: false}
	securityConfig := config.SecurityConfig{
		VerifyEmail:    false,
		CodeLength:     6,
		CodeExpiry:     10 * time.Minute,
		MaxSendPerHour: 5,
	}
	emailService := NewEmailService(verRepo, &smtpConfig, &securityConfig)
	svc := NewAuthService(userRepo, jwtManager, &jwtConfig, emailService, &securityConfig, &smtpConfig)
	return svc, db
}

func TestAuthService_Register(t *testing.T) {
	t.Run("successfully register a new user", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		req := &RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		result, appErr := svc.Register(req)
		require.Nil(t, appErr)
		require.NotNil(t, result)
		assert.Equal(t, "testuser", result.User.Username)
		assert.Equal(t, "test@example.com", result.User.Email)
		assert.NotZero(t, result.User.ID)
		assert.NotZero(t, result.User.CreatedAt)
		assert.False(t, result.NeedVerify)
	})

	t.Run("reject duplicate email", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		req1 := &RegisterRequest{
			Username: "user1",
			Email:    "same@example.com",
			Password: "password123",
		}
		_, appErr := svc.Register(req1)
		require.Nil(t, appErr)

		req2 := &RegisterRequest{
			Username: "user2",
			Email:    "same@example.com",
			Password: "password456",
		}
		_, appErr = svc.Register(req2)
		require.NotNil(t, appErr)
		assert.Equal(t, apperrors.ErrUserExists, appErr.Code)
	})

	t.Run("reject duplicate username", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		req1 := &RegisterRequest{
			Username: "duplicate",
			Email:    "email1@example.com",
			Password: "password123",
		}
		_, appErr := svc.Register(req1)
		require.Nil(t, appErr)

		req2 := &RegisterRequest{
			Username: "duplicate",
			Email:    "email2@example.com",
			Password: "password456",
		}
		_, appErr = svc.Register(req2)
		require.NotNil(t, appErr)
		assert.Equal(t, apperrors.ErrUserExists, appErr.Code)
	})

	t.Run("reject short username", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		req := &RegisterRequest{
			Username: "ab", // too short
			Email:    "test@example.com",
			Password: "password123",
		}
		_, appErr := svc.Register(req)
		require.NotNil(t, appErr)
		assert.Equal(t, apperrors.ErrInvalidUsername, appErr.Code)
	})

	t.Run("reject invalid email format", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		req := &RegisterRequest{
			Username: "testuser",
			Email:    "not-an-email",
			Password: "password123",
		}
		_, appErr := svc.Register(req)
		require.NotNil(t, appErr)
		assert.Equal(t, apperrors.ErrInvalidEmail, appErr.Code)
	})

	t.Run("reject weak password", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		req := &RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "123", // too short
		}
		_, appErr := svc.Register(req)
		require.NotNil(t, appErr)
		assert.Equal(t, apperrors.ErrWeakPassword, appErr.Code)
	})

	t.Run("normalize email to lowercase", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		req := &RegisterRequest{
			Username: "testuser",
			Email:    "Test@Example.COM",
			Password: "password123",
		}
		result, appErr := svc.Register(req)
		require.Nil(t, appErr)
		assert.Equal(t, "test@example.com", result.User.Email)
	})
}

func TestAuthService_Login(t *testing.T) {
	t.Run("successfully login with valid credentials", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		// Register first
		_, appErr := svc.Register(&RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		})
		require.Nil(t, appErr)

		// Login
		result, appErr := svc.Login(&LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		})
		require.Nil(t, appErr)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)
		assert.Equal(t, "testuser", result.User.Username)
		assert.Equal(t, "test@example.com", result.User.Email)
	})

	t.Run("reject wrong password", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		_, appErr := svc.Register(&RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "correctpassword",
		})
		require.Nil(t, appErr)

		_, appErr = svc.Login(&LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		})
		require.NotNil(t, appErr)
		assert.Equal(t, apperrors.ErrInvalidCreds, appErr.Code)
	})

	t.Run("reject non-existent email", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		_, appErr := svc.Login(&LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		})
		require.NotNil(t, appErr)
		assert.Equal(t, apperrors.ErrInvalidCreds, appErr.Code)
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	t.Run("successfully refresh token", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		// Register and login
		_, appErr := svc.Register(&RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		})
		require.Nil(t, appErr)

		loginResult, appErr := svc.Login(&LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		})
		require.Nil(t, appErr)

		// Refresh with valid refresh token
		result, appErr := svc.RefreshToken(loginResult.RefreshToken)
		require.Nil(t, appErr)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)
	})

	t.Run("reject invalid refresh token", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		_, appErr := svc.RefreshToken("invalid-token")
		require.NotNil(t, appErr)
		assert.Equal(t, apperrors.ErrTokenInvalid, appErr.Code)
	})
}

func TestAuthService_GetProfile(t *testing.T) {
	t.Run("successfully get user profile", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		result, appErr := svc.Register(&RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		})
		require.Nil(t, appErr)

		profile, appErr := svc.GetProfile(result.User.ID)
		require.Nil(t, appErr)
		require.NotNil(t, profile)
		assert.Equal(t, result.User.ID, profile.ID)
		assert.Equal(t, "testuser", profile.Username)
	})

	t.Run("return error for non-existent user", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		_, appErr := svc.GetProfile(99999)
		require.NotNil(t, appErr)
		assert.Equal(t, apperrors.ErrNotFound, appErr.Code)
	})
}

func TestAuthService_UpdateProfile(t *testing.T) {
	t.Run("successfully update username", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		result, _ := svc.Register(&RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		})

		updated, appErr := svc.UpdateProfile(result.User.ID, &UpdateProfileRequest{
			Username: "newusername",
		})
		require.Nil(t, appErr)
		assert.Equal(t, "newusername", updated.Username)
	})

	t.Run("successfully update avatar and zodiac", func(t *testing.T) {
		svc, _ := setupAuthService(t)

		result, _ := svc.Register(&RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		})

		updated, appErr := svc.UpdateProfile(result.User.ID, &UpdateProfileRequest{
			Avatar: "https://example.com/avatar.jpg",
			Zodiac: "aries",
		})
		require.Nil(t, appErr)
		assert.Equal(t, "https://example.com/avatar.jpg", updated.Avatar)
		assert.Equal(t, "aries", updated.Zodiac)
	})
}
