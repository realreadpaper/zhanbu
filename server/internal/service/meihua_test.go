package service

import (
	"testing"
	"time"

	"zhanbu/internal/model"
)

func TestMeihuaTrigramMapping(t *testing.T) {
	// Verify all 8 trigrams exist
	for i := TrigramQian; i <= TrigramKun; i++ {
		trigram, ok := meihuaTrigramByNumber[i]
		if !ok {
			t.Errorf("trigram number %d not found", i)
			continue
		}
		if trigram.Number != i {
			t.Errorf("trigram %d has wrong number: %d", i, trigram.Number)
		}
	}

	// 余0为坤
	if meihuaTrigramByNumber[TrigramKun].Name != "坤" {
		t.Errorf("remainder 0 should map to 坤, got %s", meihuaTrigramByNumber[TrigramKun].Name)
	}
}

func TestMeihuaClassicCase(t *testing.T) {
	// 经典观梅案例：辰年十二月十七日申时
	// 年数=5(辰), 月=12, 日=17, 时=9(申)
	// 上卦：(5+12+17)%8 = 34%8 = 2 -> 兑
	// 下卦：(5+12+17+9)%8 = 43%8 = 5 -> 巽? 不对，应该是3(离)
	// 等等，让我重新算：43%8 = 5... 不对
	// 43 / 8 = 5 余 3，所以下卦是3 -> 离
	// 动爻：43%6 = 1 -> 初爻动
	// 本卦：兑上离下 = 泽火革
	// 变卦：初爻变，离 -> 阳阴阳 -> 阳阴阴 = 震? 不对
	// 离 = 阳阴阳(从下到上)，初爻变 -> 阴阴阳 = 艮
	// 变卦：兑上艮下 = 泽山咸

	svc := NewMeiHuaService()

	// 手动测试 buildResult
	upperNum := TrigramDui  // 兑
	lowerNum := TrigramLi   // 离
	movingLine := 1

	result, err := svc.buildResult("观梅", mockSourceValues(), upperNum, lowerNum, movingLine, time.Now().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 验证本卦
	if result.BenGua.Name != "泽火革" {
		t.Errorf("ben gua name: want 泽火革, got %s", result.BenGua.Name)
	}

	// 验证变卦：初爻变，离(101) -> 艮(001)，兑上艮下 = 泽山咸
	if result.BianGua.Name != "泽山咸" {
		t.Errorf("bian gua name: want 泽山咸, got %s", result.BianGua.Name)
	}

	// 验证动爻
	if result.MovingLine != 1 {
		t.Errorf("moving line: want 1, got %d", result.MovingLine)
	}

	// 验证体用：动爻在下卦(1<=3)，下卦为用，上卦为体
	if result.TiYong.Ti.Name != "兑" {
		t.Errorf("ti trigram: want 兑, got %s", result.TiYong.Ti.Name)
	}
	if result.TiYong.Yong.Name != "离" {
		t.Errorf("yong trigram: want 离, got %s", result.TiYong.Yong.Name)
	}
	if result.TiYong.Relation != "用克体" {
		t.Errorf("ti-yong relation: want 用克体, got %s", result.TiYong.Relation)
	}
}

func TestMeihuaRemainderZero(t *testing.T) {
	svc := NewMeiHuaService()

	// 测试余数为0的情况
	// 上卦余0 -> 坤(TrigramKun), 下卦余0 -> 坤(TrigramKun), 动爻余0 -> 6爻
	result, err := svc.buildResult("test", mockSourceValues(), TrigramKun, TrigramKun, LineCount, time.Now().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.BenGua.Name != "坤为地" {
		t.Errorf("ben gua name: want 坤为地, got %s", result.BenGua.Name)
	}
	if result.MovingLine != LineCount {
		t.Errorf("moving line: want %d, got %d", LineCount, result.MovingLine)
	}
}

func TestMeihuaTiYongUpperMoving(t *testing.T) {
	svc := NewMeiHuaService()

	// 动爻在上卦(4<=6)，上卦为用，下卦为体
	result, err := svc.buildResult("test", mockSourceValues(), TrigramDui, TrigramLi, 5, time.Now().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 动爻5在上卦，上卦(兑)为用，下卦(离)为体
	if result.TiYong.Ti.Name != "离" {
		t.Errorf("ti trigram: want 离, got %s", result.TiYong.Ti.Name)
	}
	if result.TiYong.Yong.Name != "兑" {
		t.Errorf("yong trigram: want 兑, got %s", result.TiYong.Yong.Name)
	}
}

func TestMeihuaHuGua(t *testing.T) {
	svc := NewMeiHuaService()

	// 泽火革：兑(上)离(下)
	// 兑 = 110(从下到上), 离 = 101(从下到上)
	// 6爻从下到上：1,0,1,1,1,0
	// 互卦下卦 = 爻2,3,4 = 0,1,1 -> 巽(5)? 不对
	// 爻2=0, 爻3=1, 爻4=1 -> [0,1,1] -> 巽(5)... 不对
	// 等等，让我重新理解
	// lower = 离 = [1,0,1] (底,中,顶)
	// upper = 兑 = [1,1,0] (底,中,顶)
	// 6爻: [1,0,1, 1,1,0]
	// 互卦下卦 = 爻2,3,4 = lower[1], lower[2], upper[0] = 0,1,1 -> 巽
	// 互卦上卦 = 爻3,4,5 = lower[2], upper[0], upper[1] = 1,1,1 -> 乾
	// 互卦 = 乾上巽下 = 天风姤

	// 但传统上泽火革的互卦是天风姤
	// 让我验证一下

	huUpper, huLower := svc.calcHuGua(TrigramDui, TrigramLi) // 兑上离下
	t.Logf("huGua: upper=%d, lower=%d", huUpper, huLower)

	// 离=[1,0,1], 兑=[1,1,0]
	// lower[1]=0, lower[2]=1, upper[0]=1 -> [0,1,1] = 巽(5)
	// lower[2]=1, upper[0]=1, upper[1]=1 -> [1,1,1] = 乾(1)
	// 互卦 = 乾上巽下 = 天风姤

	if huUpper != TrigramQian {
		t.Errorf("hu gua upper: want %d(乾), got %d", TrigramQian, huUpper)
	}
	if huLower != TrigramXun {
		t.Errorf("hu gua lower: want %d(巽), got %d", TrigramXun, huLower)
	}

	// 验证卦名
	result, err := svc.buildResult("test", mockSourceValues(), TrigramDui, TrigramLi, 1, time.Now().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HuGua.Name != "天风姤" {
		t.Errorf("hu gua name: want 天风姤, got %s", result.HuGua.Name)
	}
}

func TestMeihuaNumberDivination(t *testing.T) {
	svc := NewMeiHuaService()

	// 数字起卦：12, 34
	// 前半(12) % 8 = 4 -> 震
	// 全(12+34=46) % 8 = 6 -> 坎
	// 全(46) % 6 = 4 -> 4爻动
	result, err := svc.CalculateByNumbers("事业", []int{12, 34})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Method != MethodNumber {
		t.Errorf("method: want %s, got %s", MethodNumber, result.Method)
	}
	if result.UpperTrigram.Number != TrigramZhen {
		t.Errorf("upper trigram: want 震(%d), got %s", TrigramZhen, result.UpperTrigram.Name)
	}
	if result.LowerTrigram.Number != TrigramKan {
		t.Errorf("lower trigram: want 坎(%d), got %s", TrigramKan, result.LowerTrigram.Name)
	}
	if result.MovingLine != 4 {
		t.Errorf("moving line: want 4, got %d", result.MovingLine)
	}
}

func TestMeihuaWuXingRelation(t *testing.T) {
	svc := NewMeiHuaService()

	tests := []struct {
		ti, yong, want string
	}{
		{"金", "金", RelationBiHe},
		{"金", "土", YongShengTi}, // 土生金
		{"金", "水", TiShengYong}, // 金生水
		{"金", "火", YongKeTi},    // 火克金
		{"金", "木", TiKeYong},    // 金克木
	}

	for _, tt := range tests {
		got := svc.calcWuXingRelation(tt.ti, tt.yong)
		if got != tt.want {
			t.Errorf("calcWuXingRelation(%s, %s): want %s, got %s", tt.ti, tt.yong, tt.want, got)
		}
	}
}

func mockSourceValues() model.MeiHuaSourceValues {
	return model.MeiHuaSourceValues{
		Method:     MethodTime,
		YearBranch: "辰",
		YearNumber: 5,
		LunarMonth: 12,
		LunarDay:   17,
		HourBranch: "申",
		HourNumber: 9,
	}
}
