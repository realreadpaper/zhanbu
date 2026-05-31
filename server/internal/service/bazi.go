package service

import (
	"fmt"
	"strconv"
	"strings"

	"zhanbu/internal/model"
)

// BaZiService handles bazi (八字) calculation logic.
type BaZiService struct{}

// NewBaZiService creates a new BaZiService.
func NewBaZiService() *BaZiService {
	return &BaZiService{}
}

// TianGan (Heavenly Stems) in order.
var tianGan = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}

// DiZhi (Earthly Branches) in order.
var diZhi = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

// WuXing of each TianGan.
var tianGanWuXing = []string{"木", "木", "火", "火", "土", "土", "金", "金", "水", "水"}

// WuXing of each DiZhi.
var diZhiWuXing = []string{"水", "土", "木", "木", "土", "火", "火", "土", "金", "金", "土", "水"}

// NaYin (纳音) for 60 Jiazi cycles.
var naYin = []string{
	"海中金", "海中金", "炉中火", "炉中火", "大林木", "大林木",
	"路旁土", "路旁土", "剑锋金", "剑锋金", "山头火", "山头火",
	"涧下水", "涧下水", "城头土", "城头土", "白蜡金", "白蜡金",
	"杨柳木", "杨柳木", "泉中水", "泉中水", "屋上土", "屋上土",
	"霹雳火", "霹雳火", "松柏木", "松柏木", "长流水", "长流水",
	"砂石金", "砂石金", "山下火", "山下火", "平地木", "平地木",
	"壁上土", "壁上土", "金箔金", "金箔金", "覆灯火", "覆灯火",
	"天河水", "天河水", "大驿土", "大驿土", "钗环金", "钗环金",
	"桑柘木", "桑柘木", "大溪水", "大溪水", "沙中土", "沙中土",
	"天上火", "天上火", "石榴木", "石榴木", "大海水", "大海水",
}

// HiddenGan maps each DiZhi to its hidden heavenly stems.
var hiddenGan = map[string][]string{
	"子": {"癸"},
	"丑": {"己", "癸", "辛"},
	"寅": {"甲", "丙", "戊"},
	"卯": {"乙"},
	"辰": {"戊", "乙", "癸"},
	"巳": {"丙", "庚", "戊"},
	"午": {"丁", "己"},
	"未": {"己", "丁", "乙"},
	"申": {"庚", "壬", "戊"},
	"酉": {"辛"},
	"戌": {"戊", "辛", "丁"},
	"亥": {"壬", "甲"},
}

// TenGodNames maps relationship to ten-god name.
// Key: (dayMaster_wuXing, other_wuXing) relationship.
// 生: 相生关系; 克: 相克关系
func tenGodName(dayMasterWX, otherWX string) string {
	if dayMasterWX == otherWX {
		return "比肩"
	}

	// 相生: 我生=食伤, 生我=印
	sheng := map[string]string{
		"木": "火", "火": "土", "土": "金", "金": "水", "水": "木",
	}
	// 相克: 我克=财, 克我=官
	ke := map[string]string{
		"木": "土", "火": "金", "土": "水", "金": "木", "水": "火",
	}

	if sheng[dayMasterWX] == otherWX {
		return "食神" // 我生同性
	}
	if sheng[otherWX] == dayMasterWX {
		return "偏印" // 生我同性
	}
	if ke[dayMasterWX] == otherWX {
		return "偏财" // 我克同性
	}
	if ke[otherWX] == dayMasterWX {
		return "偏官" // 克我同性
	}

	return "比肩"
}

// Calculate computes the bazi pillars from solar calendar inputs.
func (s *BaZiService) Calculate(birthDate, birthTime, gender string) (*model.BaZiResult, error) {
	// Parse date: expect "YYYY-MM-DD"
	parts := strings.Split(birthDate, "-")
	if len(parts) != 3 {
		return nil, fmt.Errorf("日期格式错误，应为 YYYY-MM-DD")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil || year < 1900 || year > 2100 {
		return nil, fmt.Errorf("无效的年份: %s", parts[0])
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil || month < 1 || month > 12 {
		return nil, fmt.Errorf("无效的月份: %s", parts[1])
	}
	day, err := strconv.Atoi(parts[2])
	if err != nil || day < 1 || day > 31 {
		return nil, fmt.Errorf("无效的日期: %s", parts[2])
	}

	// Parse time: expect "HH:MM" or "HH:MM:SS"
	timeParts := strings.Split(birthTime, ":")
	if len(timeParts) < 2 || len(timeParts) > 3 {
		return nil, fmt.Errorf("时间格式错误，应为 HH:MM")
	}
	hour, err := strconv.Atoi(timeParts[0])
	if err != nil || hour < 0 || hour > 23 {
		return nil, fmt.Errorf("无效的小时: %s", timeParts[0])
	}

	// Calculate year pillar
	yearTG, yearDZ := s.calcYearPillar(year)

	// Calculate month pillar
	monthTG, monthDZ := s.calcMonthPillar(year, month)

	// Calculate day pillar
	dayTG, dayDZ := s.calcDayPillar(year, month, day)

	// Calculate hour pillar
	hourTG, hourDZ := s.calcHourPillar(hour, dayTG)

	// Build pillars
	pillars := model.BaZiPillars{
		Year:  s.buildPillar(yearTG, yearDZ),
		Month: s.buildPillar(monthTG, monthDZ),
		Day:   s.buildPillar(dayTG, dayDZ),
		Hour:  s.buildPillar(hourTG, hourDZ),
	}

	// Five element analysis
	fiveElements := s.analyzeFiveElements(pillars, dayTG)

	// Ten gods analysis
	tenGods := s.analyzeTenGods(pillars, dayTG)

	result := &model.BaZiResult{
		Birth: model.BaZiBirthInfo{
			Solar: birthDate + " " + birthTime,
			Lunar: s.solarToLunarApprox(year, month, day),
		},
		Pillars:      pillars,
		FiveElements: fiveElements,
		TenGods:      tenGods,
	}

	_ = gender // gender can be used for additional analysis later

	return result, nil
}

// calcYearPillar calculates the year pillar (年柱) from the solar year.
// Formula: (year - 4) % 10 for TianGan, (year - 4) % 12 for DiZhi.
func (s *BaZiService) calcYearPillar(year int) (string, string) {
	tgIdx := (year - 4) % 10
	if tgIdx < 0 {
		tgIdx += 10
	}
	dzIdx := (year - 4) % 12
	if dzIdx < 0 {
		dzIdx += 12
	}
	return tianGan[tgIdx], diZhi[dzIdx]
}

// calcMonthPillar calculates the month pillar (月柱).
// The month is based on the solar term (节气). For simplicity, we approximate:
// Month 1 (寅月) starts around Feb 4, Month 2 (卯月) around Mar 6, etc.
// TianGan of month depends on year's TianGan.
func (s *BaZiService) calcMonthPillar(year, month int) (string, string) {
	// Adjust month for solar terms: month 1 (Jan) is mostly 丑月 (previous year's 12th branch)
	// Approximate: Feb=寅(2), Mar=卯(3), ..., Jan=丑(1 of previous year)
	dzMonth := month + 1 // Feb=3(寅), Mar=4(卯), ...
	if month == 1 {
		dzMonth = 13 // 丑月
	}
	if dzMonth > 12 {
		dzMonth = dzMonth - 12
	}
	dzIdx := dzMonth % 12 // 0=子, 1=丑, 2=寅, ...

	// Month DiZhi: 寅=2 is the first month
	// month 1(Jan)->丑(1), month 2(Feb)->寅(2), month 3(Mar)->卯(3), ...
	dzIdx = (month + 1) % 12
	if month == 1 {
		dzIdx = 1 // 丑
	}

	// TianGan of month: depends on year's TianGan
	// 甲己之年丙作首 (yearTG 甲/己 -> month 1 starts with 丙)
	yearTG, _ := s.calcYearPillar(year)
	yearTGIdx := indexOf(tianGan, yearTG)

	// Starting TianGan for month 1 (寅月):
	// 甲/己(0,5) -> 丙(2), 乙/庚(1,6) -> 戊(4), 丙/辛(2,7) -> 庚(6), 丁/壬(3,8) -> 壬(8), 戊/癸(4,9) -> 甲(0)
	startTG := []int{2, 4, 6, 8, 0}
	baseTG := startTG[yearTGIdx%5]

	// Month offset: month 2(寅)=0, month 3(卯)=1, ...
	monthOffset := month - 2
	if month == 1 {
		monthOffset = 11 // 丑月 is the 12th month in this cycle
	}
	tgIdx := (baseTG + monthOffset) % 10

	return tianGan[tgIdx], diZhi[dzIdx]
}

// calcDayPillar calculates the day pillar (日柱) using a known reference date.
// Reference: 1900-01-01 is 甲子日 (index 0 in the 60 Jiazi cycle).
// Actually, 1900-01-01 = 甲戌日 (index 10). We'll use a more accurate reference.
// Reference: 2000-01-01 = 甲午日 (index 30).
func (s *BaZiService) calcDayPillar(year, month, day int) (string, string) {
	// Calculate days from a reference date
	// Reference: 1900-01-01 = 庚子日 (Jiazi index 36)
	// We use the Julian Day Number approach for accuracy

	refYear, refMonth, refDay := 1900, 1, 1
	refJDN := s.gregorianToJD(refYear, refMonth, refDay)
	targetJDN := s.gregorianToJD(year, month, day)
	diff := targetJDN - refJDN

	// 1900-01-01 is 庚子日: 天干=庚(6), 地支=子(0), 六十甲子序号=36
	refIdx := 36
	jiaziIdx := (refIdx + diff) % 60
	if jiaziIdx < 0 {
		jiaziIdx += 60
	}

	tgIdx := jiaziIdx % 10
	dzIdx := jiaziIdx % 12

	return tianGan[tgIdx], diZhi[dzIdx]
}

// gregorianToJD converts a Gregorian date to Julian Day Number.
func (s *BaZiService) gregorianToJD(year, month, day int) int {
	// Simplified Julian Day Number calculation
	if month <= 2 {
		year--
		month += 12
	}
	A := year / 100
	B := 2 - A + A/4
	return int(365.25*float64(year+4716)) + int(30.6001*float64(month+1)) + day + B - 1524
}

// calcHourPillar calculates the hour pillar (时柱).
// DiZhi of hour: 子(23-1), 丑(1-3), 寅(3-5), ..., 亥(21-23).
// TianGan of hour depends on day's TianGan.
func (s *BaZiService) calcHourPillar(hour int, dayTG string) (string, string) {
	// Map hour to DiZhi index
	dzIdx := ((hour + 1) / 2) % 12
	if hour == 23 {
		dzIdx = 0 // 子时 (23:00-01:00)
	}

	// TianGan of hour: depends on day's TianGan
	// 甲己还加甲 (dayTG 甲/己 -> 子时 starts with 甲)
	dayTGIdx := indexOf(tianGan, dayTG)
	startTG := []int{0, 2, 4, 6, 8}
	baseTG := startTG[dayTGIdx%5]

	tgIdx := (baseTG + dzIdx) % 10

	return tianGan[tgIdx], diZhi[dzIdx]
}

// buildPillar constructs a BaZiPillar from TianGan and DiZhi.
func (s *BaZiService) buildPillar(tg, dz string) model.BaZiPillar {
	wx := s.getPillarWuXing(tg, dz)
	ny := s.getNaYin(tg, dz)
	hg := hiddenGan[dz]

	return model.BaZiPillar{
		TianGan:   tg,
		DiZhi:     dz,
		WuXing:    wx,
		NaYin:     ny,
		HiddenGan: hg,
	}
}

// getPillarWuXing determines the WuXing of a pillar based on its TianGan.
func (s *BaZiService) getPillarWuXing(tg, dz string) string {
	tgIdx := indexOf(tianGan, tg)
	if tgIdx >= 0 && tgIdx < len(tianGanWuXing) {
		return tianGanWuXing[tgIdx]
	}
	return "未知"
}

// getNaYin returns the NaYin for a Jiazi combination.
func (s *BaZiService) getNaYin(tg, dz string) string {
	tgIdx := indexOf(tianGan, tg)
	dzIdx := indexOf(diZhi, dz)
	if tgIdx < 0 || dzIdx < 0 {
		return "未知"
	}
	jiaziIdx := (tgIdx*6 + dzIdx*5) % 60
	if jiaziIdx < 0 || jiaziIdx >= len(naYin) {
		return "未知"
	}
	return naYin[jiaziIdx]
}

// analyzeFiveElements performs five-element analysis.
func (s *BaZiService) analyzeFiveElements(pillars model.BaZiPillars, dayTG string) model.FiveElementAnalysis {
	counts := map[string]int{
		"金": 0, "木": 0, "水": 0, "火": 0, "土": 0,
	}

	// Count from all pillars (TianGan + DiZhi)
	for _, tg := range []string{pillars.Year.TianGan, pillars.Month.TianGan, pillars.Day.TianGan, pillars.Hour.TianGan} {
		wx := tianGanWuXing[indexOf(tianGan, tg)]
		counts[wx]++
	}
	for _, dz := range []string{pillars.Year.DiZhi, pillars.Month.DiZhi, pillars.Day.DiZhi, pillars.Hour.DiZhi} {
		wx := diZhiWuXing[indexOf(diZhi, dz)]
		counts[wx]++
	}

	// Find strongest and weakest
	strongest, weakest := "金", "金"
	maxCount, minCount := 0, 999
	for wx, cnt := range counts {
		if cnt > maxCount {
			maxCount = cnt
			strongest = wx
		}
		if cnt < minCount {
			minCount = cnt
			weakest = wx
		}
	}

	// Determine Day Master element
	dayMasterWX := tianGanWuXing[indexOf(tianGan, dayTG)]

	// Determine strength (simplified)
	strength := "中和"
	if counts[dayMasterWX] >= 3 {
		strength = "偏强"
	} else if counts[dayMasterWX] <= 1 {
		strength = "偏弱"
	}

	// YongShen (用神) - simplified: use the element that controls the day master when strong
	// or supports it when weak
	sheng := map[string]string{
		"木": "火", "火": "土", "土": "金", "金": "水", "水": "木",
	}
	ke := map[string]string{
		"木": "土", "火": "金", "土": "水", "金": "木", "水": "火",
	}

	var yongShen, jiShen string
	if strength == "偏强" {
		yongShen = ke[dayMasterWX] // 泄或克
		jiShen = sheng[dayMasterWX]
	} else {
		yongShen = sheng[dayMasterWX] // 生扶
		jiShen = ke[dayMasterWX]
	}

	return model.FiveElementAnalysis{
		Metal:     counts["金"],
		Wood:      counts["木"],
		Water:     counts["水"],
		Fire:      counts["火"],
		Earth:     counts["土"],
		Strongest: strongest,
		Weakest:   weakest,
		DayMaster: dayMasterWX,
		Strength:  strength,
		YongShen:  yongShen,
		JiShen:    jiShen,
	}
}

// analyzeTenGods determines the Ten Gods for each pillar relative to the Day Master.
func (s *BaZiService) analyzeTenGods(pillars model.BaZiPillars, dayTG string) []model.TenGod {
	dayMasterWX := tianGanWuXing[indexOf(tianGan, dayTG)]

	gods := make([]model.TenGod, 0, 8)

	positions := []struct {
		position string
		tg       string
	}{
		{"年干", pillars.Year.TianGan},
		{"月干", pillars.Month.TianGan},
		{"时干", pillars.Hour.TianGan},
		// Also include DiZhi's main hidden Gan
		{"年支主气", pillars.Year.HiddenGan[0]},
		{"月支主气", pillars.Month.HiddenGan[0]},
		{"日支主气", pillars.Day.HiddenGan[0]},
		{"时支主气", pillars.Hour.HiddenGan[0]},
	}

	for _, p := range positions {
		otherWX := tianGanWuXing[indexOf(tianGan, p.tg)]
		god := tenGodName(dayMasterWX, otherWX)
		gods = append(gods, model.TenGod{
			Position: p.position,
			TianGan:  p.tg,
			God:      god,
		})
	}

	return gods
}

// solarToLunarApprox provides an approximate lunar date string.
// This is a simplified conversion; for production use, a proper lunar calendar library should be used.
func (s *BaZiService) solarToLunarApprox(year, month, day int) string {
	// Simple approximation: just format as a descriptive string
	// In production, you would use a proper solar-to-lunar conversion library
	animals := []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
	animal := animals[(year-4)%12]
	return fmt.Sprintf("农历%s年（%s年）%d月%d日", animal, animal, month, day)
}

// indexOf returns the index of a string in a slice, or -1 if not found.
func indexOf(slice []string, target string) int {
	for i, s := range slice {
		if s == target {
			return i
		}
	}
	return -1
}
