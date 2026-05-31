import { motion, AnimatePresence } from 'framer-motion'
import { useTarot } from '../hooks/useTarot'
import SpreadSelector from '../components/tarot/SpreadSelector'
import SpreadLayout from '../components/tarot/SpreadLayout'
import CardReading from '../components/tarot/CardReading'
import AIReading from '../components/common/AIReading'

export default function Tarot() {
  const {
    phase,
    selectedSpread,
    question,
    drawResult,
    revealedCards,
    selectedCardIndex,
    isLoading,
    error,
    setSelectedSpread,
    setQuestion,
    startDraw,
    revealCard,
    selectCard,
    closeCardDetail,
    reset,
    allRevealed,
  } = useTarot()

  const selectedCard =
    selectedCardIndex !== null && drawResult
      ? drawResult.cards[selectedCardIndex]
      : null

  return (
    <div className="flex-1 px-4 py-12 max-w-5xl mx-auto w-full">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center mb-10"
      >
        <h1 className="text-4xl font-bold mb-3">
          <span className="bg-gradient-to-r from-purple-400 via-pink-400 to-amber-400 bg-clip-text text-transparent">
            塔罗牌占卜
          </span>
        </h1>
        <p className="text-slate-400">静心冥想你的问题，让塔罗牌为你揭示命运的指引</p>
      </motion.div>

      {/* Error display */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="mb-6 p-4 bg-red-500/10 border border-red-500/30 rounded-xl text-red-400 text-center"
          >
            {error}
          </motion.div>
        )}
      </AnimatePresence>

      {/* Phase: Input */}
      <AnimatePresence mode="wait">
        {phase === 'input' && (
          <motion.div
            key="input"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0, x: -100 }}
            className="space-y-8"
          >
            {/* Question input */}
            <div>
              <label className="block text-amber-400 text-sm font-medium mb-3">
                ✨ 冥想你的问题（可选）
              </label>
              <input
                type="text"
                value={question}
                onChange={(e) => setQuestion(e.target.value)}
                placeholder="例如：我的事业发展方向应该如何选择？"
                className="w-full px-5 py-4 bg-slate-800/80 border border-slate-700 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500 transition-all"
              />
            </div>

            {/* Spread selection */}
            <div>
              <label className="block text-amber-400 text-sm font-medium mb-3">
                🔮 选择牌阵
              </label>
              <SpreadSelector selected={selectedSpread} onSelect={setSelectedSpread} />
            </div>

            {/* Start button */}
            <div className="text-center pt-4">
              <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                onClick={startDraw}
                disabled={isLoading}
                className="px-10 py-4 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white rounded-xl font-bold text-lg shadow-lg shadow-purple-500/20 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
              >
                {isLoading ? (
                  <span className="flex items-center gap-2">
                    <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                    </svg>
                    抽牌中...
                  </span>
                ) : (
                  '🔮 开始占卜'
                )}
              </motion.button>
            </div>
          </motion.div>
        )}

        {/* Phase: Drawing */}
        {phase === 'drawing' && drawResult && (
          <motion.div
            key="drawing"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="space-y-8"
          >
            {/* Instructions */}
            <div className="text-center">
              <p className="text-slate-300">
                {allRevealed
                  ? '所有牌已翻开，点击任意一张查看详细解读'
                  : '依次点击每张牌，揭示命运的指引'}
              </p>
              <p className="text-slate-500 text-sm mt-1">
                已翻开 {revealedCards.size} / {drawResult.cards.length} 张
              </p>
            </div>

            {/* Spread layout */}
            <SpreadLayout
              spread={selectedSpread}
              cards={drawResult.cards}
              revealedCards={revealedCards}
              onRevealCard={revealCard}
              onCardClick={selectCard}
            />

            {/* AI Reading */}
            {allRevealed && drawResult.record_id && (
              <AIReading
                type="tarot"
                resultId={drawResult.record_id}
                question={drawResult.question}
              />
            )}

            {/* Action buttons */}
            {allRevealed && (
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                className="text-center pt-4 flex justify-center gap-4"
              >
                <motion.button
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                  onClick={reset}
                  className="px-8 py-3 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white rounded-xl font-medium transition-all"
                >
                  🔄 再次占卜
                </motion.button>
              </motion.div>
            )}
          </motion.div>
        )}
      </AnimatePresence>

      {/* Card detail modal */}
      <CardReading
        drawnCard={selectedCard}
        isOpen={selectedCardIndex !== null}
        onClose={closeCardDetail}
      />
    </div>
  )
}
