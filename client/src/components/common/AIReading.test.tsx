import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, expect, test, vi } from 'vitest'
import AIReading from './AIReading'
import { useAI } from '../../hooks/useAI'

vi.mock('../../hooks/useAI', () => ({
  useAI: vi.fn(),
}))

const mockedUseAI = vi.mocked(useAI)

beforeEach(() => {
  cleanup()
  mockedUseAI.mockReset()
})

test('shows incomplete status instead of completed status for truncated readings', () => {
  mockedUseAI.mockReturnValue({
    text: '建议先稳住节奏。\n\n【系统提示：AI 输出达到长度上限，内容可能未完整生成，请稍后重新解读。】',
    isStreaming: false,
    isDone: true,
    error: null,
    start: vi.fn(),
    reset: vi.fn(),
  })

  render(<AIReading type="liuyao_v2" resultId={1} question="关系能长久吗" />)

  expect(screen.getByText('解读未完整')).toBeTruthy()
  expect(screen.queryByText('解读完成')).toBeNull()
})

test('regenerates with force when clicking reread on existing text', () => {
  const start = vi.fn()
  mockedUseAI.mockReturnValue({
    text: '旧的半截解读',
    isStreaming: false,
    isDone: true,
    error: null,
    start,
    reset: vi.fn(),
  })

  render(<AIReading type="liuyao_v2" resultId={1} question="关系能长久吗" />)
  fireEvent.click(screen.getByText('重新解读'))

  expect(start).toHaveBeenCalledWith(true)
})
