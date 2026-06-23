import { motion } from 'framer-motion'
import type { QuickQuestion } from './quickQuestionData'
import { getQuickQuestions } from './quickQuestionData'

interface QuickQuestionsProps {
  /** Questions to display */
  questions?: QuickQuestion[]
  /** Callback when a question is clicked */
  onSelect: (question: string) => void
  /** Layout variant */
  variant?: 'vertical' | 'horizontal'
  /** Optional title for horizontal suggestions */
  title?: string
}

/**
 * QuickQuestions component displays preset question buttons.
 */
export default function QuickQuestions({
  questions,
  onSelect,
  variant = 'horizontal',
  title,
}: QuickQuestionsProps) {
  const items = questions || getQuickQuestions('default')

  if (variant === 'vertical') {
    // Desktop: vertical list
    return (
      <div className="flex flex-col gap-1.5">
        {items.map((q) => (
          <motion.button
            key={q.id}
            whileHover={{ x: 2 }}
            whileTap={{ scale: 0.98 }}
            onClick={() => onSelect(q.text)}
            className="flex items-center gap-2 px-3 py-2.5 rounded-lg text-xs text-slate-400 bg-slate-800/30 border border-slate-700/30 hover:bg-slate-700/40 hover:text-violet-300 transition-all text-left"
          >
            <span className="text-sm">{q.icon}</span>
            <span>{q.label || q.text}</span>
          </motion.button>
        ))}
      </div>
    )
  }

  // Mobile: horizontal scroll
  return (
    <div className="border-t border-slate-700/30 bg-slate-900/70">
      {title && (
        <div className="px-4 pt-2 text-[11px] text-slate-500">
          {title}
        </div>
      )}
      <div className="flex gap-1.5 overflow-x-auto px-4 py-2.5 scrollbar-hide touch-pan-x">
        {items.map((q) => (
          <motion.button
            key={q.id}
            whileTap={{ scale: 0.95 }}
            onClick={() => onSelect(q.text)}
            title={q.text}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs text-slate-300 bg-slate-800/70 border border-slate-700/60 hover:bg-slate-700/70 hover:text-violet-200 hover:border-violet-500/40 transition-all whitespace-nowrap"
          >
            <span>{q.icon}</span>
            <span>{q.label || q.text}</span>
          </motion.button>
        ))}
      </div>
    </div>
  )
}
