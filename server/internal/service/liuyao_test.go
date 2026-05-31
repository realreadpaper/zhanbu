package service

import (
	"testing"
)

func TestLineSymbol(t *testing.T) {
	svc := NewLiuYaoService()

	tests := []struct {
		value  int
		symbol string
	}{
		{6, "⚋"},
		{7, "⚊"},
		{8, "⚋"},
		{9, "⚊"},
	}

	for _, tc := range tests {
		result := svc.getLineSymbol(tc.value)
		if result != tc.symbol {
			t.Errorf("value %d: expected symbol '%s', got '%s'", tc.value, tc.symbol, result)
		}
	}
}

func TestThrowOnce(t *testing.T) {
	svc := NewLiuYaoService()
	svc.rng.Seed(42) // fixed seed for reproducibility

	counts := map[int]int{}
	for i := 0; i < 1000; i++ {
		svc.rng.Seed(int64(i * 42))
		result, err := svc.throwOnce()
		if err != nil {
			t.Fatalf("throwOnce returned error: %v", err)
		}
		counts[result.Value]++
	}

	// Just check all 4 values appear
	for _, v := range []int{6, 7, 8, 9} {
		if counts[v] == 0 {
			t.Errorf("value %d never occurred in 1000 throws", v)
		}
	}
}

func TestThrowValueMapping(t *testing.T) {
	svc := NewLiuYaoService()

	// Test that the value mapping is correct:
	// 3 heads (111) -> old yang (6)
	// 2 heads (110,101,011) -> young yin (8)
	// 1 head (100,010,001) -> young yang (7)
	// 0 heads (000) -> old yin (9)

	cases := []struct {
		coins  [3]int
		value  int
		mutable bool
		typ    string
	}{
		{[3]int{1, 1, 1}, 6, true, "old_yang"},
		{[3]int{1, 1, 0}, 8, false, "young_yin"},
		{[3]int{1, 0, 0}, 7, false, "young_yang"},
		{[3]int{0, 0, 0}, 9, true, "old_yin"},
	}

	for _, tc := range cases {
		yang := tc.coins[0] + tc.coins[1] + tc.coins[2]
		result := svc.determineLine(yang)

		if result.Value != tc.value {
			t.Errorf("coins %v: expected value %d, got %d", tc.coins, tc.value, result.Value)
		}
		if result.Mutable != tc.mutable {
			t.Errorf("coins %v: expected mutable %v, got %v", tc.coins, tc.mutable, result.Mutable)
		}
		if result.Type != tc.typ {
			t.Errorf("coins %v: expected type '%s', got '%s'", tc.coins, tc.typ, result.Type)
		}
	}
}

func TestGetHexagramByBinary(t *testing.T) {
	svc := NewLiuYaoService()

	// Test lookup by binary strings
	tests := []struct {
		binary string
		id     int
		name   string
	}{
		{"111111", 1, "乾"},
		{"000000", 64, "火水未济"},
		{"111000", 8, "水地比"},
		{"010010", 46, "地风升"},
		{"101101", 19, "地泽临"},
		{"101010", 22, "山火贲"},
	}

	for _, tc := range tests {
		h := svc.GetHexagramByBinary(tc.binary)
		if h == nil {
			t.Errorf("binary %s: expected hexagram, got nil", tc.binary)
			continue
		}
		if h.ID != tc.id || h.Name != tc.name {
			t.Errorf("binary %s: expected id=%d name=%s, got id=%d name=%s",
				tc.binary, tc.id, tc.name, h.ID, h.Name)
		}
	}
}

func TestGetHexagramByID(t *testing.T) {
	svc := NewLiuYaoService()

	h := svc.GetHexagramByID(1)
	if h == nil || h.ID != 1 {
		t.Errorf("expected hexagram 1, got %v", h)
	}

	h = svc.GetHexagramByID(64)
	if h == nil || h.ID != 64 {
		t.Errorf("expected hexagram 64, got %v", h)
	}

	h = svc.GetHexagramByID(999)
	if h != nil {
		t.Errorf("expected nil for invalid ID")
	}
}

func TestThrowAllLines(t *testing.T) {
	svc := NewLiuYaoService()
	svc.rng.Seed(42)

	result, err := svc.Throw("测试问题")
	if err != nil {
		t.Fatalf("Throw returned error: %v", err)
	}

	if result.Question != "测试问题" {
		t.Errorf("expected question '测试问题', got '%s'", result.Question)
	}

	if result.BenGua == nil {
		t.Fatal("expected ben_gua, got nil")
	}

	if len(result.Lines) != 6 {
		t.Errorf("expected 6 lines, got %d", len(result.Lines))
	}

	// Lines should be indexed 0-5 (初爻 to 上爻)
	if result.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}

	// Check mutable lines
	for i, line := range result.Lines {
		if line.Value < 6 || line.Value > 9 {
			t.Errorf("line %d: unexpected value %d", i, line.Value)
		}
		if line.Mutable {
			found := false
			for _, ml := range result.MutableLines {
				if ml == i {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("line %d is mutable but not in mutable_lines", i)
			}
		}
	}

	// With fixed seed, we should get consistent results
	t.Logf("BenGua: %s (id=%d)", result.BenGua.Name, result.BenGua.ID)
	t.Logf("MutableLines: %v", result.MutableLines)
	if result.BianGua != nil {
		t.Logf("BianGua: %s (id=%d)", result.BianGua.Name, result.BianGua.ID)
	}
}

func TestThrowNoQuestion(t *testing.T) {
	svc := NewLiuYaoService()
	svc.rng.Seed(42)

	result, err := svc.Throw("")
	if err != nil {
		t.Fatalf("Throw with empty question returned error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.BenGua == nil {
		t.Fatal("expected ben_gua, got nil")
	}
}

func TestGetAllHexagrams(t *testing.T) {
	svc := NewLiuYaoService()
	hexagrams := svc.GetAllHexagrams()

	if len(hexagrams) != 64 {
		t.Errorf("expected 64 hexagrams, got %d", len(hexagrams))
	}

	ids := make(map[int]bool)
	for _, h := range hexagrams {
		ids[h.ID] = true
	}
	for i := 1; i <= 64; i++ {
		if !ids[i] {
			t.Errorf("hexagram %d missing", i)
		}
	}
}

func TestFlipLine(t *testing.T) {
	svc := NewLiuYaoService()

	tests := []struct {
		value  int
		result int
	}{
		{6, 8}, // old yang -> young yin
		{9, 7}, // old yin -> young yang
		{7, 7}, // static lines don't change
		{8, 8},
	}

	for _, tc := range tests {
		got := svc.flipLine(tc.value)
		if got != tc.result {
			t.Errorf("flipLine(%d): expected %d, got %d", tc.value, tc.result, got)
		}
	}
}

func TestGetTrigrams(t *testing.T) {
	svc := NewLiuYaoService()
	trigrams := svc.GetTrigrams()

	if len(trigrams) != 8 {
		t.Errorf("expected 8 trigrams, got %d", len(trigrams))
	}

	// Just check we have 8 trigrams with non-empty names
	for _, tg := range trigrams {
		if tg.Name == "" {
			t.Errorf("trigram has empty name")
		}
	}
}

func TestHexagram64(t *testing.T) {
	// Quick check that all 64 hexagrams have valid data
	svc := NewLiuYaoService()
	for i := 1; i <= 64; i++ {
		h := svc.GetHexagramByID(i)
		if h == nil {
			t.Errorf("hexagram %d is nil", i)
			continue
		}
		if h.Name == "" {
			t.Errorf("hexagram %d has empty name", i)
		}
		if h.Binary == "" {
			t.Errorf("hexagram %d has empty binary", i)
		}
		if len(h.LineTexts) != 6 {
			t.Errorf("hexagram %d has %d line texts, expected 6", i, len(h.LineTexts))
		}
	}
}

func TestHexagramLookup(t *testing.T) {
	svc := NewLiuYaoService()
	h := svc.GetHexagramByBinary("000101")
	if h == nil {
		t.Fatal("hexagram 000101 not found")
	}
	t.Logf("Found: %s (binary=%s)", h.Name, h.Binary)
	total := len(svc.GetAllHexagrams())
	t.Logf("Total hexagrams loaded: %d", total)
}

func TestMutableProbability(t *testing.T) {
	svc := NewLiuYaoService()

	mutableCount := 0
	total := 6000 // 1000 throws * 6 lines
	for i := 0; i < 1000; i++ {
		result, err := svc.Throw("")
		if err != nil {
			t.Fatalf("throw %d failed: %v", i, err)
		}
		if result == nil {
			t.Fatal("result is nil")
		}
		for _, line := range result.Lines {
			if line.Mutable {
				mutableCount++
			}
		}
	}

	ratio := float64(mutableCount) / float64(total)
	// Expected: 25% (1/4 chance per line since 6 and 9 each have 1/8 prob, total 2/8 = 1/4)
	if ratio < 0.20 || ratio > 0.30 {
		t.Errorf("mutable ratio %.3f is too far from 0.25", ratio)
	}
}
