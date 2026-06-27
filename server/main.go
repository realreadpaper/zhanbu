package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"zhanbu/config"
	"zhanbu/internal/database"
	"zhanbu/internal/router"
)

func main() {
	// Initialize logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		logger.Warn().Err(err).Msg("failed to load config file, using defaults")
	}
	// Always apply environment variable overrides (Load uses Viper naming conventions
	// that may not match the ZHANBU_DB_* env vars set in docker-compose)
	cfg = config.LoadFromEnv()

	// Production environment check: reject default JWT secret
	if cfg.Server.Mode == "release" && cfg.JWT.Secret == "dev-secret-key-change-in-production" {
		logger.Fatal().Msg("refusing to start in production mode with default JWT secret. Set ZHANBU_JWT_SECRET environment variable.")
	}

	// Set log level
	if cfg.Server.Mode == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Ensure data directory exists
	if _, err := config.GetDataDir(cfg); err != nil {
		logger.Fatal().Err(err).Msg("failed to create data directory")
	}

	// Initialize database
	db, err := database.Init(&cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize database")
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		logger.Fatal().Err(err).Msg("failed to run database migrations")
	}
	logger.Info().Msg("database migrations completed")

	// Seed data
	if err := database.SeedData(db); err != nil {
		logger.Fatal().Err(err).Msg("failed to seed data")
	}
	logger.Info().Msg("seed data completed")

	// Load prompt profiles
	logger.Info().Msg("loading prompt profiles...")
	profiles, err := config.LoadProfiles("")
	if err != nil {
		logger.Warn().Err(err).Msg("failed to load prompt profiles, using empty config")
		profiles = &config.ProfilesConfig{
			Profiles:        make(map[string]config.ProfileConfig),
			DefaultBindings: make(map[string]string),
		}
	}
	logger.Info().Msgf("prompt profiles: %s", profiles.FormatLoadedSummary())

	// Setup router
	r := router.SetupRouter(db, cfg, logger, profiles)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info().Str("addr", addr).Msg("starting server")
	if err := r.Run(addr); err != nil {
		logger.Fatal().Err(err).Msg("failed to start server")
	}
}
