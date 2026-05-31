import { useState } from 'react'
import { motion } from 'framer-motion'
import type { DrawnCard } from '../../types/tarot'

interface TarotCardComponentProps {
  drawnCard: DrawnCard
  isRevealed: boolean
  onReveal?: () => void
  onClick?: () => void
  layoutStyle?: React.CSSProperties
  delay?: number
}

const cardBackDesign = (
  <div className="absolute inset-0 bg-gradient-to-br from-purple-900 to-indigo-900 rounded-xl flex items-center justify-center border-2 border-purple-500/30">
    <div className="relative">
      <div className="w-16 h-16 border-2 border-amber-400/50 rounded-full flex items-center justify-center">
        <div className="w-10 h-10 border border-amber-400/30 rounded-full flex items-center justify-center">
          <span className="text-amber-400 text-lg">✦</span>
        </div>
      </div>
      <div className="absolute -top-1 -left-1 w-18 h-18 border border-purple-400/20 rounded-full" />
    </div>
    {/* Decorative corners */}
    <div className="absolute top-2 left-2 w-4 h-4 border-t-2 border-l-2 border-amber-400/30 rounded-tl" />
    <div className="absolute top-2 right-2 w-4 h-4 border-t-2 border-r-2 border-amber-400/30 rounded-tr" />
    <div className="absolute bottom-2 left-2 w-4 h-4 border-b-2 border-l-2 border-amber-400/30 rounded-bl" />
    <div className="absolute bottom-2 right-2 w-4 h-4 border-b-2 border-r-2 border-amber-400/30 rounded-br" />
  </div>
)

export default function TarotCardComponent({
  drawnCard,
  isRevealed,
  onReveal,
  onClick,
  layoutStyle,
  delay = 0,
}: TarotCardComponentProps) {
  const [isFlipped, setIsFlipped] = useState(false)
  const { card, orientation, position_name } = drawnCard
  const isReversed = orientation === 'reversed'
  const position = { name: position_name }

  const handleClick = () => {
    if (!isRevealed) {
      setIsFlipped(true)
      onReveal?.()
    } else {
      onClick?.()
    }
  }

  return (
    <motion.div
      className="relative cursor-pointer"
      style={layoutStyle}
      initial={{ opacity: 0, scale: 0.8 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ delay, duration: 0.4 }}
    >
      {/* Position label */}
      <div className="text-center mb-2">
        <span className="text-xs text-purple-300/70 font-medium">{position.name}</span>
      </div>

      <div
        className="relative w-32 h-52 sm:w-36 sm:h-56"
        style={{ perspective: '1000px' }}
        onClick={handleClick}
      >
        <motion.div
          className="absolute inset-0 w-full h-full"
          style={{ transformStyle: 'preserve-3d' }}
          animate={{ rotateY: isFlipped || isRevealed ? 180 : 0 }}
          transition={{ duration: 0.6, ease: 'easeInOut' }}
        >
          {/* Card Back */}
          <div
            className="absolute inset-0 w-full h-full"
            style={{ backfaceVisibility: 'hidden' }}
          >
            {cardBackDesign}
          </div>

          {/* Card Front */}
          <div
            className="absolute inset-0 w-full h-full rounded-xl overflow-hidden"
            style={{
              backfaceVisibility: 'hidden',
              transform: `rotateY(180deg) ${isReversed ? 'rotate(180deg)' : ''}`,
            }}
          >
            <div className="w-full h-full bg-gradient-to-b from-slate-800 to-slate-900 border-2 border-amber-500/40 rounded-xl flex flex-col items-center justify-center p-3">
              {/* Card image placeholder */}
              <div className="w-20 h-28 bg-gradient-to-br from-purple-800/50 to-indigo-800/50 rounded-lg flex items-center justify-center mb-3 border border-purple-500/20">
                <span className="text-3xl">
                  {card.type === 'major' ? '✦' : card.suit === 'wands' ? '🪄' : card.suit === 'cups' ? '🏆' : card.suit === 'swords' ? '⚔️' : '💰'}
                </span>
              </div>
              <h4 className="text-amber-300 text-xs font-bold text-center leading-tight">
                {card.name}
              </h4>
              <p className="text-slate-400 text-[10px] mt-1">{card.name_en}</p>
              {isReversed && (
                <span className="mt-1 text-[10px] text-red-400 font-medium">逆位</span>
              )}
            </div>
          </div>
        </motion.div>
      </div>

      {!isRevealed && (
        <motion.p
          className="text-center text-xs text-slate-500 mt-2"
          animate={{ opacity: [0.5, 1, 0.5] }}
          transition={{ repeat: Infinity, duration: 2 }}
        >
          点击翻牌
        </motion.p>
      )}
    </motion.div>
  )
}
