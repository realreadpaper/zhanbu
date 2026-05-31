import api from './api'
import type { BaZiResult } from '../types/bazi'

interface ApiResponse<T> {
  code: number
  data: T
  message: string
}

export async function calculateBaZi(
  birthDate: string,
  birthTime: string,
  gender: string = ''
): Promise<BaZiResult> {
  const { data } = await api.post<ApiResponse<BaZiResult>>('/bazi/calculate', {
    birth_date: birthDate,
    birth_time: birthTime,
    gender,
  })
  return data.data
}
