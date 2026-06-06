import { useState, useCallback } from 'react'
import { drawCards } from '../services/tarot'
import type { DrawResult } from '../types/tarot'

export type DivinationPhase = 'input' | 'drawing' | 'result'

interface UseTarotReturn {
  phase: DivinationPhase
  selectedSpread: string
  question: string
  drawResult: DrawResult | null
  revealedCards: Set<number>
  selectedCardIndex: number | null
  isLoading: boolean
  error: string | null
  setSelectedSpread: (spread: string) => void
  setQuestion: (question: string) => void
  startDraw: () => Promise<void>
  revealCard: (index: number) => void
  selectCard: (index: number) => void
  closeCardDetail: () => void
  reset: () => void
  allRevealed: boolean
}

export function useTarot(): UseTarotReturn {
  const [phase, setPhase] = useState<DivinationPhase>('input')
  const [selectedSpread, setSelectedSpread] = useState<string>('single')
  const [question, setQuestion] = useState('')
  const [drawResult, setDrawResult] = useState<DrawResult | null>(null)
  const [revealedCards, setRevealedCards] = useState<Set<number>>(new Set())
  const [selectedCardIndex, setSelectedCardIndex] = useState<number | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const startDraw = useCallback(async () => {
    setIsLoading(true)
    setError(null)
    setPhase('drawing')

    try {
      const result = await drawCards(selectedSpread, question || undefined)
      setDrawResult(result)
      setRevealedCards(new Set())
      setPhase('drawing')
    } catch (err) {
      setError('抽牌失败，请稍后重试')
      setPhase('input')
      console.error('Draw error:', err)
    } finally {
      setIsLoading(false)
    }
  }, [selectedSpread, question])

  const revealCard = useCallback((index: number) => {
    setRevealedCards((prev) => {
      const next = new Set(prev)
      next.add(index)
      return next
    })
  }, [])

  const selectCard = useCallback((index: number) => {
    setSelectedCardIndex(index)
  }, [])

  const closeCardDetail = useCallback(() => {
    setSelectedCardIndex(null)
  }, [])

  const reset = useCallback(() => {
    setPhase('input')
    setDrawResult(null)
    setRevealedCards(new Set())
    setSelectedCardIndex(null)
    setError(null)
    setQuestion('')
  }, [])

  const allRevealed = drawResult
    ? revealedCards.size === drawResult.cards.length
    : false

  return {
    phase,
    selectedSpread,
    question,
    drawResult,
    revealedCards,
    selectedCardIndex,
    isLoading,
    error,
    setSelectedSpread,
    setQuestion,
    startDraw,
    revealCard,
    selectCard,
    closeCardDetail,
    reset,
    allRevealed,
  }
}
