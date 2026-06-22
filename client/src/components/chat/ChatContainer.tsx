import { useState, useRef, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useChat } from '../../hooks/useChat'
import ChatMessage from './ChatMessage'
import ChatInput from './ChatInput'
import QuickQuestions from './QuickQuestions'
import { getQuickQuestions } from './quickQuestionData'

interface ChatContainerProps {
  /** The divination record ID */
  recordId?: number
  /** Initial divination type */
  divinationType?: string
  /** Initial question */
  question?: string
  /** Initial AI reading (if already exists) */
  initialReading?: string
  /** Callback when back button is clicked */
  onBack?: () => void
  /** Callback for history button */
  onHistory?: () => void
  /** Called when this chat has active content */
  onActiveChange?: (active: boolean) => void
}

/**
 * ChatContainer is the main chat interface component.
 * Handles the full chat experience including welcome, messages, and input.
 */
export default function ChatContainer({
  recordId,
  divinationType = 'tarot',
  question,
  onBack,
  onHistory,
  onActiveChange,
}: ChatContainerProps) {
  const selectedType = divinationType
  const {
    session,
    messages,
    isStreaming,
    isLoading,
    error,
    initSession,
    startModeSession,
    sendMessage,
    reset,
  } = useChat()

  const messagesEndRef = useRef<HTMLDivElement>(null)
  const chatAreaRef = useRef<HTMLDivElement>(null)
  const [showScrollButton, setShowScrollButton] = useState(false)

  // Initialize session on mount
  useEffect(() => {
    if (recordId) {
      initSession(recordId)
    }
  }, [recordId, initSession])

  useEffect(() => {
    onActiveChange?.(Boolean(session || messages.length > 0 || isLoading))
  }, [isLoading, messages.length, onActiveChange, session])

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' })
    }
  }, [messages])

  // Handle scroll to show/hide scroll button
  const handleScroll = () => {
    if (chatAreaRef.current) {
      const { scrollTop, scrollHeight, clientHeight } = chatAreaRef.current
      setShowScrollButton(scrollHeight - scrollTop - clientHeight > 100)
    }
  }

  // Scroll to bottom
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  const handleSend = (content: string) => {
    if (session) {
      sendMessage(content)
      return
    }
    startModeSession(selectedType, content)
  }

  const handleSuggestedQuestion = (content: string) => {
    if (session) {
      sendMessage(content)
      return
    }
    startModeSession(selectedType, content)
  }

  // Get welcome content based on divination type
  const getWelcomeContent = () => {
    const contents: Record<string, { icon: string; title: string; desc: string; hint: string }> = {
      tarot: {
        icon: '🎴',
        title: '塔罗牌占卜',
        desc: '通过塔罗牌的象征图案，洞察您内心的答案。\n支持单牌、三牌、凯尔特十字等多种牌阵。',
        hint: '💡 请输入您的问题，我将为您抽牌解读',
      },
      liuyao: {
        icon: '☯️',
        title: '六爻占卜',
        desc: '基于《高岛易断》的深度卦象解读系统。\n支持蓍草法和铜钱法起卦。',
        hint: '💡 请输入您想问的事情，我将为您起卦',
      },
      bazi: {
        icon: '📋',
        title: '八字排盘',
        desc: '根据出生时间排盘，分析命理五行。\n涵盖十神、用神忌神等深度分析。',
        hint: '💡 请告诉我您的出生年月日时',
      },
      horoscope: {
        icon: '⭐',
        title: '星座运势',
        desc: '每日/每周/每月星座运势解读。\n涵盖爱情、事业、财运、健康。',
        hint: '💡 请告诉我您的星座',
      },
    }
    return contents[selectedType] || contents.tarot
  }

  // Get placeholder based on state
  const getPlaceholder = () => {
    if (isStreaming) return 'AI正在回复中...'
    if (!session && selectedType === 'bazi') return '输入出生年月日时...'
    if (!session) return '输入您想问的问题...'
    return '继续提问...'
  }

  const welcome = getWelcomeContent()

  return (
    <div className="flex flex-col h-full bg-slate-900/50 relative">
      {/* Chat header */}
      <div className="px-4 sm:px-6 py-3 border-b border-slate-700/50 flex items-center justify-between">
        <div>
          <div className="text-sm font-semibold text-slate-200 flex items-center gap-2">
            🔮 AI 占卜师
          </div>
          <div className="text-xs text-slate-500">
            {session ? '基于您的占卜结果进行深度解读' : '选择占卜方式开始'}
          </div>
        </div>
        <div className="flex gap-2">
          {session && (
            <button
              onClick={() => {
                const firstUserMessage = messages.find((msg) => msg.role === 'user')?.content || question || '请重新解读'
                reset()
                startModeSession(selectedType, firstUserMessage)
              }}
              className="px-3 py-1.5 text-xs text-slate-400 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:bg-slate-700/50 transition-colors"
            >
              重新解读
            </button>
          )}
          {onHistory && (
            <button
              onClick={onHistory}
              className="px-3 py-1.5 text-xs text-slate-400 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:bg-slate-700/50 transition-colors"
            >
              历史记录
            </button>
          )}
          {onBack && (
            <button
              onClick={onBack}
              className="px-3 py-1.5 text-xs text-slate-400 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:bg-slate-700/50 transition-colors"
            >
              ← 返回
            </button>
          )}
        </div>
      </div>

      {/* Messages area */}
      <div
        ref={chatAreaRef}
        onScroll={handleScroll}
        className="flex-1 overflow-y-auto px-4 sm:px-6 py-4 space-y-4"
      >
        {/* Welcome card (show when no messages) */}
        {messages.length === 0 && !isLoading && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="flex flex-col items-center justify-center py-8"
          >
            <div className="bg-gradient-to-br from-violet-500/10 to-indigo-500/5 border border-violet-500/20 rounded-2xl p-6 sm:p-8 max-w-md text-center">
              <div className="text-5xl mb-4">{welcome.icon}</div>
              <h3 className="text-lg font-semibold text-slate-200 mb-2">{welcome.title}</h3>
              <p className="text-sm text-slate-400 leading-relaxed whitespace-pre-line mb-4">
                {welcome.desc}
              </p>
              <div className="bg-violet-500/10 rounded-xl px-4 py-3 text-xs text-indigo-300">
                {welcome.hint}
              </div>
            </div>
          </motion.div>
        )}

        {/* Loading state */}
        {isLoading && (
          <div className="flex justify-center py-8">
            <div className="w-8 h-8 border-2 border-violet-500/30 border-t-violet-500 rounded-full animate-spin" />
          </div>
        )}

        {/* Messages */}
        {messages.map((msg, index) => (
          <ChatMessage
            key={msg.id || index}
            message={msg}
            isStreaming={isStreaming && index === messages.length - 1 && msg.role === 'assistant'}
          />
        ))}

        {/* Error display */}
        {error && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="flex justify-center"
          >
            <div className="bg-red-500/10 border border-red-500/30 rounded-xl px-4 py-3 text-sm text-red-400">
              {error}
            </div>
          </motion.div>
        )}

        <div ref={messagesEndRef} />
      </div>

      {/* Scroll to bottom button */}
      <AnimatePresence>
        {showScrollButton && (
          <motion.button
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 10 }}
            onClick={scrollToBottom}
            className="absolute bottom-32 right-6 w-8 h-8 bg-slate-700 border border-slate-600 rounded-full flex items-center justify-center text-slate-300 shadow-lg hover:bg-slate-600 transition-colors"
          >
            ↓
          </motion.button>
        )}
      </AnimatePresence>

      <QuickQuestions
        questions={getQuickQuestions(selectedType)}
        onSelect={handleSuggestedQuestion}
        variant="horizontal"
        title={session ? '可继续追问' : '也可以从这些问题开始'}
      />

      {/* Input area */}
      <ChatInput
        placeholder={getPlaceholder()}
        disabled={isStreaming || isLoading}
        onSend={handleSend}
      />
    </div>
  )
}
