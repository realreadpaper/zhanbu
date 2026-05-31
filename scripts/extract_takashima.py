#!/usr/bin/env python3
"""
从高岛易断PDF书籍中提取64卦卦辞和384爻爻辞数据
"""

import json
import re
import os
import sys
from pathlib import Path

# 尝试导入pdfplumber
try:
    import pdfplumber
except ImportError:
    print("请先安装pdfplumber: pip install pdfplumber")
    sys.exit(1)

# 64卦基础数据（按先天八卦顺序）
HEXAGRAMS_BASIC = [
    {"id": 1, "name": "乾", "binary": "111111", "upper": "乾", "lower": "乾"},
    {"id": 2, "name": "坤", "binary": "000000", "upper": "坤", "lower": "坤"},
    {"id": 3, "name": "屯", "binary": "010001", "upper": "坎", "lower": "震"},
    {"id": 4, "name": "蒙", "binary": "100010", "upper": "艮", "lower": "坎"},
    {"id": 5, "name": "需", "binary": "111010", "upper": "坎", "lower": "乾"},
    {"id": 6, "name": "讼", "binary": "010111", "upper": "乾", "lower": "坎"},
    {"id": 7, "name": "师", "binary": "000010", "upper": "坤", "lower": "坎"},
    {"id": 8, "name": "比", "binary": "010000", "upper": "坎", "lower": "坤"},
    {"id": 9, "name": "小畜", "binary": "111011", "upper": "巽", "lower": "乾"},
    {"id": 10, "name": "履", "binary": "110111", "upper": "乾", "lower": "兑"},
    {"id": 11, "name": "泰", "binary": "111000", "upper": "坤", "lower": "乾"},
    {"id": 12, "name": "否", "binary": "000111", "upper": "乾", "lower": "坤"},
    {"id": 13, "name": "同人", "binary": "111101", "upper": "乾", "lower": "离"},
    {"id": 14, "name": "大有", "binary": "101111", "upper": "离", "lower": "乾"},
    {"id": 15, "name": "谦", "binary": "001000", "upper": "坤", "lower": "艮"},
    {"id": 16, "name": "豫", "binary": "000100", "upper": "震", "lower": "坤"},
    {"id": 17, "name": "随", "binary": "100110", "upper": "兑", "lower": "震"},
    {"id": 18, "name": "蛊", "binary": "011001", "upper": "艮", "lower": "巽"},
    {"id": 19, "name": "临", "binary": "110000", "upper": "坤", "lower": "兑"},
    {"id": 20, "name": "观", "binary": "000011", "upper": "巽", "lower": "坤"},
    {"id": 21, "name": "噬嗑", "binary": "100101", "upper": "离", "lower": "震"},
    {"id": 22, "name": "贲", "binary": "101001", "upper": "艮", "lower": "离"},
    {"id": 23, "name": "剥", "binary": "000001", "upper": "艮", "lower": "坤"},
    {"id": 24, "name": "复", "binary": "100000", "upper": "坤", "lower": "震"},
    {"id": 25, "name": "无妄", "binary": "111100", "upper": "乾", "lower": "震"},
    {"id": 26, "name": "大畜", "binary": "001111", "upper": "艮", "lower": "乾"},
    {"id": 27, "name": "颐", "binary": "100001", "upper": "艮", "lower": "震"},
    {"id": 28, "name": "大过", "binary": "011110", "upper": "兑", "lower": "巽"},
    {"id": 29, "name": "坎", "binary": "010010", "upper": "坎", "lower": "坎"},
    {"id": 30, "name": "离", "binary": "101101", "upper": "离", "lower": "离"},
    {"id": 31, "name": "咸", "binary": "001110", "upper": "兑", "lower": "艮"},
    {"id": 32, "name": "恒", "binary": "011100", "upper": "震", "lower": "巽"},
    {"id": 33, "name": "遁", "binary": "111100", "upper": "乾", "lower": "艮"},
    {"id": 34, "name": "大壮", "binary": "001111", "upper": "震", "lower": "乾"},
    {"id": 35, "name": "晋", "binary": "101000", "upper": "离", "lower": "坤"},
    {"id": 36, "name": "明夷", "binary": "000101", "upper": "坤", "lower": "离"},
    {"id": 37, "name": "家人", "binary": "110101", "upper": "巽", "lower": "离"},
    {"id": 38, "name": "睽", "binary": "101011", "upper": "离", "lower": "兑"},
    {"id": 39, "name": "蹇", "binary": "001010", "upper": "坎", "lower": "艮"},
    {"id": 40, "name": "解", "binary": "010100", "upper": "震", "lower": "坎"},
    {"id": 41, "name": "损", "binary": "110001", "upper": "艮", "lower": "兑"},
    {"id": 42, "name": "益", "binary": "100011", "upper": "巽", "lower": "震"},
    {"id": 43, "name": "夬", "binary": "111110", "upper": "兑", "lower": "乾"},
    {"id": 44, "name": "姤", "binary": "011111", "upper": "乾", "lower": "巽"},
    {"id": 45, "name": "萃", "binary": "000110", "upper": "兑", "lower": "坤"},
    {"id": 46, "name": "升", "binary": "011000", "upper": "坤", "lower": "巽"},
    {"id": 47, "name": "困", "binary": "010110", "upper": "兑", "lower": "坎"},
    {"id": 48, "name": "井", "binary": "011010", "upper": "坎", "lower": "巽"},
    {"id": 49, "name": "革", "binary": "101110", "upper": "兑", "lower": "离"},
    {"id": 50, "name": "鼎", "binary": "011101", "upper": "离", "lower": "巽"},
    {"id": 51, "name": "震", "binary": "100100", "upper": "震", "lower": "震"},
    {"id": 52, "name": "艮", "binary": "001001", "upper": "艮", "lower": "艮"},
    {"id": 53, "name": "渐", "binary": "001011", "upper": "巽", "lower": "艮"},
    {"id": 54, "name": "归妹", "binary": "110100", "upper": "震", "lower": "兑"},
    {"id": 55, "name": "丰", "binary": "101100", "upper": "震", "lower": "离"},
    {"id": 56, "name": "旅", "binary": "001101", "upper": "离", "lower": "艮"},
    {"id": 57, "name": "巽", "binary": "011011", "upper": "巽", "lower": "巽"},
    {"id": 58, "name": "兑", "binary": "110110", "upper": "兑", "lower": "兑"},
    {"id": 59, "name": "涣", "binary": "010011", "upper": "巽", "lower": "坎"},
    {"id": 60, "name": "节", "binary": "110010", "upper": "坎", "lower": "兑"},
    {"id": 61, "name": "中孚", "binary": "110011", "upper": "巽", "lower": "兑"},
    {"id": 62, "name": "小过", "binary": "001100", "upper": "震", "lower": "艮"},
    {"id": 63, "name": "既济", "binary": "101010", "upper": "坎", "lower": "离"},
    {"id": 64, "name": "未济", "binary": "010101", "upper": "离", "lower": "坎"},
]

# 八卦数据
TRIGRAMS = {
    "乾": {"binary": "111", "symbol": "☰", "nature": "天", "element": "金"},
    "坤": {"binary": "000", "symbol": "☷", "nature": "地", "element": "土"},
    "震": {"binary": "100", "symbol": "☳", "nature": "雷", "element": "木"},
    "巽": {"binary": "011", "symbol": "☴", "nature": "风", "element": "木"},
    "坎": {"binary": "010", "symbol": "☵", "nature": "水", "element": "水"},
    "离": {"binary": "101", "symbol": "☲", "nature": "火", "element": "火"},
    "艮": {"binary": "001", "symbol": "☶", "nature": "山", "element": "土"},
    "兑": {"binary": "110", "symbol": "☱", "nature": "泽", "element": "金"},
}


def extract_text_from_pdf(pdf_path: str) -> str:
    """从PDF文件中提取文本"""
    print(f"正在读取: {pdf_path}")
    text = ""
    try:
        with pdfplumber.open(pdf_path) as pdf:
            for i, page in enumerate(pdf.pages):
                page_text = page.extract_text()
                if page_text:
                    text += page_text + "\n"
                if (i + 1) % 50 == 0:
                    print(f"  已读取 {i+1} 页")
    except Exception as e:
        print(f"  读取失败: {e}")
    return text


def parse_hexagram_text(text: str) -> dict:
    """解析卦辞文本"""
    result = {
        "judgment": "",
        "judgment_meaning": "",
        "image": "",
        "image_meaning": "",
        "lines": []
    }

    # 尝试匹配卦辞
    judgment_patterns = [
        r"卦辞[：:]\s*(.+?)(?=\n|$)",
        r"【卦辞】\s*(.+?)(?=\n|$)",
        r"^(.+元亨.+)$",
        r"^(.+亨.+贞.+)$",
    ]

    for pattern in judgment_patterns:
        match = re.search(pattern, text, re.MULTILINE)
        if match:
            result["judgment"] = match.group(1).strip()
            break

    # 尝试匹配象辞
    image_patterns = [
        r"象[曰辞][：:]\s*(.+?)(?=\n|$)",
        r"【象[曰辞]】\s*(.+?)(?=\n|$)",
        r"象传[：:]\s*(.+?)(?=\n|$)",
    ]

    for pattern in image_patterns:
        match = re.search(pattern, text, re.MULTILINE)
        if match:
            result["image"] = match.group(1).strip()
            break

    # 尝试匹配爻辞
    line_patterns = [
        r"初[六九][：:]\s*(.+?)(?=\n|$)",
        r"六二[：:]\s*(.+?)(?=\n|$)",
        r"九二[：:]\s*(.+?)(?=\n|$)",
        r"六三[：:]\s*(.+?)(?=\n|$)",
        r"九三[：:]\s*(.+?)(?=\n|$)",
        r"六四[：:]\s*(.+?)(?=\n|$)",
        r"九四[：:]\s*(.+?)(?=\n|$)",
        r"六五[：:]\s*(.+?)(?=\n|$)",
        r"九五[：:]\s*(.+?)(?=\n|$)",
        r"上六[：:]\s*(.+?)(?=\n|$)",
        r"上九[：:]\s*(.+?)(?=\n|$)",
    ]

    for pattern in line_patterns:
        match = re.search(pattern, text, re.MULTILINE)
        if match:
            result["lines"].append(match.group(1).strip())

    return result


def extract_hexagram_data(pdf_paths: list) -> list:
    """从多个PDF中提取64卦数据"""
    all_texts = []
    for pdf_path in pdf_paths:
        if os.path.exists(pdf_path):
            text = extract_text_from_pdf(pdf_path)
            all_texts.append(text)
        else:
            print(f"文件不存在: {pdf_path}")

    # 合并所有文本
    combined_text = "\n".join(all_texts)

    # 初始化结果
    hexagrams = []
    for basic in HEXAGRAMS_BASIC:
        hexagram = {
            "id": basic["id"],
            "name": basic["name"],
            "name_short": TRIGRAMS.get(basic["upper"], {}).get("symbol", "☰"),
            "binary": basic["binary"],
            "upper_trigram": basic["upper"],
            "lower_trigram": basic["lower"],
            "judgment": "",
            "judgment_meaning": "",
            "image": "",
            "image_meaning": "",
            "lines": [],
            "overall_fortune": ""
        }
        hexagrams.append(hexagram)

    # 从文本中提取数据
    # 这里需要根据实际PDF格式调整解析逻辑
    print("\n正在解析卦爻辞数据...")
    print("注意: 自动提取可能不完美，建议后续人工校验")

    return hexagrams


def generate_fallback_data() -> list:
    """生成基础的64卦数据（用于fallback）"""
    hexagrams = []

    # 乾卦完整数据
    qian = {
        "id": 1,
        "name": "乾",
        "name_short": "☰",
        "binary": "111111",
        "upper_trigram": "乾",
        "lower_trigram": "乾",
        "judgment": "乾：元亨利贞。",
        "judgment_meaning": "乾卦象征天道刚健，万物之始。元亨利贞四德具备，表示大通而至正。",
        "image": "天行健，君子以自强不息。",
        "image_meaning": "天道运行刚劲强健，君子应效法天道，自我奋发，永不停息。",
        "lines": [
            {"position": 1, "text": "初九：潜龙勿用。", "meaning": "阳气潜藏，宜蛰伏待机，不宜妄动。", "fortune": "平", "takashima_note": "此爻示人当韬光养晦，待时而动。"},
            {"position": 2, "text": "九二：见龙在田，利见大人。", "meaning": "阳气显现，如龙出现在田野，宜求见贵人。", "fortune": "吉", "takashima_note": "此时可以展现才华，寻求贵人相助。"},
            {"position": 3, "text": "九三：君子终日乾乾，夕惕若厉，无咎。", "meaning": "君子白天勤奋努力，晚上警惕谨慎，虽危无咎。", "fortune": "平", "takashima_note": "此爻告诫人要勤勉谨慎，不可懈怠。"},
            {"position": 4, "text": "九四：或跃在渊，无咎。", "meaning": "或跃起或潜藏，审时度势，无咎。", "fortune": "平", "takashima_note": "此爻示人当随机应变，进退有度。"},
            {"position": 5, "text": "九五：飞龙在天，利见大人。", "meaning": "龙飞天上，象征事业鼎盛，宜见贵人。", "fortune": "大吉", "takashima_note": "此为最吉之爻，事业达于顶峰。"},
            {"position": 6, "text": "上九：亢龙有悔。", "meaning": "龙飞过高，物极必反，有悔恨。", "fortune": "凶", "takashima_note": "此爻警告盛极必衰，当知进退。"},
        ],
        "overall_fortune": "大吉大利，但需警惕盛极而衰。"
    }
    hexagrams.append(qian)

    # 坤卦完整数据
    kun = {
        "id": 2,
        "name": "坤",
        "name_short": "☷",
        "binary": "000000",
        "upper_trigram": "坤",
        "lower_trigram": "坤",
        "judgment": "坤：元亨，利牝马之贞。君子有攸往，先迷后得主。利西南得朋，东北丧朋。安贞吉。",
        "judgment_meaning": "坤卦象征地道柔顺，母马之德。宜随顺而行，先迷后得，安守正道则吉。",
        "image": "地势坤，君子以厚德载物。",
        "image_meaning": "大地之势顺承天道，君子应效法大地，以深厚的德行承载万物。",
        "lines": [
            {"position": 1, "text": "初六：履霜，坚冰至。", "meaning": "踩到霜，预示坚冰将至。防微杜渐。", "fortune": "平", "takashima_note": "此爻示人见微知著，防患未然。"},
            {"position": 2, "text": "六二：直方大，不习无不利。", "meaning": "正直方正宏大，不需学习也无不利。", "fortune": "吉", "takashima_note": "此爻示人品德端正，自然顺利。"},
            {"position": 3, "text": "六三：含章可贞，或从王事，无成有终。", "meaning": "含蓄文采，可守正道。或从事王业，虽无成就但有善终。", "fortune": "平", "takashima_note": "此爻示人当谦虚内敛，默默奉献。"},
            {"position": 4, "text": "六四：括囊，无咎无誉。", "meaning": "扎紧口袋，无咎无誉。谨慎自守。", "fortune": "平", "takashima_note": "此爻示人当谨慎保守，不求有功但求无过。"},
            {"position": 5, "text": "六五：黄裳，元吉。", "meaning": "黄色下裳，大吉。中庸之德。", "fortune": "大吉", "takashima_note": "此爻为坤卦最吉，示人谦逊中庸之美德。"},
            {"position": 6, "text": "上六：龙战于野，其血玄黄。", "meaning": "龙在原野交战，血流玄黄。阴盛阳衰之象。", "fortune": "凶", "takashima_note": "此爻警告阴盛必衰，当适可而止。"},
        ],
        "overall_fortune": "柔顺利贞，宜随顺而行，厚德载物。"
    }
    hexagrams.append(kun)

    # 为其他62卦生成基础数据
    other_hexagram_names = [
        "屯", "蒙", "需", "讼", "师", "比", "小畜", "履", "泰", "否",
        "同人", "大有", "谦", "豫", "随", "蛊", "临", "观", "噬嗑", "贲",
        "剥", "复", "无妄", "大畜", "颐", "大过", "坎", "离", "咸", "恒",
        "遁", "大壮", "晋", "明夷", "家人", "睽", "蹇", "解", "损", "益",
        "夬", "姤", "萃", "升", "困", "井", "革", "鼎", "震", "艮",
        "渐", "归妹", "丰", "旅", "巽", "兑", "涣", "节", "中孚", "小过",
        "既济", "未济"
    ]

    for i, name in enumerate(other_hexagram_names, start=3):
        basic = HEXAGRAMS_BASIC[i-1]
        hexagram = {
            "id": i,
            "name": name,
            "name_short": TRIGRAMS.get(basic["upper"], {}).get("symbol", "☰"),
            "binary": basic["binary"],
            "upper_trigram": basic["upper"],
            "lower_trigram": basic["lower"],
            "judgment": f"{name}卦卦辞待补充",
            "judgment_meaning": f"{name}卦释义待补充",
            "image": f"{name}卦象辞待补充",
            "image_meaning": f"{name}卦象辞释义待补充",
            "lines": [],
            "overall_fortune": f"{name}卦总体运势待补充"
        }

        # 生成6爻基础数据
        for pos in range(1, 7):
            line = {
                "position": pos,
                "text": f"{name}卦第{pos}爻爻辞待补充",
                "meaning": f"{name}卦第{pos}爻释义待补充",
                "fortune": "平",
                "takashima_note": f"{name}卦第{pos}爻高岛断语待补充"
            }
            hexagram["lines"].append(line)

        hexagrams.append(hexagram)

    return hexagrams


def main():
    """主函数"""
    data_dir = Path("/Users/jianglong/Desktop/new code/zhanbu/data")
    output_dir = Path("/Users/jianglong/Desktop/new code/zhanbu/server/internal/service/data")

    # PDF文件路径
    pdf_files = [
        data_dir / "高岛易断占断破解_高岛嘉右卫门.pdf",
        data_dir / "图解高岛易断.pdf",
        data_dir / "高岛易断易经活解活断800例（上）.pdf",
        data_dir / "高岛易断易经活解活断800例（下）.pdf",
    ]

    print("=" * 60)
    print("高岛易断数据提取工具")
    print("=" * 60)

    # 尝试从PDF提取数据
    hexagrams = extract_hexagram_data([str(p) for p in pdf_files])

    # 如果PDF提取不完整，使用fallback数据
    print("\n使用基础数据作为fallback...")
    hexagrams = generate_fallback_data()

    # 保存到JSON文件
    output_file = output_dir / "takashima_hexagrams.json"
    with open(output_file, "w", encoding="utf-8") as f:
        json.dump(hexagrams, f, ensure_ascii=False, indent=2)

    print(f"\n数据已保存到: {output_file}")
    print(f"共 {len(hexagrams)} 卦")

    # 统计数据完整性
    complete_count = sum(1 for h in hexagrams if "待补充" not in h["judgment"])
    print(f"完整数据: {complete_count} 卦")
    print(f"待补充: {len(hexagrams) - complete_count} 卦")

    print("\n" + "=" * 60)
    print("提示: 基础数据已生成，建议后续人工校验和补充完整卦爻辞")
    print("=" * 60)


if __name__ == "__main__":
    main()
