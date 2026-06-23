package service

import (
	"fmt"
	"time"

	"github.com/6tail/lunar-go/calendar"
	"zhanbu/internal/model"
)

// ── 先天八卦常量 ──────────────────────────────────────────────────

const (
	TrigramQian  = 1 // 乾
	TrigramDui   = 2 // 兑
	TrigramLi    = 3 // 离
	TrigramZhen  = 4 // 震
	TrigramXun   = 5 // 巽
	TrigramKan   = 6 // 坎
	TrigramGen   = 7 // 艮
	TrigramKun   = 8 // 坤
)

// TrigramCount 八卦总数，取模基数。
const TrigramCount = 8

// LineCount 六爻总数，动爻取模基数。
const LineCount = 6

// LowerTrigramLineBoundary 下卦爻号上限（1~3为下卦，4~6为上卦）。
const LowerTrigramLineBoundary = 3

// DiZhiCount 地支总数。
const DiZhiCount = 12

// FallbackYearNumber 当年支无法解析时的默认序数（辰=5）。
const FallbackYearNumber = 5

// ── 起卦方式 ──────────────────────────────────────────────────────

const (
	MethodTime   = "time"   // 时间起卦
	MethodNumber = "number" // 数字起卦
)

// ── 体用位置 ──────────────────────────────────────────────────────

const (
	HexPositionUpper = "upper"
	HexPositionLower = "lower"
)

// ── 五行体用关系 ──────────────────────────────────────────────────

const (
	RelationBiHe    = "比和"
	YongShengTi     = "用生体"
	TiShengYong     = "体生用"
	YongKeTi        = "用克体"
	TiKeYong        = "体克用"
)

// ── 爻值 ──────────────────────────────────────────────────────────

const (
	LineYin = 0 // 阴爻
	LineYang = 1 // 阳爻
)

// ── 先天八卦数映射：余数 -> 卦 ────────────────────────────────────

var meihuaTrigramByNumber = map[int]model.MeiHuaTrigram{
	TrigramQian: {Number: TrigramQian, Name: "乾", Symbol: "☰", Element: "金", YinYang: "阳"},
	TrigramDui:  {Number: TrigramDui, Name: "兑", Symbol: "☱", Element: "金", YinYang: "阴"},
	TrigramLi:   {Number: TrigramLi, Name: "离", Symbol: "☲", Element: "火", YinYang: "阴"},
	TrigramZhen: {Number: TrigramZhen, Name: "震", Symbol: "☳", Element: "木", YinYang: "阳"},
	TrigramXun:  {Number: TrigramXun, Name: "巽", Symbol: "☴", Element: "木", YinYang: "阴"},
	TrigramKan:  {Number: TrigramKan, Name: "坎", Symbol: "☵", Element: "水", YinYang: "阳"},
	TrigramGen:  {Number: TrigramGen, Name: "艮", Symbol: "☶", Element: "土", YinYang: "阳"},
	TrigramKun:  {Number: TrigramKun, Name: "坤", Symbol: "☷", Element: "土", YinYang: "阴"},
}

// ── 地支序数映射 ──────────────────────────────────────────────────

var diZhiNumbers = map[string]int{
	"子": 1, "丑": 2, "寅": 3, "卯": 4,
	"辰": 5, "巳": 6, "午": 7, "未": 8,
	"申": 9, "酉": 10, "戌": 11, "亥": 12,
}

var diZhiNames = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

// 卦名映射：上下卦组合 -> 卦名
var hexagramNames = map[string]map[string]string{
	"乾": {"乾": "乾为天", "兑": "天泽履", "离": "天火同人", "震": "天雷无妄", "巽": "天风姤", "坎": "天水讼", "艮": "天山遁", "坤": "天地否"},
	"兑": {"乾": "泽天夬", "兑": "兑为泽", "离": "泽火革", "震": "泽雷随", "巽": "泽风大过", "坎": "泽水困", "艮": "泽山咸", "坤": "泽地萃"},
	"离": {"乾": "火天大有", "兑": "火泽睽", "离": "离为火", "震": "火雷噬嗑", "巽": "火风鼎", "坎": "火水未济", "艮": "火山旅", "坤": "火地晋"},
	"震": {"乾": "雷天大壮", "兑": "雷泽归妹", "离": "雷火丰", "震": "震为雷", "巽": "雷风恒", "坎": "雷水解", "艮": "雷山小过", "坤": "雷地豫"},
	"巽": {"乾": "风天小畜", "兑": "风泽中孚", "离": "风火家人", "震": "风雷益", "巽": "巽为风", "坎": "风水涣", "艮": "风山渐", "坤": "风地观"},
	"坎": {"乾": "水天需", "兑": "水泽节", "离": "水火既济", "震": "水雷屯", "巽": "水风井", "坎": "坎为水", "艮": "水山蹇", "坤": "水地比"},
	"艮": {"乾": "山天大畜", "兑": "山泽损", "离": "山火贲", "震": "山雷颐", "巽": "山风蛊", "坎": "山水蒙", "艮": "艮为山", "坤": "山地剥"},
	"坤": {"乾": "地天泰", "兑": "地泽临", "离": "地火明夷", "震": "地雷复", "巽": "地风升", "坎": "地水师", "艮": "地山谦", "坤": "坤为地"},
}

// 五行生克关系
var wuXingSheng = map[string]string{
	"木": "火", "火": "土", "土": "金", "金": "水", "水": "木",
}

var wuXingKe = map[string]string{
	"木": "土", "土": "水", "水": "火", "火": "金", "金": "木",
}

// MeiHuaService implements the MeiHua (Plum Blossom) divination algorithm.
type MeiHuaService struct{}

// NewMeiHuaService creates a new MeiHuaService.
func NewMeiHuaService() *MeiHuaService {
	return &MeiHuaService{}
}

// CalculateByTime performs divination based on the current time.
func (s *MeiHuaService) CalculateByTime(question string, t time.Time, timezone string) (*model.MeiHuaResult, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.FixedZone("CST", 8*3600) // default to China Standard Time
	}
	localTime := t.In(loc)

	// Convert to lunar calendar
	lunar := calendar.NewSolarFromDate(localTime)
	lunarDate := lunar.GetLunar()

	yearGanZhi := lunarDate.GetYearInGanZhi()
	// 地支是后两个字符
	runes := []rune(yearGanZhi)
	yearBranch := string(runes[len(runes)-2:])

	lunarMonth := lunarDate.GetMonth()
	lunarDay := lunarDate.GetDay()

	// 时辰：23-1子时，1-3丑时，...
	hour := localTime.Hour()
	hourBranchIndex := (hour + 1) / 2
	if hourBranchIndex >= DiZhiCount {
		hourBranchIndex = 0 // 23点是子时
	}
	hourBranch := diZhiNames[hourBranchIndex]
	hourNumber := diZhiNumbers[hourBranch]

	// 年数：按地支序数
	yearNumber := 0
	for _, ch := range yearBranch {
		if n, ok := diZhiNumbers[string(ch)]; ok {
			yearNumber = n
			break
		}
	}
	if yearNumber == 0 {
		yearNumber = FallbackYearNumber
	}

	// 处理农历月为负数（闰月）的情况
	absMonth := lunarMonth
	if absMonth < 0 {
		absMonth = -absMonth
	}

	sourceValues := model.MeiHuaSourceValues{
		Method:     MethodTime,
		YearBranch: yearBranch,
		YearNumber: yearNumber,
		LunarMonth: absMonth,
		LunarDay:   lunarDay,
		HourBranch: hourBranch,
		HourNumber: hourNumber,
	}

	// 上卦：(年 + 月 + 日) % 8
	upperSum := yearNumber + absMonth + lunarDay
	upperRemainder := upperSum % TrigramCount
	if upperRemainder == 0 {
		upperRemainder = TrigramKun
	}

	// 下卦：(年 + 月 + 日 + 时) % 8
	lowerSum := upperSum + hourNumber
	lowerRemainder := lowerSum % TrigramCount
	if lowerRemainder == 0 {
		lowerRemainder = TrigramKun
	}

	// 动爻：(年 + 月 + 日 + 时) % 6
	movingLineRemainder := lowerSum % LineCount
	if movingLineRemainder == 0 {
		movingLineRemainder = LineCount
	}

	return s.buildResult(question, sourceValues, upperRemainder, lowerRemainder, movingLineRemainder, t.Format(time.RFC3339))
}

// MinNumbersRequired 数字起卦最少需要的数字个数。
const MinNumbersRequired = 2

// CalculateByNumbers performs divination based on user-provided numbers.
func (s *MeiHuaService) CalculateByNumbers(question string, numbers []int) (*model.MeiHuaResult, error) {
	if len(numbers) < MinNumbersRequired {
		return nil, fmt.Errorf("at least %d numbers required", MinNumbersRequired)
	}

	sourceValues := model.MeiHuaSourceValues{
		Method:  MethodNumber,
		Numbers: numbers,
	}

	// 按梅花易数数字起卦法：
	// 前半部分数字之和 % 8 = 上卦
	// 全部数字之和 % 8 = 下卦
	// 全部数字之和 % 6 = 动爻
	mid := len(numbers) / 2
	if mid == 0 {
		mid = 1
	}

	upperSum := 0
	for _, n := range numbers[:mid] {
		upperSum += n
	}

	totalSum := 0
	for _, n := range numbers {
		totalSum += n
	}

	upperRemainder := upperSum % TrigramCount
	if upperRemainder == 0 {
		upperRemainder = TrigramKun
	}

	lowerRemainder := totalSum % TrigramCount
	if lowerRemainder == 0 {
		lowerRemainder = TrigramKun
	}

	movingLineRemainder := totalSum % LineCount
	if movingLineRemainder == 0 {
		movingLineRemainder = LineCount
	}

	return s.buildResult(question, sourceValues, upperRemainder, lowerRemainder, movingLineRemainder, time.Now().Format(time.RFC3339))
}

// buildResult constructs the full MeiHuaResult from the calculated remainders.
func (s *MeiHuaService) buildResult(
	question string,
	sourceValues model.MeiHuaSourceValues,
	upperNum, lowerNum, movingLine int,
	timestamp string,
) (*model.MeiHuaResult, error) {
	upper := meihuaTrigramByNumber[upperNum]
	lower := meihuaTrigramByNumber[lowerNum]

	benGua := s.buildHexagram(upper, lower)

	// 互卦：取本卦2、3、4爻为下卦，3、4、5爻为上卦
	huUpper, huLower := s.calcHuGua(upperNum, lowerNum)
	huGua := s.buildHexagram(meihuaTrigramByNumber[huUpper], meihuaTrigramByNumber[huLower])

	// 变卦：动爻变后的新卦
	bianUpper, bianLower := s.calcBianGua(upperNum, lowerNum, movingLine)
	bianGua := s.buildHexagram(meihuaTrigramByNumber[bianUpper], meihuaTrigramByNumber[bianLower])

	// 体用判断：动爻所在卦为用，另一卦为体
	var tiTrigram, yongTrigram model.MeiHuaTrigram
	tiHex, yongHex := HexPositionUpper, HexPositionLower
	if movingLine <= LowerTrigramLineBoundary {
		// 动爻在下卦，下卦为用，上卦为体
		tiTrigram = upper
		yongTrigram = lower
		tiHex, yongHex = HexPositionUpper, HexPositionLower
	} else {
		// 动爻在上卦，上卦为用，下卦为体
		tiTrigram = lower
		yongTrigram = upper
		tiHex, yongHex = HexPositionLower, HexPositionUpper
	}

	relation := s.calcWuXingRelation(tiTrigram.Element, yongTrigram.Element)

	tiYong := model.MeiHuaTiYong{
		Ti:       tiTrigram,
		Yong:     yongTrigram,
		TiHex:    tiHex,
		YongHex:  yongHex,
		Relation: relation,
	}

	return &model.MeiHuaResult{
		Question:     question,
		Method:       sourceValues.Method,
		SourceValues: sourceValues,
		UpperTrigram: upper,
		LowerTrigram: lower,
		BenGua:       benGua,
		HuGua:        huGua,
		BianGua:      bianGua,
		MovingLine:   movingLine,
		TiYong:       tiYong,
		Timestamp:    timestamp,
	}, nil
}

// buildHexagram constructs a MeiHuaHexagram from upper and lower trigrams.
func (s *MeiHuaService) buildHexagram(upper, lower model.MeiHuaTrigram) model.MeiHuaHexagram {
	name := hexagramNames[upper.Name][lower.Name]
	nameShort := ""
	if len([]rune(name)) >= 2 {
		// 取最后两个字作为简称，如"泽火革" -> "革"
		runes := []rune(name)
		nameShort = string(runes[len(runes)-1:])
	}

	return model.MeiHuaHexagram{
		Upper:     upper,
		Lower:     lower,
		Name:      name,
		NameShort: nameShort,
	}
}

// calcHuGua calculates the互卦 trigram numbers.
// 互卦：取本卦2、3、4爻（从下往上数）为下卦，3、4、5爻为上卦。
// 本卦的六爻从下到上：1爻(下卦底)、2爻(下卦中)、3爻(下卦顶)、4爻(上卦底)、5爻(上卦中)、6爻(上卦顶)
// 互卦下卦 = 2、3、4爻 -> 对应下卦的上中 + 上卦的下
// 互卦上卦 = 3、4、5爻 -> 对应下卦的上 + 上卦的下中
func (s *MeiHuaService) calcHuGua(upperNum, lowerNum int) (huUpper, huLower int) {
	// 将上下卦展开为6爻（阴阳）
	// 每卦3爻，从下到上
	lowerLines := trigramToLines(lowerNum) // [底, 中, 顶]
	upperLines := trigramToLines(upperNum) // [底, 中, 顶]

	// 6爻从下到上：[lowerLines[0], lowerLines[1], lowerLines[2], upperLines[0], upperLines[1], upperLines[2]]
	// 互卦下卦 = 爻2、3、4 = lowerLines[1], lowerLines[2], upperLines[0]
	// 互卦上卦 = 爻3、4、5 = lowerLines[2], upperLines[0], upperLines[1]
	huLowerLines := [3]int{lowerLines[1], lowerLines[2], upperLines[0]}
	huUpperLines := [3]int{lowerLines[2], upperLines[0], upperLines[1]}

	huLower = linesToTrigramNum(huLowerLines)
	huUpper = linesToTrigramNum(huUpperLines)
	return
}

// calcBianGua calculates the变卦 trigram numbers after the moving line changes.
func (s *MeiHuaService) calcBianGua(upperNum, lowerNum, movingLine int) (bianUpper, bianLower int) {
	lowerLines := trigramToLines(lowerNum)
	upperLines := trigramToLines(upperNum)

	// 动爻1-6，从下往上，阴阳互变
	if movingLine <= LowerTrigramLineBoundary {
		lowerLines[movingLine-1] = flipLine(lowerLines[movingLine-1])
	} else {
		upperLines[movingLine-LowerTrigramLineBoundary-1] = flipLine(upperLines[movingLine-LowerTrigramLineBoundary-1])
	}

	bianLower = linesToTrigramNum(lowerLines)
	bianUpper = linesToTrigramNum(upperLines)
	return
}

// flipLine 阴阳互变：阴变阳，阳变阴。
func flipLine(line int) int {
	if line == LineYang {
		return LineYin
	}
	return LineYang
}

// trigramToLines converts a trigram number to its 3 lines (from bottom to top).
func trigramToLines(num int) [3]int {
	// 先天八卦对应的爻象（从下到上）
	switch num {
	case TrigramQian: // ☰ 阳阳阳
		return [3]int{LineYang, LineYang, LineYang}
	case TrigramDui: // ☱ 阳阳阴
		return [3]int{LineYang, LineYang, LineYin}
	case TrigramLi: // ☲ 阳阴阳
		return [3]int{LineYang, LineYin, LineYang}
	case TrigramZhen: // ☳ 阳阴阴
		return [3]int{LineYang, LineYin, LineYin}
	case TrigramXun: // ☴ 阴阳阳
		return [3]int{LineYin, LineYang, LineYang}
	case TrigramKan: // ☵ 阴阳阴
		return [3]int{LineYin, LineYang, LineYin}
	case TrigramGen: // ☶ 阴阴阳
		return [3]int{LineYin, LineYin, LineYang}
	case TrigramKun: // ☷ 阴阴阴
		return [3]int{LineYin, LineYin, LineYin}
	default:
		return [3]int{LineYin, LineYin, LineYin}
	}
}

// linesToTrigramNum converts 3 lines (bottom to top) to a trigram number.
func linesToTrigramNum(lines [3]int) int {
	switch {
	case lines == [3]int{LineYang, LineYang, LineYang}:
		return TrigramQian
	case lines == [3]int{LineYang, LineYang, LineYin}:
		return TrigramDui
	case lines == [3]int{LineYang, LineYin, LineYang}:
		return TrigramLi
	case lines == [3]int{LineYang, LineYin, LineYin}:
		return TrigramZhen
	case lines == [3]int{LineYin, LineYang, LineYang}:
		return TrigramXun
	case lines == [3]int{LineYin, LineYang, LineYin}:
		return TrigramKan
	case lines == [3]int{LineYin, LineYin, LineYang}:
		return TrigramGen
	case lines == [3]int{LineYin, LineYin, LineYin}:
		return TrigramKun
	default:
		return TrigramKun
	}
}

// calcWuXingRelation determines the 生克 relation between ti and yong elements.
func (s *MeiHuaService) calcWuXingRelation(tiElement, yongElement string) string {
	if tiElement == yongElement {
		return RelationBiHe
	}
	if wuXingSheng[yongElement] == tiElement {
		return YongShengTi
	}
	if wuXingSheng[tiElement] == yongElement {
		return TiShengYong
	}
	if wuXingKe[yongElement] == tiElement {
		return YongKeTi
	}
	if wuXingKe[tiElement] == yongElement {
		return TiKeYong
	}
	return RelationBiHe
}

// MockMeiHuaService is a mock for testing.
type MockMeiHuaService struct {
	Result *model.MeiHuaResult
	Err    error
}
