import api from './api'
import type { LiuYaoResult, HexagramBrief, LiuYaoV2Result, TakashimaHexagram, LiuYaoV2Config } from '../types/liuyao'

interface ApiResponse<T> {
  code: number
  data: T
  message: string
}

// ==================== v1 接口 ====================

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

// ==================== v2 接口 (高岛易断) ====================

export async function throwHexagramsV2(question?: string, method?: string): Promise<LiuYaoV2Result> {
  const { data } = await api.post<ApiResponse<LiuYaoV2Result>>('/liuyao/v2/throw', {
    question: question || '',
    method: method || '',
  })
  return data.data
}

export async function fetchHexagramsV2(): Promise<TakashimaHexagram[]> {
  const { data } = await api.get<ApiResponse<TakashimaHexagram[]>>('/liuyao/v2/hexagrams')
  return data.data
}

export async function fetchHexagramByIdV2(id: number): Promise<TakashimaHexagram> {
  const { data } = await api.get<ApiResponse<TakashimaHexagram>>(`/liuyao/v2/hexagrams/${id}`)
  return data.data
}

export async function fetchLiuYaoConfig(): Promise<LiuYaoV2Config> {
  const { data } = await api.get<ApiResponse<LiuYaoV2Config>>('/liuyao/v2/config')
  return data.data
}
