import { useState, useEffect } from 'react'
import { useParams, useNavigate, useSearchParams } from 'react-router-dom'
import ChatContainer from '../components/chat/ChatContainer'
import DivinationSelector from '../components/chat/DivinationSelector'
import { getDivinationPersona } from '../components/chat/divinationPersona'
import { fetchHistoryDetail } from '../services/history'

interface DivinationRecord {
  id: number
  type: string
  question: string
  result: string
  ai_reading: string
  created_at: string
}

interface TarotResultCard {
  name?: string
  position?: number
  is_reversed?: boolean
}

interface TarotResult {
  cards?: TarotResultCard[]
}

const CHAT_TYPES = ['tarot', 'liuyao', 'liuyao_v2', 'bazi', 'horoscope'] as const

const getInitialType = (type: string | null) => {
  return type && CHAT_TYPES.includes(type as (typeof CHAT_TYPES)[number]) ? type : 'liuyao'
}

/**
 * ChatPage is the main chat page with responsive layout.
 * Desktop: left-right split (info + chat)
 * Mobile: single column chat
 */
export default function ChatPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const queryType = getInitialType(searchParams.get('type'))
  const [record, setRecord] = useState<DivinationRecord | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedType, setSelectedType] = useState(queryType)
  const [chatActive, setChatActive] = useState(false)

  // Load divination record
  useEffect(() => {
    const loadRecord = async () => {
      if (!id) {
        setRecord(null)
        setChatActive(false)
        setLoading(false)
        return
      }

      try {
        setLoading(true)
        const data = await fetchHistoryDetail(parseInt(id))
        setRecord(data)
        setSelectedType(data.type)
        setChatActive(true)
      } catch (err) {
        setError((err as Error).message || '加载失败')
      } finally {
        setLoading(false)
      }
    }

    loadRecord()
  }, [id])

  // Get divination type icon
  const getTypeIcon = (type: string) => {
    const icons: Record<string, string> = {
      tarot: '🎴',
      liuyao: '☯️',
      liuyao_v2: '☯️',
      bazi: '📋',
      horoscope: '⭐',
    }
    return icons[type] || '🔮'
  }

  // Format date
  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleString('zh-CN', {
      month: 'numeric',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  // Parse result to get card info (for tarot)
  const getCardInfo = () => {
    if (!record || record.type !== 'tarot') return null

    try {
      const result = JSON.parse(record.result) as TarotResult
      if (result.cards && Array.isArray(result.cards)) {
        return result.cards.map((card) => ({
          name: card.name,
          position: card.position,
          isReversed: card.is_reversed,
        }))
      }
    } catch {
      // Ignore parse errors
    }
    return null
  }

  // Loading state
  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen bg-slate-900">
        <div className="text-center">
          <div className="w-12 h-12 border-2 border-violet-500/30 border-t-violet-500 rounded-full animate-spin mx-auto mb-4" />
          <p className="text-slate-400">加载中...</p>
        </div>
      </div>
    )
  }

  // Error state
  if (error || (id && !record)) {
    return (
      <div className="flex items-center justify-center h-screen bg-slate-900">
        <div className="text-center">
          <div className="text-4xl mb-4">😵</div>
          <p className="text-red-400 mb-4">{error || '记录不存在'}</p>
          <button
            onClick={() => navigate('/history')}
            className="px-4 py-2 bg-violet-600 text-white rounded-lg hover:bg-violet-700 transition-colors"
          >
            返回历史记录
          </button>
        </div>
      </div>
    )
  }

  const cardInfo = getCardInfo()
  const activeType = record?.type || selectedType || queryType
  const persona = getDivinationPersona(activeType)

  return (
    <div className="h-screen flex flex-col lg:flex-row bg-slate-900">
      {/* Left panel - Desktop only */}
      <div className="hidden lg:flex lg:w-80 xl:w-96 flex-col border-r border-slate-700/50 bg-slate-900/80 overflow-y-auto">
        {/* Divination selector */}
        <div className="p-4 border-b border-slate-700/30">
          <div className="text-xs text-slate-500 font-semibold uppercase tracking-wider mb-3">
            选择占卜方式
          </div>
          <DivinationSelector
            selected={selectedType}
            onSelect={(type) => {
              if (record && !window.confirm('切换占卜方式会开始新的聊天，当前聊天可在历史记录中继续查看。是否继续？')) {
                return
              }
              setRecord(null)
              setChatActive(false)
              setSelectedType(type)
              navigate('/chat')
            }}
            variant="grid"
          />
        </div>

        {/* Current divination info */}
        <div className="p-4 border-b border-slate-700/30">
          <div className="text-xs text-slate-500 font-semibold uppercase tracking-wider mb-3">
            {record ? '当前占卜' : '使用说明'}
          </div>
          <div className="bg-gradient-to-br from-violet-500/10 to-indigo-500/5 border border-violet-500/20 rounded-xl p-4">
            <div className="flex items-center gap-3 mb-3">
              <div className="w-10 h-10 bg-gradient-to-br from-violet-600 to-indigo-600 rounded-xl flex items-center justify-center text-lg">
                {getTypeIcon(activeType)}
              </div>
              <div>
                <div className="text-sm font-semibold text-slate-200">
                  {persona.name}
                </div>
                <div className="text-xs text-slate-500">
                  {record ? formatDate(record.created_at) : persona.subtitle}
                </div>
              </div>
            </div>
            <div className="text-xs text-slate-400 mb-3">
              <strong className="text-violet-300">问题：</strong>
              {record?.question || '输入您想问的问题，AI 将在聊天中完成解读'}
            </div>
            {/* Card display for tarot */}
            {cardInfo && (
              <div className="flex gap-2 justify-center mt-3">
                {cardInfo.slice(0, 3).map((card, index) => (
                  <div
                    key={index}
                    className="w-14 h-20 bg-gradient-to-br from-violet-500/15 to-indigo-500/10 border border-violet-500/25 rounded-lg flex flex-col items-center justify-center text-[10px] text-violet-300 text-center p-1"
                  >
                    <div className="text-xl mb-0.5">
                      {index === 0 ? '👸' : index === 1 ? '🏎️' : '☀️'}
                    </div>
                    <div className="leading-tight">
                      {card.name}
                      <br />
                      <span className="text-slate-500">
                        {card.isReversed ? '逆位' : '正位'}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Back button */}
        <div className="p-4 mt-auto border-t border-slate-700/30">
          <button
            onClick={() => navigate('/history')}
            className="w-full px-4 py-2.5 text-sm text-slate-400 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:bg-slate-700/50 transition-colors flex items-center justify-center gap-2"
          >
            ← 返回历史记录
          </button>
        </div>
      </div>

      {/* Right panel - Chat area */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Mobile header */}
        <div className="lg:hidden flex items-center justify-between px-4 py-3 border-b border-slate-700/50 bg-slate-900/80">
          <button
            onClick={() => navigate('/history')}
            className="text-indigo-400 text-lg"
          >
            ‹
          </button>
          <div className="text-center flex-1">
            <div className="text-sm font-semibold text-slate-200">
              {persona.icon} {persona.name}
            </div>
            <div className="text-xs text-slate-500">
              {record?.question ? record.question.slice(0, 20) + '...' : persona.welcomeTitle}
            </div>
          </div>
          <div className="w-8" />
        </div>

        {/* Mobile divination selector */}
        <div className="lg:hidden">
          <DivinationSelector
            selected={selectedType}
            onSelect={(type) => {
              if (record && !window.confirm('切换占卜方式会开始新的聊天，当前聊天可在历史记录中继续查看。是否继续？')) {
                return
              }
              setRecord(null)
              setChatActive(false)
              setSelectedType(type)
              navigate('/chat')
            }}
            variant="scroll"
            compact={chatActive}
          />
        </div>

        {/* Chat container */}
        <ChatContainer
          recordId={record?.id}
          divinationType={activeType}
          question={record?.question}
          initialReading={record?.ai_reading}
          onActiveChange={setChatActive}
        />
      </div>
    </div>
  )
}
