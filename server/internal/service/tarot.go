package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"zhanbu/internal/model"
)

// TarotCardReader 塔罗牌数据读取接口（便于测试mock）
type TarotCardReader interface {
	FindAll() ([]model.TarotCard, error)
	FindByID(id uint) (*model.TarotCard, error)
	FindByType(cardType string) ([]model.TarotCard, error)
	FindBySuit(suit string) ([]model.TarotCard, error)
}

// DivinationRecordSaver 保存占卜记录的接口（便于测试mock）
type DivinationRecordSaver interface {
	Create(record *model.DivinationRecord) error
}

// TarotService 塔罗牌服务
type TarotService struct {
	repo      TarotCardReader
	saver     DivinationRecordSaver
	rng       *rand.Rand
}

// NewTarotService 创建塔罗牌服务
func NewTarotService(repo TarotCardReader) *TarotService {
	return &TarotService{
		repo: repo,
		rng:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SetRecordSaver 设置记录保存器
func (s *TarotService) SetRecordSaver(saver DivinationRecordSaver) {
	s.saver = saver
}

// GetSpreads 获取所有可用牌阵
func (s *TarotService) GetSpreads() []model.Spread {
	return []model.Spread{
		{
			ID:    "single",
			Name:  "单牌抽取",
			Count: 1,
			Positions: []model.SpreadPosition{
				{Index: 0, Name: "核心信息", Description: "代表当前核心信息和指引"},
			},
		},
		{
			ID:    "three",
			Name:  "三牌阵",
			Count: 3,
			Positions: []model.SpreadPosition{
				{Index: 0, Name: "过去", Description: "过去的影响力"},
				{Index: 1, Name: "现在", Description: "当前状况"},
				{Index: 2, Name: "未来", Description: "未来趋势"},
			},
		},
		{
			ID:    "celtic",
			Name:  "凯尔特十字阵",
			Count: 10,
			Positions: []model.SpreadPosition{
				{Index: 0, Name: "当前状况", Description: "核心问题或当前状况"},
				{Index: 1, Name: "挑战", Description: "面临的挑战或障碍"},
				{Index: 2, Name: "潜意识", Description: "潜意识或根基"},
				{Index: 3, Name: "过去", Description: "过去的影响力"},
				{Index: 4, Name: "意识", Description: "意识或目标"},
				{Index: 5, Name: "近期未来", Description: "近期未来趋势"},
				{Index: 6, Name: "自我认知", Description: "自我认知"},
				{Index: 7, Name: "外部环境", Description: "外部环境影响"},
				{Index: 8, Name: "希望与恐惧", Description: "内心深处的希望与恐惧"},
				{Index: 9, Name: "最终结果", Description: "最终结果"},
			},
		},
		{
			ID:    "love",
			Name:  "爱情十字阵",
			Count: 5,
			Positions: []model.SpreadPosition{
				{Index: 0, Name: "你的现状", Description: "你当前的感情状态"},
				{Index: 1, Name: "对方的现状", Description: "对方当前的状态"},
				{Index: 2, Name: "关系的挑战", Description: "关系中面临的挑战"},
				{Index: 3, Name: "你的期望", Description: "你对这段关系的期望"},
				{Index: 4, Name: "关系的走向", Description: "关系未来的发展方向"},
			},
		},
	}
}

// GetSpreadByID 根据ID获取牌阵
func (s *TarotService) GetSpreadByID(id string) (*model.Spread, error) {
	for _, spread := range s.GetSpreads() {
		if spread.ID == id {
			return &spread, nil
		}
	}
	return nil, fmt.Errorf("未知的牌阵类型: %s", id)
}

// FisherYatesShuffle Fisher-Yates 洗牌算法
func (s *TarotService) FisherYatesShuffle(cards []model.TarotCard) []model.TarotCard {
	shuffled := make([]model.TarotCard, len(cards))
	copy(shuffled, cards)
	for i := len(shuffled) - 1; i > 0; i-- {
		j := s.rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}

// GetAllCards 获取所有牌
func (s *TarotService) GetAllCards() ([]model.TarotCard, error) {
	return s.repo.FindAll()
}

// GetCardByID 根据ID获取牌
func (s *TarotService) GetCardByID(id uint) (*model.TarotCard, error) {
	return s.repo.FindByID(id)
}

// DrawCards 抽牌。userID > 0 时自动保存记录。
func (s *TarotService) DrawCards(spreadID string, question string, userID ...uint) (*model.DrawResult, error) {
	spread, err := s.GetSpreadByID(spreadID)
	if err != nil {
		return nil, err
	}

	cards, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("获取牌组失败: %w", err)
	}

	if len(cards) == 0 {
		return nil, fmt.Errorf("牌组为空，请先导入种子数据")
	}

	shuffled := s.FisherYatesShuffle(cards)

	drawnCards := make([]model.DrawnCard, spread.Count)
	for i := 0; i < spread.Count; i++ {
		orientation := "upright"
		if s.rng.Intn(2) == 0 {
			orientation = "reversed"
		}
		drawnCards[i] = model.DrawnCard{
			Position:     i + 1,
			PositionName: spread.Positions[i].Name,
			Card:         shuffled[i],
			Orientation:  orientation,
		}
	}

	result := &model.DrawResult{
		Spread:    spreadID,
		Question:  question,
		Cards:     drawnCards,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 自动保存记录（如果有saver且传入了userID）
	if len(userID) > 0 && userID[0] > 0 && s.saver != nil {
		resultJSON, _ := json.Marshal(result)
		record := &model.DivinationRecord{
			UserID:   userID[0],
			Type:     "tarot",
			Question: question,
			Result:   string(resultJSON),
		}
		if err := s.saver.Create(record); err == nil {
			result.RecordID = record.ID
		}
	}

	return result, nil
}
