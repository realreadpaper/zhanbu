package repository

import (
	"zhanbu/internal/model"

	"gorm.io/gorm"
)

// ChatRepository handles chat session and message database operations.
type ChatRepository struct {
	db *gorm.DB
}

// NewChatRepository creates a new ChatRepository.
func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// CreateSession creates a new chat session.
func (r *ChatRepository) CreateSession(session *model.ChatSession) error {
	return r.db.Create(session).Error
}

// FindSessionByID finds a chat session by ID with messages.
func (r *ChatRepository) FindSessionByID(id uint) (*model.ChatSession, error) {
	var session model.ChatSession
	err := r.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("chat_messages.created_at ASC")
	}).First(&session, id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// FindSessionByUserAndID finds a chat session by user ID and session ID.
func (r *ChatRepository) FindSessionByUserAndID(userID uint, sessionID uint) (*model.ChatSession, error) {
	var session model.ChatSession
	err := r.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("chat_messages.created_at ASC")
	}).Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// FindSessionByUserAndRecord finds a chat session by user ID and record ID.
func (r *ChatRepository) FindSessionByUserAndRecord(userID uint, recordID uint) (*model.ChatSession, error) {
	var session model.ChatSession
	err := r.db.Where("user_id = ? AND record_id = ?", userID, recordID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// ListSessionsByUser lists chat sessions for a user with pagination.
func (r *ChatRepository) ListSessionsByUser(userID uint, page, size int) ([]model.ChatSession, int64, error) {
	var sessions []model.ChatSession
	var total int64

	// Count total
	if err := r.db.Model(&model.ChatSession{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with last message
	offset := (page - 1) * size
	err := r.db.Where("user_id = ?", userID).
		Order("updated_at DESC").
		Offset(offset).
		Limit(size).
		Find(&sessions).Error
	if err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

// UpdateSessionTitle updates a chat session's title.
func (r *ChatRepository) UpdateSessionTitle(id uint, title string) error {
	return r.db.Model(&model.ChatSession{}).Where("id = ?", id).Update("title", title).Error
}

// UpdateSessionTimestamp updates a chat session's updated_at timestamp.
func (r *ChatRepository) UpdateSessionTimestamp(id uint) error {
	return r.db.Model(&model.ChatSession{}).Where("id = ?", id).Update("updated_at", gorm.Expr("NOW()")).Error
}

// DeleteSession deletes a chat session and its messages.
func (r *ChatRepository) DeleteSession(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete messages first
		if err := tx.Where("session_id = ?", id).Delete(&model.ChatMessage{}).Error; err != nil {
			return err
		}
		// Delete session
		return tx.Delete(&model.ChatSession{}, id).Error
	})
}

// CreateMessage creates a new chat message.
func (r *ChatRepository) CreateMessage(message *model.ChatMessage) error {
	return r.db.Create(message).Error
}

// GetRecentMessages gets recent messages for a session.
func (r *ChatRepository) GetRecentMessages(sessionID uint, limit int) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage
	err := r.db.Where("session_id = ?", sessionID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	// Reverse to chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// UpdateSessionRecordID updates a chat session's record ID（重新占卜时切换记录）.
func (r *ChatRepository) UpdateSessionRecordID(sessionID uint, recordID uint) error {
	return r.db.Model(&model.ChatSession{}).Where("id = ?", sessionID).Update("record_id", recordID).Error
}

// CountMessagesBySession counts messages in a session.
func (r *ChatRepository) CountMessagesBySession(sessionID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.ChatMessage{}).Where("session_id = ?", sessionID).Count(&count).Error
	return count, err
}

// CountSessionsByUser counts sessions for a user.
func (r *ChatRepository) CountSessionsByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.ChatSession{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
