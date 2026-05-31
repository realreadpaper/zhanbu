package config

import "testing"

func TestLoadFromEnv_ReadsAIConfig(t *testing.T) {
	t.Setenv("ZHANBU_AI_API_KEY", "sk-test")
	t.Setenv("ZHANBU_AI_MODEL", "test-model")
	t.Setenv("ZHANBU_AI_BASE_URL", "https://ai.example.test/v1")

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
}
