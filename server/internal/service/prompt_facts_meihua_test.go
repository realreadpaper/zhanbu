package service

import (
	"testing"
)

func TestParseMeihuaFacts_Success(t *testing.T) {
	resultJSON := `{
		"method": "time",
		"source_values": {
			"year_branch": "午",
			"lunar_month": 5,
			"lunar_day": 9,
			"lunar_display": "丙午年五月初九子时",
			"hour_branch": "子",
			"hour_number": 1
		},
		"ben_gua": {
			"upper": {"name": "兑", "element": "金"},
			"lower": {"name": "离", "element": "火"},
			"name": "泽火革"
		},
		"hu_gua": {
			"upper": {"name": "乾", "element": "金"},
			"lower": {"name": "巽", "element": "木"},
			"name": "天风姤"
		},
		"bian_gua": {
			"upper": {"name": "兑", "element": "金"},
			"lower": {"name": "艮", "element": "土"},
			"name": "泽山咸"
		},
		"moving_line": 1,
		"ti_yong": {
			"ti": {"name": "兑", "element": "金"},
			"yong": {"name": "离", "element": "火"},
			"relation": "用克体"
		}
	}`

	facts, err := parseMeihuaFacts(resultJSON, "最近事业如何？")
	if err != nil {
		t.Fatalf("parseMeihuaFacts: %v", err)
	}

	expected := map[string]string{
		"Method":         "时间起卦（丙午年五月初九子时）",
		"BenGua":         "泽火革（兑上离下）",
		"HuGua":          "天风姤（乾上巽下）",
		"BianGua":        "泽山咸（兑上艮下）",
		"MovingLine":     "初爻动",
		"Ti":             "兑（金）",
		"Yong":           "离（火）",
		"TiYongRelation": "用克体",
	}

	for key, expectedVal := range expected {
		actualVal, ok := facts[key]
		if !ok {
			t.Errorf("missing fact key: %s", key)
			continue
		}
		if actualVal != expectedVal {
			t.Errorf("fact %s: expected %q, got %q", key, expectedVal, actualVal)
		}
	}
}

func TestParseMeihuaFacts_NumberMethod(t *testing.T) {
	resultJSON := `{
		"method": "number",
		"source_values": {
			"numbers": [12, 34]
		},
		"ben_gua": {
			"upper": {"name": "乾", "element": "金"},
			"lower": {"name": "兑", "element": "金"},
			"name": "天泽履"
		},
		"hu_gua": {
			"upper": {"name": "巽", "element": "木"},
			"lower": {"name": "离", "element": "火"},
			"name": "风火家人"
		},
		"bian_gua": {
			"upper": {"name": "乾", "element": "金"},
			"lower": {"name": "离", "element": "火"},
			"name": "天火同人"
		},
		"moving_line": 3,
		"ti_yong": {
			"ti": {"name": "乾", "element": "金"},
			"yong": {"name": "兑", "element": "金"},
			"relation": "比和"
		}
	}`

	facts, err := parseMeihuaFacts(resultJSON, "数字 12 34，问事业")
	if err != nil {
		t.Fatalf("parseMeihuaFacts: %v", err)
	}

	if facts["Method"] != "数字起卦（12、34）" {
		t.Errorf("expected method 数字起卦（12、34）, got %s", facts["Method"])
	}
	if facts["MovingLine"] != "三爻动" {
		t.Errorf("expected moving line 三爻动, got %s", facts["MovingLine"])
	}
	if facts["TiYongRelation"] != "比和" {
		t.Errorf("expected relation 比和, got %s", facts["TiYongRelation"])
	}
}

func TestParseMeihuaFacts_InvalidJSON(t *testing.T) {
	_, err := parseMeihuaFacts("not json", "test")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseMeihuaFacts_AllFieldsPresent(t *testing.T) {
	resultJSON := `{
		"method": "time",
		"source_values": {"year_branch": "午", "lunar_month": 5, "lunar_day": 9, "hour_branch": "子"},
		"ben_gua": {"upper": {"name": "兑", "element": "金"}, "lower": {"name": "离", "element": "火"}, "name": "泽火革"},
		"hu_gua": {"upper": {"name": "乾", "element": "金"}, "lower": {"name": "巽", "element": "木"}, "name": "天风姤"},
		"bian_gua": {"upper": {"name": "兑", "element": "金"}, "lower": {"name": "艮", "element": "土"}, "name": "泽山咸"},
		"moving_line": 2,
		"ti_yong": {"ti": {"name": "兑", "element": "金"}, "yong": {"name": "离", "element": "火"}, "relation": "用克体"}
	}`

	facts, err := parseMeihuaFacts(resultJSON, "test")
	if err != nil {
		t.Fatalf("parseMeihuaFacts: %v", err)
	}

	requiredKeys := []string{"Method", "BenGua", "HuGua", "BianGua", "MovingLine", "Ti", "Yong", "TiYongRelation"}
	for _, key := range requiredKeys {
		if _, ok := facts[key]; !ok {
			t.Errorf("missing required fact key: %s", key)
		}
	}
}
