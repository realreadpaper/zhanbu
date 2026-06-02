#!/usr/bin/env python3
"""
Extract Takashima I Ching content from the full PDF into traceable data files.

This script intentionally avoids invented fallback text. Empty fields stay empty
and are reported in the quality report.
"""

from __future__ import annotations

import argparse
import json
import re
import subprocess
from dataclasses import dataclass
from pathlib import Path
from typing import Any


DEFAULT_PDF = Path("data/高岛易断占断破解_高岛嘉右卫门.pdf")
DEFAULT_OUT = Path("output/takashima")
DEFAULT_SERVER_DATA = Path("server/internal/service/data")

NATURE_TO_TRIGRAM = {
    "天": ("乾", "111", "☰", "金"),
    "泽": ("兑", "110", "☱", "金"),
    "火": ("离", "101", "☲", "火"),
    "雷": ("震", "100", "☳", "木"),
    "风": ("巽", "011", "☴", "木"),
    "水": ("坎", "010", "☵", "水"),
    "山": ("艮", "001", "☶", "土"),
    "地": ("坤", "000", "☷", "土"),
}

SPECIAL_FULL_NAMES = {
    "乾为天": ("乾", "天", "天"),
    "坤为地": ("坤", "地", "地"),
    "坎为水": ("坎", "水", "水"),
    "离为火": ("离", "火", "火"),
    "震为雷": ("震", "雷", "雷"),
    "艮为山": ("艮", "山", "山"),
    "巽为风": ("巽", "风", "风"),
    "兑为泽": ("兑", "泽", "泽"),
}

LINE_NAMES = ["初九", "初六", "九二", "六二", "九三", "六三", "九四", "六四", "九五", "六五", "上九", "上六", "用九", "用六"]
LINE_NAME_RE = "|".join(LINE_NAMES)
LABEL_SEPARATOR_RE = r"[:：；;]"
LINE_LABEL_SEPARATOR_RE = r"[:：；;，,]"
PLACEHOLDER_RE = re.compile(r"第\d爻(?:含义|高岛断语)|总体运势|象征意义|待补充")


@dataclass
class LocatedMatch:
    page: int
    start: int
    end: int
    text: str


def run_pdftotext(pdf: Path) -> str:
    if not pdf.exists():
        raise FileNotFoundError(pdf)
    proc = subprocess.run(
        ["pdftotext", "-layout", "-enc", "UTF-8", str(pdf), "-"],
        check=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    return proc.stdout.decode("utf-8", errors="replace")


def normalize_text(text: str) -> str:
    text = text.replace("\u3000", " ")
    text = re.sub(r"[ \t]+\n", "\n", text)
    text = re.sub(r"\n{3,}", "\n\n", text)
    return text.strip()


def clean_pdf_artifacts(text: str) -> str:
    """Remove OCR/layout artifacts while preserving the extracted wording."""
    if not text:
        return ""

    text = strip_page_markers(text)
    text = normalize_text(text)
    text = re.sub(r"(?m)^\s*[“”\"']+\s*$\n?", "", text)
    text = re.sub(r"曰\s*大\s*[“”\"']?\s*有\s*。", "曰大有。", text)

    previous = None
    while previous != text:
        previous = text
        text = re.sub(r"([\u4e00-\u9fff])\s+([\u4e00-\u9fff])", r"\1\2", text)

    text = re.sub(r"([\u4e00-\u9fff])\s+([，。！？；：、])", r"\1\2", text)
    text = re.sub(r"([，。！？；：、])\s+([\u4e00-\u9fff])", r"\1\2", text)
    text = re.sub(r"《\s+", "《", text)
    text = re.sub(r"\s+》", "》", text)
    text = re.sub(r"\n{3,}", "\n\n", text)
    return text.strip()


def clean_short_text(text: str) -> str:
    text = clean_pdf_artifacts(text)
    text = re.sub(r"\s*\n\s*", "", text)
    return text.strip()


def split_pages(raw_text: str) -> list[dict[str, Any]]:
    raw_pages = raw_text.split("\f")
    if raw_pages and raw_pages[-1].strip() == "":
        raw_pages = raw_pages[:-1]
    return [
        {
            "page": i + 1,
            "text": page.strip("\n"),
            "char_count": len(page.strip()),
        }
        for i, page in enumerate(raw_pages)
    ]


def make_global_text(pages: list[dict[str, Any]]) -> tuple[str, list[tuple[int, int, int]]]:
    chunks: list[str] = []
    spans: list[tuple[int, int, int]] = []
    cursor = 0
    for page in pages:
        marker = f"\n\n[[PAGE {page['page']:04d}]]\n"
        chunks.append(marker)
        cursor += len(marker)
        text = page["text"]
        start = cursor
        chunks.append(text)
        cursor += len(text)
        spans.append((page["page"], start, cursor))
    return "".join(chunks), spans


def page_for_offset(offset: int, spans: list[tuple[int, int, int]]) -> int:
    last_page = spans[0][0] if spans else 1
    for page, start, end in spans:
        if start <= offset <= end:
            return page
        if offset >= start:
            last_page = page
    return last_page


def pages_for_range(start: int, end: int, spans: list[tuple[int, int, int]]) -> list[int]:
    found = [page for page, p_start, p_end in spans if p_start <= end and p_end >= start]
    return found or [page_for_offset(start, spans)]


def infer_hexagram_name(full_name: str) -> tuple[str, str, str]:
    if full_name in SPECIAL_FULL_NAMES:
        return SPECIAL_FULL_NAMES[full_name]
    if len(full_name) >= 3:
        upper_nature = full_name[0]
        lower_nature = full_name[1]
        short_name = full_name[2:]
        if upper_nature in NATURE_TO_TRIGRAM and lower_nature in NATURE_TO_TRIGRAM:
            return short_name, upper_nature, lower_nature
    return full_name, "", ""


def trigram_meta(nature: str) -> tuple[str, str, str, str]:
    return NATURE_TO_TRIGRAM.get(nature, ("", "", "", ""))


def find_toc_titles(pages: list[dict[str, Any]]) -> list[dict[str, Any]]:
    toc_text = "\n".join(page["text"] for page in pages[:25])
    titles: dict[int, str] = {}
    for match in re.finditer(r"(?m)^\s*(\d{2})\s+([^\s，。:：]+)\s*$", toc_text):
        idx = int(match.group(1))
        name = match.group(2).strip()
        if 1 <= idx <= 64 and idx not in titles:
            titles[idx] = name
    return [{"id": idx, "full_name": titles[idx]} for idx in sorted(titles)]


def find_body_starts(global_text: str, spans: list[tuple[int, int, int]], titles: list[dict[str, Any]]) -> list[LocatedMatch]:
    starts: list[LocatedMatch] = []
    for title in titles:
        idx = title["id"]
        full_name = title["full_name"]
        pattern = re.compile(rf"(?m)^\s*{idx:02d}\s+{re.escape(full_name)}\s*$")
        candidates = []
        for match in pattern.finditer(global_text):
            page = page_for_offset(match.start(), spans)
            if page >= 45:
                candidates.append(LocatedMatch(page, match.start(), match.end(), match.group(0)))
        if candidates:
            starts.append(candidates[0])
    starts.sort(key=lambda m: m.start)
    return starts


def extract_quoted_after(marker: str, text: str) -> str:
    idx = text.find(marker)
    if idx < 0:
        return ""
    tail = text[idx + len(marker):].strip()
    tail = re.sub(r"^[：:；;]\s*", "", tail)
    enders = ["\n  ", "\n《", "\n【", "\n\n"]
    ends = [tail.find(e) for e in enders if tail.find(e) > 0]
    end = min(ends) if ends else min(len(tail), 500)
    return clean_short_text(tail[:end])


def extract_first_after(markers: list[str], text: str) -> str:
    candidates: list[tuple[int, str]] = []
    for marker in markers:
        idx = text.find(marker)
        if idx >= 0:
            candidates.append((idx, marker))
    if not candidates:
        return ""
    _, marker = min(candidates, key=lambda item: item[0])
    return extract_quoted_after(marker, text)


def extract_image_after(markers: list[str], text: str) -> str:
    candidates: list[tuple[int, str]] = []
    for marker in markers:
        idx = text.find(marker)
        if idx >= 0:
            candidates.append((idx, marker))
    if not candidates:
        return ""

    idx, marker = min(candidates, key=lambda item: item[0])
    tail = text[idx + len(marker):]
    tail = re.sub(r"^[：:；;]\s*", "", tail.strip())
    parts: list[str] = []
    for line in tail.splitlines():
        stripped = line.strip()
        if not stripped:
            if parts:
                break
            continue
        if stripped.startswith("[[PAGE"):
            continue
        if re.fullmatch(r"[“”\"' ]+", stripped):
            continue
        if stripped.startswith("《") or stripped.startswith("【"):
            break
        parts.append(stripped)
        if stripped.endswith(("。", "！", "？")):
            break
    return clean_short_text("".join(parts))


def first_nonempty_lines(text: str, max_lines: int = 12) -> list[str]:
    return [line.strip() for line in text.splitlines() if line.strip()][:max_lines]


def compact_label(label: str) -> str:
    label = re.sub(r"［.*?］|\[.*?\]|\(.*?\)|（.*?）", "", label)
    label = re.sub(r"[\s《》“”\"'、.．。]", "", label)
    return label


def line_position(name: str) -> int:
    if name.startswith("初"):
        return 1
    if name.endswith("二"):
        return 2
    if name.endswith("三"):
        return 3
    if name.endswith("四"):
        return 4
    if name.endswith("五"):
        return 5
    if name.startswith("上"):
        return 6
    if name in ("用九", "用六"):
        return 7
    return 0


def strip_page_markers(text: str) -> str:
    return re.sub(r"\n?\[\[PAGE\s+\d{4}\]\]\n?", "\n", text)


def ordered_line_matches(matches: list[re.Match[str]]) -> list[re.Match[str]]:
    selected: list[re.Match[str]] = []
    cursor = 0
    for expected_position in range(1, 7):
        next_match = None
        for match in matches:
            if match.start() < cursor:
                continue
            if line_position(match.group(1)) == expected_position:
                next_match = match
                break
        if next_match is None:
            continue
        selected.append(next_match)
        cursor = next_match.end()

    for match in matches:
        if match.start() >= cursor and match.group(1) in ("用九", "用六"):
            selected.append(match)
            break

    return selected


def find_judgment(raw: str, short_name: str, full_name: str) -> tuple[str, int | None]:
    aliases = {compact_label(short_name), compact_label(full_name)}
    if short_name == "坎":
        aliases.add("习坎")

    lines = raw.splitlines()
    char_offset = 0
    for i, line in enumerate(lines[:120]):
        stripped = line.strip()
        match = re.match(rf"^\s*(.+?){LABEL_SEPARATOR_RE}\s*(.+)$", line)
        if not match:
            char_offset += len(line) + 1
            continue

        label = compact_label(match.group(1))
        if label not in aliases and not any(label.endswith(alias) for alias in aliases if alias):
            char_offset += len(line) + 1
            continue

        text_parts = [match.group(2).strip()]
        for next_line in lines[i + 1:i + 6]:
            next_stripped = next_line.strip()
            if not next_stripped:
                break
            if (
                next_stripped.startswith("[[PAGE")
                or re.match(r"^\d+(?:\s+\d+)*$", next_stripped)
                or "篆书" in next_stripped
                or "甲骨文" in next_stripped
                or next_stripped.startswith("《")
            ):
                break
            if re.match(rf"^\s*({LINE_NAME_RE}){LINE_LABEL_SEPARATOR_RE}", next_line):
                break
            text_parts.append(next_stripped)

        return clean_short_text("\n".join(text_parts)), char_offset + match.start()

    return "", None


def split_line_blocks(block_text: str, block_start: int, spans: list[tuple[int, int, int]]) -> list[dict[str, Any]]:
    matches = ordered_line_matches(list(re.finditer(rf"(?m)^\s*({LINE_NAME_RE}){LINE_LABEL_SEPARATOR_RE}\s*", block_text)))
    lines: list[dict[str, Any]] = []
    for i, match in enumerate(matches):
        name = match.group(1)
        start = match.start()
        end = matches[i + 1].start() if i + 1 < len(matches) else len(block_text)
        line_text = block_text[start:end].strip()
        absolute_start = block_start + start
        absolute_end = block_start + end
        source_pages = pages_for_range(absolute_start, absolute_end, spans)

        original = ""
        after_label = block_text[match.end():end].strip()
        original_end_candidates = []
        for marker in ["《象传》曰", "【占】", "【例】", "\n  "]:
            pos = after_label.find(marker)
            if pos > 0:
                original_end_candidates.append(pos)
        original_end = min(original_end_candidates) if original_end_candidates else min(len(after_label), 240)
        original = clean_short_text(after_label[:original_end])

        lines.append({
            "position": line_position(name),
            "name": name,
            "original": original,
            "commentary": clean_short_text(extract_quoted_after("《象传》曰", line_text)),
            "takashima_analysis": clean_pdf_artifacts(line_text),
            "source_pages": source_pages,
            "char_count": len(line_text),
        })
    return lines


def parse_hexagram_block(
    title: dict[str, Any],
    start: LocatedMatch,
    end_offset: int,
    global_text: str,
    spans: list[tuple[int, int, int]],
) -> dict[str, Any]:
    full_name = title["full_name"]
    short_name, upper_nature, lower_nature = infer_hexagram_name(full_name)
    upper_trigram, upper_binary, upper_symbol, upper_element = trigram_meta(upper_nature)
    lower_trigram, lower_binary, lower_symbol, lower_element = trigram_meta(lower_nature)
    raw = global_text[start.end:end_offset].strip()
    source_pages = pages_for_range(start.start, end_offset, spans)

    judgment, judgment_offset = find_judgment(raw, short_name, full_name)

    first_line = re.search(rf"(?m)^\s*({LINE_NAME_RE}){LABEL_SEPARATOR_RE}\s*", raw)
    overview = raw[:first_line.start()] if first_line else raw

    tuan = extract_first_after([
        "《彖传》曰：",
        "《彖传》曰:",
        "《彖传》曰；",
        "《彖传》曰;",
        "《彖传》，曰：",
        "《彖传》，曰:",
        "《象传》曰：",
        "《象传》曰:",
        "《象传》曰；",
        "《象传》曰;",
    ], overview)
    image = extract_image_after([
        "《大象》曰：",
        "《大象》曰:",
        "《大象》曰；",
        "《大象》曰;",
        "《大彖》曰：",
        "《大彖》曰:",
        "《大彖》曰；",
        "《大彖》曰;",
        "《象》曰：",
        "《象》曰:",
    ], overview)
    lines = split_line_blocks(raw, start.end, spans)

    case_matches = []
    for m in re.finditer(r"【例】", raw):
        absolute = start.end + m.start()
        case_matches.append({
            "source_page": page_for_offset(absolute, spans),
            "offset": m.start(),
        })

    warnings: list[str] = []
    if not upper_trigram or not lower_trigram:
        warnings.append("missing_trigram")
    if not judgment:
        warnings.append("missing_judgment")
    if len([line for line in lines if line["name"] not in ("用九", "用六")]) < 6:
        warnings.append("fewer_than_six_lines")
    if PLACEHOLDER_RE.search(raw):
        warnings.append("placeholder_text_detected")

    return {
        "id": title["id"],
        "name": short_name,
        "full_name": full_name,
        "binary": upper_binary + lower_binary if upper_binary and lower_binary else "",
        "upper_trigram": upper_trigram,
        "lower_trigram": lower_trigram,
        "upper_nature": upper_nature,
        "lower_nature": lower_nature,
        "upper_symbol": upper_symbol,
        "lower_symbol": lower_symbol,
        "elements": {
            "upper": upper_element,
            "lower": lower_element,
        },
        "source": {
            "book": "高岛易断:占断破解",
            "pdf": str(DEFAULT_PDF),
            "start_page": min(source_pages),
            "end_page": max(source_pages),
            "pages": source_pages,
        },
        "judgment": {
            "text": judgment,
            "source_page": page_for_offset(start.end + judgment_offset, spans) if judgment_offset is not None else None,
        },
        "tuan": {
            "text": tuan,
        },
        "image": {
            "text": image,
        },
        "lines": lines,
        "cases": case_matches,
        "raw_text": clean_pdf_artifacts(raw),
        "quality": {
            "needs_review": bool(warnings),
            "warnings": warnings,
            "char_count": len(raw),
        },
    }


def build_sections(pages: list[dict[str, Any]], hexagrams: list[dict[str, Any]]) -> dict[str, Any]:
    def text_between(start_page: int, end_page: int) -> str:
        return normalize_text("\n".join(page["text"] for page in pages[start_page - 1:end_page]))

    return {
        "book": {
            "title": "高岛易断:占断破解",
            "total_pages": len(pages),
        },
        "front_matter": [
            {"title": "目录", "pages": [2, 22], "text": text_between(2, 22)},
            {"title": "序与编校说明", "pages": [23, 34], "text": text_between(23, 34)},
            {"title": "高岛吞象传略", "pages": [35, 45], "text": text_between(35, 45)},
            {"title": "高岛吞象易经筮法揭秘", "pages": [46, 51], "text": text_between(46, 51)},
        ],
        "hexagrams": [
            {
                "id": h["id"],
                "full_name": h["full_name"],
                "pages": [h["source"]["start_page"], h["source"]["end_page"]],
                "char_count": h["quality"]["char_count"],
                "needs_review": h["quality"]["needs_review"],
            }
            for h in hexagrams
        ],
    }


def build_interpretation_rules(pages: list[dict[str, Any]], hexagrams: list[dict[str, Any]]) -> dict[str, Any]:
    method_text = normalize_text("\n".join(page["text"] for page in pages[45:51]))
    case_count = sum(len(h["cases"]) for h in hexagrams)
    line_count = sum(len(h["lines"]) for h in hexagrams)
    return {
        "source": {
            "book": "高岛易断:占断破解",
            "method_pages": [46, 51],
        },
        "method_raw_text": method_text,
        "divination_methods": [
            {
                "name": "略筮法",
                "summary": "五十根竹签或蓍草先抽一根为太极；静默祈祷后分蓍。男子先取左手数除八为上卦，再取右手数除八为下卦；女子左右相反。再合四十九根分蓍，男子取左手数除六、女子取右手数除六为动爻。",
                "source_pages": [47, 48],
            },
            {
                "name": "八卦竹筒法",
                "summary": "一个竹筒置乾、兑、离、震、巽、坎、艮、坤八签，另一个竹筒置一至六爻。求占者静心默祷后抽上卦、下卦、动爻，男子左先右后，女子反之。",
                "source_pages": [48],
            },
            {
                "name": "六十四卦签法",
                "summary": "将六十四卦卦象与卦名刻在竹签上，静默祈祷后直接抽卦，再按男左女右抽动爻。",
                "source_pages": [48],
            },
            {
                "name": "三百八十四爻签法",
                "summary": "将三百八十四爻直接做成签，祈祷后抽出卦象与动爻，并依据爻辞直接索解。",
                "source_pages": [48, 49],
            },
            {
                "name": "时间起卦法",
                "summary": "以年月日时数取上卦、下卦与变爻；同一时辰可加入姓氏笔画或自由报数。书中也举出用时分起卦的简化方式。",
                "source_pages": [50, 51],
            },
        ],
        "interpretation_workflow": [
            {
                "step": 1,
                "name": "定心立问",
                "rule": "占前消除杂念，专注默想所问之事，书中称关键在至诚无息。",
                "source_pages": [46, 49],
            },
            {
                "step": 2,
                "name": "得本卦与动爻",
                "rule": "高岛占例大多只有一爻发动；取得本卦后重点定位动爻。",
                "source_pages": [46, 48],
            },
            {
                "step": 3,
                "name": "先读卦辞爻辞",
                "rule": "找出对应卦辞、动爻爻辞，并以经文作为判断核心。",
                "source_pages": [48, 49],
            },
            {
                "step": 4,
                "name": "以事拟象",
                "rule": "抓住当前问题重点，依据爻辞进行模拟，使经文、卦象与所问人事相互对应。",
                "source_pages": [49],
            },
            {
                "step": 5,
                "name": "参考占例迁移",
                "rule": "每卦每爻的【占】与【例】作为高岛断法案例库，按问题类型迁移其推断路径。",
                "source_pages": [54, 671],
            },
        ],
        "extracted_principles": [
            {
                "name": "至诚起占",
                "rule": "起占前须静心诚意，以所问之事专一其念；脚本保留原文，具体断法仍以书中筮法段落为准。",
                "source_pages": [46, 51],
            },
            {
                "name": "卦辞爻辞为核心证据",
                "rule": "解读必须先引用本卦、动爻与变卦原文，再展开象数、人事、占例推断；不得脱离经文凭空发挥。",
                "source_pages": [46, 51],
            },
            {
                "name": "按问题迁移占例",
                "rule": "每卦每爻中的【占】和【例】是高岛式判断的案例库；现代问题应先匹配问题类型，再借鉴对应占断逻辑。",
                "source_pages": [46, 51],
            },
            {
                "name": "动爻优先",
                "rule": "有动爻时，以动爻爻辞及其高岛分析为判断重点，再用本卦总论与变卦趋势校正。",
                "source_pages": [46, 51],
            },
        ],
        "corpus_stats": {
            "hexagram_count": len(hexagrams),
            "line_block_count": line_count,
            "case_marker_count": case_count,
        },
    }


def quality_report(pages: list[dict[str, Any]], titles: list[dict[str, Any]], starts: list[LocatedMatch], hexagrams: list[dict[str, Any]]) -> dict[str, Any]:
    warnings: list[dict[str, Any]] = []
    for h in hexagrams:
        for warning in h["quality"]["warnings"]:
            warnings.append({"id": h["id"], "full_name": h["full_name"], "warning": warning})
    return {
        "pages": len(pages),
        "toc_title_count": len(titles),
        "body_start_count": len(starts),
        "hexagram_count": len(hexagrams),
        "warnings_count": len(warnings),
        "warnings": warnings,
        "placeholder_hits": sum(1 for h in hexagrams if PLACEHOLDER_RE.search(json.dumps(h, ensure_ascii=False))),
        "line_blocks_total": sum(len(h["lines"]) for h in hexagrams),
        "case_markers_total": sum(len(h["cases"]) for h in hexagrams),
    }


def write_json(path: Path, data: Any) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--pdf", type=Path, default=DEFAULT_PDF)
    parser.add_argument("--out", type=Path, default=DEFAULT_OUT)
    parser.add_argument("--server-data", type=Path, default=DEFAULT_SERVER_DATA)
    args = parser.parse_args()

    raw_text = run_pdftotext(args.pdf)
    pages = split_pages(raw_text)
    global_text, spans = make_global_text(pages)
    titles = find_toc_titles(pages)
    starts = find_body_starts(global_text, spans, titles)

    start_by_id = {i + 1: start for i, start in enumerate(starts)}
    hexagrams = []
    for idx, title in enumerate(titles):
        start = start_by_id.get(title["id"])
        if start is None:
            continue
        if idx + 1 < len(titles):
            next_start = start_by_id.get(titles[idx + 1]["id"])
            end_offset = next_start.start if next_start else len(global_text)
        else:
            end_offset = len(global_text)
        hexagrams.append(parse_hexagram_block(title, start, end_offset, global_text, spans))

    sections = build_sections(pages, hexagrams)
    rules = build_interpretation_rules(pages, hexagrams)
    report = quality_report(pages, titles, starts, hexagrams)

    write_json(args.out / "book_full_text.json", pages)
    write_json(args.out / "book_sections.json", sections)
    write_json(args.out / "takashima_hexagrams.json", hexagrams)
    write_json(args.out / "takashima_interpretation_rules.json", rules)
    write_json(args.out / "quality_report.json", report)
    write_json(args.server_data / "takashima_hexagrams.json", hexagrams)
    write_json(args.server_data / "takashima_interpretation_rules.json", rules)

    print(json.dumps(report, ensure_ascii=False, indent=2))
    return 0 if report["hexagram_count"] == 64 and report["placeholder_hits"] == 0 else 1


if __name__ == "__main__":
    raise SystemExit(main())
