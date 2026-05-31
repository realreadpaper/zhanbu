import { motion, AnimatePresence } from 'framer-motion'
import { useAI } from '../../hooks/useAI'

export interface AIReadingProps {
  /** The divination type (tarot, horoscope, liuyao, bazi) */
  type: string
  /** The result ID to interpret (for saved records) */
  resultId?: number
  /** The result JSON directly (for non-saved results) */
  result?: string
  /** Optional question to pass to the AI */
  question?: string
  /** Whether to show the button */
  show?: boolean
}

/**
 * AIReading component displays an AI interpretation button and streaming result.
 */
export default function AIReading({ type, resultId, result, question, show = true }: AIReadingProps) {
  const { text, isStreaming, isDone, error, start, reset } = useAI({
    type,
    resultId,
    result,
    question,
  })

  if (!show) return null

  const hasEmptyDone = isDone && !text.trim()

  return (
    <div className="mt-8">
      {/* AI Reading Button */}
      {!isStreaming && !isDone && !text && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-center"
        >
          <motion.button
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            onClick={start}
            disabled={isStreaming || (!resultId && !result)}
            className="px-8 py-4 bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-700 hover:to-indigo-700 text-white rounded-xl font-bold text-lg shadow-lg shadow-violet-500/20 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            <span className="flex items-center gap-3">
              <svg
                className="w-6 h-6"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M12 2L2 7l10 5 10-5-10-5z" />
                <path d="M2 17l10 5 10-5" />
                <path d="M2 12l10 5 10-5" />
              </svg>
              ✨ AI深度解读
            </span>
          </motion.button>
          <p className="text-slate-500 text-sm mt-3">
            点击获取AI大模型的个性化解读
          </p>
        </motion.div>
      )}

      {/* Error Display */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="p-4 bg-red-500/10 border border-red-500/30 rounded-xl text-red-400 text-center"
          >
            <p>{error}</p>
            <button
              onClick={start}
              className="mt-3 px-4 py-2 bg-red-500/20 hover:bg-red-500/30 text-red-300 rounded-lg text-sm transition-colors"
            >
              重试
            </button>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Empty result display */}
      <AnimatePresence>
        {hasEmptyDone && !error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="p-4 bg-amber-500/10 border border-amber-500/30 rounded-xl text-amber-300 text-center"
          >
            <p>AI未返回解读内容，请检查提示词或稍后重试</p>
            <button
              onClick={start}
              className="mt-3 px-4 py-2 bg-amber-500/20 hover:bg-amber-500/30 text-amber-200 rounded-lg text-sm transition-colors"
            >
              重试
            </button>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Streaming Content */}
      <AnimatePresence>
        {(isStreaming || text || (isDone && text.trim())) && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className="relative"
          >
            {/* Header */}
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-violet-400 font-semibold flex items-center gap-2">
                <svg
                  className="w-5 h-5"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M12 2L2 7l10 5 10-5-10-5z" />
                  <path d="M2 17l10 5 10-5" />
                  <path d="M2 12l10 5 10-5" />
                </svg>
                AI深度解读
              </h3>

              {isDone && (
                <button
                  onClick={reset}
                  className="text-slate-400 hover:text-white text-sm transition-colors"
                >
                  重新解读
                </button>
              )}
            </div>

            {/* Content */}
            <div className="p-6 bg-gradient-to-br from-violet-500/5 to-indigo-500/5 border border-violet-500/20 rounded-xl">
              <div className="text-slate-300 leading-relaxed whitespace-pre-line">
                {text}
                {isStreaming && (
                  <motion.span
                    animate={{ opacity: [0, 1, 0] }}
                    transition={{ repeat: Infinity, duration: 1 }}
                    className="inline-block w-2 h-5 bg-violet-400 ml-0.5 align-middle"
                  />
                )}
              </div>

              {/* Streaming indicator */}
              {isStreaming && (
                <div className="mt-4 flex items-center gap-2 text-violet-400 text-sm">
                  <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24">
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                      fill="none"
                    />
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                    />
                  </svg>
                  <span>AI正在解读中...</span>
                </div>
              )}

              {/* Done indicator */}
              {isDone && (
                <motion.div
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  className="mt-4 flex items-center gap-2 text-emerald-400 text-sm"
                >
                  <svg
                    className="w-4 h-4"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <polyline points="20 6 9 17 4 12" />
                  </svg>
                  <span>解读完成</span>
                </motion.div>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}
