import api from './api'
import type { TarotCard, Spread, DrawResult } from '../types/tarot'

interface ApiResponse<T> {
  code: number
  data: T
}

export async function fetchCards(): Promise<TarotCard[]> {
  const { data } = await api.get<ApiResponse<TarotCard[]>>('/tarot/cards')
  return data.data
}

export async function fetchCardById(id: number): Promise<TarotCard> {
  const { data } = await api.get<ApiResponse<TarotCard>>(`/tarot/cards/${id}`)
  return data.data
}

export async function fetchSpreads(): Promise<Spread[]> {
  const { data } = await api.get<ApiResponse<Spread[]>>('/tarot/spreads')
  return data.data
}

export async function drawCards(
  spread: string,
  question?: string
): Promise<DrawResult> {
  const { data } = await api.post<ApiResponse<DrawResult>>('/tarot/draw', {
    spread,
    question,
  })
  return data.data
}
