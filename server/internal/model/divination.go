package model

import "time"

// DivinationRecord represents a user's divination record.
type DivinationRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Type      string    `gorm:"type:text;not null;index" json:"type"` // tarot | horoscope | liuyao | bazi
	Question  string    `gorm:"type:text;default:''" json:"question"`
	Result    string    `gorm:"type:text;not null" json:"result"`       // JSON string
	AIReading string    `gorm:"type:text;default:''" json:"ai_reading"` // AI interpretation
	CreatedAt time.Time `json:"created_at"`

	// Association
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName specifies the table name for DivinationRecord.
func (DivinationRecord) TableName() string {
	return "divination_records"
}
