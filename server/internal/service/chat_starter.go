package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"zhanbu/internal/model"
	apperrors "zhanbu/pkg/errors"
)

// ChatModeStarter creates divination records for sessions started directly
// from the chat page.
type ChatModeStarter struct {
	recordRepo DivinationRecordSaver
	tarot      *TarotService
	liuyaoV2   *LiuYaoV2Service
	bazi       *BaZiService
	horoscope  *HoroscopeService
}

// NewChatModeStarter returns a starter that delegates to existing divination services.
func NewChatModeStarter(
	recordRepo DivinationRecordSaver,
	tarot *TarotService,
	liuyaoV2 *LiuYaoV2Service,
	bazi *BaZiService,
	horoscope *HoroscopeService,
) *ChatModeStarter {
	return &ChatModeStarter{
		recordRepo: recordRepo,
		tarot:      tarot,
		liuyaoV2:   liuyaoV2,
		bazi:       bazi,
		horoscope:  horoscope,
	}
}

// Start creates a real divination record for the selected mode.
func (s *ChatModeStarter) Start(userID uint, divinationType string, question string) (*model.DivinationRecord, error) {
	switch normalizeChatDivinationType(divinationType) {
	case "tarot":
		return s.startTarot(userID, question)
	case "liuyao_v2":
		return s.startLiuYaoV2(userID, question)
	case "bazi":
		return s.startBaZi(userID, question)
	case "horoscope":
		return s.startHoroscope(userID, question)
	default:
		return nil, apperrors.New(apperrors.ErrBadRequest, "invalid divination type")
	}
}

func (s *ChatModeStarter) startTarot(userID uint, question string) (*model.DivinationRecord, error) {
	if s.tarot == nil {
		return nil, apperrors.New(apperrors.ErrAIServiceUnavailable, "tarot service is not configured")
	}
	result, err := s.tarot.DrawCards("single", question)
	if err != nil {
		return nil, err
	}
	return s.saveResult(userID, "tarot", question, result)
}

func (s *ChatModeStarter) startLiuYaoV2(userID uint, question string) (*model.DivinationRecord, error) {
	if s.liuyaoV2 == nil {
		return nil, apperrors.New(apperrors.ErrAIServiceUnavailable, "liuyao service is not configured")
	}
	result, err := s.liuyaoV2.Throw(question, "")
	if err != nil {
		return nil, err
	}
	return s.saveResult(userID, "liuyao_v2", question, result)
}

func (s *ChatModeStarter) startBaZi(userID uint, question string) (*model.DivinationRecord, error) {
	if s.bazi == nil {
		return nil, apperrors.New(apperrors.ErrAIServiceUnavailable, "bazi service is not configured")
	}
	birthDate, birthTime, gender, ok := parseBaZiInput(question)
	if !ok {
		return nil, apperrors.New(apperrors.ErrBadRequest, "请在问题中包含出生日期和时间，例如：1990-05-12 08:30 女，看看事业")
	}
	result, err := s.bazi.Calculate(birthDate, birthTime, gender)
	if err != nil {
		return nil, err
	}
	return s.saveResult(userID, "bazi", question, result)
}

func (s *ChatModeStarter) startHoroscope(userID uint, question string) (*model.DivinationRecord, error) {
	if s.horoscope == nil {
		return nil, apperrors.New(apperrors.ErrAIServiceUnavailable, "horoscope service is not configured")
	}
	zodiac := parseZodiac(question)
	if zodiac == "" {
		return nil, apperrors.New(apperrors.ErrBadRequest, "请在问题中包含星座，例如：天蝎座今天事业运如何")
	}
	result, err := s.horoscope.Generate(zodiac, "daily", "")
	if err != nil {
		return nil, err
	}
	return s.saveResult(userID, "horoscope", question, result)
}

func (s *ChatModeStarter) saveResult(userID uint, divinationType string, question string, result any) (*model.DivinationRecord, error) {
	if s.recordRepo == nil {
		return nil, apperrors.New(apperrors.ErrAIServiceUnavailable, "record repository is not configured")
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal divination result: %w", err)
	}
	record := &model.DivinationRecord{
		UserID:   userID,
		Type:     divinationType,
		Question: question,
		Result:   string(resultJSON),
	}
	if err := s.recordRepo.Create(record); err != nil {
		return nil, err
	}
	return record, nil
}

var baziInputPattern = regexp.MustCompile(`(\d{4})[-/.年](\d{1,2})[-/.月](\d{1,2})日?\s+(\d{1,2}):(\d{2})`)

func parseBaZiInput(input string) (birthDate string, birthTime string, gender string, ok bool) {
	matches := baziInputPattern.FindStringSubmatch(input)
	if len(matches) != 6 {
		return "", "", "unknown", false
	}

	gender = "unknown"
	if strings.Contains(input, "女") || strings.Contains(strings.ToLower(input), "female") {
		gender = "female"
	} else if strings.Contains(input, "男") || strings.Contains(strings.ToLower(input), "male") {
		gender = "male"
	}

	return fmt.Sprintf("%s-%02s-%02s", matches[1], matches[2], matches[3]),
		fmt.Sprintf("%02s:%s", matches[4], matches[5]),
		gender,
		true
}

func parseZodiac(input string) string {
	normalized := strings.ToLower(strings.TrimSpace(input))
	aliases := map[string]string{
		"白羊":          "aries",
		"白羊座":         "aries",
		"aries":       "aries",
		"金牛":          "taurus",
		"金牛座":         "taurus",
		"taurus":      "taurus",
		"双子":          "gemini",
		"双子座":         "gemini",
		"gemini":      "gemini",
		"巨蟹":          "cancer",
		"巨蟹座":         "cancer",
		"cancer":      "cancer",
		"狮子":          "leo",
		"狮子座":         "leo",
		"leo":         "leo",
		"处女":          "virgo",
		"处女座":         "virgo",
		"virgo":       "virgo",
		"天秤":          "libra",
		"天秤座":         "libra",
		"libra":       "libra",
		"天蝎":          "scorpio",
		"天蝎座":         "scorpio",
		"scorpio":     "scorpio",
		"射手":          "sagittarius",
		"射手座":         "sagittarius",
		"sagittarius": "sagittarius",
		"摩羯":          "capricorn",
		"摩羯座":         "capricorn",
		"魔羯":          "capricorn",
		"魔羯座":         "capricorn",
		"capricorn":   "capricorn",
		"水瓶":          "aquarius",
		"水瓶座":         "aquarius",
		"aquarius":    "aquarius",
		"双鱼":          "pisces",
		"双鱼座":         "pisces",
		"pisces":      "pisces",
	}
	for alias, zodiac := range aliases {
		if strings.Contains(normalized, alias) {
			return zodiac
		}
	}
	return ""
}
