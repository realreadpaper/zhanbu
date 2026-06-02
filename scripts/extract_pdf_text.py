#!/usr/bin/env python3
"""
提取PDF文本内容，用于后续分析
"""

import pdfplumber
import os
from pathlib import Path

data_dir = Path("/Users/jianglong/Desktop/new code/zhanbu/data")
output_dir = Path("/Users/jianglong/Desktop/new code/zhanbu/scripts/extracted_texts")
output_dir.mkdir(exist_ok=True)

pdf_files = [
    "高岛易断占断破解_高岛嘉右卫门.pdf",
    "图解高岛易断.pdf",
    "高岛易断易经活解活断800例（上）.pdf",
    "高岛易断易经活解活断800例（下）.pdf",
]

for pdf_name in pdf_files:
    pdf_path = data_dir / pdf_name
    output_file = output_dir / f"{pdf_name.replace('.pdf', '.txt')}"

    print(f"\n{'='*60}")
    print(f"正在处理: {pdf_name}")
    print(f"{'='*60}")

    try:
        with pdfplumber.open(pdf_path) as pdf:
            total_pages = len(pdf.pages)
            print(f"总页数: {total_pages}")

            text_content = []
            for i, page in enumerate(pdf.pages):
                try:
                    page_text = page.extract_text()
                    if page_text:
                        text_content.append(f"\n--- 第 {i+1} 页 ---\n")
                        text_content.append(page_text)
                except Exception as e:
                    print(f"  第 {i+1} 页提取失败: {e}")

                if (i + 1) % 20 == 0:
                    print(f"  已处理 {i+1}/{total_pages} 页")

            # 保存到文件
            with open(output_file, "w", encoding="utf-8") as f:
                f.write("\n".join(text_content))

            print(f"已保存到: {output_file}")
            print(f"提取文本长度: {len(''.join(text_content))} 字符")

    except Exception as e:
        print(f"处理失败: {e}")

print("\n" + "="*60)
print("所有PDF处理完成！")
print("="*60)
