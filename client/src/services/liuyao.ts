import api from './api'
import type { LiuYaoResult, HexagramBrief } from '../types/liuyao'

interface ApiResponse<T> {
  code: number
  data: T
  message: string
}

export async function throwHexagrams(question?: string): Promise<LiuYaoResult> {
  const { data } = await api.post<ApiResponse<LiuYaoResult>>('/liuyao/throw', {
    question: question || '',
  })
  return data.data
}

export async function fetchHexagrams(): Promise<HexagramBrief[]> {
  const { data } = await api.get<ApiResponse<HexagramBrief[]>>('/liuyao/hexagrams')
  return data.data
}
