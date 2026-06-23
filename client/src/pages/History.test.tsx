import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
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

  render(
    <MemoryRouter>
      <History />
    </MemoryRouter>
  )

  const item = await screen.findByText('昨天认识一个女生，关系能长久吗')
  fireEvent.click(item)

  await waitFor(() => {
    expect(mockedFetchHistoryDetail).toHaveBeenCalledWith(7)
  })
  expect(await screen.findByText('记录详情')).toBeTruthy()
  expect(screen.getAllByText(/泽山咸.*雷风恒/).length).toBeGreaterThan(0)
  expect(screen.getByText('这段关系需要顺势推进，先稳住节奏。')).toBeTruthy()
})

test('displays meihua history detail with chinese lunar time', async () => {
  mockedFetchHistory.mockResolvedValue({
    total: 1,
    page: 1,
    page_size: 10,
    items: [
      {
        id: 9,
        type: 'meihua',
        type_cn: '梅花易数',
        question: '最近事业如何',
        summary: '泽火革',
        created_at: '2026-06-23 23:30:00',
      },
    ],
  })
  mockedFetchHistoryDetail.mockResolvedValue({
    id: 9,
    user_id: 1,
    type: 'meihua',
    question: '最近事业如何',
    result: JSON.stringify({
      source_values: { lunar_display: '丙午年五月初九子时' },
      ben_gua: { name: '泽火革' },
      hu_gua: { name: '天风姤' },
      bian_gua: { name: '泽山咸' },
      moving_line: 1,
      ti_yong: {
        ti: { name: '兑' },
        yong: { name: '离' },
        relation: '用克体',
      },
    }),
    ai_reading: '此卦宜先稳后动。',
    created_at: '2026-06-23T23:30:00Z',
  })

  render(
    <MemoryRouter>
      <History />
    </MemoryRouter>
  )

  const item = await screen.findByText('最近事业如何')
  fireEvent.click(item)

  expect(await screen.findByText('记录详情')).toBeTruthy()
  expect(screen.getByText(/起卦时间：丙午年五月初九子时/)).toBeTruthy()
  expect(screen.getByText(/本卦：泽火革/)).toBeTruthy()
  expect(screen.getByText(/动爻：初爻动/)).toBeTruthy()
  expect(screen.getByText('此卦宜先稳后动。')).toBeTruthy()
})
