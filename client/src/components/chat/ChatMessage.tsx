import { motion } from 'framer-motion'
import type { ChatMessage as ChatMessageType } from '../../services/chat'
import { getDivinationPersona } from './divinationPersona'

interface ChatMessageProps {
  message: ChatMessageType
  isStreaming?: boolean
  /** Current divination type for persona display */
  divinationType?: string
}

/**
 * ChatMessage component displays a single chat message bubble.
 */
export default function ChatMessage({ message, isStreaming, divinationType }: ChatMessageProps) {
  const isUser = message.role === 'user'
  const isAssistant = message.role === 'assistant'
  const persona = divinationType ? getDivinationPersona(divinationType) : null

  // Format time
  const formatTime = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }

  // Simple markdown-like rendering
  const renderContent = (content: string) => {
    if (!content) return null

    const paragraphs = content.split('\n').filter(p => p.trim())

    return paragraphs.map((p, i) => {
      const parts = p.split(/(\*\*.*?\*\*)/g)

      return (
        <p key={i} className="mb-2 last:mb-0">
          {parts.map((part, index) => {
            if (part.startsWith('**') && part.endsWith('**') && part.length > 4) {
              return <strong key={index}>{part.slice(2, -2)}</strong>
            }
            return <span key={index}>{part}</span>
          })}
        </p>
      )
    })
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      className={`flex gap-2 sm:gap-3 ${isUser ? 'flex-row-reverse' : 'flex-row'}`}
    >
      {/* Avatar */}
      <div
        className={`w-8 h-8 sm:w-9 sm:h-9 rounded-lg flex items-center justify-center text-sm flex-shrink-0 ${
          isAssistant
            ? 'bg-gradient-to-br from-violet-500/30 to-indigo-500/30'
            : 'bg-gradient-to-br from-emerald-500/30 to-teal-500/30'
        }`}
      >
        {isAssistant ? (persona?.icon || '🔮') : '👤'}
      </div>

      {/* Message content */}
      <div className={`flex-1 min-w-0 ${isUser ? 'text-right' : ''}`}>
        {/* Name */}
        <div className={`text-xs text-slate-500 mb-1 px-1 ${isUser ? 'text-right' : ''}`}>
          {isAssistant ? (persona?.name || 'AI 占卜师') : '我'}
        </div>

        {/* Bubble */}
        <div
          className={`inline-block text-left px-4 py-3 rounded-2xl text-sm leading-relaxed max-w-full ${
            isAssistant
              ? 'bg-slate-800/80 border border-slate-700/50 rounded-tl-md'
              : 'bg-gradient-to-br from-violet-500/20 to-indigo-500/20 border border-violet-500/30 rounded-tr-md'
          }`}
        >
          {renderContent(message.content)}

          {/* Streaming cursor */}
          {isStreaming && isAssistant && message.content && (
            <motion.span
              animate={{ opacity: [0, 1, 0] }}
              transition={{ repeat: Infinity, duration: 1 }}
              className="inline-block w-0.5 h-4 bg-violet-400 ml-0.5 align-middle"
            />
          )}

          {/* Typing indicator */}
          {isStreaming && isAssistant && !message.content && (
            <div className="flex gap-1.5 py-1">
              <motion.div
                animate={{ y: [0, -6, 0] }}
                transition={{ repeat: Infinity, duration: 1.4, delay: 0 }}
                className="w-1.5 h-1.5 bg-indigo-400 rounded-full"
              />
              <motion.div
                animate={{ y: [0, -6, 0] }}
                transition={{ repeat: Infinity, duration: 1.4, delay: 0.2 }}
                className="w-1.5 h-1.5 bg-indigo-400 rounded-full"
              />
              <motion.div
                animate={{ y: [0, -6, 0] }}
                transition={{ repeat: Infinity, duration: 1.4, delay: 0.4 }}
                className="w-1.5 h-1.5 bg-indigo-400 rounded-full"
              />
            </div>
          )}
        </div>

        {/* Time */}
        <div className={`text-[10px] text-slate-600 mt-1 px-1 ${isUser ? 'text-right' : ''}`}>
          {formatTime(message.created_at)}
        </div>
      </div>
    </motion.div>
  )
}
