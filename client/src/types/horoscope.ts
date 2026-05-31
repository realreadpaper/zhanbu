export interface HoroscopeDetail {
  love: string
  career: string
  wealth: string
  health: string
}

export interface HoroscopeResult {
  zodiac: string
  zodiac_cn: string
  period: 'daily' | 'weekly' | 'monthly'
  date: string
  overall: number
  love: number
  career: number
  wealth: number
  health: number
  lucky_number: number
  lucky_color: string
  summary: string
  detail: HoroscopeDetail
}

export interface ZodiacSign {
  name: string
  name_cn: string
  symbol: string
  element: 'fire' | 'earth' | 'air' | 'water'
  dateRange: string
}

export const ZODIAC_SIGNS: ZodiacSign[] = [
  { name: 'aries', name_cn: '白羊座', symbol: '♈', element: 'fire', dateRange: '3.21 - 4.19' },
  { name: 'taurus', name_cn: '金牛座', symbol: '♉', element: 'earth', dateRange: '4.20 - 5.20' },
  { name: 'gemini', name_cn: '双子座', symbol: '♊', element: 'air', dateRange: '5.21 - 6.21' },
  { name: 'cancer', name_cn: '巨蟹座', symbol: '♋', element: 'water', dateRange: '6.22 - 7.22' },
  { name: 'leo', name_cn: '狮子座', symbol: '♌', element: 'fire', dateRange: '7.23 - 8.22' },
  { name: 'virgo', name_cn: '处女座', symbol: '♍', element: 'earth', dateRange: '8.23 - 9.22' },
  { name: 'libra', name_cn: '天秤座', symbol: '♎', element: 'air', dateRange: '9.23 - 10.23' },
  { name: 'scorpio', name_cn: '天蝎座', symbol: '♏', element: 'water', dateRange: '10.24 - 11.22' },
  { name: 'sagittarius', name_cn: '射手座', symbol: '♐', element: 'fire', dateRange: '11.23 - 12.21' },
  { name: 'capricorn', name_cn: '摩羯座', symbol: '♑', element: 'earth', dateRange: '12.22 - 1.19' },
  { name: 'aquarius', name_cn: '水瓶座', symbol: '♒', element: 'air', dateRange: '1.20 - 2.18' },
  { name: 'pisces', name_cn: '双鱼座', symbol: '♓', element: 'water', dateRange: '2.19 - 3.20' },
]

export const ELEMENT_COLORS: Record<string, string> = {
  fire: '#ef4444',
  earth: '#a3e635',
  air: '#38bdf8',
  water: '#818cf8',
}
