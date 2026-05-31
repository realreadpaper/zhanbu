import api from './api'
import type { HoroscopeResult } from '../types/horoscope'

interface ApiResponse<T> {
  code: number
  data: T
}

export async function fetchHoroscope(
  zodiac: string,
  period: 'daily' | 'weekly' | 'monthly' = 'daily',
  date?: string
): Promise<HoroscopeResult> {
  const params: Record<string, string> = { period }
  if (date) {
    params.date = date
  }
  const { data } = await api.get<ApiResponse<HoroscopeResult>>(
    `/horoscope/${zodiac}`,
    { params }
  )
  return data.data
}
