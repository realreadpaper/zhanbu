package repository

import (
	"time"

	"zhanbu/internal/model"

	"gorm.io/gorm"
)

// VerificationRepository handles database operations for email verifications.
type VerificationRepository struct {
	db *gorm.DB
}

// NewVerificationRepository creates a new VerificationRepository.
func NewVerificationRepository(db *gorm.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

// Create inserts a new verification record.
func (r *VerificationRepository) Create(v *model.EmailVerification) error {
	return r.db.Create(v).Error
}

// FindLatest finds the latest verification for email+purpose.
func (r *VerificationRepository) FindLatest(email, purpose string) (*model.EmailVerification, error) {
	var v model.EmailVerification
	err := r.db.Where("email = ? AND purpose = ?", email, purpose).
		Order("created_at DESC").
		First(&v).Error
	return &v, err
}

// MarkUsed marks a verification code as used.
func (r *VerificationRepository) MarkUsed(id uint) error {
	return r.db.Model(&model.EmailVerification{}).Where("id = ?", id).Update("used", true).Error
}

// CreateSendLog inserts a send log for rate limiting.
func (r *VerificationRepository) CreateSendLog(log *model.SendLog) error {
	return r.db.Create(log).Error
}

// CountRecentSends counts how many emails were sent recently for rate limiting.
func (r *VerificationRepository) CountRecentSends(email, purpose string, within time.Duration) (int64, error) {
	var count int64
	since := time.Now().Add(-within)
	err := r.db.Model(&model.SendLog{}).
		Where("email = ? AND purpose = ? AND sent_at > ?", email, purpose, since).
		Count(&count).Error
	return count, err
}

// VerificationRepoReader interface for testing.
type VerificationRepoReader interface {
	Create(v *model.EmailVerification) error
	FindLatest(email, purpose string) (*model.EmailVerification, error)
	MarkUsed(id uint) error
	CreateSendLog(log *model.SendLog) error
	CountRecentSends(email, purpose string, within time.Duration) (int64, error)
}
