package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "prompt_profiles.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write temp yaml: %v", err)
	}
	return path
}

func TestLoadProfiles_Success(t *testing.T) {
	yaml := `
prompt_profiles:
  shaoyong_meihua:
    version: "v1"
    divination_type: "meihua"
    name: "康节先生"
    title: "梅花易数解读者"
    subtitle: "师承邵雍一脉"
    icon: "🌸"
    description: "test"
    enabled: true
    system_identity: "你是康节先生"
    reasoning_framework:
      - "先看体用生克"
    voice_style:
      - "沉稳通透"
    output_structure:
      - "卦象总览"
    guardrails:
      - "不得编造"

default_bindings:
  meihua: "shaoyong_meihua"
`
	path := writeTempYAML(t, yaml)

	cfg, err := LoadProfiles(path)
	if err != nil {
		t.Fatalf("LoadProfiles: %v", err)
	}

	if len(cfg.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(cfg.Profiles))
	}
	if len(cfg.DefaultBindings) != 1 {
		t.Fatalf("expected 1 binding, got %d", len(cfg.DefaultBindings))
	}
}

func TestLoadProfiles_FileNotFound_FallsBackToEmbedded(t *testing.T) {
	// When file is missing, the embedded YAML is used as fallback.
	// This ensures profiles work in Docker where config files aren't copied.
	cfg, err := LoadProfiles("/nonexistent/path.yaml")
	if err != nil {
		t.Fatalf("expected fallback to embedded YAML, got error: %v", err)
	}
	if len(cfg.Profiles) == 0 {
		t.Fatal("expected profiles from embedded YAML")
	}
	if cfg.DefaultBindings["meihua"] != "shaoyong_meihua" {
		t.Fatalf("expected meihua→shaoyong_meihua binding, got %q", cfg.DefaultBindings["meihua"])
	}
}

func TestDefaultProfile_Success(t *testing.T) {
	cfg := &ProfilesConfig{
		Profiles: map[string]ProfileConfig{
			"shaoyong_meihua": {
				Version:        "v1",
				DivinationType: "meihua",
				Name:           "康节先生",
				Enabled:        true,
			},
		},
		DefaultBindings: map[string]string{
			"meihua": "shaoyong_meihua",
		},
	}

	profile, err := cfg.DefaultProfile("meihua")
	if err != nil {
		t.Fatalf("DefaultProfile: %v", err)
	}
	if profile.Name != "康节先生" {
		t.Fatalf("expected name 康节先生, got %s", profile.Name)
	}
	if profile.Version != "v1" {
		t.Fatalf("expected version v1, got %s", profile.Version)
	}
}

func TestDefaultProfile_NoBinding(t *testing.T) {
	cfg := &ProfilesConfig{
		Profiles:        map[string]ProfileConfig{},
		DefaultBindings: map[string]string{},
	}

	_, err := cfg.DefaultProfile("unknown_type")
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestDefaultProfile_ProfileNotFound(t *testing.T) {
	cfg := &ProfilesConfig{
		Profiles: map[string]ProfileConfig{},
		DefaultBindings: map[string]string{
			"meihua": "nonexistent",
		},
	}

	_, err := cfg.DefaultProfile("meihua")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestDefaultProfile_Disabled(t *testing.T) {
	cfg := &ProfilesConfig{
		Profiles: map[string]ProfileConfig{
			"shaoyong_meihua": {
				Version:        "v1",
				DivinationType: "meihua",
				Name:           "康节先生",
				Enabled:        false,
			},
		},
		DefaultBindings: map[string]string{
			"meihua": "shaoyong_meihua",
		},
	}

	_, err := cfg.DefaultProfile("meihua")
	if err == nil {
		t.Fatal("expected error for disabled profile")
	}
}

func TestProfile_Success(t *testing.T) {
	cfg := &ProfilesConfig{
		Profiles: map[string]ProfileConfig{
			"shaoyong_meihua": {
				Version: "v1",
				Name:    "康节先生",
			},
		},
	}

	profile, err := cfg.Profile("shaoyong_meihua")
	if err != nil {
		t.Fatalf("Profile: %v", err)
	}
	if profile.Name != "康节先生" {
		t.Fatalf("expected name 康节先生, got %s", profile.Name)
	}
}

func TestProfile_NotFound(t *testing.T) {
	cfg := &ProfilesConfig{
		Profiles: map[string]ProfileConfig{},
	}

	_, err := cfg.Profile("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestLoadProfiles_P0OnlyMeihuaProfileEnabled(t *testing.T) {
	// Load the actual YAML file from the config directory
	// This test verifies the real config file loads correctly
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	cfg, err := LoadProfiles(filepath.Join(cwd, "prompt_profiles.yaml"))
	if err != nil {
		t.Fatalf("LoadProfiles: %v", err)
	}

	expectedProfiles := []string{
		"shaoyong_meihua",
	}
	if len(cfg.Profiles) != len(expectedProfiles) {
		t.Errorf("expected %d profiles, got %d", len(expectedProfiles), len(cfg.Profiles))
	}
	for _, id := range expectedProfiles {
		if _, ok := cfg.Profiles[id]; !ok {
			t.Errorf("missing profile: %s", id)
		}
	}

	expectedBindings := map[string]string{
		"meihua": "shaoyong_meihua",
	}
	if len(cfg.DefaultBindings) != len(expectedBindings) {
		t.Errorf("expected %d default bindings, got %d", len(expectedBindings), len(cfg.DefaultBindings))
	}
	for typ, expectedID := range expectedBindings {
		actualID, ok := cfg.DefaultBindings[typ]
		if !ok {
			t.Errorf("missing binding for type: %s", typ)
			continue
		}
		if actualID != expectedID {
			t.Errorf("binding for %s: expected %s, got %s", typ, expectedID, actualID)
		}
	}
}

func TestLoadProfiles_MeihuaProfileContent(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	cfg, err := LoadProfiles(filepath.Join(cwd, "prompt_profiles.yaml"))
	if err != nil {
		t.Fatalf("LoadProfiles: %v", err)
	}

	profile, err := cfg.DefaultProfile("meihua")
	if err != nil {
		t.Fatalf("DefaultProfile(meihua): %v", err)
	}

	if profile.Name != "康节先生" {
		t.Errorf("expected name 康节先生, got %s", profile.Name)
	}
	if profile.DivinationType != "meihua" {
		t.Errorf("expected divination_type meihua, got %s", profile.DivinationType)
	}
	if !profile.Enabled {
		t.Error("expected profile to be enabled")
	}
	if len(profile.ReasoningFramework) == 0 {
		t.Error("expected non-empty reasoning_framework")
	}
	if len(profile.VoiceStyle) == 0 {
		t.Error("expected non-empty voice_style")
	}
	if len(profile.OutputStructure) == 0 {
		t.Error("expected non-empty output_structure")
	}
	if len(profile.Guardrails) == 0 {
		t.Error("expected non-empty guardrails")
	}
}
