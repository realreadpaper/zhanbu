import TarotCardComponent from './TarotCardComponent'
import type { DrawnCard } from '../../types/tarot'

interface SpreadLayoutProps {
  spread: string
  cards: DrawnCard[]
  revealedCards: Set<number>
  onRevealCard: (index: number) => void
  onCardClick: (index: number) => void
}

const layoutConfigs: Record<string, { positions: { x: number; y: number }[]; label: string }> = {
  single: {
    label: '单牌占卜',
    positions: [{ x: 50, y: 50 }],
  },
  three: {
    label: '三牌阵',
    positions: [
      { x: 17, y: 50 },
      { x: 50, y: 50 },
      { x: 83, y: 50 },
    ],
  },
  celtic: {
    label: '凯尔特十字',
    positions: [
      { x: 35, y: 50 }, // 1 - 当前情况
      { x: 35, y: 50 }, // 2 - 障碍
      { x: 35, y: 15 }, // 3 - 目标/意识
      { x: 35, y: 85 }, // 4 - 根基/潜意识
      { x: 10, y: 50 }, // 5 - 过去
      { x: 60, y: 50 }, // 6 - 未来
      { x: 80, y: 85 }, // 7 - 自我认知
      { x: 80, y: 65 }, // 8 - 环境
      { x: 80, y: 45 }, // 9 - 希望/恐惧
      { x: 80, y: 25 }, // 10 - 最终结果
    ],
  },
  love: {
    label: '爱情十字',
    positions: [
      { x: 50, y: 15 }, // 1 - 你的感受
      { x: 50, y: 50 }, // 2 - 对方的感受
      { x: 20, y: 50 }, // 3 - 关系基础
      { x: 80, y: 50 }, // 4 - 过去影响
      { x: 50, y: 85 }, // 5 - 关系走向
    ],
  },
}

export default function SpreadLayout({
  spread,
  cards,
  revealedCards,
  onRevealCard,
  onCardClick,
}: SpreadLayoutProps) {
  const config = layoutConfigs[spread] || layoutConfigs.single

  return (
    <div className="relative w-full max-w-3xl mx-auto">
      <h3 className="text-center text-amber-400 text-lg font-semibold mb-6">
        {config.label}
      </h3>

      <div className="relative" style={{ minHeight: spread === 'celtic' ? '480px' : '300px' }}>
        {cards.map((drawnCard, index) => {
          const pos = config.positions[index] || { x: 50, y: 50 }
          const isCelticCross = spread === 'celtic' && (index === 0 || index === 1)

          return (
            <div
              key={index}
              className="absolute"
              style={{
                left: `${pos.x}%`,
                top: `${pos.y}%`,
                transform: 'translate(-50%, -50%)',
                zIndex: isCelticCross && index === 1 ? 10 : 1,
              }}
            >
              <TarotCardComponent
                drawnCard={drawnCard}
                isRevealed={revealedCards.has(index)}
                onReveal={() => onRevealCard(index)}
                onClick={() => onCardClick(index)}
                delay={index * 0.15}
              />
            </div>
          )
        })}
      </div>
    </div>
  )
}
