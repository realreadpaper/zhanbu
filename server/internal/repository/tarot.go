package repository

import (
	"zhanbu/internal/model"

	"gorm.io/gorm"
)

// TarotRepository 塔罗牌数据访问层
type TarotRepository struct {
	db *gorm.DB
}

// NewTarotRepository 创建塔罗牌Repository
func NewTarotRepository(db *gorm.DB) *TarotRepository {
	return &TarotRepository{db: db}
}

// FindByID 根据ID查找牌
func (r *TarotRepository) FindByID(id uint) (*model.TarotCard, error) {
	var card model.TarotCard
	err := r.db.First(&card, id).Error
	return &card, err
}

// FindAll 查找所有牌
func (r *TarotRepository) FindAll() ([]model.TarotCard, error) {
	var cards []model.TarotCard
	err := r.db.Order("id").Find(&cards).Error
	return cards, err
}

// FindByType 根据类型查找牌
func (r *TarotRepository) FindByType(cardType string) ([]model.TarotCard, error) {
	var cards []model.TarotCard
	err := r.db.Where("type = ?", cardType).Order("id").Find(&cards).Error
	return cards, err
}

// FindBySuit 根据花色查找牌
func (r *TarotRepository) FindBySuit(suit string) ([]model.TarotCard, error) {
	var cards []model.TarotCard
	err := r.db.Where("suit = ?", suit).Order("id").Find(&cards).Error
	return cards, err
}

// Count 返回牌总数
func (r *TarotRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.TarotCard{}).Count(&count).Error
	return count, err
}
