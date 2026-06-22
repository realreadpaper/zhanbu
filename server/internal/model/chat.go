package model

import "time"

// ChatSession represents a chat session for a divination record.
type ChatSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	RecordID  uint      `gorm:"not null;index" json:"record_id"`
	Title     string    `gorm:"type:varchar(200);default:''" json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Associations
	User               User               `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	DivinationRecord   DivinationRecord   `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"-"`
	Messages           []ChatMessage      `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE" json:"messages,omitempty"`
}

// TableName specifies the table name for ChatSession.
func (ChatSession) TableName() string {
	return "chat_sessions"
}

// ChatMessage represents a message in a chat session.
type ChatMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID uint      `gorm:"not null;index" json:"session_id"`
	Role      string    `gorm:"type:varchar(20);not null" json:"role"` // 'user' | 'assistant'
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the table name for ChatMessage.
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// ChatSessionResponse is the API response for a chat session.
type ChatSessionResponse struct {
	ID        uint                  `json:"id"`
	RecordID  uint                  `json:"record_id"`
	Title     string                `json:"title"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	Messages  []ChatMessageResponse `json:"messages,omitempty"`
}

// ChatMessageResponse is the API response for a chat message.
type ChatMessageResponse struct {
	ID        uint      `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a ChatSession to ChatSessionResponse.
func (s *ChatSession) ToResponse() ChatSessionResponse {
	resp := ChatSessionResponse{
		ID:        s.ID,
		RecordID:  s.RecordID,
		Title:     s.Title,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
	if len(s.Messages) > 0 {
		resp.Messages = make([]ChatMessageResponse, len(s.Messages))
		for i, msg := range s.Messages {
			resp.Messages[i] = msg.ToResponse()
		}
	}
	return resp
}

// ToResponse converts a ChatMessage to ChatMessageResponse.
func (m *ChatMessage) ToResponse() ChatMessageResponse {
	return ChatMessageResponse{
		ID:        m.ID,
		Role:      m.Role,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
	}
}
