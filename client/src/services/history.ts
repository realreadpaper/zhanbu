import api from './api'

interface ApiResponse<T> {
  code: number
  data: T
  message: string
}

export interface HistoryItem {
  id: number
  type: string
  type_cn: string
  question: string
  summary: string
  created_at: string
}

export interface HistoryListResponse {
  total: number
  page: number
  page_size: number
  items: HistoryItem[]
}

export interface DivinationRecord {
  id: number
  user_id: number
  type: string
  question: string
  result: string
  ai_reading: string
  created_at: string
  prompt_profile_id?: string
  prompt_profile_name?: string
  prompt_profile_version?: string
}

export async function fetchHistory(
  type?: string,
  page: number = 1,
  pageSize: number = 10,
  startDate?: string,
  endDate?: string
): Promise<HistoryListResponse> {
  const params: Record<string, string | number> = { page, page_size: pageSize }
  if (type) {
    params.type = type
  }
  if (startDate) {
    params.start_date = startDate
  }
  if (endDate) {
    params.end_date = endDate
  }
  const { data } = await api.get<ApiResponse<HistoryListResponse>>('/history', { params })
  return data.data
}

export async function fetchHistoryDetail(id: number): Promise<DivinationRecord> {
  const { data } = await api.get<ApiResponse<DivinationRecord>>(`/history/${id}`)
  return data.data
}

export async function deleteHistory(id: number): Promise<void> {
  await api.delete(`/history/${id}`)
}
