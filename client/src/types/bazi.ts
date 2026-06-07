export interface BaZiPillar {
  tian_gan: string
  di_zhi: string
  wu_xing: string
  na_yin: string
  hidden_gan?: string[]
}

export interface BaZiPillars {
  year: BaZiPillar
  month: BaZiPillar
  day: BaZiPillar
  hour: BaZiPillar
}

export interface FiveElementAnalysis {
  metal: number
  wood: number
  water: number
  fire: number
  earth: number
  strongest: string
  weakest: string
  day_master: string
  strength: string
  yong_shen: string
  ji_shen: string
}

export interface TenGod {
  position: string
  tian_gan: string
  god: string
}

export interface BaZiBirthInfo {
  solar: string
  lunar: string
}

export interface BaZiResult {
  record_id?: number
  birth: BaZiBirthInfo
  pillars: BaZiPillars
  five_elements: FiveElementAnalysis
  ten_gods: TenGod[]
}
