package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"zhanbu/config"
)

var DB *gorm.DB

// Init initializes the PostgreSQL database connection.
func Init(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	logLevel := logger.Silent

	dsn := cfg.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings for PostgreSQL
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	DB = db
	return db, nil
}

// GetDB returns the global database instance.
func GetDB() *gorm.DB {
	return DB
}
