package database

import (
	"encoding/json"
	"fmt"
	"os"

	"zhanbu/internal/model"

	"gorm.io/gorm"
)

// SeedData 导入种子数据
func SeedData(db *gorm.DB) error {
	// 导入塔罗牌数据
	if err := seedTarotCards(db); err != nil {
		return fmt.Errorf("导入塔罗牌数据失败: %w", err)
	}
	return nil
}

func seedTarotCards(db *gorm.DB) error {
	// 检查是否已有数据
	var count int64
	db.Model(&model.TarotCard{}).Count(&count)
	if count > 0 {
		fmt.Printf("塔罗牌数据已存在 (%d张)，跳过导入\n", count)
		return nil
	}

	// 导入大阿尔卡纳
	if err := importTarotJSON(db, "data/tarot/major_arcana.json"); err != nil {
		return fmt.Errorf("导入大阿尔卡纳失败: %w", err)
	}

	// 导入小阿尔卡纳
	if err := importTarotJSON(db, "data/tarot/minor_arcana.json"); err != nil {
		return fmt.Errorf("导入小阿尔卡纳失败: %w", err)
	}

	var total int64
	db.Model(&model.TarotCard{}).Count(&total)
	fmt.Printf("塔罗牌数据导入完成，共 %d 张\n", total)
	return nil
}

func importTarotJSON(db *gorm.DB, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("读取文件 %s 失败: %w", filepath, err)
	}

	var cards []model.TarotCard
	if err := json.Unmarshal(data, &cards); err != nil {
		return fmt.Errorf("解析JSON失败: %w", err)
	}

	for i := range cards {
		if err := db.Create(&cards[i]).Error; err != nil {
			return fmt.Errorf("插入牌 %d 失败: %w", cards[i].ID, err)
		}
	}

	fmt.Printf("成功导入 %d 张牌 from %s\n", len(cards), filepath)
	return nil
}
