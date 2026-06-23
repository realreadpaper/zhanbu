import { useState, useRef } from 'react'
import type { KeyboardEvent } from 'react'
import { motion } from 'framer-motion'

interface ChatInputProps {
  /** Placeholder text */
  placeholder?: string
  /** Whether the input is disabled (AI is streaming) */
  disabled?: boolean
  /** Callback when message is sent */
  onSend: (content: string) => void
}

/**
 * ChatInput component with auto-resize textarea and send button.
 */
export default function ChatInput({ placeholder = '输入您的问题...', disabled, onSend }: ChatInputProps) {
  const [content, setContent] = useState('')
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  const handleSend = () => {
    const trimmed = content.trim()
    if (!trimmed || disabled) return

    onSend(trimmed)
    setContent('')

    // Reset textarea height
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto'
    }
  }

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    // Enter to send, Shift+Enter for new line
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const handleInput = () => {
    // Auto-resize textarea
    const textarea = textareaRef.current
    if (textarea) {
      textarea.style.height = 'auto'
      textarea.style.height = Math.min(textarea.scrollHeight, 120) + 'px'
    }
  }

  return (
    <div className="p-3 pb-[max(0.75rem,env(safe-area-inset-bottom))] sm:p-4 border-t border-slate-700/50 bg-slate-900/50">
      <div className="flex gap-2 sm:gap-3 items-end">
        <textarea
          ref={textareaRef}
          value={content}
          onChange={(e) => setContent(e.target.value)}
          onKeyDown={handleKeyDown}
          onInput={handleInput}
          placeholder={placeholder}
          disabled={disabled}
          rows={1}
          className="flex-1 bg-slate-800/80 border border-slate-600/50 rounded-xl px-4 py-2.5 text-base text-slate-200 placeholder-slate-500 outline-none focus:border-violet-500/50 transition-colors resize-none min-h-[40px] max-h-[120px] disabled:opacity-50 sm:text-sm"
        />
        <motion.button
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
          onClick={handleSend}
          disabled={disabled || !content.trim()}
          className="w-10 h-10 bg-gradient-to-br from-violet-600 to-indigo-600 rounded-xl flex items-center justify-center text-white text-lg shadow-lg shadow-violet-500/20 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex-shrink-0"
        >
          ↑
        </motion.button>
      </div>
      <div className="text-[11px] text-slate-600 mt-2 text-center">
        按 Enter 发送，Shift + Enter 换行
      </div>
    </div>
  )
}
