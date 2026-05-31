package model

import "time"

// User represents a user in the system.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"type:text;not null;uniqueIndex" json:"username"`
	Email        string    `gorm:"type:text;not null;uniqueIndex" json:"email"`
	PasswordHash  string    `gorm:"type:text;not null" json:"-"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	Avatar        string    `gorm:"type:text;default:''" json:"avatar"`
	Zodiac       string    `gorm:"type:text;default:''" json:"zodiac"`
	BirthDate    string    `gorm:"type:text;default:''" json:"birth_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserResponse is the public user info returned in API responses.
type UserResponse struct {
	ID            uint      `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Avatar        string    `json:"avatar"`
	Zodiac        string    `json:"zodiac"`
	BirthDate     string    `json:"birth_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ToResponse converts a User to a UserResponse (hides password hash).
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		Username:      u.Username,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		Avatar:        u.Avatar,
		Zodiac:        u.Zodiac,
		BirthDate:     u.BirthDate,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
