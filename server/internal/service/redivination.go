package service

import (
	"regexp"
	"strings"
)

// redivinationPatterns 正面模式：触发重新占卜的关键词。
var redivinationPatterns = []*regexp.Regexp{
	regexp.MustCompile(`重新(占卜|起卦|算|测|卜|来|摇|问)`),
	regexp.MustCompile(`再(占|算|卜|起|测|问|摇).{0,4}(一次|一卦|一下|卦)`),
	regexp.MustCompile(`另起.{0,2}卦`),
	regexp.MustCompile(`起(个|一|新)卦`),
	regexp.MustCompile(`(新|换|另).{0,2}(占|算|卜|测)`),
	regexp.MustCompile(`再(来|试|玩|弄).{0,2}(一次|下)`),
	regexp.MustCompile(`(重新|再).{0,1}(起|占|算|测)`),
	regexp.MustCompile(`数字\s*[\d\s]+`), // 数字起卦模式
	regexp.MustCompile(`用时间起卦`),
	regexp.MustCompile(`按时间起卦`),
	regexp.MustCompile(`用当前时间(起卦)?`),
}

// negativePatterns 负面模式：追问解读，不应触发重新占卜。
var negativePatterns = []*regexp.Regexp{
	regexp.MustCompile(`重新(解读|解释|分析|看|说)`),
	regexp.MustCompile(`换个(角度|说法|思路|方式)`),
	regexp.MustCompile(`详细(说说|解释|分析|讲)`),
	regexp.MustCompile(`再(解释|分析|说|讲|看|读).{0,2}(一下|一遍|一次)`),
	regexp.MustCompile(`具体(说说|解释|分析|讲)`),
	regexp.MustCompile(`展开(说说|讲)`),
	regexp.MustCompile(`(怎么说|怎么看|什么含义|什么意思)`),
	regexp.MustCompile(`能.{0,3}(详细|具体|仔细|再)`),
}

// questionPrefixPattern 问题提取：先去掉重新占卜关键词，再提取问/测后的内容。
var questionCleanPattern = regexp.MustCompile(`(^|[，,。\s])(?:问|测|算|占|卜)(?:一下|一卦|一次)?[：:\s]*(.+)`)

// ReDivinationIntent 表示检测到的重新占卜意图。
type ReDivinationIntent struct {
	ShouldReDivine bool   // 是否需要重新占卜
	Question       string // 用于新占卜的问题，如果为空则使用原始消息
}

// ReDivinationDetector 检测消息中是否包含重新占卜意图。
type ReDivinationDetector struct{}

// Detect 检测消息中是否包含重新占卜意图。
// 仅对 meihua 类型进行检测；其他类型直接返回不占卜。
func (d *ReDivinationDetector) Detect(divinationType string, message string) *ReDivinationIntent {
	if divinationType != "meihua" {
		return &ReDivinationIntent{ShouldReDivine: false}
	}

	normalized := strings.TrimSpace(message)
	if normalized == "" {
		return &ReDivinationIntent{ShouldReDivine: false}
	}

	// 先检查负面模式（追问解读）
	for _, p := range negativePatterns {
		if p.MatchString(normalized) {
			return &ReDivinationIntent{ShouldReDivine: false}
		}
	}

	// 再检查正面模式（重新占卜）
	matched := false
	for _, p := range redivinationPatterns {
		if p.MatchString(normalized) {
			matched = true
			break
		}
	}

	if !matched {
		return &ReDivinationIntent{ShouldReDivine: false}
	}

	// 提取问题
	question := extractQuestion(normalized)

	return &ReDivinationIntent{
		ShouldReDivine: true,
		Question:       question,
	}
}

// extractQuestion 从消息中提取占卜问题。
func extractQuestion(message string) string {
	// Step 1: 去除已知的重新占卜关键词，避免它们干扰问题提取
	cleaned := message
	cleaned = regexp.MustCompile(`重新(占卜|起卦|算|测|卜|来|摇|问)`).ReplaceAllString(cleaned, "")
	cleaned = regexp.MustCompile(`再(占|算|卜|起|测|问|摇).{0,4}(一次|一卦|一下|卦)`).ReplaceAllString(cleaned, "")
	cleaned = regexp.MustCompile(`另起.{0,2}卦`).ReplaceAllString(cleaned, "")
	cleaned = regexp.MustCompile(`起(个|一|新)卦`).ReplaceAllString(cleaned, "")
	cleaned = regexp.MustCompile(`(?:用|按|当前)时间起卦`).ReplaceAllString(cleaned, "")
	cleaned = regexp.MustCompile(`数字\s*[\d\s]+`).ReplaceAllString(cleaned, "")

	// Step 2: 在清理后的文本中查找 "问/测/算 X" 模式
	matches := questionCleanPattern.FindStringSubmatch(cleaned)
	if len(matches) >= 3 {
		q := strings.TrimSpace(matches[2])
		if q != "" {
			return q
		}
	}

	// Step 3: 去除多余标点，返回清理后的文本
	cleaned = regexp.MustCompile(`^[，,。；;！!、\s]+|[，,。；;！!、\s]+$`).ReplaceAllString(cleaned, "")
	cleaned = regexp.MustCompile(`[，,。；;！!、\s]+`).ReplaceAllString(cleaned, " ")
	cleaned = strings.TrimSpace(cleaned)

	if cleaned == "" {
		return message
	}
	return cleaned
}
