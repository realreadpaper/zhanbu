package service

import (
	"fmt"
	"strings"

	"zhanbu/config"
)

// InterpretationInput 是单次占卜解读的输入数据。
type InterpretationInput struct {
	DivinationType string            `json:"divination_type"`
	Question       string            `json:"question"`
	ResultJSON     string            `json:"result_json"`
	ResultFacts    map[string]string `json:"result_facts"`
	Mode           string            `json:"mode"` // "direct" or "chat"
}

// Messages 是 Compose 的输出，包含 system 和 user 两条消息。
type Messages struct {
	System string `json:"system"`
	User   string `json:"user"`
}

// Compose 根据人格配置和解读输入，组装 AI 请求所需的 system/user 消息。
// 它是纯函数：不访问数据库、全局状态、网络或日志。
func Compose(profile *config.ProfileConfig, input InterpretationInput) (*Messages, error) {
	if profile == nil {
		return nil, fmt.Errorf("profile is nil")
	}
	if profile.DivinationType == "" {
		return nil, fmt.Errorf("profile divination type is empty")
	}
	if profile.DivinationType != input.DivinationType {
		return nil, fmt.Errorf("profile type %s does not match input type %s", profile.DivinationType, input.DivinationType)
	}
	if strings.TrimSpace(input.Question) == "" {
		return nil, fmt.Errorf("question is empty")
	}

	// Validate required facts for meihua
	if input.DivinationType == "meihua" {
		required := []string{"Method", "BenGua", "HuGua", "BianGua", "MovingLine", "Ti", "Yong", "TiYongRelation"}
		for _, key := range required {
			if strings.TrimSpace(input.ResultFacts[key]) == "" {
				return nil, fmt.Errorf("missing required fact: %s", key)
			}
		}
	}

	system := buildSystem(profile)
	user := buildUser(profile, input)
	return &Messages{System: system, User: user}, nil
}

// buildSystem 构建 system prompt：身份、推理框架、表达风格、边界规则。
func buildSystem(profile *config.ProfileConfig) string {
	var b strings.Builder

	// System identity
	b.WriteString(strings.TrimSpace(profile.SystemIdentity))
	b.WriteString("\n\n")

	// Reasoning framework
	if len(profile.ReasoningFramework) > 0 {
		b.WriteString("【推理框架】\n")
		for i, step := range profile.ReasoningFramework {
			fmt.Fprintf(&b, "%d. %s\n", i+1, step)
		}
		b.WriteString("\n")
	}

	// Voice style
	if len(profile.VoiceStyle) > 0 {
		b.WriteString("【表达风格】\n")
		for _, style := range profile.VoiceStyle {
			fmt.Fprintf(&b, "- %s\n", style)
		}
		b.WriteString("\n")
	}

	// Guardrails
	if len(profile.Guardrails) > 0 {
		b.WriteString("【对话规则】\n")
		for i, rule := range profile.Guardrails {
			fmt.Fprintf(&b, "%d. %s\n", i+1, rule)
		}
		fmt.Fprintf(&b, "%d. 回答要简洁明了，每次回复控制在 200-400 字\n", len(profile.Guardrails)+1)
		fmt.Fprintf(&b, "%d. 使用 Markdown 格式让回答更清晰（如加粗、列表等）\n", len(profile.Guardrails)+2)
	}

	return b.String()
}

// buildUser 构建 user prompt：问题、占卜事实、输出结构。
func buildUser(profile *config.ProfileConfig, input InterpretationInput) string {
	var b strings.Builder

	// Question
	fmt.Fprintf(&b, "【用户问题】\n%s\n\n", input.Question)

	// Result facts
	if len(input.ResultFacts) > 0 {
		b.WriteString("【卦象信息】\n")
		for _, key := range factDisplayOrder(input.DivinationType) {
			if val, ok := input.ResultFacts[key]; ok && val != "" {
				fmt.Fprintf(&b, "%s：%s\n", factDisplayName(key), val)
			}
		}
		b.WriteString("\n")
	}

	// Output structure
	if len(profile.OutputStructure) > 0 {
		b.WriteString("请按以下结构输出：\n")
		for i, section := range profile.OutputStructure {
			fmt.Fprintf(&b, "%d. **%s**\n", i+1, section)
		}
	}

	// Mode-specific instructions
	if input.Mode == "chat" {
		b.WriteString("\n【对话规则】\n")
		b.WriteString("1. 始终以上述身份回答\n")
		b.WriteString("2. 可以深入分析某一卦、某一爻的具体含义\n")
		b.WriteString("3. 可以结合用户的具体生活情况进行个性化解读\n")
	}

	return b.String()
}

// factDisplayOrder 返回不同占卜类型的事实显示顺序。
func factDisplayOrder(divinationType string) []string {
	switch divinationType {
	case "meihua":
		return []string{"Method", "BenGua", "HuGua", "BianGua", "MovingLine", "Ti", "Yong", "TiYongRelation"}
	default:
		return nil
	}
}

// factDisplayName 返回事实字段的中文显示名。
func factDisplayName(key string) string {
	names := map[string]string{
		"Method":         "起卦方式",
		"BenGua":         "本卦",
		"HuGua":          "互卦",
		"BianGua":        "变卦",
		"MovingLine":     "动爻",
		"Ti":             "体卦",
		"Yong":           "用卦",
		"TiYongRelation": "体用关系",
	}
	if name, ok := names[key]; ok {
		return name
	}
	return key
}
