export interface TarotCard {
  id: number
  name: string
  name_en: string
  type: 'major' | 'minor'
  suit?: 'wands' | 'cups' | 'swords' | 'pentacles'
  number?: number
  image?: string
  keywords_up: string
  keywords_down: string
  meaning_up: string
  meaning_down: string
  description: string
}

export interface SpreadPosition {
  index: number
  name: string
  description: string
}

export interface Spread {
  id: string
  name: string
  count: number
  positions: SpreadPosition[]
}

export interface DrawnCard {
  position: number
  position_name: string
  card: TarotCard
  orientation: string // 'upright' | 'reversed'
}

export interface DrawResult {
  spread: string
  question?: string
  cards: DrawnCard[]
  timestamp: string
  record_id?: number
}
