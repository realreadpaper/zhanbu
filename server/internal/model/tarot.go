package model

// TarotCard 塔罗牌模型
type TarotCard struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"not null" json:"name"`
	NameEn       string `gorm:"not null" json:"name_en"`
	Type         string `gorm:"not null" json:"type"` // major/minor
	Suit         string `json:"suit,omitempty"`        // wands/cups/swords/pentacles
	Number       int    `json:"number"`
	Image        string `json:"image"`
	KeywordsUp   string `json:"keywords_up"`   // JSON数组
	KeywordsDown string `json:"keywords_down"` // JSON数组
	MeaningUp    string `json:"meaning_up"`
	MeaningDown  string `json:"meaning_down"`
	Description  string `json:"description"`
}

// Spread 牌阵定义
type Spread struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Count     int            `json:"count"`
	Positions []SpreadPosition `json:"positions"`
}

// SpreadPosition 牌阵位置
type SpreadPosition struct {
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DrawResult 抽牌结果
type DrawResult struct {
	Spread    string       `json:"spread"`
	Question  string       `json:"question,omitempty"`
	Cards     []DrawnCard  `json:"cards"`
	Timestamp string       `json:"timestamp"`
	RecordID  uint         `json:"record_id,omitempty"` // 关联的记录ID（已登录用户）
}

// DrawnCard 已抽取的牌
type DrawnCard struct {
	Position     int       `json:"position"`
	PositionName string    `json:"position_name"`
	Card         TarotCard `json:"card"`
	Orientation  string    `json:"orientation"` // upright/reversed
}
