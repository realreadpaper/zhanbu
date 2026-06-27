package service

import (
	"strings"
	"testing"

	"zhanbu/config"
)

var testMeihuaProfile = &config.ProfileConfig{
	Version:        "v1",
	DivinationType: "meihua",
	Name:           "康节先生",
	Title:          "梅花易数解读者",
	Subtitle:       "师承邵雍一脉，观物取象",
	Icon:           "🌸",
	Description:    "以邵雍梅花易数体系为核心的 AI 解读者",
	Enabled:        true,
	SystemIdentity: "你是「康节先生」人格，一位以邵雍梅花易数体系为核心的 AI 解读者。\n你不自称历史上的邵雍本人，而是以邵雍一脉的观物取象、体用生克、数理成卦方法进行判断。",
	ReasoningFramework: []string{
		"先以体用生克定大方向（用生体吉、体生用泄、用克体凶、体克用可控、比和稳）",
		"再看本卦现状：上下卦象征与卦义",
		"再看互卦过程：内在因素与隐含趋势",
		"再看动爻与变卦：变化焦点与最终走向",
		"如有时间判断，参考卦数推断应期",
	},
	VoiceStyle: []string{
		"文言与白话交融，沉稳通透",
		"善用比喻，不空泛玄谈",
		"语言现代、清楚、直接",
	},
	OutputStructure: []string{
		"卦象总览",
		"本卦解读",
		"互卦分析",
		"动爻与变卦",
		"体用生克",
		"结论与建议",
	},
	Guardrails: []string{
		"必须基于起卦结果解读，不得编造卦象信息",
		"不得自称历史上的邵雍本人",
		"用户问与占卜无关问题时，礼貌引导回占卜话题",
	},
}

var testMeihuaFacts = map[string]string{
	"Method":         "时间起卦（丙午年五月初九子时）",
	"BenGua":         "泽火革（兑上离下）",
	"HuGua":          "天风姤（乾上巽下）",
	"BianGua":        "泽山咸（兑上艮下）",
	"MovingLine":     "初爻动",
	"Ti":             "兑（金）",
	"Yong":           "离（火）",
	"TiYongRelation": "用克体",
}

func TestCompose_Success(t *testing.T) {
	input := InterpretationInput{
		DivinationType: "meihua",
		Question:       "最近事业如何？",
		ResultFacts:    testMeihuaFacts,
		Mode:           "direct",
	}

	msgs, err := Compose(testMeihuaProfile, input)
	if err != nil {
		t.Fatalf("Compose: %v", err)
	}

	// System must contain persona name
	if !strings.Contains(msgs.System, "康节先生") {
		t.Error("system should contain 康节先生")
	}

	// System must contain guardrails about not claiming to be real Shao Yong
	if !strings.Contains(msgs.System, "不得自称历史上的邵雍本人") {
		t.Error("system should contain guardrail about not claiming to be Shao Yong")
	}

	// System must contain reasoning framework
	if !strings.Contains(msgs.System, "推理框架") {
		t.Error("system should contain 推理框架")
	}

	// TiYong must appear before BenGua in reasoning framework (body-use is primary)
	tiYongIdx := strings.Index(msgs.System, "体用生克")
	benGuaIdx := strings.Index(msgs.System, "本卦")
	if tiYongIdx < 0 || benGuaIdx < 0 {
		t.Fatal("system should contain both 体用生克 and 本卦")
	}
	if tiYongIdx > benGuaIdx {
		t.Error("reasoning framework should mention 体用生克 before 本卦")
	}

	// User must contain question
	if !strings.Contains(msgs.User, "最近事业如何？") {
		t.Error("user should contain the question")
	}

	// User must contain all facts
	for key, val := range testMeihuaFacts {
		if !strings.Contains(msgs.User, val) {
			t.Errorf("user should contain fact %s = %s", key, val)
		}
	}

	// User must contain output structure
	if !strings.Contains(msgs.User, "卦象总览") {
		t.Error("user should contain output structure section 卦象总览")
	}

	// No raw JSON dump in output
	if strings.Contains(msgs.System, "{") || strings.Contains(msgs.User, "{") {
		// Small chance of false positive, but generally facts should be formatted
		t.Log("warning: output may contain raw JSON")
	}
}

func TestCompose_ChatMode(t *testing.T) {
	input := InterpretationInput{
		DivinationType: "meihua",
		Question:       "最近事业如何？",
		ResultFacts:    testMeihuaFacts,
		Mode:           "chat",
	}

	msgs, err := Compose(testMeihuaProfile, input)
	if err != nil {
		t.Fatalf("Compose: %v", err)
	}

	// Chat mode should include conversation rules
	if !strings.Contains(msgs.User, "对话规则") {
		t.Error("chat mode user prompt should contain 对话规则")
	}
}

func TestCompose_NilProfile(t *testing.T) {
	input := InterpretationInput{
		DivinationType: "meihua",
		Question:       "test",
	}

	_, err := Compose(nil, input)
	if err == nil {
		t.Fatal("expected error for nil profile")
	}
}

func TestCompose_EmptyDivinationType(t *testing.T) {
	profile := &config.ProfileConfig{
		DivinationType: "",
	}

	input := InterpretationInput{
		DivinationType: "meihua",
		Question:       "test",
	}

	_, err := Compose(profile, input)
	if err == nil {
		t.Fatal("expected error for empty divination type")
	}
}

func TestCompose_TypeMismatch(t *testing.T) {
	input := InterpretationInput{
		DivinationType: "tarot",
		Question:       "test",
	}

	_, err := Compose(testMeihuaProfile, input)
	if err == nil {
		t.Fatal("expected error for type mismatch")
	}
}

func TestCompose_EmptyQuestion(t *testing.T) {
	input := InterpretationInput{
		DivinationType: "meihua",
		Question:       "",
		ResultFacts:    testMeihuaFacts,
	}

	_, err := Compose(testMeihuaProfile, input)
	if err == nil {
		t.Fatal("expected error for empty question")
	}
}

func TestCompose_MissingRequiredFacts(t *testing.T) {
	input := InterpretationInput{
		DivinationType: "meihua",
		Question:       "test",
		ResultFacts: map[string]string{
			"Method": "时间起卦",
			// Missing other required fields
		},
	}

	_, err := Compose(testMeihuaProfile, input)
	if err == nil {
		t.Fatal("expected error for missing required facts")
	}
}

func TestCompose_StableOutput(t *testing.T) {
	input := InterpretationInput{
		DivinationType: "meihua",
		Question:       "最近事业如何？",
		ResultFacts:    testMeihuaFacts,
		Mode:           "direct",
	}

	msgs1, err := Compose(testMeihuaProfile, input)
	if err != nil {
		t.Fatalf("first Compose: %v", err)
	}

	msgs2, err := Compose(testMeihuaProfile, input)
	if err != nil {
		t.Fatalf("second Compose: %v", err)
	}

	if msgs1.System != msgs2.System {
		t.Error("Compose should produce stable system output")
	}
	if msgs1.User != msgs2.User {
		t.Error("Compose should produce stable user output")
	}
}

func TestCompose_NoGlobalProfileDependency(t *testing.T) {
	// Verify Compose uses only the passed-in profile, not any global state
	profile1 := &config.ProfileConfig{
		Version:        "v1",
		DivinationType: "meihua",
		Name:           "Profile A",
		SystemIdentity: "Identity A",
		ReasoningFramework: []string{"Step A1", "Step A2"},
		OutputStructure: []string{"Section A"},
	}
	profile2 := &config.ProfileConfig{
		Version:        "v1",
		DivinationType: "meihua",
		Name:           "Profile B",
		SystemIdentity: "Identity B",
		ReasoningFramework: []string{"Step B1", "Step B2"},
		OutputStructure: []string{"Section B"},
	}

	input := InterpretationInput{
		DivinationType: "meihua",
		Question:       "test",
		ResultFacts:    testMeihuaFacts,
	}

	msgs1, _ := Compose(profile1, input)
	msgs2, _ := Compose(profile2, input)

	if msgs1.System == msgs2.System {
		t.Error("different profiles should produce different system prompts")
	}
	if !strings.Contains(msgs1.System, "Identity A") {
		t.Error("first profile should use Identity A")
	}
	if !strings.Contains(msgs2.System, "Identity B") {
		t.Error("second profile should use Identity B")
	}
}
