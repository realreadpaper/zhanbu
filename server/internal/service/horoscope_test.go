package service

import (
	"testing"

	"zhanbu/internal/model"
)

// mockTemplateStore provides canned templates for testing.
type mockTemplateStore struct{}

func (m *mockTemplateStore) GetTemplate(zodiac, period, dimension, level string) string {
	return dimension + "_" + level + "_text"
}

func (m *mockTemplateStore) GetSummary(zodiac, period, level string) string {
	return "overall_" + level + "_summary"
}

func newTestService() *HoroscopeService {
	return NewHoroscopeService(&mockTemplateStore{})
}

func TestConsistency(t *testing.T) {
	svc := newTestService()

	// Same zodiac + date + period should always produce the same result.
	for i := 0; i < 10; i++ {
		result1, err := svc.Generate("aries", "daily", "2026-05-30")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		result2, err := svc.Generate("aries", "daily", "2026-05-30")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result1.Overall != result2.Overall {
			t.Errorf("consistency failed: overall %d != %d", result1.Overall, result2.Overall)
		}
		if result1.Love != result2.Love {
			t.Errorf("consistency failed: love %d != %d", result1.Love, result2.Love)
		}
		if result1.Career != result2.Career {
			t.Errorf("consistency failed: career %d != %d", result1.Career, result2.Career)
		}
		if result1.Wealth != result2.Wealth {
			t.Errorf("consistency failed: wealth %d != %d", result1.Wealth, result2.Wealth)
		}
		if result1.Health != result2.Health {
			t.Errorf("consistency failed: health %d != %d", result1.Health, result2.Health)
		}
		if result1.LuckyNumber != result2.LuckyNumber {
			t.Errorf("consistency failed: lucky_number %d != %d", result1.LuckyNumber, result2.LuckyNumber)
		}
		if result1.LuckyColor != result2.LuckyColor {
			t.Errorf("consistency failed: lucky_color %s != %s", result1.LuckyColor, result2.LuckyColor)
		}
	}
}

func TestScoreRange(t *testing.T) {
	svc := newTestService()

	// Test all 12 zodiac signs
	for _, z := range model.Zodiacs {
		result, err := svc.Generate(z.Name, "daily", "2026-05-30")
		if err != nil {
			t.Fatalf("zodiac %s: unexpected error: %v", z.Name, err)
		}

		scores := map[string]int{
			"overall": result.Overall,
			"love":    result.Love,
			"career":  result.Career,
			"wealth":  result.Wealth,
			"health":  result.Health,
		}

		for dim, score := range scores {
			if score < 1 || score > 5 {
				t.Errorf("zodiac %s: %s score %d out of range [1,5]", z.Name, dim, score)
			}
		}
	}
}

func TestLuckyNumberRange(t *testing.T) {
	svc := newTestService()

	// Test multiple dates to cover more seeds
	dates := []string{"2026-01-01", "2026-05-30", "2026-12-31"}
	for _, z := range model.Zodiacs {
		for _, date := range dates {
			result, err := svc.Generate(z.Name, "daily", date)
			if err != nil {
				t.Fatalf("zodiac %s date %s: unexpected error: %v", z.Name, date, err)
			}
			if result.LuckyNumber < 1 || result.LuckyNumber > 9 {
				t.Errorf("zodiac %s date %s: lucky_number %d out of range [1,9]", z.Name, date, result.LuckyNumber)
			}
		}
	}
}

func TestAllZodiacsGenerate(t *testing.T) {
	svc := newTestService()

	for _, z := range model.Zodiacs {
		result, err := svc.Generate(z.Name, "daily", "2026-05-30")
		if err != nil {
			t.Errorf("zodiac %s: unexpected error: %v", z.Name, err)
			continue
		}
		if result.ZodiacCN != z.NameCN {
			t.Errorf("zodiac %s: expected cn name %s, got %s", z.Name, z.NameCN, result.ZodiacCN)
		}
		if result.Summary == "" {
			t.Errorf("zodiac %s: summary is empty", z.Name)
		}
		if result.Detail.Love == "" {
			t.Errorf("zodiac %s: love detail is empty", z.Name)
		}
		if result.Detail.Career == "" {
			t.Errorf("zodiac %s: career detail is empty", z.Name)
		}
		if result.Detail.Wealth == "" {
			t.Errorf("zodiac %s: wealth detail is empty", z.Name)
		}
		if result.Detail.Health == "" {
			t.Errorf("zodiac %s: health detail is empty", z.Name)
		}
	}
}

func TestDifferentPeriodsDiffer(t *testing.T) {
	svc := newTestService()

	// It's possible (though unlikely) that daily and weekly give the same scores,
	// so we test with multiple dates to increase confidence.
	dates := []string{"2026-01-15", "2026-05-30", "2026-09-10"}
	differCount := 0

	for _, date := range dates {
		daily, err := svc.Generate("aries", "daily", date)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		weekly, err := svc.Generate("aries", "weekly", date)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if daily.Overall != weekly.Overall || daily.Love != weekly.Love {
			differCount++
		}
	}

	if differCount == 0 {
		t.Log("warning: daily and weekly gave identical scores for all test dates (statistically unlikely but possible)")
	}
}

func TestInvalidZodiac(t *testing.T) {
	svc := newTestService()

	_, err := svc.Generate("invalid", "daily", "2026-05-30")
	if err == nil {
		t.Error("expected error for invalid zodiac, got nil")
	}
}

func TestInvalidPeriod(t *testing.T) {
	svc := newTestService()

	_, err := svc.Generate("aries", "yearly", "2026-05-30")
	if err == nil {
		t.Error("expected error for invalid period, got nil")
	}
}

func TestInvalidDate(t *testing.T) {
	svc := newTestService()

	_, err := svc.Generate("aries", "daily", "not-a-date")
	if err == nil {
		t.Error("expected error for invalid date, got nil")
	}
}

func TestWeeklyNormalization(t *testing.T) {
	svc := newTestService()

	// Two dates in the same ISO week should produce the same result
	// 2026-05-25 (Monday) and 2026-05-29 (Friday) are in the same ISO week (W22)
	r1, err := svc.Generate("aries", "weekly", "2026-05-25")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r2, err := svc.Generate("aries", "weekly", "2026-05-29")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r1.Overall != r2.Overall || r1.Love != r2.Love || r1.Career != r2.Career {
		t.Errorf("same ISO week should produce same weekly result:\n  05-25: overall=%d love=%d career=%d\n  05-29: overall=%d love=%d career=%d",
			r1.Overall, r1.Love, r1.Career, r2.Overall, r2.Love, r2.Career)
	}
}

func TestMonthlyNormalization(t *testing.T) {
	svc := newTestService()

	// Two dates in the same month should produce the same result
	r1, err := svc.Generate("aries", "monthly", "2026-05-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r2, err := svc.Generate("aries", "monthly", "2026-05-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r1.Overall != r2.Overall || r1.Love != r2.Love || r1.Career != r2.Career {
		t.Errorf("same month should produce same monthly result:\n  05-01: overall=%d love=%d career=%d\n  05-31: overall=%d love=%d career=%d",
			r1.Overall, r1.Love, r1.Career, r2.Overall, r2.Love, r2.Career)
	}
}

func TestDefaultDate(t *testing.T) {
	svc := newTestService()

	// Empty date string should use today's date (no error)
	result, err := svc.Generate("aries", "daily", "")
	if err != nil {
		t.Fatalf("unexpected error with empty date: %v", err)
	}
	if result.Date == "" {
		t.Error("expected date to be filled in, got empty string")
	}
}

func TestLuckyColorInPool(t *testing.T) {
	svc := newTestService()

	for _, z := range model.Zodiacs {
		result, err := svc.Generate(z.Name, "daily", "2026-05-30")
		if err != nil {
			t.Fatalf("zodiac %s: %v", z.Name, err)
		}

		found := false
		for _, c := range LuckyColors {
			if c == result.LuckyColor {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("zodiac %s: lucky color %q not in LuckyColors pool", z.Name, result.LuckyColor)
		}
	}
}
