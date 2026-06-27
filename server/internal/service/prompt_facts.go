package service

import "fmt"

// FactParser 从占卜结果 JSON 中提取关键事实。
// 单方法接口以 er 结尾命名，符合 Go 惯例。
type FactParser interface {
	Parse(resultJSON, question string) (map[string]string, error)
}

// factParsers 是各占卜类型的 facts 解析器注册表。
var factParsers = map[string]func(resultJSON, question string) (map[string]string, error){
	"meihua": parseMeihuaFacts,
}

// ParseFacts 根据占卜类型解析结果 JSON 为结构化事实。
func ParseFacts(divinationType, resultJSON, question string) (map[string]string, error) {
	parser, ok := factParsers[divinationType]
	if !ok {
		return nil, fmt.Errorf("no fact parser for type: %s", divinationType)
	}
	return parser(resultJSON, question)
}
