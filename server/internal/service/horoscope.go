package service

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"zhanbu/internal/model"
)

// HoroscopeTemplateStore provides access to horoscope templates.
// Interface allows mocking in tests.
type HoroscopeTemplateStore interface {
	GetTemplate(zodiac, period, dimension, level string) string
	GetSummary(zodiac, period, level string) string
}

// JSONTemplateStore loads templates from JSON files.
type JSONTemplateStore struct {
	daily   map[string]map[string]string
	weekly  map[string]map[string]string
	monthly map[string]map[string]string
}

// templateFileData is the structure of a single template JSON file.
// Key = zodiac english name, value = map of template_key -> text.
type templateFileData map[string]map[string]string

// NewJSONTemplateStore creates a store by reading JSON files from dataDir.
func NewJSONTemplateStore(dataDir string) (*JSONTemplateStore, error) {
	store := &JSONTemplateStore{}

	var err error
	store.daily, err = loadTemplateFile(filepath.Join(dataDir, "daily.json"))
	if err != nil {
		return nil, fmt.Errorf("load daily.json: %w", err)
	}
	store.weekly, err = loadTemplateFile(filepath.Join(dataDir, "weekly.json"))
	if err != nil {
		return nil, fmt.Errorf("load weekly.json: %w", err)
	}
	store.monthly, err = loadTemplateFile(filepath.Join(dataDir, "monthly.json"))
	if err != nil {
		return nil, fmt.Errorf("load monthly.json: %w", err)
	}

	return store, nil
}

func loadTemplateFile(path string) (map[string]map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result templateFileData
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// periodMap returns the template map for a period.
func (s *JSONTemplateStore) periodMap(period string) map[string]map[string]string {
	switch period {
	case "daily":
		return s.daily
	case "weekly":
		return s.weekly
	case "monthly":
		return s.monthly
	default:
		return s.daily
	}
}

// GetTemplate returns the template text for a zodiac/period/dimension/level combination.
// dimension = love|career|wealth|health, level = high|mid|low.
func (s *JSONTemplateStore) GetTemplate(zodiac, period, dimension, level string) string {
	pm := s.periodMap(period)
	if pm == nil {
		return ""
	}
	zodiacTemplates, ok := pm[zodiac]
	if !ok {
		return ""
	}
	key := dimension + "_" + level
	return zodiacTemplates[key]
}

// GetSummary returns the overall summary template for a zodiac/period/level.
func (s *JSONTemplateStore) GetSummary(zodiac, period, level string) string {
	pm := s.periodMap(period)
	if pm == nil {
		return ""
	}
	zodiacTemplates, ok := pm[zodiac]
	if !ok {
		return ""
	}
	key := "overall_" + level
	return zodiacTemplates[key]
}

// LuckyColors is the pool of possible lucky colors.
var LuckyColors = []string{
	"红色", "橙色", "黄色", "绿色", "蓝色",
	"紫色", "粉色", "白色", "金色", "银色",
	"青色", "棕色",
}

// HoroscopeService generates horoscope readings.
type HoroscopeService struct {
	templates HoroscopeTemplateStore
}

// NewHoroscopeService creates a new horoscope service.
func NewHoroscopeService(templates HoroscopeTemplateStore) *HoroscopeService {
	return &HoroscopeService{templates: templates}
}

// Generate produces a horoscope result for the given zodiac, period, and optional date.
// If dateStr is empty, today's date is used.
// The result is deterministic: same zodiac+date+period always yields the same output.
func (s *HoroscopeService) Generate(zodiac, period, dateStr string) (*model.HoroscopeResult, error) {
	zodiacInfo := model.ZodiacLookup(zodiac)
	if zodiacInfo == nil {
		return nil, fmt.Errorf("未知的星座: %s", zodiac)
	}

	// Validate period
	switch period {
	case "daily", "weekly", "monthly":
		// ok
	default:
		return nil, fmt.Errorf("无效的运势周期: %s (可选: daily, weekly, monthly)", period)
	}

	// Resolve date
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	// For weekly/monthly, normalize the date string to ensure consistency
	// weekly: use ISO week string (e.g. "2026-W22")
	// monthly: use year-month (e.g. "2026-05")
	normalizedDate := dateStr
	switch period {
	case "weekly":
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("无效的日期格式，应为 YYYY-MM-DD: %s", dateStr)
		}
		year, week := t.ISOWeek()
		normalizedDate = fmt.Sprintf("%d-W%02d", year, week)
	case "monthly":
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("无效的日期格式，应为 YYYY-MM-DD: %s", dateStr)
		}
		normalizedDate = t.Format("2006-01")
	case "daily":
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("无效的日期格式，应为 YYYY-MM-DD: %s", dateStr)
		}
		normalizedDate = t.Format("2006-01-02")
	}

	// Generate deterministic seed
	seed := generateSeed(zodiac, normalizedDate, period)

	// Generate scores (1-5)
	overall := intRangeFromSeed(seed, 0, 1, 5)
	love := intRangeFromSeed(seed, 1, 1, 5)
	career := intRangeFromSeed(seed, 2, 1, 5)
	wealth := intRangeFromSeed(seed, 3, 1, 5)
	health := intRangeFromSeed(seed, 4, 1, 5)

	// Generate lucky elements
	luckyNumber := intRangeFromSeed(seed, 5, 1, 9)
	luckyColorIdx := intRangeFromSeed(seed, 6, 0, len(LuckyColors)-1)
	luckyColor := LuckyColors[luckyColorIdx]

	// Map scores to template levels
	overallLevel := scoreToLevel(overall)
	loveLevel := scoreToLevel(love)
	careerLevel := scoreToLevel(career)
	wealthLevel := scoreToLevel(wealth)
	healthLevel := scoreToLevel(health)

	// Get template texts
	summary := s.templates.GetSummary(zodiac, period, overallLevel)
	loveText := s.templates.GetTemplate(zodiac, period, "love", loveLevel)
	careerText := s.templates.GetTemplate(zodiac, period, "career", careerLevel)
	wealthText := s.templates.GetTemplate(zodiac, period, "wealth", wealthLevel)
	healthText := s.templates.GetTemplate(zodiac, period, "health", healthLevel)

	result := &model.HoroscopeResult{
		Zodiac:      zodiacInfo.Name,
		ZodiacCN:    zodiacInfo.NameCN,
		Period:      period,
		Date:        dateStr,
		Overall:     overall,
		Love:        love,
		Career:      career,
		Wealth:      wealth,
		Health:      health,
		LuckyNumber: luckyNumber,
		LuckyColor:  luckyColor,
		Summary:     summary,
		Detail: model.HoroscopeDetail{
			Love:   loveText,
			Career: careerText,
			Wealth: wealthText,
			Health: healthText,
		},
	}

	return result, nil
}

// generateSeed computes a deterministic uint64 seed from zodiac+date+period using SHA-256.
func generateSeed(zodiac, date, period string) uint64 {
	input := zodiac + date + period
	hash := sha256.Sum256([]byte(input))
	// Take the first 8 bytes as a uint64
	return binary.BigEndian.Uint64(hash[:8])
}

// intRangeFromSeed extracts a deterministic integer in [min, max] from the seed.
// offset selects which "lane" of the seed to use (0-7).
func intRangeFromSeed(seed uint64, offset, min, max int) int {
	// Use different parts of the seed for different dimensions
	// XOR with offset to get variation
	shifted := seed >> (offset * 8)
	// Mix bits
	mixed := shifted ^ (seed >> 32)
	// Ensure positive and within range
	val := int(mixed & 0xFFFFFFFF)
	if val < 0 {
		val = -val
	}
	return min + val%(max-min+1)
}

// scoreToLevel maps a 1-5 score to a template level.
func scoreToLevel(score int) string {
	switch score {
	case 1, 2:
		return "low"
	case 3:
		return "mid"
	case 4, 5:
		return "high"
	default:
		return "mid"
	}
}
