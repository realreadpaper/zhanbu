package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaZiCalculate(t *testing.T) {
	svc := NewBaZiService()

	result, err := svc.Calculate("1990-05-15", "14:30", "male")
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Check birth info
	assert.Equal(t, "1990-05-15 14:30", result.Birth.Solar)
	assert.NotEmpty(t, result.Birth.Lunar)

	// Check pillars
	assert.NotEmpty(t, result.Pillars.Year.TianGan)
	assert.NotEmpty(t, result.Pillars.Year.DiZhi)
	assert.NotEmpty(t, result.Pillars.Month.TianGan)
	assert.NotEmpty(t, result.Pillars.Month.DiZhi)
	assert.NotEmpty(t, result.Pillars.Day.TianGan)
	assert.NotEmpty(t, result.Pillars.Day.DiZhi)
	assert.NotEmpty(t, result.Pillars.Hour.TianGan)
	assert.NotEmpty(t, result.Pillars.Hour.DiZhi)

	// Check WuXing
	assert.NotEmpty(t, result.Pillars.Year.WuXing)
	assert.NotEmpty(t, result.Pillars.Month.WuXing)
	assert.NotEmpty(t, result.Pillars.Day.WuXing)
	assert.NotEmpty(t, result.Pillars.Hour.WuXing)

	// Check NaYin
	assert.NotEmpty(t, result.Pillars.Year.NaYin)
	assert.NotEmpty(t, result.Pillars.Month.NaYin)
	assert.NotEmpty(t, result.Pillars.Day.NaYin)
	assert.NotEmpty(t, result.Pillars.Hour.NaYin)

	// Check five elements
	assert.GreaterOrEqual(t, result.FiveElements.Metal, 0)
	assert.GreaterOrEqual(t, result.FiveElements.Wood, 0)
	assert.NotEmpty(t, result.FiveElements.DayMaster)
	assert.NotEmpty(t, result.FiveElements.Strength)
	assert.NotEmpty(t, result.FiveElements.YongShen)
	assert.NotEmpty(t, result.FiveElements.JiShen)

	// Check ten gods
	assert.NotEmpty(t, result.TenGods)
}

func TestBaZiInvalidDate(t *testing.T) {
	svc := NewBaZiService()

	_, err := svc.Calculate("invalid", "14:30", "male")
	assert.Error(t, err)

	_, err = svc.Calculate("1990-13-15", "14:30", "male")
	assert.Error(t, err)

	_, err = svc.Calculate("1990-05-15", "25:30", "male")
	assert.Error(t, err)
}

func TestCalcYearPillar(t *testing.T) {
	svc := NewBaZiService()

	// 2024 is 甲辰年
	tg, dz := svc.calcYearPillar(2024)
	assert.Equal(t, "甲", tg)
	assert.Equal(t, "辰", dz)

	// 2025 is 乙巳年
	tg, dz = svc.calcYearPillar(2025)
	assert.Equal(t, "乙", tg)
	assert.Equal(t, "巳", dz)

	// 1990 is 庚午年
	tg, dz = svc.calcYearPillar(1990)
	assert.Equal(t, "庚", tg)
	assert.Equal(t, "午", dz)
}

func TestCalcDayPillar(t *testing.T) {
	svc := NewBaZiService()

	// Verify consecutive days produce consecutive Jiazi indices
	// We test by checking that the TianGan advances by 1 and DiZhi advances by 1
	tg1, dz1 := svc.calcDayPillar(2024, 1, 1)
	tg2, dz2 := svc.calcDayPillar(2024, 1, 2)

	tgIdx1 := indexOf(tianGan, tg1)
	tgIdx2 := indexOf(tianGan, tg2)
	dzIdx1 := indexOf(diZhi, dz1)
	dzIdx2 := indexOf(diZhi, dz2)

	// TianGan should advance by 1 (mod 10)
	assert.Equal(t, 1, (tgIdx2-tgIdx1+10)%10)
	// DiZhi should advance by 1 (mod 12)
	assert.Equal(t, 1, (dzIdx2-dzIdx1+12)%12)
}

func TestTenGodName(t *testing.T) {
	// Same element -> 比肩
	assert.Equal(t, "比肩", tenGodName("木", "木"))

	// 我生 (木生火)
	assert.Equal(t, "食神", tenGodName("木", "火"))

	// 生我 (水生木)
	assert.Equal(t, "偏印", tenGodName("木", "水"))

	// 我克 (木克土)
	assert.Equal(t, "偏财", tenGodName("木", "土"))

	// 克我 (金克木)
	assert.Equal(t, "偏官", tenGodName("木", "金"))
}
