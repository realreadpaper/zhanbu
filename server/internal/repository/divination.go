package repository

import (
	"zhanbu/internal/model"

	"gorm.io/gorm"
)

// DivinationRepository handles database operations for divination records.
type DivinationRepository struct {
	db *gorm.DB
}

// NewDivinationRepository creates a new DivinationRepository.
func NewDivinationRepository(db *gorm.DB) *DivinationRepository {
	return &DivinationRepository{db: db}
}

// Create inserts a new divination record.
func (r *DivinationRepository) Create(record *model.DivinationRecord) error {
	return r.db.Create(record).Error
}

// FindByID finds a divination record by ID.
func (r *DivinationRepository) FindByID(id uint) (*model.DivinationRecord, error) {
	var record model.DivinationRecord
	err := r.db.First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindByUserIDAndID finds a record by ID, scoped to a user.
func (r *DivinationRepository) FindByUserIDAndID(userID uint, id uint) (*model.DivinationRecord, error) {
	var record model.DivinationRecord
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// ListByUserID returns paginated records for a user, optionally filtered by type.
func (r *DivinationRepository) ListByUserID(userID uint, divinationType string, page, pageSize int) ([]model.DivinationRecord, int64, error) {
	var records []model.DivinationRecord
	var total int64

	query := r.db.Model(&model.DivinationRecord{}).Where("user_id = ?", userID)
	if divinationType != "" {
		query = query.Where("type = ?", divinationType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// DeleteByUserIDAndID deletes a record by ID, scoped to a user.
// Returns an error if the record doesn't exist or doesn't belong to the user.
func (r *DivinationRepository) DeleteByUserIDAndID(userID uint, id uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.DivinationRecord{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpdateAIReading updates the AI reading for a divination record.
func (r *DivinationRepository) UpdateAIReading(id uint, reading string) error {
	return r.db.Model(&model.DivinationRecord{}).Where("id = ?", id).Update("ai_reading", reading).Error
}
