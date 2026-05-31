import { useState, useCallback } from 'react'
import { fetchHoroscope } from '../services/horoscope'
import type { HoroscopeResult } from '../types/horoscope'

type Period = 'daily' | 'weekly' | 'monthly'

interface UseHoroscopeReturn {
  selectedZodiac: string | null
  period: Period
  result: HoroscopeResult | null
  isLoading: boolean
  error: string | null
  selectZodiac: (zodiac: string) => void
  setPeriod: (period: Period) => void
  refresh: () => void
}

export function useHoroscope(): UseHoroscopeReturn {
  const [selectedZodiac, setSelectedZodiac] = useState<string | null>(null)
  const [period, setPeriod] = useState<Period>('daily')
  const [result, setResult] = useState<HoroscopeResult | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const doFetch = useCallback(async (zodiac: string, p: Period) => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await fetchHoroscope(zodiac, p)
      setResult(data)
    } catch (err) {
      setError('获取运势失败，请稍后重试')
      console.error('Horoscope fetch error:', err)
    } finally {
      setIsLoading(false)
    }
  }, [])

  const selectZodiac = useCallback((zodiac: string) => {
    setSelectedZodiac(zodiac)
    doFetch(zodiac, period)
  }, [period, doFetch])

  const setPeriodAndFetch = useCallback((p: Period) => {
    setPeriod(p)
    if (selectedZodiac) {
      doFetch(selectedZodiac, p)
    }
  }, [selectedZodiac, doFetch])

  const refresh = useCallback(() => {
    if (selectedZodiac) {
      doFetch(selectedZodiac, period)
    }
  }, [selectedZodiac, period, doFetch])

  return {
    selectedZodiac,
    period,
    result,
    isLoading,
    error,
    selectZodiac,
    setPeriod: setPeriodAndFetch,
    refresh,
  }
}
