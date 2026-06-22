package service

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"

	"zhanbu/internal/model"
	"zhanbu/internal/repository"
	apperrors "zhanbu/pkg/errors"
)

// DivinationRecordReader defines the read operations for divination records (for testing).
type DivinationRecordReader interface {
	Create(record *model.DivinationRecord) error
	FindByUserIDAndID(userID uint, id uint) (*model.DivinationRecord, error)
	ListByUserID(userID uint, divinationType string, page, pageSize int) ([]model.DivinationRecord, int64, error)
	DeleteByUserIDAndID(userID uint, id uint) error
	UpdateAIReading(id uint, reading string) error
}

// HistoryService handles history business logic.
type HistoryService struct {
	repo DivinationRecordReader
}

// NewHistoryService creates a new HistoryService.
func NewHistoryService(repo *repository.DivinationRepository) *HistoryService {
	return &HistoryService{repo: repo}
}

// NewHistoryServiceWithReader creates a HistoryService with a custom reader (for testing).
func NewHistoryServiceWithReader(repo DivinationRecordReader) *HistoryService {
	return &HistoryService{repo: repo}
}

// HistoryItem is a single history list item.
type HistoryItem struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	TypeCN    string `json:"type_cn"`
	Question  string `json:"question"`
	Summary   string `json:"summary"`
	CreatedAt string `json:"created_at"`
}

// HistoryListResponse is the paginated history response.
type HistoryListResponse struct {
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Items    []HistoryItem `json:"items"`
}

// typeCNMap maps divination types to Chinese names.
var typeCNMap = map[string]string{
	"tarot":     "塔罗牌",
	"horoscope": "星座运势",
	"liuyao":    "六爻",
	"liuyao_v2": "高岛易断",
	"bazi":      "八字",
}

// buildSummary creates a short summary from the result JSON.
func buildSummary(result string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(result), &m); err != nil {
		if len(result) > 100 {
			return result[:100] + "..."
		}
		return result
	}

	if cards, ok := m["cards"].([]interface{}); ok && len(cards) > 0 {
		var names []string
		for _, c := range cards {
			if cardMap, ok := c.(map[string]interface{}); ok {
				if card, ok := cardMap["card"].(map[string]interface{}); ok {
					if name, ok := card["name"].(string); ok {
						names = append(names, name)
					}
				}
			}
		}
		if len(names) > 0 {
			summary := ""
			for i, n := range names {
				if i > 0 {
					summary += ", "
				}
				summary += n
				if i >= 2 {
					summary += "..."
					break
				}
			}
			return summary
		}
	}

	if len(result) > 100 {
		return result[:100] + "..."
	}
	return result
}

// List returns paginated history for a user.
func (s *HistoryService) List(userID uint, divinationType string, page, pageSize int) (*HistoryListResponse, *apperrors.AppError) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	records, total, err := s.repo.ListByUserID(userID, divinationType, page, pageSize)
	if err != nil {
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "failed to query history", err)
	}

	items := make([]HistoryItem, len(records))
	for i, r := range records {
		typeCN := typeCNMap[r.Type]
		if typeCN == "" {
			typeCN = r.Type
		}
		summary := buildSummary(r.Result)
		if r.AIReading != "" {
			summary = r.AIReading
			if len(summary) > 100 {
				summary = summary[:100] + "..."
			}
		}
		items[i] = HistoryItem{
			ID:        r.ID,
			Type:      r.Type,
			TypeCN:    typeCN,
			Question:  r.Question,
			Summary:   summary,
			CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &HistoryListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    items,
	}, nil
}

// GetDetail returns the full record detail for a user.
func (s *HistoryService) GetDetail(userID uint, recordID uint) (*model.DivinationRecord, *apperrors.AppError) {
	record, err := s.repo.FindByUserIDAndID(userID, recordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(apperrors.ErrNotFound, "record not found")
		}
		return nil, apperrors.NewWithErr(apperrors.ErrInternalServer, "database error", err)
	}
	return record, nil
}

// Delete deletes a record for a user.
func (s *HistoryService) Delete(userID uint, recordID uint) *apperrors.AppError {
	err := s.repo.DeleteByUserIDAndID(userID, recordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New(apperrors.ErrNotFound, "record not found")
		}
		return apperrors.NewWithErr(apperrors.ErrInternalServer, "failed to delete record", err)
	}
	return nil
}
