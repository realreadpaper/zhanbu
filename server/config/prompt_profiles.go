package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

//go:embed prompt_profiles.yaml
var embeddedProfilesYAML []byte

// ProfileConfig 定义单个 AI 解读人格。
type ProfileConfig struct {
	Version            string   `mapstructure:"version"`
	DivinationType     string   `mapstructure:"divination_type"`
	Name               string   `mapstructure:"name"`
	Title              string   `mapstructure:"title"`
	Subtitle           string   `mapstructure:"subtitle"`
	Icon               string   `mapstructure:"icon"`
	Description        string   `mapstructure:"description"`
	Enabled            bool     `mapstructure:"enabled"`
	SystemIdentity     string   `mapstructure:"system_identity"`
	ReasoningFramework []string `mapstructure:"reasoning_framework"`
	VoiceStyle         []string `mapstructure:"voice_style"`
	OutputStructure    []string `mapstructure:"output_structure"`
	Guardrails         []string `mapstructure:"guardrails"`
}

// ProfilesConfig 是 prompt_profiles.yaml 的顶层结构。
type ProfilesConfig struct {
	Profiles        map[string]ProfileConfig `mapstructure:"prompt_profiles"`
	DefaultBindings map[string]string        `mapstructure:"default_bindings"`
}

// LoadProfiles 从 YAML 文件或嵌入内容加载所有人格配置。
// 优先尝试文件系统（便于本地开发热编辑），失败时回退到编译时嵌入的 YAML。
func LoadProfiles(configPath string) (*ProfilesConfig, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	var source string
	if configPath != "" {
		v.SetConfigFile(configPath)
		source = fmt.Sprintf("explicit path: %s", configPath)
	} else {
		v.SetConfigName("prompt_profiles")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("../config")
		source = "search paths: ., ./config, ../config"
	}

	log.Info().Str("source", source).Msg("loading prompt profiles")

	if err := v.ReadInConfig(); err != nil {
		log.Warn().Err(err).Str("source", source).Msg("file-based profiles not found, falling back to embedded YAML")

		// 回退：使用编译时嵌入的 YAML
		if len(embeddedProfilesYAML) == 0 {
			return nil, fmt.Errorf("read prompt profiles config: %w (embedded fallback is empty)", err)
		}

		log.Info().Int("bytes", len(embeddedProfilesYAML)).Msg("reading prompt profiles from embedded YAML")
		if err := v.ReadConfig(bytes.NewReader(embeddedProfilesYAML)); err != nil {
			return nil, fmt.Errorf("read embedded prompt profiles: %w", err)
		}
	} else {
		log.Info().Str("file", v.ConfigFileUsed()).Msg("loaded prompt profiles from file")
	}

	var cfg ProfilesConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal prompt profiles: %w", err)
	}

	// 详细日志：已加载的人格与绑定
	profileIDs := make([]string, 0, len(cfg.Profiles))
	for id, p := range cfg.Profiles {
		profileIDs = append(profileIDs, fmt.Sprintf("%s(%s %s enabled=%v)", id, p.Name, p.Version, p.Enabled))
	}
	log.Info().
		Int("profile_count", len(cfg.Profiles)).
		Strs("profiles", profileIDs).
		Int("binding_count", len(cfg.DefaultBindings)).
		Interface("default_bindings", cfg.DefaultBindings).
		Msg("prompt profiles loaded successfully")

	return &cfg, nil
}

// DefaultProfile 获取指定占卜类型的默认人格配置。
// 返回指针便于调用方判断 nil；未找到或已禁用时返回非 nil error。
func (c *ProfilesConfig) DefaultProfile(divinationType string) (*ProfileConfig, error) {
	log.Debug().
		Str("divination_type", divinationType).
		Interface("available_bindings", c.DefaultBindings).
		Int("available_profiles", len(c.Profiles)).
		Msg("resolving default profile")

	profileID, ok := c.DefaultBindings[divinationType]
	if !ok {
		log.Warn().
			Str("divination_type", divinationType).
			Strs("known_types", bindingKeys(c.DefaultBindings)).
			Msg("no default profile binding")
		return nil, fmt.Errorf("no default profile for type: %s", divinationType)
	}

	log.Debug().
		Str("divination_type", divinationType).
		Str("profile_id", profileID).
		Msg("resolved profile binding")

	profile, ok := c.Profiles[profileID]
	if !ok {
		log.Warn().
			Str("profile_id", profileID).
			Strs("known_profile_ids", profileKeys(c.Profiles)).
			Msg("profile not found in registry")
		return nil, fmt.Errorf("profile not found: %s", profileID)
	}

	if !profile.Enabled {
		log.Warn().
			Str("profile_id", profileID).
			Str("profile_name", profile.Name).
			Msg("profile is disabled")
		return nil, fmt.Errorf("profile disabled: %s", profileID)
	}

	log.Debug().
		Str("profile_id", profileID).
		Str("profile_name", profile.Name).
		Str("profile_version", profile.Version).
		Msg("default profile resolved")

	return &profile, nil
}

// Profile 按 ID 获取人格配置。
func (c *ProfilesConfig) Profile(profileID string) (*ProfileConfig, error) {
	profile, ok := c.Profiles[profileID]
	if !ok {
		return nil, fmt.Errorf("profile not found: %s", profileID)
	}
	return &profile, nil
}

// bindingKeys 返回 bindings map 的所有 key，供日志使用。
func bindingKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// profileKeys 返回 profiles map 的所有 key，供日志使用。
func profileKeys(m map[string]ProfileConfig) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// FormatLoadedSummary 返回已加载人格配置的简要摘要，供启动日志使用。
func (c *ProfilesConfig) FormatLoadedSummary() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d profiles, %d bindings", len(c.Profiles), len(c.DefaultBindings))
	for typ, profileID := range c.DefaultBindings {
		if p, ok := c.Profiles[profileID]; ok {
			fmt.Fprintf(&b, "\n  %s → %s (%s)", typ, profileID, p.Name)
		}
	}
	return b.String()
}
