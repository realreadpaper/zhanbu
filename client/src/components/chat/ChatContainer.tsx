import { useState, useRef, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useChat } from '../../hooks/useChat'
import ChatMessage from './ChatMessage'
import ChatInput from './ChatInput'
import QuickQuestions from './QuickQuestions'
import { getQuickQuestions } from './quickQuestionData'
import DivinationRitual from './DivinationRitual'
import DivinationResultCard from './DivinationResultCard'
import { getDivinationPersona } from './divinationPersona'

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
    record,
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
  const persona = getDivinationPersona(selectedType)
  const isInitialRitual = isLoading && !session && messages.length === 0

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

  // Get placeholder based on state
  const getPlaceholder = () => {
    if (isStreaming) return `${persona.name}正在解读...`
    if (!session && selectedType === 'bazi') return '输入出生年月日时，例如：1990-05-12 08:30 女...'
    if (!session && selectedType === 'horoscope') return '输入星座和问题，例如：天蝎座今天事业运如何...'
    if (!session) return '输入您想问的问题...'
    return '继续提问...'
  }

  return (
    <div className="flex flex-1 min-h-0 flex-col bg-slate-900/50 relative">
      {/* Chat header — 手机端隐藏,由 ChatPage mobile header 替代 */}
      <div className="hidden sm:flex px-4 sm:px-6 py-3 border-b border-slate-700/50 items-center justify-between">
        <div>
          <div className="text-sm font-semibold text-slate-200 flex items-center gap-2">
            <span>{persona.icon}</span>
            <span>{persona.name}</span>
          </div>
          <div className="text-xs text-slate-500">
            {session ? persona.subtitle : '选择占卜方式开始'}
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
        className="flex-1 overflow-y-auto overscroll-contain px-4 sm:px-6 py-4 space-y-4"
      >
        {/* Welcome card (show when no messages) */}
        {messages.length === 0 && !isLoading && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="flex flex-col items-center justify-center py-8"
          >
            <div className="bg-gradient-to-br from-violet-500/10 to-indigo-500/5 border border-violet-500/20 rounded-2xl p-6 sm:p-8 max-w-md text-center">
              <div className="text-5xl mb-4">{persona.icon}</div>
              <h3 className="text-lg font-semibold text-slate-200 mb-1">{persona.welcomeTitle}</h3>
              <div className="text-sm text-violet-200 mb-3">{persona.title}</div>
              <p className="text-sm text-slate-400 leading-relaxed whitespace-pre-line mb-4">
                {persona.welcomeDescription}
              </p>
              <div className="bg-violet-500/10 rounded-xl px-4 py-3 text-xs text-indigo-300">
                {persona.welcomeHint}
              </div>
            </div>
          </motion.div>
        )}

        <AnimatePresence>
          {isInitialRitual && (
            <div className="py-8">
              <DivinationRitual persona={persona} />
            </div>
          )}
        </AnimatePresence>

        {isLoading && !isInitialRitual && (
          <div className="flex justify-center py-8">
            <div className={`h-8 w-8 rounded-full border-2 border-slate-700 border-t-violet-400 animate-spin`} />
          </div>
        )}

        {/* Messages */}
        {messages.map((msg, index) => {
          const showResultAfterMessage = record && index === 0 && msg.role === 'user'

          return (
            <div key={msg.id || index} className="space-y-4">
              <ChatMessage
                message={msg}
                isStreaming={isStreaming && index === messages.length - 1 && msg.role === 'assistant'}
              />
              {showResultAfterMessage && <DivinationResultCard record={record} />}
            </div>
          )
        })}

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
