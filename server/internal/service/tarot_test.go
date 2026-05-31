package service

import (
	"testing"

	"zhanbu/internal/model"
)

// mockTarotRepo 用于测试的mock仓库
type mockTarotRepo struct {
	cards []model.TarotCard
}

func (m *mockTarotRepo) FindAll() ([]model.TarotCard, error) {
	return m.cards, nil
}

func (m *mockTarotRepo) FindByID(id uint) (*model.TarotCard, error) {
	for _, c := range m.cards {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, nil
}

func (m *mockTarotRepo) FindByType(t string) ([]model.TarotCard, error) {
	var result []model.TarotCard
	for _, c := range m.cards {
		if c.Type == t {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *mockTarotRepo) FindBySuit(s string) ([]model.TarotCard, error) {
	var result []model.TarotCard
	for _, c := range m.cards {
		if c.Suit == s {
			result = append(result, c)
		}
	}
	return result, nil
}

// generateTestCards 生成78张测试牌
func generateTestCards() []model.TarotCard {
	var cards []model.TarotCard
	id := uint(1)
	// 22张大阿尔卡纳
	for i := 0; i < 22; i++ {
		cards = append(cards, model.TarotCard{ID: id, Name: "牌", NameEn: "Card", Type: "major", Number: i})
		id++
	}
	// 56张小阿尔卡纳
	suits := []string{"wands", "cups", "swords", "pentacles"}
	for _, s := range suits {
		for i := 1; i <= 14; i++ {
			cards = append(cards, model.TarotCard{ID: id, Name: "牌", NameEn: "Card", Type: "minor", Suit: s, Number: i})
			id++
		}
	}
	return cards
}

func TestFisherYatesShuffle(t *testing.T) {
	cards := generateTestCards()
	repo := &mockTarotRepo{cards: cards}
	svc := NewTarotService(repo)

	shuffled := svc.FisherYatesShuffle(cards)

	if len(shuffled) != len(cards) {
		t.Fatalf("expected %d cards, got %d", len(cards), len(shuffled))
	}

	originalIDs := make(map[uint]bool)
	for _, c := range cards {
		originalIDs[c.ID] = true
	}
	for _, c := range shuffled {
		if !originalIDs[c.ID] {
			t.Errorf("shuffled contains unexpected card ID %d", c.ID)
		}
	}
}

func TestDrawSingle(t *testing.T) {
	cards := generateTestCards()
	repo := &mockTarotRepo{cards: cards}
	svc := NewTarotService(repo)

	result, err := svc.DrawCards("single", "测试问题")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Cards) != 1 {
		t.Fatalf("expected 1 card, got %d", len(result.Cards))
	}
	if result.Cards[0].PositionName != "核心信息" {
		t.Errorf("expected position '核心信息', got '%s'", result.Cards[0].PositionName)
	}
	if result.Cards[0].Orientation != "upright" && result.Cards[0].Orientation != "reversed" {
		t.Errorf("unexpected orientation: %s", result.Cards[0].Orientation)
	}
}

func TestDrawThree(t *testing.T) {
	cards := generateTestCards()
	repo := &mockTarotRepo{cards: cards}
	svc := NewTarotService(repo)

	result, err := svc.DrawCards("three", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Cards) != 3 {
		t.Fatalf("expected 3 cards, got %d", len(result.Cards))
	}
	expected := []string{"过去", "现在", "未来"}
	for i, name := range expected {
		if result.Cards[i].PositionName != name {
			t.Errorf("card %d: expected '%s', got '%s'", i, name, result.Cards[i].PositionName)
		}
	}
}

func TestDrawCeltic(t *testing.T) {
	cards := generateTestCards()
	repo := &mockTarotRepo{cards: cards}
	svc := NewTarotService(repo)

	result, err := svc.DrawCards("celtic", "事业如何")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Cards) != 10 {
		t.Fatalf("expected 10 cards, got %d", len(result.Cards))
	}
	if result.Cards[0].PositionName != "当前状况" {
		t.Errorf("expected '当前状况', got '%s'", result.Cards[0].PositionName)
	}
	if result.Cards[9].PositionName != "最终结果" {
		t.Errorf("expected '最终结果', got '%s'", result.Cards[9].PositionName)
	}
}

func TestDrawLove(t *testing.T) {
	cards := generateTestCards()
	repo := &mockTarotRepo{cards: cards}
	svc := NewTarotService(repo)

	result, err := svc.DrawCards("love", "感情")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Cards) != 5 {
		t.Fatalf("expected 5 cards, got %d", len(result.Cards))
	}
	expected := []string{"你的现状", "对方的现状", "关系的挑战", "你的期望", "关系的走向"}
	for i, name := range expected {
		if result.Cards[i].PositionName != name {
			t.Errorf("card %d: expected '%s', got '%s'", i, name, result.Cards[i].PositionName)
		}
	}
}

func TestUprightReversedRatio(t *testing.T) {
	cards := generateTestCards()
	repo := &mockTarotRepo{cards: cards}
	svc := NewTarotService(repo)

	upright := 0
	total := 1000
	for i := 0; i < total; i++ {
		result, _ := svc.DrawCards("single", "")
		if result.Cards[0].Orientation == "upright" {
			upright++
		}
	}

	ratio := float64(upright) / float64(total)
	if ratio < 0.4 || ratio > 0.6 {
		t.Errorf("upright ratio %.2f is not close to 0.5 (expected 0.4-0.6)", ratio)
	}
}

func TestDrawInvalidSpread(t *testing.T) {
	cards := generateTestCards()
	repo := &mockTarotRepo{cards: cards}
	svc := NewTarotService(repo)

	_, err := svc.DrawCards("invalid", "")
	if err == nil {
		t.Error("expected error for invalid spread, got nil")
	}
}
