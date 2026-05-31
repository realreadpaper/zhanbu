package model

import "time"

// EmailVerification stores email verification codes.
type EmailVerification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"type:text;not null;index" json:"email"`
	Code      string    `gorm:"type:text;not null" json:"code"`
	Purpose   string    `gorm:"type:text;not null" json:"purpose"` // register / reset_password
	Used      bool      `gorm:"default:false" json:"used"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// IsExpired checks if the code has expired.
func (v *EmailVerification) IsExpired() bool {
	return time.Now().After(v.ExpiresAt)
}

// SendLog tracks email sending for rate limiting.
type SendLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"type:text;not null;index" json:"email"`
	Purpose   string    `gorm:"type:text;not null" json:"purpose"`
	SentAt    time.Time `gorm:"not null;default:current_timestamp" json:"sent_at"`
}
