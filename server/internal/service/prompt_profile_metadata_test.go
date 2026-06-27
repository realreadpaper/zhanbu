package service

import (
	"testing"

	"zhanbu/config"
	"zhanbu/internal/model"
)

func TestApplyDefaultPromptProfilePopulatesRecord(t *testing.T) {
	profiles := &config.ProfilesConfig{
		Profiles: map[string]config.ProfileConfig{
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
	record := &model.DivinationRecord{Type: "meihua"}

	ApplyDefaultPromptProfile(record, profiles)

	if record.PromptProfileID != "shaoyong_meihua" {
		t.Fatalf("expected profile id shaoyong_meihua, got %s", record.PromptProfileID)
	}
	if record.PromptProfileName != "康节先生" {
		t.Fatalf("expected profile name 康节先生, got %s", record.PromptProfileName)
	}
	if record.PromptProfileVersion != "v1" {
		t.Fatalf("expected profile version v1, got %s", record.PromptProfileVersion)
	}
}

func TestApplyDefaultPromptProfileSkipsUnboundType(t *testing.T) {
	profiles := &config.ProfilesConfig{
		Profiles: map[string]config.ProfileConfig{
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
	record := &model.DivinationRecord{Type: "liuyao_v2"}

	ApplyDefaultPromptProfile(record, profiles)

	if record.PromptProfileID != "" || record.PromptProfileName != "" || record.PromptProfileVersion != "" {
		t.Fatalf("expected no profile metadata for unbound type, got %#v", record)
	}
}
