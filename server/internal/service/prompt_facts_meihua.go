package service

import (
	"encoding/json"
	"fmt"

	"zhanbu/internal/model"
)

// parseMeihuaFacts 从梅花易数结果 JSON 中提取关键事实。
func parseMeihuaFacts(resultJSON, question string) (map[string]string, error) {
	var reading struct {
		Method     string               `json:"method"`
		BenGua     model.MeiHuaHexagram `json:"ben_gua"`
		HuGua      model.MeiHuaHexagram `json:"hu_gua"`
		BianGua    model.MeiHuaHexagram `json:"bian_gua"`
		MovingLine int                  `json:"moving_line"`
		TiYong     model.MeiHuaTiYong   `json:"ti_yong"`
		SourceValues struct {
			YearBranch   string `json:"year_branch"`
			LunarMonth   int    `json:"lunar_month"`
			LunarDay     int    `json:"lunar_day"`
			LunarDisplay string `json:"lunar_display"`
			HourBranch   string `json:"hour_branch"`
			Numbers      []int  `json:"numbers"`
		} `json:"source_values"`
	}
	if err := json.Unmarshal([]byte(resultJSON), &reading); err != nil {
		return nil, fmt.Errorf("parse meihua result: %w", err)
	}

	methodDesc := formatMeihuaMethodDesc(reading.Method, resultJSON)

	facts := map[string]string{
		"Method":         methodDesc,
		"BenGua":         formatMeihuaHexagram(reading.BenGua),
		"HuGua":          formatMeihuaHexagram(reading.HuGua),
		"BianGua":        formatMeihuaHexagram(reading.BianGua),
		"MovingLine":     formatMovingLine(reading.MovingLine),
		"Ti":             fmt.Sprintf("%s（%s）", reading.TiYong.Ti.Name, reading.TiYong.Ti.Element),
		"Yong":           fmt.Sprintf("%s（%s）", reading.TiYong.Yong.Name, reading.TiYong.Yong.Element),
		"TiYongRelation": reading.TiYong.Relation,
	}

	return facts, nil
}
