package service

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"zhanbu/internal/model"
)

//go:embed data/hexagrams.json
var hexagramsJSON embed.FS

// LiuYaoService handles the hexagram divination logic.
type LiuYaoService struct {
	hexagrams map[int]*model.Hexagram
	trigrams  []model.Trigram
	binaryIdx map[string]*model.Hexagram
	rng       *rand.Rand
}

// NewLiuYaoService creates a new LiuYaoService with hexagram data loaded.
func NewLiuYaoService() *LiuYaoService {
	svc := &LiuYaoService{
		hexagrams: make(map[int]*model.Hexagram),
		binaryIdx: make(map[string]*model.Hexagram),
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	svc.initTrigrams()
	svc.initHexagrams()
	return svc
}

//go:embed data/trigrams.json
var trigramsJSON embed.FS

func (s *LiuYaoService) initTrigrams() {
	// Load from embedded file
	data, err := trigramsJSON.ReadFile("data/trigrams.json")
	if err == nil {
		var tris []model.Trigram
		if json.Unmarshal(data, &tris) == nil && len(tris) > 0 {
			s.trigrams = tris
			return
		}
	}

	// Fallback embedded trigrams data
	s.trigrams = []model.Trigram{
		{ID: 7, Name: "乾", Symbol: "☰", Binary: 7, Nature: "天", Element: "金", YinYang: "阳"},
		{ID: 6, Name: "兑", Symbol: "☱", Binary: 6, Nature: "泽", Element: "金", YinYang: "阴"},
		{ID: 5, Name: "离", Symbol: "☲", Binary: 5, Nature: "火", Element: "火", YinYang: "阴"},
		{ID: 4, Name: "震", Symbol: "☳", Binary: 4, Nature: "雷", Element: "木", YinYang: "阳"},
		{ID: 3, Name: "巽", Symbol: "☴", Binary: 3, Nature: "风", Element: "木", YinYang: "阴"},
		{ID: 2, Name: "坎", Symbol: "☵", Binary: 2, Nature: "水", Element: "水", YinYang: "阳"},
		{ID: 1, Name: "艮", Symbol: "☶", Binary: 1, Nature: "山", Element: "土", YinYang: "阳"},
		{ID: 0, Name: "坤", Symbol: "☷", Binary: 0, Nature: "地", Element: "土", YinYang: "阴"},
	}
}

func (s *LiuYaoService) initHexagrams() {
	// Load hexagram data from embedded file
	data, err := hexagramsJSON.ReadFile("data/hexagrams.json")
	_ = err
	if err == nil {
		var hxList []model.Hexagram
		if werr := json.Unmarshal(data, &hxList); werr == nil {
			for i := range hxList {
				h := &hxList[i]
				s.hexagrams[h.ID] = h
				s.binaryIdx[h.Binary] = h
			}
			if len(s.hexagrams) >= 64 {
				return
			}
		}
	}

	// Embedded 64 hexagrams data (fallback)
	raw := getEmbeddedHexagrams()
	for i := range raw {
		h := &raw[i]
		s.hexagrams[h.ID] = h
		s.binaryIdx[h.Binary] = h
	}
}

// getLineSymbol returns the symbol for a line value.
func (s *LiuYaoService) getLineSymbol(value int) string {
	if value == 7 || value == 9 { // yang lines
		return "⚊"
	}
	return "⚋" // yin lines (6, 8)
}

// determineLine converts the number of yang coins to a LineResult.
func (s *LiuYaoService) determineLine(yangCount int) model.LineResult {
	switch yangCount {
	case 3: // 老阳
		return model.LineResult{Value: 6, Type: "old_yang", Mutable: true, Symbol: "⚋"}
	case 2: // 少阴
		return model.LineResult{Value: 8, Type: "young_yin", Mutable: false, Symbol: "⚋"}
	case 1: // 少阳
		return model.LineResult{Value: 7, Type: "young_yang", Mutable: false, Symbol: "⚊"}
	default: // 0 = 老阴
		return model.LineResult{Value: 9, Type: "old_yin", Mutable: true, Symbol: "⚊"}
	}
}

// throwOnce simulates throwing three coins once.
func (s *LiuYaoService) throwOnce() (model.LineResult, error) {
	yangCount := 0
	for i := 0; i < 3; i++ {
		if s.rng.Intn(2) == 1 {
			yangCount++
		}
	}
	return s.determineLine(yangCount), nil
}

// flipLine flips a mutable line: old_yang(6)->young_yin(8), old_yin(9)->young_yang(7).
func (s *LiuYaoService) flipLine(value int) int {
	switch value {
	case 6:
		return 8 // old yang becomes yin
	case 9:
		return 7 // old yin becomes yang
	default:
		return value // static lines don't change
	}
}

// linesToBinary converts 6 line values to a binary string.
// Yao order: index 0=初爻(bottom), index 5=上爻(top)
// Binary: reading from top to bottom as in trigrams
func linesToBinary(lines [6]model.LineResult) string {
	b := ""
	for i := 5; i >= 0; i-- {
		if lines[i].Value == 7 || lines[i].Value == 9 {
			b += "1" // yang
		} else {
			b += "0" // yin
		}
	}
	return b
}

// Throw performs a full 6-line coin throw divination.
func (s *LiuYaoService) Throw(question string) (*model.LiuYaoResult, error) {
	var lines [6]model.LineResult
	var mutableLines []int

	for i := 0; i < 6; i++ {
		line, err := s.throwOnce()
		if err != nil {
			return nil, fmt.Errorf("第%d爻投掷失败: %w", i+1, err)
		}
		lines[i] = line
		if line.Mutable {
			mutableLines = append(mutableLines, i)
		}
	}

	// Calculate ben_gua (本卦) binary
	benBinary := linesToBinary(lines)
	benGua := s.GetHexagramByBinary(benBinary)
	if benGua == nil {
		return nil, fmt.Errorf("无法找到本卦，二进制: %s", benBinary)
	}

	result := &model.LiuYaoResult{
		Question:     question,
		Lines:        lines,
		BenGua:       benGua,
		MutableLines: mutableLines,
		Timestamp:    time.Now().Format(time.RFC3339),
	}

	// Calculate bian_gua (变卦) if there are mutable lines
	if len(mutableLines) > 0 {
		var flippedLines [6]model.LineResult
		copy(flippedLines[:], lines[:])
		for _, idx := range mutableLines {
			flippedLines[idx].Value = s.flipLine(lines[idx].Value)
		}
		bianBinary := linesToBinary(flippedLines)
		result.BianGua = s.GetHexagramByBinary(bianBinary)
	}

	return result, nil
}

// GetHexagramByID returns a hexagram by its ID (1-64).
func (s *LiuYaoService) GetHexagramByID(id int) *model.Hexagram {
	return s.hexagrams[id]
}

// GetHexagramByBinary returns a hexagram by its binary string.
func (s *LiuYaoService) GetHexagramByBinary(binary string) *model.Hexagram {
	return s.binaryIdx[binary]
}

// GetAllHexagrams returns all 64 hexagrams.
func (s *LiuYaoService) GetAllHexagrams() []model.Hexagram {
	result := make([]model.Hexagram, 0, 64)
	for id := 1; id <= 64; id++ {
		if h, ok := s.hexagrams[id]; ok {
			result = append(result, *h)
		}
	}
	return result
}

// GetTrigrams returns all 8 trigrams.
func (s *LiuYaoService) GetTrigrams() []model.Trigram {
	return s.trigrams
}
