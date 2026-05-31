package service

import "zhanbu/internal/model"

// getEmbeddedHexagrams returns all 64 hexagrams as embedded data.
func getEmbeddedHexagrams() []model.Hexagram {
	return []model.Hexagram{
		{ID: 1, Name: "乾为天", NameShort: "乾", UpperTrigram: "乾", LowerTrigram: "乾", Binary: "111111", Judgment: "乾，元亨利贞。", Image: "天行健，君子以自强不息。", LineTexts: []string{"初九：潜龙勿用。", "九二：见龙在田，利见大人。", "九三：君子终日乾乾，夕惕若，厉无咎。", "九四：或跃在渊，无咎。", "九五：飞龙在天，利见大人。", "上九：亢龙有悔。"}, Description: "乾卦象征天，纯阳之卦，刚健中正。"},
		{ID: 2, Name: "坤为地", NameShort: "坤", UpperTrigram: "坤", LowerTrigram: "坤", Binary: "000000", Judgment: "坤，元亨，利牝马之贞。", Image: "地势坤，君子以厚德载物。", LineTexts: []string{"初六：履霜，坚冰至。", "六二：直方大，不习无不利。", "六三：含章可贞，或从王事。", "六四：括囊，无咎无誉。", "六五：黄裳，元吉。", "上六：龙战于野，其血玄黄。"}, Description: "坤卦象征地，纯阴之卦。柔顺包容，厚德载物。"},
	}
}
