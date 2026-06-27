package service

import "testing"

func TestReDivinationDetector_PositivePatterns(t *testing.T) {
	detector := &ReDivinationDetector{}

	tests := []struct {
		message string
		desc    string
	}{
		{"重新占卜", "直接命令"},
		{"重新起卦", "直接命令-起卦"},
		{"再占一次", "再占一次"},
		{"再算一卦", "再算一卦"},
		{"起新卦", "起新卦"},
		{"帮我再算一下", "帮我再算"},
		{"重新算", "重新算"},
		{"再来一次", "再来一次"},
		{"数字 5 8", "数字起卦"},
		{"数字 12 34 56，问事业", "数字起卦带问题"},
		{"用时间起卦", "时间起卦"},
		{"按时间起卦", "按时间起卦"},
		{"重新占卜，问感情如何", "带问题"},
		{"再占一卦，测一下事业运", "带问题-测"},
		{"起新卦 问一下财运", "起新卦带问题"},
		{"重新摇一卦", "重新摇卦"},
		{"另起一卦", "另起一卦"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			intent := detector.Detect("meihua", tt.message)
			if !intent.ShouldReDivine {
				t.Errorf("expected ShouldReDivine=true for %q (%s)", tt.message, tt.desc)
			}
		})
	}
}

func TestReDivinationDetector_NegativePatterns(t *testing.T) {
	detector := &ReDivinationDetector{}

	tests := []struct {
		message string
		desc    string
	}{
		{"重新解读一下", "重新解读"},
		{"换个角度分析", "换个角度"},
		{"详细说说这个卦", "详细说说"},
		{"再解释一下变卦的含义", "再解释"},
		{"具体分析一下互卦", "具体分析"},
		{"怎么看这个结果", "怎么看"},
		{"能详细说说吗", "能详细说说"},
		{"那我的财运呢", "追问财运"},
		{"感情方面怎么样", "追问感情"},
		{"展开讲讲", "展开讲讲"},
		{"重新分析", "重新分析"},
		{"能不能再分析一下", "再分析"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			intent := detector.Detect("meihua", tt.message)
			if intent.ShouldReDivine {
				t.Errorf("expected ShouldReDivine=false for %q (%s)", tt.message, tt.desc)
			}
		})
	}
}

func TestReDivinationDetector_NonMeihuaType(t *testing.T) {
	detector := &ReDivinationDetector{}

	// Should not detect re-divination for non-meihua types even with keywords
	if detector.Detect("tarot", "重新占卜").ShouldReDivine {
		t.Error("should not detect re-divination for tarot type")
	}
	if detector.Detect("liuyao_v2", "再算一卦").ShouldReDivine {
		t.Error("should not detect re-divination for liuyao_v2 type")
	}
	if detector.Detect("bazi", "重新起卦").ShouldReDivine {
		t.Error("should not detect re-divination for bazi type")
	}
}

func TestReDivinationDetector_EmptyMessage(t *testing.T) {
	detector := &ReDivinationDetector{}
	intent := detector.Detect("meihua", "")
	if intent.ShouldReDivine {
		t.Error("should not detect re-divination for empty message")
	}
}

func TestReDivinationDetector_QuestionExtraction(t *testing.T) {
	detector := &ReDivinationDetector{}

	tests := []struct {
		message string
		wantQ   string
	}{
		{"重新占卜，问事业如何", "事业如何"},
		{"再算一卦，测一下感情运", "感情运"},
		{"起新卦问财运", "财运"},
		{"重新占卜感情方面", "感情方面"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			intent := detector.Detect("meihua", tt.message)
			if !intent.ShouldReDivine {
				t.Errorf("expected ShouldReDivine=true for %q", tt.message)
				return
			}
			if intent.Question != tt.wantQ {
				t.Errorf("expected question %q, got %q", tt.wantQ, intent.Question)
			}
		})
	}
}
