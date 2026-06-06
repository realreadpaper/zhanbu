package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application.
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	AI        AIConfig        `mapstructure:"ai"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	CORS      CORSConfig      `mapstructure:"cors"`
	SMTP      SMTPConfig      `mapstructure:"smtp"`
	Security  SecurityConfig  `mapstructure:"security"`
	LiuYao    LiuYaoConfig    `mapstructure:"liuyao"`
}

// SMTPConfig holds email configuration.
type SMTPConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
	SSL      bool   `mapstructure:"ssl"`
}

// SecurityConfig holds security-related config.
type SecurityConfig struct {
	VerifyEmail    bool          `mapstructure:"verify_email"`
	CodeLength     int           `mapstructure:"code_length"`
	CodeExpiry     time.Duration `mapstructure:"code_expiry"`
	MaxSendPerHour int           `mapstructure:"max_send_per_hour"`
}

// LiuYaoConfig holds LiuYao divination configuration.
type LiuYaoConfig struct {
	Version string `mapstructure:"version"` // v1=traditional, v2=takashima
	Method  string `mapstructure:"method"`  // yarrow=蓍草法, coin=铜钱法, both=两种都支持
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig holds database-related configuration.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	Timezone string `mapstructure:"timezone"`
}

// DSN returns the PostgreSQL connection string.
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		d.Host, d.User, d.Password, d.DBName, d.Port, d.SSLMode, d.Timezone,
	)
}

// JWTConfig holds JWT-related configuration.
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	AccessTTL  time.Duration `mapstructure:"access_ttl"`
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
}

// AIConfig holds AI-related configuration.
type AIConfig struct {
	Provider    string  `mapstructure:"provider"`
	APIKey      string  `mapstructure:"api_key"`
	Model       string  `mapstructure:"model"`
	BaseURL     string  `mapstructure:"base_url"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float64 `mapstructure:"temperature"`
}

// RateLimitConfig holds rate limit configuration.
type RateLimitConfig struct {
	AIPerMinute  int `mapstructure:"ai_per_minute"`
	APIPerMinute int `mapstructure:"api_per_minute"`
}

// CORSConfig holds CORS configuration.
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

// Load reads configuration from file and environment variables.
func Load(configPath string) (*Config, error) {
	v := viper.New()
	setDefaults(v)

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("../config")
	}

	v.SetEnvPrefix("ZHANBU")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// LoadFromEnv creates a Config from environment variables only.
func LoadFromEnv() *Config {
	cfg := &Config{
		Server:    ServerConfig{Port: 18080, Mode: "debug"},
		Database:  DatabaseConfig{Host: "localhost", Port: 15432, User: "zhanbu", Password: "zhanbu_secret", DBName: "zhanbu", SSLMode: "disable", Timezone: "Asia/Shanghai"},
		JWT:       JWTConfig{Secret: "dev-secret-key-change-in-production", AccessTTL: time.Hour, RefreshTTL: 7 * 24 * time.Hour},
		AI:        AIConfig{Provider: "openai", APIKey: "", Model: "gpt-4", BaseURL: "https://api.openai.com/v1", MaxTokens: 3000, Temperature: 0.7},
		RateLimit: RateLimitConfig{AIPerMinute: 5, APIPerMinute: 60},
		CORS:      CORSConfig{AllowedOrigins: []string{"http://localhost:5173"}},
		SMTP:      SMTPConfig{Enabled: false, Host: "smtp.qq.com", Port: 465, SSL: true},
		Security:  SecurityConfig{VerifyEmail: false, CodeLength: 6, CodeExpiry: 10 * time.Minute, MaxSendPerHour: 5},
		LiuYao:    LiuYaoConfig{Version: "v1", Method: "both"},
	}

	if v := os.Getenv("PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("ZHANBU_SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("ZHANBU_SERVER_MODE"); v != "" {
		cfg.Server.Mode = v
	}
	if v := os.Getenv("ZHANBU_DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("ZHANBU_DB_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Database.Port)
	}
	if v := os.Getenv("ZHANBU_DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("ZHANBU_DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("ZHANBU_DB_NAME"); v != "" {
		cfg.Database.DBName = v
	}
	if v := os.Getenv("ZHANBU_DB_SSLMODE"); v != "" {
		cfg.Database.SSLMode = v
	}
	if v := os.Getenv("ZHANBU_JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("ZHANBU_SMTP_HOST"); v != "" {
		cfg.SMTP.Enabled = true
		cfg.SMTP.Host = v
	}
	if v := os.Getenv("ZHANBU_SMTP_ENABLED"); v != "" {
		cfg.SMTP.Enabled = v == "true" || v == "1"
	}
	if v := os.Getenv("ZHANBU_SMTP_USERNAME"); v != "" {
		cfg.SMTP.Username = v
	}
	if v := os.Getenv("ZHANBU_SMTP_PASSWORD"); v != "" {
		cfg.SMTP.Password = v
	}
	if v := os.Getenv("ZHANBU_SMTP_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.SMTP.Port)
	}
	if v := os.Getenv("ZHANBU_SMTP_FROM"); v != "" {
		cfg.SMTP.From = v
	}
	if v := os.Getenv("ZHANBU_SECURITY_VERIFY_EMAIL"); v == "true" || v == "1" {
		cfg.Security.VerifyEmail = true
	}
	if v := os.Getenv("ZHANBU_AI_PROVIDER"); v != "" {
		cfg.AI.Provider = v
	}
	if v := os.Getenv("ZHANBU_AI_API_KEY"); v != "" {
		cfg.AI.APIKey = v
	}
	if v := os.Getenv("ZHANBU_AI_MODEL"); v != "" {
		cfg.AI.Model = v
	}
	if v := os.Getenv("ZHANBU_AI_BASE_URL"); v != "" {
		cfg.AI.BaseURL = v
	}
	if v := os.Getenv("ZHANBU_AI_MAX_TOKENS"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.AI.MaxTokens)
	}
	if v := os.Getenv("ZHANBU_AI_TEMPERATURE"); v != "" {
		fmt.Sscanf(v, "%f", &cfg.AI.Temperature)
	}
	if v := os.Getenv("ZHANBU_CORS_ALLOWED_ORIGINS"); v != "" {
		cfg.CORS.AllowedOrigins = splitCommaSeparated(v)
	}
	if v := os.Getenv("ZHANBU_LIUYAO_VERSION"); v != "" {
		cfg.LiuYao.Version = v
	}
	if v := os.Getenv("ZHANBU_LIUYAO_METHOD"); v != "" {
		cfg.LiuYao.Method = v
	}

	return cfg
}

func splitCommaSeparated(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}

// GetDataDir returns the data directory path for static data files.
func GetDataDir(cfg *Config) (string, error) {
	dir := "data"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create data directory: %w", err)
		}
	}
	return dir, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 18080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 15432)
	v.SetDefault("database.user", "zhanbu")
	v.SetDefault("database.password", "zhanbu_secret")
	v.SetDefault("database.dbname", "zhanbu")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.timezone", "Asia/Shanghai")
	v.SetDefault("jwt.secret", "dev-secret-key-change-in-production")
	v.SetDefault("jwt.access_ttl", "1h")
	v.SetDefault("jwt.refresh_ttl", "168h")
	v.SetDefault("ai.provider", "openai")
	v.SetDefault("ai.api_key", "")
	v.SetDefault("ai.model", "gpt-4")
	v.SetDefault("ai.base_url", "https://api.openai.com/v1")
	v.SetDefault("ai.max_tokens", 3000)
	v.SetDefault("ai.temperature", 0.7)
	v.SetDefault("rate_limit.ai_per_minute", 5)
	v.SetDefault("rate_limit.api_per_minute", 60)
	v.SetDefault("cors.allowed_origins", []string{"http://localhost:5173"})
	v.SetDefault("smtp.enabled", false)
	v.SetDefault("smtp.host", "smtp.qq.com")
	v.SetDefault("smtp.port", 465)
	v.SetDefault("smtp.ssl", true)
	v.SetDefault("security.verify_email", false)
	v.SetDefault("security.code_length", 6)
	v.SetDefault("security.code_expiry", "10m")
	v.SetDefault("security.max_send_per_hour", 5)
	v.SetDefault("liuyao.version", "v1")
	v.SetDefault("liuyao.method", "both")
}
