import { motion } from 'framer-motion'

export interface SpreadOption {
  id: string
  name: string
  description: string
  cardCount: number
  icon: string
}

const spreadOptions: SpreadOption[] = [
  {
    id: 'single',
    name: '单牌占卜',
    description: '最简洁的方式，一牌定乾坤，快速获得指引',
    cardCount: 1,
    icon: '🃏',
  },
  {
    id: 'three',
    name: '三牌阵',
    description: '过去、现在、未来，揭示事物发展的脉络',
    cardCount: 3,
    icon: '🔮',
  },
  {
    id: 'celtic',
    name: '凯尔特十字',
    description: '最经典的牌阵，全方位深度解读问题',
    cardCount: 10,
    icon: '✦',
  },
  {
    id: 'love',
    name: '爱情十字',
    description: '专为感情问题设计，解读爱情的方方面面',
    cardCount: 5,
    icon: '💕',
  },
]

interface SpreadSelectorProps {
  selected: string
  onSelect: (spreadId: string) => void
}

function SpreadPreview({ cardCount }: { cardCount: number }) {
  const layouts: Record<string, { x: number; y: number }[]> = {
    single: [{ x: 50, y: 50 }],
    three: [
      { x: 20, y: 50 },
      { x: 50, y: 50 },
      { x: 80, y: 50 },
    ],
    celtic: [
      { x: 50, y: 50 },
      { x: 50, y: 50 },
      { x: 50, y: 20 },
      { x: 50, y: 80 },
      { x: 20, y: 50 },
      { x: 80, y: 50 },
      { x: 85, y: 20 },
      { x: 85, y: 40 },
      { x: 85, y: 60 },
      { x: 85, y: 80 },
    ],
    love: [
      { x: 50, y: 20 },
      { x: 50, y: 50 },
      { x: 20, y: 50 },
      { x: 80, y: 50 },
      { x: 50, y: 80 },
    ],
  }

  const positions = layouts[cardCount === 1 ? 'single' : cardCount === 3 ? 'three' : cardCount === 10 ? 'celtic' : 'love']

  return (
    <div className="relative w-full h-24">
      {positions.map((pos, i) => (
        <div
          key={i}
          className="absolute w-4 h-6 bg-purple-600/40 border border-purple-400/30 rounded-sm"
          style={{
            left: `${pos.x}%`,
            top: `${pos.y}%`,
            transform: 'translate(-50%, -50%)',
          }}
        />
      ))}
    </div>
  )
}

export default function SpreadSelector({ selected, onSelect }: SpreadSelectorProps) {
  return (
    <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
      {spreadOptions.map((option) => (
        <motion.button
          key={option.id}
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          onClick={() => onSelect(option.id)}
          className={`relative p-4 rounded-xl border-2 transition-all duration-300 text-left ${
            selected === option.id
              ? 'border-purple-500 bg-purple-500/10 shadow-lg shadow-purple-500/10'
              : 'border-slate-700 bg-slate-800/50 hover:border-slate-600'
          }`}
        >
          <div className="text-2xl mb-2">{option.icon}</div>
          <h3 className="text-white font-semibold text-sm mb-1">{option.name}</h3>
          <p className="text-slate-400 text-xs leading-relaxed mb-3">{option.description}</p>
          <SpreadPreview cardCount={option.cardCount} />
          <div className="mt-2 text-xs text-slate-500">{option.cardCount} 张牌</div>
        </motion.button>
      ))}
    </div>
  )
}
