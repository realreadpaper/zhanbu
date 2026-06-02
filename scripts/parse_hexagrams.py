#!/usr/bin/env python3
"""
从提取的文本中解析64卦数据
"""

import json
import re
from pathlib import Path

# 64卦名称和对应关系
HEXAGRAM_NAMES = [
    ("乾", "乾为天"), ("坤", "坤为地"), ("屯", "水雷屯"), ("蒙", "山水蒙"),
    ("需", "水天需"), ("讼", "天水讼"), ("师", "地水师"), ("比", "水地比"),
    ("小畜", "风天小畜"), ("履", "天泽履"), ("泰", "地天泰"), ("否", "天地否"),
    ("同人", "天火同人"), ("大有", "火天大有"), ("谦", "地山谦"), ("豫", "雷地豫"),
    ("随", "泽雷随"), ("蛊", "山风蛊"), ("临", "地泽临"), ("观", "风地观"),
    ("噬嗑", "火雷噬嗑"), ("贲", "山火贲"), ("剥", "山地剥"), ("复", "地雷复"),
    ("无妄", "天雷无妄"), ("大畜", "山天大畜"), ("颐", "山雷颐"), ("大过", "泽风大过"),
    ("坎", "坎为水"), ("离", "离为火"), ("咸", "泽山咸"), ("恒", "雷风恒"),
    ("遁", "天山遁"), ("大壮", "雷天大壮"), ("晋", "火地晋"), ("明夷", "地火明夷"),
    ("家人", "风火家人"), ("睽", "火泽睽"), ("蹇", "水山蹇"), ("解", "雷水解"),
    ("损", "山泽损"), ("益", "风雷益"), ("夬", "泽天夬"), ("姤", "天风姤"),
    ("萃", "泽地萃"), ("升", "地风升"), ("困", "泽水困"), ("井", "水风井"),
    ("革", "泽火革"), ("鼎", "火风鼎"), ("震", "震为雷"), ("艮", "艮为山"),
    ("渐", "风山渐"), ("归妹", "雷泽归妹"), ("丰", "雷火丰"), ("旅", "火山旅"),
    ("巽", "巽为风"), ("兑", "兑为泽"), ("涣", "风水涣"), ("节", "水泽节"),
    ("中孚", "风泽中孚"), ("小过", "雷山小过"), ("既济", "水火既济"), ("未济", "火水未济"),
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

# 64卦二进制映射（上卦下卦组合）
HEXAGRAM_BINARY = {
    "乾为天": "111111", "坤为地": "000000", "水雷屯": "010001", "山水蒙": "100010",
    "水天需": "111010", "天水讼": "010111", "地水师": "000010", "水地比": "010000",
    "风天小畜": "111011", "天泽履": "110111", "地天泰": "111000", "天地否": "000111",
    "天火同人": "111101", "火天大有": "101111", "地山谦": "001000", "雷地豫": "000100",
    "泽雷随": "100110", "山风蛊": "011001", "地泽临": "110000", "风地观": "000011",
    "火雷噬嗑": "100101", "山火贲": "101001", "山地剥": "000001", "地雷复": "100000",
    "天雷无妄": "111100", "山天大畜": "001111", "山雷颐": "100001", "泽风大过": "011110",
    "坎为水": "010010", "离为火": "101101", "泽山咸": "001110", "雷风恒": "011100",
    "天山遁": "111100", "雷天大壮": "001111", "火地晋": "101000", "地火明夷": "000101",
    "风火家人": "110101", "火泽睽": "101011", "水山蹇": "001010", "雷水解": "010100",
    "山泽损": "110001", "风雷益": "100011", "泽天夬": "111110", "天风姤": "011111",
    "泽地萃": "000110", "地风升": "011000", "泽水困": "010110", "水风井": "011010",
    "泽火革": "101110", "火风鼎": "011101", "震为雷": "100100", "艮为山": "001001",
    "风山渐": "001011", "雷泽归妹": "110100", "雷火丰": "101100", "火山旅": "001101",
    "巽为风": "011011", "兑为泽": "110110", "风水涣": "010011", "水泽节": "110010",
    "风泽中孚": "110011", "雷山小过": "001100", "水火既济": "101010", "火水未济": "010101",
}


def parse_text_file(file_path: str) -> str:
    """读取文本文件"""
    with open(file_path, "r", encoding="utf-8") as f:
        return f.read()


def extract_hexagram_sections(text: str) -> list:
    """提取每个卦的文本段落"""
    sections = []

    # 匹配卦名模式：中文名 + 三位数字编号
    pattern = r'([一-龥]+)\s*\n(\d{2,3})\s*\n'

    matches = list(re.finditer(pattern, text))

    for i, match in enumerate(matches):
        hexagram_name = match.group(1).strip()
        hexagram_number = int(match.group(2))

        # 验证是否是有效的卦名
        valid_names = [name for _, name in HEXAGRAM_NAMES]
        if hexagram_name in valid_names and 1 <= hexagram_number <= 64:
            start = match.end()
            end = matches[i + 1].start() if i + 1 < len(matches) else len(text)
            section_text = text[start:end].strip()

            sections.append({
                "name": hexagram_name,
                "number": hexagram_number,
                "text": section_text
            })

    return sections


def parse_hexagram_section(section: dict) -> dict:
    """解析单个卦的数据"""
    name = section["name"]
    number = section["number"]
    text = section["text"]

    # 查找对应的短名
    short_name = None
    for short, full in HEXAGRAM_NAMES:
        if full == name:
            short_name = short
            break

    # 获取二进制
    binary = HEXAGRAM_BINARY.get(name, "000000")

    # 获取上下卦
    upper_trigram = ""
    lower_trigram = ""
    for tri_name, tri_data in TRIGRAMS.items():
        if tri_name in name:
            if not upper_trigram:
                upper_trigram = tri_name
            else:
                lower_trigram = tri_name

    # 解析卦辞（第一行非空文本，直到遇到爻辞）
    lines = text.split('\n')
    judgment = ""
    judgment_lines = []

    for line in lines:
        line = line.strip()
        if not line:
            continue
        # 如果遇到爻辞，停止
        if re.match(r'^(初|九|六|上)[六九一二三四五]', line):
            break
        judgment_lines.append(line)

    judgment = " ".join(judgment_lines).strip()

    # 解析爻辞
    line_texts = []
    current_line = ""

    for line in lines:
        line = line.strip()
        if not line:
            continue

        # 匹配爻辞
        line_match = re.match(r'^(初[六九]|[六九][二三四五]|上[六九])[：:](.+)$', line)
        if line_match:
            if current_line:
                line_texts.append(current_line)
            current_line = line_match.group(0)
        elif current_line:
            current_line += " " + line

    if current_line:
        line_texts.append(current_line)

    # 确保有6个爻辞
    while len(line_texts) < 6:
        line_texts.append(f"{name}卦第{len(line_texts) + 1}爻爻辞待补充")

    # 构建结果
    result = {
        "id": number,
        "name": short_name or name[:2],
        "name_short": TRIGRAMS.get(upper_trigram, {}).get("symbol", "☰"),
        "binary": binary,
        "upper_trigram": upper_trigram,
        "lower_trigram": lower_trigram,
        "judgment": judgment,
        "judgment_meaning": f"{name}卦卦辞释义待补充",
        "image": f"{name}卦象辞待补充",
        "image_meaning": f"{name}卦象辞释义待补充",
        "lines": [],
        "overall_fortune": f"{name}卦总体运势待补充"
    }

    # 解析每个爻
    line_positions = ["初", "二", "三", "四", "五", "上"]
    for i, line_text in enumerate(line_texts[:6]):
        # 提取爻辞内容
        content = re.sub(r'^[初九六上][六九一二三四五]+[：:]\s*', '', line_text)

        line_data = {
            "position": i + 1,
            "text": line_text,
            "meaning": content,
            "fortune": "平",  # 默认
            "takashima_note": f"{name}卦第{i + 1}爻高岛断语待补充"
        }
        result["lines"].append(line_data)

    return result


def main():
    """主函数"""
    input_file = Path("/Users/jianglong/Desktop/new code/zhanbu/scripts/extracted_texts/高岛易断占断破解_高岛嘉右卫门.txt")
    output_file = Path("/Users/jianglong/Desktop/new code/zhanbu/server/internal/service/data/takashima_hexagrams.json")

    print("=" * 60)
    print("解析高岛易断64卦数据")
    print("=" * 60)

    # 读取文本
    text = parse_text_file(str(input_file))
    print(f"文本长度: {len(text)} 字符")

    # 提取卦段落
    sections = extract_hexagram_sections(text)
    print(f"提取到 {len(sections)} 个卦")

    # 解析每个卦
    hexagrams = []
    for section in sections:
        hexagram = parse_hexagram_section(section)
        hexagrams.append(hexagram)
        print(f"  {hexagram['id']:2d}. {hexagram['name']} - 卦辞: {hexagram['judgment'][:30]}...")

    # 按ID排序
    hexagrams.sort(key=lambda x: x["id"])

    # 保存到文件
    with open(output_file, "w", encoding="utf-8") as f:
        json.dump(hexagrams, f, ensure_ascii=False, indent=2)

    print(f"\n数据已保存到: {output_file}")
    print(f"共 {len(hexagrams)} 卦")

    # 统计
    complete = sum(1 for h in hexagrams if "待补充" not in h["judgment"])
    print(f"完整卦辞: {complete} 卦")
    print(f"待补充: {len(hexagrams) - complete} 卦")

    print("\n" + "=" * 60)


if __name__ == "__main__":
    main()
