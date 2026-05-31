package model

import "time"

// Trigram represents one of the 8 basic trigrams.
type Trigram struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Binary  int    `json:"binary"`
	Nature  string `json:"nature"`
	Element string `json:"element"`
	YinYang string `json:"yin_yang"`
}

// Hexagram represents one of the 64 hexagrams.
type Hexagram struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	NameShort    string   `json:"name_short"`
	UpperTrigram string   `json:"upper_trigram"`
	LowerTrigram string   `json:"lower_trigram"`
	Binary       string   `json:"binary"`
	Judgment     string   `json:"judgment"`
	Image        string   `json:"image"`
	LineTexts    []string `json:"line_texts"`
	Description  string   `json:"description"`
}

// HexagramBrief is a light version for list display.
type HexagramBrief struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	NameShort    string `json:"name_short"`
	UpperTrigram string `json:"upper_trigram"`
	LowerTrigram string `json:"lower_trigram"`
	Binary       string `json:"binary"`
	Description  string `json:"description"`
}

// LineResult represents a single line throw.
type LineResult struct {
	Value   int    `json:"value"`   // 6/7/8/9
	Type    string `json:"type"`    // old_yang/young_yang/old_yin/young_yin
	Mutable bool   `json:"mutable"` // is mutable line?
	Symbol  string `json:"symbol"`  // ⚊ or ⚋
}

// LiuYaoResult represents the full throw result.
type LiuYaoResult struct {
	Question     string       `json:"question,omitempty"`
	Lines        [6]LineResult `json:"lines"`
	BenGua       *Hexagram    `json:"ben_gua"`
	BianGua      *Hexagram    `json:"bian_gua,omitempty"`
	MutableLines []int         `json:"mutable_lines"`
	Timestamp    string       `json:"timestamp"`
}

// BaZiPillar represents one pillar (year/month/day/hour).
type BaZiPillar struct {
	TianGan    string   `json:"tian_gan"`
	DiZhi      string   `json:"di_zhi"`
	WuXing     string   `json:"wu_xing"`
	NaYin      string   `json:"na_yin"`
	HiddenGan  []string `json:"hidden_gan,omitempty"`
}

// BaZiPillars contains all four pillars.
type BaZiPillars struct {
	Year  BaZiPillar `json:"year"`
	Month BaZiPillar `json:"month"`
	Day   BaZiPillar `json:"day"`
	Hour  BaZiPillar `json:"hour"`
}

// FiveElementAnalysis represents five element distribution.
type FiveElementAnalysis struct {
	Metal     int    `json:"metal"`
	Wood      int    `json:"wood"`
	Water     int    `json:"water"`
	Fire      int    `json:"fire"`
	Earth     int    `json:"earth"`
	Strongest string `json:"strongest"`
	Weakest   string `json:"weakest"`
	DayMaster string `json:"day_master"`
	Strength  string `json:"strength"`
	YongShen  string `json:"yong_shen"`
	JiShen    string `json:"ji_shen"`
}

// TenGod represents one ten-god analysis.
type TenGod struct {
	Position string `json:"position"`
	TianGan  string `json:"tian_gan"`
	God      string `json:"god"`
}

// BaZiBirthInfo contains birth information.
type BaZiBirthInfo struct {
	Solar string `json:"solar"`
	Lunar string `json:"lunar"`
}

// BaZiResult represents the full bazi calculation result.
type BaZiResult struct {
	Birth        BaZiBirthInfo       `json:"birth"`
	Pillars      BaZiPillars         `json:"pillars"`
	FiveElements FiveElementAnalysis `json:"five_elements"`
	TenGods      []TenGod            `json:"ten_gods"`
}

// LiuYaoRequest is the request for liuyao throw.
type LiuYaoRequest struct {
	Question string `json:"question"`
}

// BaZiRequest is the request for bazi calculation.
type BaZiRequest struct {
	BirthDate string `json:"birth_date"`
	BirthTime string `json:"birth_time"`
	Gender    string `json:"gender"`
}

// ToBrief converts a Hexagram to a brief version.
func (h *Hexagram) ToBrief() HexagramBrief {
	return HexagramBrief{
		ID:           h.ID,
		Name:         h.Name,
		NameShort:    h.NameShort,
		UpperTrigram: h.UpperTrigram,
		LowerTrigram: h.LowerTrigram,
		Binary:       h.Binary,
		Description:  h.Description,
	}
}

// variable for timestamp
var _ = time.Now
