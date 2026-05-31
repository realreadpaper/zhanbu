import { motion, AnimatePresence } from 'framer-motion'
import type { DrawnCard } from '../../types/tarot'

interface CardReadingProps {
  drawnCard: DrawnCard | null
  isOpen: boolean
  onClose: () => void
}

export default function CardReading({ drawnCard, isOpen, onClose }: CardReadingProps) {
  if (!drawnCard) return null

  const { card, orientation, position_name } = drawnCard
  const isReversed = orientation === 'reversed'
  const meaning = isReversed ? card.meaning_down : card.meaning_up
  const positionLabel = isReversed ? '逆位' : '正位'
  const position = { name: position_name, description: '' }

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-black/70 backdrop-blur-sm z-40"
          />

          {/* Modal */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.9, y: 20 }}
            className="fixed inset-0 z-50 flex items-center justify-center p-4"
          >
            <div className="bg-gradient-to-b from-slate-800 to-slate-900 rounded-2xl border border-purple-500/30 max-w-lg w-full max-h-[85vh] overflow-y-auto shadow-2xl shadow-purple-500/10">
              {/* Header */}
              <div className="relative p-6 border-b border-slate-700/50">
                <button
                  onClick={onClose}
                  className="absolute top-4 right-4 text-slate-400 hover:text-white transition-colors"
                >
                  <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>

                <div className="flex items-start gap-4">
                  <div className="w-16 h-24 bg-gradient-to-br from-purple-800/50 to-indigo-800/50 rounded-lg flex items-center justify-center border border-purple-500/20 flex-shrink-0">
                    <span className="text-3xl">
                      {card.type === 'major' ? '✦' : card.suit === 'wands' ? '🪄' : card.suit === 'cups' ? '🏆' : card.suit === 'swords' ? '⚔️' : '💰'}
                    </span>
                  </div>
                  <div>
                    <p className="text-purple-300 text-xs font-medium mb-1">{position.name}</p>
                    <h2 className="text-2xl font-bold text-white">{card.name}</h2>
                    <p className="text-slate-400 text-sm">{card.name_en}</p>
                    <div className="flex items-center gap-2 mt-2">
                      <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
                        isReversed
                          ? 'bg-red-500/20 text-red-400'
                          : 'bg-green-500/20 text-green-400'
                      }`}>
                        {positionLabel}
                      </span>
                      <span className="px-2 py-0.5 rounded-full text-xs font-medium bg-purple-500/20 text-purple-400">
                        {card.type === 'major' ? '大阿卡纳' : '小阿卡纳'}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Content */}
              <div className="p-6 space-y-5">
                {/* Position meaning */}
                <div>
                  <h3 className="text-amber-400 text-sm font-semibold mb-2 flex items-center gap-2">
                    <span>📍</span> 牌位含义
                  </h3>
                  <p className="text-slate-300 text-sm leading-relaxed">{position.description}</p>
                </div>

                {/* Card meaning */}
                <div>
                  <h3 className="text-amber-400 text-sm font-semibold mb-2 flex items-center gap-2">
                    <span>🔮</span> {positionLabel}解读
                  </h3>
                  <p className="text-slate-300 text-sm leading-relaxed">{meaning}</p>
                </div>

                {/* Description */}
                <div>
                  <h3 className="text-amber-400 text-sm font-semibold mb-2 flex items-center gap-2">
                    <span>📖</span> 牌面描述
                  </h3>
                  <p className="text-slate-300 text-sm leading-relaxed">{card.description}</p>
                </div>
              </div>

              {/* Footer */}
              <div className="p-4 border-t border-slate-700/50">
                <button
                  onClick={onClose}
                  className="w-full py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg font-medium transition-colors"
                >
                  关闭
                </button>
              </div>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  )
}
