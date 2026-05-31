package database

import (
	"gorm.io/gorm"

	"zhanbu/internal/model"
)

// Migrate runs auto-migration for all models.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.DivinationRecord{},
		&model.TarotCard{},
		&model.EmailVerification{},
		&model.SendLog{},
	)
}
