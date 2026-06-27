package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zhanbu/config"
	"zhanbu/pkg/utils"
)

func TestSetupRouter_AIInterpretRouteExistsWhenProviderUnavailable(t *testing.T) {
	cfg := &config.Config{
		Server:    config.ServerConfig{Mode: "test"},
		JWT:       config.JWTConfig{Secret: "test-secret", AccessTTL: time.Hour, RefreshTTL: time.Hour},
		AI:        config.AIConfig{APIKey: ""},
		RateLimit: config.RateLimitConfig{APIPerMinute: 60},
	}
	r := SetupRouter(nil, cfg, zerolog.Nop(), &config.ProfilesConfig{
		Profiles:        make(map[string]config.ProfileConfig),
		DefaultBindings: make(map[string]string),
	})

	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	tokenPair, err := jwtManager.GenerateTokenPair(1, "tester")
	require.NoError(t, err)

	body := bytes.NewBufferString(`{"type":"bazi","result":"{}","question":"八字排盘"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/ai/interpret", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Contains(t, rec.Body.String(), "AI service is not configured")
}

func TestSetupRouter_RefreshTokenRouteDoesNotRequireAccessToken(t *testing.T) {
	cfg := &config.Config{
		Server:    config.ServerConfig{Mode: "test"},
		JWT:       config.JWTConfig{Secret: "test-secret", AccessTTL: time.Hour, RefreshTTL: time.Hour},
		AI:        config.AIConfig{APIKey: ""},
		RateLimit: config.RateLimitConfig{APIPerMinute: 60},
	}
	r := SetupRouter(nil, cfg, zerolog.Nop(), &config.ProfilesConfig{
		Profiles:        make(map[string]config.ProfileConfig),
		DefaultBindings: make(map[string]string),
	})

	body := bytes.NewBufferString(`{}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Contains(t, rec.Body.String(), "refresh_token is required")
}
