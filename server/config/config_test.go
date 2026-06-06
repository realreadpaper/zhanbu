package config

import "testing"

func TestLoadFromEnv_ReadsAIConfig(t *testing.T) {
	t.Setenv("ZHANBU_AI_API_KEY", "sk-test")
	t.Setenv("ZHANBU_AI_MODEL", "test-model")
	t.Setenv("ZHANBU_AI_BASE_URL", "https://ai.example.test/v1")
	t.Setenv("ZHANBU_AI_MAX_TOKENS", "4096")

	cfg := LoadFromEnv()

	if cfg.AI.APIKey != "sk-test" {
		t.Fatalf("expected AI API key from env, got %q", cfg.AI.APIKey)
	}
	if cfg.AI.Model != "test-model" {
		t.Fatalf("expected AI model from env, got %q", cfg.AI.Model)
	}
	if cfg.AI.BaseURL != "https://ai.example.test/v1" {
		t.Fatalf("expected AI base URL from env, got %q", cfg.AI.BaseURL)
	}
	if cfg.AI.MaxTokens != 4096 {
		t.Fatalf("expected AI max tokens from env, got %d", cfg.AI.MaxTokens)
	}
}

func TestLoadFromEnv_ReadsRenderPortAndCORSOrigins(t *testing.T) {
	t.Setenv("PORT", "10000")
	t.Setenv("ZHANBU_CORS_ALLOWED_ORIGINS", "https://zhanbu.vercel.app, https://example.fun")

	cfg := LoadFromEnv()

	if cfg.Server.Port != 10000 {
		t.Fatalf("expected server port from PORT env, got %d", cfg.Server.Port)
	}

	expected := []string{"https://zhanbu.vercel.app", "https://example.fun"}
	if len(cfg.CORS.AllowedOrigins) != len(expected) {
		t.Fatalf("expected %d CORS origins, got %#v", len(expected), cfg.CORS.AllowedOrigins)
	}
	for i := range expected {
		if cfg.CORS.AllowedOrigins[i] != expected[i] {
			t.Fatalf("expected CORS origin %d to be %q, got %q", i, expected[i], cfg.CORS.AllowedOrigins[i])
		}
	}
}

func TestLoadFromEnv_ZhanbuServerPortOverridesRenderPort(t *testing.T) {
	t.Setenv("PORT", "10000")
	t.Setenv("ZHANBU_SERVER_PORT", "18080")

	cfg := LoadFromEnv()

	if cfg.Server.Port != 18080 {
		t.Fatalf("expected ZHANBU_SERVER_PORT to override PORT, got %d", cfg.Server.Port)
	}
}
