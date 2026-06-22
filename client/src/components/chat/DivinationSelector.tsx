import { motion } from 'framer-motion'

export interface DivinationType {
  id: string
  name: string
  icon: string
  description: string
}

interface DivinationSelectorProps {
  /** Available divination types */
  types?: DivinationType[]
  /** Currently selected type */
  selected: string
  /** Callback when type is selected */
  onSelect: (typeId: string) => void
  /** Layout variant */
  variant?: 'grid' | 'scroll'
  /** Compact mobile layout after a chat has started */
  compact?: boolean
}

const DEFAULT_TYPES: DivinationType[] = [
  { id: 'liuyao', name: '六爻', icon: '☯️', description: '卦象分析' },
  { id: 'bazi', name: '八字', icon: '📋', description: '命理排盘' },
  { id: 'tarot', name: '塔罗牌', icon: '🎴', description: '抽牌解读' },
  { id: 'horoscope', name: '星座', icon: '⭐', description: '运势解读' },
]

/**
 * DivinationSelector component for choosing divination type.
 * Supports grid layout (desktop) and scroll layout (mobile).
 */
export default function DivinationSelector({
  types = DEFAULT_TYPES,
  selected,
  onSelect,
  variant = 'grid',
  compact = false,
}: DivinationSelectorProps) {
  if (variant === 'scroll') {
    // Mobile: horizontal scroll
    return (
      <div className={`px-4 border-b border-slate-700/30 transition-all ${compact ? 'py-2' : 'py-3'}`}>
        {!compact && <div className="text-xs text-slate-500 mb-2">选择占卜方式</div>}
        <div className="flex gap-2 overflow-x-auto pb-1 scrollbar-hide">
          {types.map((type) => (
            <motion.button
              key={type.id}
              whileTap={{ scale: 0.95 }}
              onClick={() => onSelect(type.id)}
              className={`flex items-center gap-1.5 rounded-full text-xs whitespace-nowrap transition-all ${
                selected === type.id
                  ? 'bg-gradient-to-r from-violet-500/25 to-indigo-500/20 border border-violet-500/40 text-violet-300 font-semibold shadow-lg shadow-violet-500/10'
                  : 'bg-slate-800/50 border border-slate-700/50 text-slate-400 hover:bg-slate-700/50'
              } ${compact ? 'px-2.5 py-1.5' : 'px-3 py-2'}`}
            >
              <span className={compact ? 'text-sm' : 'text-base'}>{type.icon}</span>
              <span>{type.name}</span>
            </motion.button>
          ))}
        </div>
      </div>
    )
  }

  // Desktop: 2x2 grid
  return (
    <div className="grid grid-cols-2 gap-2">
      {types.map((type) => (
        <motion.button
          key={type.id}
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          onClick={() => onSelect(type.id)}
          className={`p-3 rounded-xl text-center transition-all ${
            selected === type.id
              ? 'bg-gradient-to-br from-violet-500/20 to-indigo-500/15 border border-violet-500/35 shadow-lg shadow-violet-500/15'
              : 'bg-slate-800/30 border border-slate-700/30 hover:bg-slate-700/30'
          }`}
        >
          <div className="text-2xl mb-1">{type.icon}</div>
          <div className={`text-xs font-semibold ${selected === type.id ? 'text-violet-300' : 'text-slate-300'}`}>
            {type.name}
          </div>
          <div className="text-[10px] text-slate-500 mt-0.5">{type.description}</div>
        </motion.button>
      ))}
    </div>
  )
}
