export interface Trigram {
  id: number
  name: string
  symbol: string
  binary: number
  nature: string
  element: string
  yin_yang: string
}

export interface Hexagram {
  id: number
  name: string
  name_short: string
  upper_trigram: string
  lower_trigram: string
  binary: string
  judgment: string
  image: string
  line_texts: string[]
  description: string
}

export interface HexagramBrief {
  id: number
  name: string
  name_short: string
  upper_trigram: string
  lower_trigram: string
  binary: string
  description: string
}

export interface LineResult {
  value: number   // 6/7/8/9
  type: string    // old_yang/young_yang/old_yin/young_yin
  mutable: boolean
  symbol: string  // ⚊ or ⚋
}

export interface LiuYaoResult {
  question?: string
  lines: LineResult[]
  ben_gua: Hexagram
  bian_gua?: Hexagram
  mutable_lines: number[]
  timestamp: string
}

// ==================== 高岛易断 v2 类型 ====================

export interface TakashimaLine {
  position: number
  name: string
  original: string
  commentary: string
  takashima_analysis: string
  source_pages: number[]
}

export interface TakashimaTextSource {
  text: string
  source_page?: number
}

export type TakashimaText = string | TakashimaTextSource

export interface TakashimaSource {
  book?: string
  pdf?: string
  start_page?: number
  end_page?: number
  pages?: number[]
}

export interface TakashimaEvidenceSnippet {
  kind: string
  title: string
  text: string
  source_pages: number[]
  score: number
}

export interface TakashimaBookEvidence {
  query_terms: string[]
  snippets: TakashimaEvidenceSnippet[]
  method_rules?: TakashimaEvidenceSnippet[]
}

export interface TakashimaHexagram {
  id: number
  name: string
  name_short?: string
  full_name: string
  binary: string
  upper_trigram: string
  lower_trigram: string
  upper_nature: string
  lower_nature: string
  upper_symbol: string
  lower_symbol: string
  source: TakashimaSource
  judgment: TakashimaText
  tuan: TakashimaText
  image: TakashimaText
  lines: TakashimaLine[]
  raw_text: string
}

export interface LiuYaoV2Result {
  question?: string
  lines: LineResult[]
  ben_gua: TakashimaHexagram
  bian_gua?: TakashimaHexagram
  mutable_lines: number[]
  book_evidence?: TakashimaBookEvidence
  method: string
  timestamp: string
}

export interface LiuYaoV2Config {
  version: string
  method: string
}
