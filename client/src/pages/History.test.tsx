import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { beforeEach, expect, test, vi } from 'vitest'
import History from './History'
import { fetchHistory, fetchHistoryDetail, deleteHistory } from '../services/history'

vi.mock('../services/history', () => ({
  fetchHistory: vi.fn(),
  fetchHistoryDetail: vi.fn(),
  deleteHistory: vi.fn(),
}))

const mockedFetchHistory = vi.mocked(fetchHistory)
const mockedFetchHistoryDetail = vi.mocked(fetchHistoryDetail)
const mockedDeleteHistory = vi.mocked(deleteHistory)

beforeEach(() => {
  mockedFetchHistory.mockReset()
  mockedFetchHistoryDetail.mockReset()
  mockedDeleteHistory.mockReset()
})

test('loads and displays full history detail when a record is clicked', async () => {
  mockedFetchHistory.mockResolvedValue({
    total: 1,
    page: 1,
    page_size: 10,
    items: [
      {
        id: 7,
        type: 'liuyao_v2',
        type_cn: '高岛易断',
        question: '昨天认识一个女生，关系能长久吗',
        summary: '{"question":"昨天认识一个女生"}',
        created_at: '2026-06-07 10:44:20',
      },
    ],
  })
  mockedFetchHistoryDetail.mockResolvedValue({
    id: 7,
    user_id: 1,
    type: 'liuyao_v2',
    question: '昨天认识一个女生，关系能长久吗',
    result: JSON.stringify({
      ben_gua: { name: '泽山咸', name_short: '咸' },
      bian_gua: { name: '雷风恒', name_short: '恒' },
      mutable_lines: [2],
    }),
    ai_reading: '这段关系需要顺势推进，先稳住节奏。',
    created_at: '2026-06-07T10:44:20Z',
  })

  render(<History />)

  const item = await screen.findByText('昨天认识一个女生，关系能长久吗')
  fireEvent.click(item)

  await waitFor(() => {
    expect(mockedFetchHistoryDetail).toHaveBeenCalledWith(7)
  })
  expect(await screen.findByText('记录详情')).toBeTruthy()
  expect(screen.getAllByText(/泽山咸.*雷风恒/).length).toBeGreaterThan(0)
  expect(screen.getByText('这段关系需要顺势推进，先稳住节奏。')).toBeTruthy()
})
