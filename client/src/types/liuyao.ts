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
