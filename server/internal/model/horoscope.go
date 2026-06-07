package model

// ZodiacInfo holds static zodiac data.
type ZodiacInfo struct {
	Name    string `json:"name"`    // e.g. "aries"
	NameCN  string `json:"name_cn"` // e.g. "白羊座"
	Element string `json:"element"` // fire/earth/air/water
	Planet  string `json:"planet"`  // ruling planet
}

// Zodiacs is the ordered list of all 12 zodiac signs.
var Zodiacs = []ZodiacInfo{
	{Name: "aries", NameCN: "白羊座", Element: "fire", Planet: "火星"},
	{Name: "taurus", NameCN: "金牛座", Element: "earth", Planet: "金星"},
	{Name: "gemini", NameCN: "双子座", Element: "air", Planet: "水星"},
	{Name: "cancer", NameCN: "巨蟹座", Element: "water", Planet: "月亮"},
	{Name: "leo", NameCN: "狮子座", Element: "fire", Planet: "太阳"},
	{Name: "virgo", NameCN: "处女座", Element: "earth", Planet: "水星"},
	{Name: "libra", NameCN: "天秤座", Element: "air", Planet: "金星"},
	{Name: "scorpio", NameCN: "天蝎座", Element: "water", Planet: "冥王星"},
	{Name: "sagittarius", NameCN: "射手座", Element: "fire", Planet: "木星"},
	{Name: "capricorn", NameCN: "摩羯座", Element: "earth", Planet: "土星"},
	{Name: "aquarius", NameCN: "水瓶座", Element: "air", Planet: "天王星"},
	{Name: "pisces", NameCN: "双鱼座", Element: "water", Planet: "海王星"},
}

// HoroscopeResult is the complete horoscope response.
type HoroscopeResult struct {
	RecordID    uint            `json:"record_id,omitempty"`
	Zodiac      string          `json:"zodiac"`
	ZodiacCN    string          `json:"zodiac_cn"`
	Period      string          `json:"period"`
	Date        string          `json:"date"`
	Overall     int             `json:"overall"`
	Love        int             `json:"love"`
	Career      int             `json:"career"`
	Wealth      int             `json:"wealth"`
	Health      int             `json:"health"`
	LuckyNumber int             `json:"lucky_number"`
	LuckyColor  string          `json:"lucky_color"`
	Summary     string          `json:"summary"`
	Detail      HoroscopeDetail `json:"detail"`
}

// HoroscopeDetail holds dimension-specific text.
type HoroscopeDetail struct {
	Love   string `json:"love"`
	Career string `json:"career"`
	Wealth string `json:"wealth"`
	Health string `json:"health"`
}

// ZodiacLookup looks up a zodiac by English name (case-insensitive).
// Returns nil if not found.
func ZodiacLookup(name string) *ZodiacInfo {
	for i := range Zodiacs {
		if Zodiacs[i].Name == name {
			return &Zodiacs[i]
		}
	}
	return nil
}
