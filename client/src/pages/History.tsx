import { useState, useEffect, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useNavigate } from 'react-router-dom'
import { DayPicker, type DateRange } from 'react-day-picker'
import { format } from 'date-fns'
import { zhCN } from 'date-fns/locale/zh-CN'
import {
  fetchHistory,
  fetchHistoryDetail,
  deleteHistory,
  type DivinationRecord,
  type HistoryListResponse,
} from '../services/history'
import { listSessions, type ChatSession } from '../services/chat'
import 'react-day-picker/style.css'

const typeFilters = [
  { value: '', label: '全部' },
  { value: 'liuyao', label: '六爻' },
  { value: 'liuyao_v2', label: '高岛易断' },
  { value: 'bazi', label: '八字' },
  { value: 'tarot', label: '塔罗牌' },
  { value: 'horoscope', label: '星座' },
]

function parseResult(result: string): unknown {
  try {
    return JSON.parse(result)
  } catch {
    return result
  }
}

function valueAsRecord(value: unknown): Record<string, unknown> | null {
  return value && typeof value === 'object' && !Array.isArray(value)
    ? value as Record<string, unknown>
    : null
}

function getNestedString(value: unknown, path: string[]): string {
  let current: unknown = value
  for (const key of path) {
    const record = valueAsRecord(current)
    if (!record) return ''
    current = record[key]
  }
  return typeof current === 'string' ? current : ''
}

function getResultTitle(record: DivinationRecord): string {
  const result = parseResult(record.result)
  const benGua = getNestedString(result, ['ben_gua', 'name']) || getNestedString(result, ['ben_gua', 'name_short'])
  const bianGua = getNestedString(result, ['bian_gua', 'name']) || getNestedString(result, ['bian_gua', 'name_short'])
  if (benGua && bianGua) return `${benGua} → ${bianGua}`
  if (benGua) return benGua

  const zodiac = getNestedString(result, ['zodiac_cn'])
  const period = getNestedString(result, ['period'])
  if (zodiac) return period ? `${zodiac} · ${period}` : zodiac

  const birth = getNestedString(result, ['birth', 'solar'])
  if (birth) return birth

  const cards = valueAsRecord(result)?.cards
  if (Array.isArray(cards)) {
    const names = cards
      .map((item) => getNestedString(item, ['card', 'name']))
      .filter(Boolean)
      .slice(0, 4)
    if (names.length > 0) return names.join(' · ')
  }

  return record.type
}

function getResultSummary(record: DivinationRecord): string {
  const result = parseResult(record.result)
  const summary = getNestedString(result, ['summary'])
  if (summary) return summary

  const judgment = getNestedString(result, ['ben_gua', 'judgment', 'text']) || getNestedString(result, ['ben_gua', 'judgment'])
  if (judgment) return judgment

  if (typeof result === 'string') return result
  return record.result
}

function useIsMobile(breakpoint = 640) {
  const getInitialValue = () =>
    typeof window !== 'undefined' && window.matchMedia(`(max-width: ${breakpoint - 1}px)`).matches
  const [isMobile, setIsMobile] = useState(getInitialValue)

  useEffect(() => {
    const mq = window.matchMedia(`(max-width: ${breakpoint - 1}px)`)
    const handler = (e: MediaQueryListEvent) => setIsMobile(e.matches)
    mq.addEventListener('change', handler)
    return () => mq.removeEventListener('change', handler)
  }, [breakpoint])

  return isMobile
}

export default function History() {
  const navigate = useNavigate()
  const isMobile = useIsMobile()
  const [data, setData] = useState<HistoryListResponse | null>(null)
  const [activeType, setActiveType] = useState('')
  const [page, setPage] = useState(1)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [dateRange, setDateRange] = useState<DateRange | undefined>(undefined)
  const [showCalendar, setShowCalendar] = useState(false)
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set())
  const [activeDetailId, setActiveDetailId] = useState<number | null>(null)
  const [detail, setDetail] = useState<DivinationRecord | null>(null)
  const [isDetailLoading, setIsDetailLoading] = useState(false)
  const [chatSessions, setChatSessions] = useState<Map<number, ChatSession>>(new Map())

  const loadData = useCallback(async (type: string, p: number, range?: DateRange) => {
    setIsLoading(true)
    setError(null)
    try {
      const startDate = range?.from ? format(range.from, 'yyyy-MM-dd') : undefined
      const endDate = range?.to ? format(range.to, 'yyyy-MM-dd') : undefined
      const result = await fetchHistory(type || undefined, p, 10, startDate, endDate)
      setData(result)
      setSelectedIds(new Set())
      setActiveDetailId(null)
      setDetail(null)

      // Load chat sessions to check which records have sessions
      try {
        const sessionsResult = await listSessions(1, 100)
        const sessionMap = new Map<number, ChatSession>()
        sessionsResult.sessions.forEach(s => sessionMap.set(s.record_id, s))
        setChatSessions(sessionMap)
      } catch {
        // Ignore chat session load errors
      }
    } catch (err) {
      setError('加载历史记录失败')
      console.error('History load error:', err)
    } finally {
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      loadData(activeType, page, dateRange)
    }, 0)

    return () => {
      window.clearTimeout(timeoutId)
    }
  }, [activeType, page, dateRange, loadData])

  const handleFilterChange = (type: string) => {
    setActiveType(type)
    setPage(1)
  }

  const handleDateRangeChange = (range: DateRange | undefined) => {
    setDateRange(range)
    setPage(1)
  }

  const handleClearDates = () => {
    setDateRange(undefined)
    setPage(1)
  }

  const handleSelectAll = () => {
    if (!data) return
    if (selectedIds.size === data.items.length) {
      setSelectedIds(new Set())
    } else {
      setSelectedIds(new Set(data.items.map(item => item.id)))
    }
  }

  const handleToggleSelect = (id: number) => {
    setSelectedIds(prev => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        next.add(id)
      }
      return next
    })
  }

  const handleOpenDetail = async (id: number) => {
    if (activeDetailId === id && detail) {
      setActiveDetailId(null)
      setDetail(null)
      return
    }

    setActiveDetailId(id)
    setIsDetailLoading(true)
    setError(null)
    try {
      const result = await fetchHistoryDetail(id)
      setDetail(result)
    } catch (err) {
      setError('加载详情失败')
      console.error('History detail error:', err)
      setActiveDetailId(null)
      setDetail(null)
    } finally {
      setIsDetailLoading(false)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这条记录吗？')) return

    try {
      await deleteHistory(id)
      if (activeDetailId === id) {
        setActiveDetailId(null)
        setDetail(null)
      }
      loadData(activeType, page, dateRange)
    } catch (err) {
      setError('删除失败')
      console.error('Delete error:', err)
    }
  }

  const typeIcons: Record<string, string> = {
    tarot: '🔮',
    liuyao: '☯',
    liuyao_v2: '☯',
    bazi: '📜',
    horoscope: '✨',
  }

  return (
    <div className="flex-1 px-4 py-12 max-w-4xl mx-auto w-full">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center mb-10"
      >
        <h1 className="text-2xl sm:text-4xl font-bold mb-3">
          <span className="bg-gradient-to-r from-amber-400 via-orange-400 to-red-400 bg-clip-text text-transparent">
            历史记录
          </span>
        </h1>
        <p className="text-slate-400">查看你的占卜历史</p>
      </motion.div>

      {/* Filters */}
      <div className="flex gap-2 mb-4 overflow-x-auto scrollbar-hide touch-pan-x sm:flex-wrap sm:justify-center">
        {typeFilters.map((filter) => (
          <button
            key={filter.value}
            onClick={() => handleFilterChange(filter.value)}
            className={`flex-shrink-0 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              activeType === filter.value
                ? 'bg-amber-500/20 text-amber-400 border border-amber-500/30'
                : 'bg-slate-800/50 text-slate-400 border border-slate-700 hover:border-slate-600'
            }`}
          >
            {filter.label}
          </button>
        ))}
      </div>

      {/* Date Range Filter */}
      <div className="flex flex-col items-center gap-3 mb-8">
        <div className="flex items-center gap-3">
          <button
            onClick={() => setShowCalendar(!showCalendar)}
            className="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium bg-slate-800/50 text-slate-400 border border-slate-700 hover:border-slate-600 transition-all"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            {dateRange?.from
              ? `${format(dateRange.from, 'MM/dd')}${dateRange.to ? ` - ${format(dateRange.to, 'MM/dd')}` : ''}`
              : '按日期筛选'
            }
          </button>
          {dateRange?.from && (
            <button
              onClick={handleClearDates}
              className="px-3 py-2 rounded-lg text-sm text-slate-500 hover:text-red-400 transition-colors"
            >
              清除日期
            </button>
          )}
        </div>

        {showCalendar && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="bg-slate-800/80 border border-slate-700 rounded-xl p-4 backdrop-blur-sm"
          >
            <DayPicker
              mode="range"
              selected={dateRange}
              onSelect={handleDateRangeChange}
              locale={zhCN}
              numberOfMonths={isMobile ? 1 : 2}
              className="text-slate-300"
              classNames={{
                months: 'flex gap-4',
                month: 'flex flex-col gap-3',
                caption_label: 'text-slate-300 text-sm font-medium',
                nav: 'flex items-center gap-1',
                button_previous: 'p-1 rounded hover:bg-slate-700 text-slate-400 transition-colors',
                button_next: 'p-1 rounded hover:bg-slate-700 text-slate-400 transition-colors',
                weekday: 'text-slate-500 text-xs font-medium w-9 py-1',
                day: 'w-9 h-9 flex items-center justify-center rounded-lg text-sm transition-all hover:bg-slate-700/50',
                selected: 'bg-amber-500/20 text-amber-400',
                range_start: 'bg-amber-500/20 text-amber-400 rounded-r-none',
                range_end: 'bg-amber-500/20 text-amber-400 rounded-l-none',
                range_middle: 'bg-amber-500/10 text-amber-300 rounded-none',
                today: 'text-amber-500 font-bold',
                disabled: 'text-slate-600 opacity-50',
                outside: 'text-slate-600',
              }}
            />
          </motion.div>
        )}
      </div>

      {/* Select All Button */}
      {data && data.items.length > 0 && (
        <div className="flex justify-end mb-4">
          <button
            onClick={handleSelectAll}
            className="flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs font-medium bg-slate-800/50 text-slate-400 border border-slate-700 hover:border-slate-600 transition-all"
          >
            <span className={`w-4 h-4 rounded border flex items-center justify-center transition-all ${
              selectedIds.size === data.items.length
                ? 'bg-amber-500/20 border-amber-500/30'
                : 'border-slate-600'
            }`}>
              {selectedIds.size === data.items.length && (
                <svg className="w-3 h-3 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                </svg>
              )}
            </span>
            {selectedIds.size === data.items.length ? '取消全选' : '全选当前页'}
            {selectedIds.size > 0 && selectedIds.size < data.items.length && (
              <span className="text-amber-400">({selectedIds.size})</span>
            )}
          </button>
        </div>
      )}

      {/* Error display */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="mb-6 p-4 bg-red-500/10 border border-red-500/30 rounded-xl text-red-400 text-center"
          >
            {error}
          </motion.div>
        )}
      </AnimatePresence>

      {/* Loading */}
      {isLoading && (
        <div className="text-center py-20">
          <svg className="animate-spin h-8 w-8 mx-auto text-amber-400" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
        </div>
      )}

      {/* History list */}
      {!isLoading && data && data.items.length === 0 && (
        <div className="text-center py-20">
          <div className="text-5xl mb-4">📋</div>
          <p className="text-slate-400">暂无历史记录</p>
          <p className="text-slate-500 text-sm mt-2">完成占卜后，记录将自动保存</p>
        </div>
      )}

      {!isLoading && data && data.items.length > 0 && (
        <div className="space-y-4">
          {data.items.map((item, index) => (
            <motion.div
              key={item.id}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.05 }}
              role="button"
              tabIndex={0}
              onClick={() => handleOpenDetail(item.id)}
              onKeyDown={(event) => {
                if (event.key === 'Enter' || event.key === ' ') {
                  event.preventDefault()
                  handleOpenDetail(item.id)
                }
              }}
              className={`bg-slate-800/50 border rounded-xl p-5 transition-all group cursor-pointer ${
                activeDetailId === item.id
                  ? 'border-amber-500/40 bg-slate-800/70'
                  : 'border-slate-700/50 hover:border-slate-600/50'
              }`}
            >
              <div className="flex items-start justify-between">
                <div className="flex items-start gap-3 flex-1 min-w-0">
                  <button
                    onClick={(event) => {
                      event.stopPropagation()
                      handleToggleSelect(item.id)
                    }}
                    className={`w-5 h-5 rounded border flex-shrink-0 flex items-center justify-center transition-all mt-0.5 ${
                      selectedIds.has(item.id)
                        ? 'bg-amber-500/20 border-amber-500/30'
                        : 'border-slate-600 hover:border-slate-500'
                    }`}
                  >
                    {selectedIds.has(item.id) && (
                      <svg className="w-3 h-3 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                      </svg>
                    )}
                  </button>
                  <span className="text-2xl flex-shrink-0">{typeIcons[item.type] || '🔮'}</span>
                  <div className="min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-white font-medium">{item.type_cn}</span>
                      <span className="text-slate-500 text-xs">{item.created_at}</span>
                    </div>
                    {item.question && (
                      <p className="text-slate-400 text-sm truncate">{item.question}</p>
                    )}
                    <p className="text-slate-500 text-xs mt-1 line-clamp-2">{item.summary}</p>
                  </div>
                </div>
                <button
                  onClick={(event) => {
                    event.stopPropagation()
                    handleDelete(item.id)
                  }}
                  className="text-slate-600 hover:text-red-400 transition-colors sm:opacity-0 sm:group-hover:opacity-100 ml-4 flex-shrink-0"
                  title="删除"
                >
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>

              <AnimatePresence>
                {activeDetailId === item.id && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: 'auto' }}
                    exit={{ opacity: 0, height: 0 }}
                    className="overflow-hidden"
                    onClick={(event) => event.stopPropagation()}
                  >
                    <div className="mt-5 border-t border-slate-700/60 pt-5">
                      {isDetailLoading && (
                        <p className="text-sm text-slate-500">详情加载中...</p>
                      )}

                      {!isDetailLoading && detail && detail.id === item.id && (
                        <div className="space-y-5">
                          <div>
                            <p className="text-xs font-medium text-amber-400 mb-1">记录详情</p>
                            <h2 className="text-lg font-semibold text-white">{getResultTitle(detail)}</h2>
                            {detail.question && (
                              <p className="mt-2 text-sm text-slate-400">{detail.question}</p>
                            )}
                          </div>

                          <div className="rounded-lg border border-slate-700/60 bg-slate-900/40 p-4">
                            <p className="mb-2 text-xs font-medium text-slate-500">占卜结果</p>
                            <p className="whitespace-pre-wrap text-sm leading-6 text-slate-300 line-clamp-6">
                              {getResultSummary(detail)}
                            </p>
                          </div>

                          <div className="rounded-lg border border-violet-500/20 bg-violet-500/5 p-4">
                            <p className="mb-2 text-xs font-medium text-violet-300">AI解读</p>
                            {detail.ai_reading ? (
                              <p className="whitespace-pre-wrap text-sm leading-6 text-slate-300">
                                {detail.ai_reading}
                              </p>
                            ) : (
                              <p className="text-sm text-slate-500">暂无 AI 解读</p>
                            )}
                          </div>

                          {/* Chat button */}
                          <div className="flex gap-2">
                            <button
                              onClick={(e) => {
                                e.stopPropagation()
                                navigate(`/chat/${item.id}`)
                              }}
                              className="flex-1 px-4 py-2.5 bg-gradient-to-r from-violet-600 to-indigo-600 rounded-lg text-white text-sm font-medium hover:from-violet-700 hover:to-indigo-700 transition-all flex items-center justify-center gap-2 shadow-lg shadow-violet-500/20"
                            >
                              💬 {chatSessions.has(item.id) ? '继续对话' : '进入对话'}
                            </button>
                          </div>
                        </div>
                      )}
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </motion.div>
          ))}
        </div>
      )}

      {/* Pagination */}
      {data && data.total > data.page_size && (
        <div className="flex justify-center items-center gap-4 mt-8">
          <button
            onClick={() => setPage(p => Math.max(1, p - 1))}
            disabled={page <= 1}
            className="px-4 py-2 bg-slate-800 border border-slate-700 rounded-lg text-slate-300 disabled:opacity-30 disabled:cursor-not-allowed hover:border-slate-600 transition-all"
          >
            上一页
          </button>
          <span className="text-slate-400 text-sm">
            第 {page} 页 / 共 {Math.ceil(data.total / data.page_size)} 页
          </span>
          <button
            onClick={() => setPage(p => p + 1)}
            disabled={page >= Math.ceil(data.total / data.page_size)}
            className="px-4 py-2 bg-slate-800 border border-slate-700 rounded-lg text-slate-300 disabled:opacity-30 disabled:cursor-not-allowed hover:border-slate-600 transition-all"
          >
            下一页
          </button>
        </div>
      )}
    </div>
  )
}
