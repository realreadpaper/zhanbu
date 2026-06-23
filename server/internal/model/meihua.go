package model

// MeiHuaTrigram represents a trigram in MeiHua divination.
type MeiHuaTrigram struct {
	Number  int    `json:"number"`   // 先天八卦数 1-8
	Name    string `json:"name"`     // 卦名：乾兑离震巽坎艮坤
	Symbol  string `json:"symbol"`   // 卦符：☰☱☲☳☴☵☶☷
	Element string `json:"element"`  // 五行
	YinYang string `json:"yin_yang"` // 阴阳
}

// MeiHuaHexagram represents a hexagram composed of upper and lower trigrams.
type MeiHuaHexagram struct {
	Upper     MeiHuaTrigram `json:"upper"`
	Lower     MeiHuaTrigram `json:"lower"`
	Name      string        `json:"name"`       // 卦名如"泽火革"
	NameShort string        `json:"name_short"` // 简称如"革"
}

// MeiHuaTiYong represents the ti-yong (体用) analysis.
type MeiHuaTiYong struct {
	Ti       MeiHuaTrigram `json:"ti"`        // 体卦
	Yong     MeiHuaTrigram `json:"yong"`      // 用卦
	TiHex    string        `json:"ti_hex"`    // 体卦所在：upper or lower
	YongHex  string        `json:"yong_hex"`  // 用卦所在：upper or lower
	Relation string        `json:"relation"`  // 生克关系：体克用/用克体/体生用/用生体/比和
}

// MeiHuaSourceValues stores the raw input values used for calculation.
type MeiHuaSourceValues struct {
	Method       string `json:"method"`        // "time" or "number"
	YearBranch   string `json:"year_branch"`   // 地支名
	YearNumber   int    `json:"year_number"`   // 地支序数
	LunarMonth   int    `json:"lunar_month"`   // 农历月
	LunarDay     int    `json:"lunar_day"`     // 农历日
	HourBranch   string `json:"hour_branch"`   // 时辰地支
	HourNumber   int    `json:"hour_number"`   // 时辰序数
	Numbers      []int  `json:"numbers,omitempty"` // 数字起卦时的原始数字
}

// MeiHuaResult represents the full MeiHua divination result.
type MeiHuaResult struct {
	RecordID      uint              `json:"record_id,omitempty"`
	Question      string            `json:"question,omitempty"`
	Method        string            `json:"method"`         // "time" or "number"
	SourceValues  MeiHuaSourceValues `json:"source_values"`
	UpperTrigram  MeiHuaTrigram     `json:"upper_trigram"`
	LowerTrigram  MeiHuaTrigram     `json:"lower_trigram"`
	BenGua        MeiHuaHexagram    `json:"ben_gua"`
	HuGua         MeiHuaHexagram    `json:"hu_gua"`
	BianGua       MeiHuaHexagram    `json:"bian_gua"`
	MovingLine    int               `json:"moving_line"`    // 1-6
	TiYong        MeiHuaTiYong      `json:"ti_yong"`
	Timestamp     string            `json:"timestamp"`
}

// MeiHuaTimeRequest is the request for time-based MeiHua divination.
type MeiHuaTimeRequest struct {
	Question  string `json:"question"`
	Timezone  string `json:"timezone"` // e.g. "Asia/Shanghai"
}

// MeiHuaNumberRequest is the request for number-based MeiHua divination.
type MeiHuaNumberRequest struct {
	Question string `json:"question"`
	Numbers  []int  `json:"numbers"` // user-provided numbers
}
