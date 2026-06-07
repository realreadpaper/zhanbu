import { useState, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { throwHexagramsV2 } from '../services/liuyao'
import { useAI } from '../hooks/useAI'
import AIReading from '../components/common/AIReading'
import type { LiuYaoV2Result, LineResult, TakashimaHexagram, TakashimaText } from '../types/liuyao'

const CASTING_RITUAL_MS = 1800

const wait = (ms: number) => new Promise((resolve) => window.setTimeout(resolve, ms))

function takashimaText(value?: TakashimaText): string {
  if (!value) return ''
  if (typeof value === 'string') return value
  return value.text || ''
}

function CastingRitual() {
  const yarrowStalks = ['蓍', '草', '筮']

  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -12 }}
      className="mx-auto max-w-xl rounded-xl border border-amber-500/20 bg-slate-900/60 p-6 shadow-xl shadow-amber-950/30"
    >
      <div className="relative mx-auto mb-5 flex h-32 w-32 items-center justify-center">
        <motion.div
          animate={{ rotate: 360 }}
          transition={{ repeat: Infinity, duration: 5, ease: 'linear' }}
          className="absolute inset-0 rounded-full border border-dashed border-amber-400/40"
        />
        <motion.div
          animate={{ scale: [0.9, 1.08, 0.9], opacity: [0.35, 0.8, 0.35] }}
          transition={{ repeat: Infinity, duration: 1.6, ease: 'easeInOut' }}
          className="absolute h-20 w-20 rounded-full bg-amber-400/10 blur-sm"
        />
        {yarrowStalks.map((stalk, index) => (
          <motion.div
            key={stalk}
            animate={{
              rotateY: [0, 180, 360],
              y: [0, -12, 0],
              x: [0, index === 1 ? 8 : -8, 0],
            }}
            transition={{
              repeat: Infinity,
              duration: 0.9 + index * 0.16,
              ease: 'easeInOut',
            }}
            className="relative mx-1 flex h-10 w-10 items-center justify-center rounded-full border border-amber-300/50 bg-gradient-to-br from-amber-200 to-amber-700 text-sm font-bold text-slate-950 shadow-lg shadow-amber-900/30"
          >
            {stalk}
          </motion.div>
        ))}
      </div>

      <div className="space-y-2">
        {['初', '二', '三', '四', '五', '上'].map((label, index) => (
          <motion.div
            key={label}
            initial={{ opacity: 0.25 }}
            animate={{ opacity: [0.25, 1, 0.25] }}
            transition={{ repeat: Infinity, duration: 1.2, delay: index * 0.12 }}
            className="grid grid-cols-[2.5rem_1fr] items-center gap-3 text-sm"
          >
            <span className="text-amber-300">{label}爻</span>
            <span className="h-2 rounded-full bg-gradient-to-r from-amber-500/20 via-yellow-300/70 to-amber-500/20" />
          </motion.div>
        ))}
      </div>

      <p className="mt-5 text-center text-sm text-slate-400">蓍草揲筮，卦象生成...</p>
    </motion.div>
  )
}

function LineDisplay({ line, index }: { line: LineResult; index: number }) {
  const isYang = line.value === 7 || line.value === 9
  const labels = ['初爻', '二爻', '三爻', '四爻', '五爻', '上爻']

  return (
    <motion.div
      initial={{ opacity: 0, x: -20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: index * 0.1 }}
      className="flex items-center gap-4 py-2"
    >
      <span className="text-slate-500 text-sm w-12">{labels[index]}</span>
      <div className="flex-1 flex items-center gap-2">
        <span className="text-2xl font-mono tracking-widest">
          {isYang ? '⚊' : '⚋'}
        </span>
        <span className="text-slate-400 text-sm">
          {line.type === 'old_yang' ? '老阳 ▲' :
           line.type === 'young_yang' ? '少阳' :
           line.type === 'old_yin' ? '老阴 ▼' : '少阴'}
        </span>
        {line.mutable && (
          <span className="px-2 py-0.5 bg-amber-500/20 text-amber-400 text-xs rounded-full">动</span>
        )}
      </div>
    </motion.div>
  )
}

function TakashimaHexagramDisplay({ hexagram, label }: { hexagram: TakashimaHexagram; label: string }) {
  const shortName = hexagram.name_short || hexagram.name
  const fullName = hexagram.full_name || hexagram.name
  const judgment = takashimaText(hexagram.judgment)
  const tuan = takashimaText(hexagram.tuan)
  const image = takashimaText(hexagram.image)

  return (
    <div className="bg-slate-800/50 border border-amber-500/20 rounded-xl p-5">
      <h3 className="text-amber-400 text-sm font-medium mb-2">{label}</h3>
      <div className="flex items-baseline gap-3 mb-2">
        <span className="text-3xl font-bold text-white">{shortName}</span>
        <span className="text-lg text-slate-300">{fullName}</span>
      </div>

      {/* 卦辞 */}
      <div className="mt-3 p-3 bg-slate-900/50 rounded-lg">
        <p className="text-amber-300 text-sm font-medium mb-1">卦辞</p>
        <p className="text-slate-300 text-sm leading-relaxed">{judgment}</p>
      </div>

      {/* 彖传 */}
      {tuan && (
        <div className="mt-2 p-3 bg-slate-900/30 rounded-lg">
          <p className="text-amber-200 text-xs font-medium mb-1">彖传</p>
          <p className="text-slate-400 text-xs leading-relaxed">{tuan}</p>
        </div>
      )}

      {/* 象辞 */}
      {image && (
        <div className="mt-2 p-3 bg-slate-900/30 rounded-lg">
          <p className="text-amber-200 text-xs font-medium mb-1">象辞</p>
          <p className="text-slate-400 text-xs leading-relaxed">{image}</p>
        </div>
      )}
    </div>
  )
}

export default function LiuYaoV2() {
  const [question, setQuestion] = useState('')
  const [result, setResult] = useState<LiuYaoV2Result | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [method, setMethod] = useState<string>('yarrow')

  const ai = useAI({
    type: 'liuyao_v2',
    resultId: result?.record_id,
    result: result && !result.record_id ? JSON.stringify(result) : undefined,
    question: result?.question,
  })

  const handleThrow = useCallback(async () => {
    setIsLoading(true)
    setError(null)

    try {
      const [data] = await Promise.all([
        throwHexagramsV2(question || undefined, method),
        wait(CASTING_RITUAL_MS),
      ])
      setResult(data)
    } catch (err) {
      setError('投掷失败，请稍后重试')
      console.error('LiuYao v2 throw error:', err)
    } finally {
      setIsLoading(false)
    }
  }, [question, method])

  const handleReset = useCallback(() => {
    setResult(null)
    setQuestion('')
    setError(null)
    ai.reset()
  }, [ai])

  return (
    <div className="flex-1 px-4 py-12 max-w-3xl mx-auto w-full">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center mb-10"
      >
        <h1 className="text-4xl font-bold mb-3">
          <span className="bg-gradient-to-r from-amber-400 via-orange-400 to-red-400 bg-clip-text text-transparent">
            高岛易断
          </span>
        </h1>
        <p className="text-slate-400">传承高岛嘉右卫门易学智慧</p>
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

      <AnimatePresence mode="wait">
        {!result ? (
          /* Input Phase */
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
                ☯ 请输入你的问题（可选）
              </label>
              <input
                type="text"
                value={question}
                onChange={(e) => setQuestion(e.target.value)}
                placeholder="例如：此事能否成功？近期运势如何？"
                className="w-full px-5 py-4 bg-slate-800/80 border border-slate-700 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-amber-500 focus:ring-1 focus:ring-amber-500 transition-all"
              />
            </div>

            {/* Method selection */}
            <div>
              <label className="block text-amber-400 text-sm font-medium mb-3">
                起卦方法
              </label>
              <div className="flex gap-4">
                <button
                  onClick={() => setMethod('yarrow')}
                  className={`flex-1 py-3 px-4 rounded-xl border transition-all ${
                    method === 'yarrow'
                      ? 'border-amber-500 bg-amber-500/10 text-amber-400'
                      : 'border-slate-700 bg-slate-800/50 text-slate-400 hover:border-slate-600'
                  }`}
                >
                  <div className="text-lg mb-1">蓍草揲筮法</div>
                  <div className="text-xs opacity-70">古法正宗，大衍之数</div>
                </button>
                <button
                  onClick={() => setMethod('coin')}
                  className={`flex-1 py-3 px-4 rounded-xl border transition-all ${
                    method === 'coin'
                      ? 'border-amber-500 bg-amber-500/10 text-amber-400'
                      : 'border-slate-700 bg-slate-800/50 text-slate-400 hover:border-slate-600'
                  }`}
                >
                  <div className="text-lg mb-1">铜钱摇卦法</div>
                  <div className="text-xs opacity-70">三枚铜钱，简便快捷</div>
                </button>
              </div>
            </div>

            {/* Throw button */}
            <div className="text-center pt-4">
              <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                onClick={handleThrow}
                disabled={isLoading}
                className="px-10 py-4 bg-gradient-to-r from-amber-600 to-orange-600 hover:from-amber-700 hover:to-orange-700 text-white rounded-xl font-bold text-lg shadow-lg shadow-amber-500/20 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
              >
                {isLoading ? (
                  <span className="flex items-center gap-2">
                    <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                    </svg>
                    掷蓍草中...
                  </span>
                ) : (
                  method === 'yarrow' ? '🌿 操筮起卦' : '🪙 掷铜钱'
                )}
              </motion.button>
              <p className="text-slate-500 text-sm mt-3">
                {method === 'yarrow' ? '大衍之数五十，其用四十有九' : '模拟三枚铜钱投掷六次'}
              </p>
            </div>

            <AnimatePresence>
              {isLoading && <CastingRitual />}
            </AnimatePresence>
          </motion.div>
        ) : (
          /* Result Phase */
          <motion.div
            key="result"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="space-y-8"
          >
            {/* Method display */}
            <div className="text-center">
              <span className="px-3 py-1 bg-amber-500/10 border border-amber-500/30 rounded-full text-amber-400 text-sm">
                {result.method === 'yarrow' ? '蓍草揲筮法' : '铜钱摇卦法'}
              </span>
            </div>

            {/* Question display */}
            {result.question && (
              <div className="text-center">
                <p className="text-slate-400 text-sm">所问之事</p>
                <p className="text-white text-lg mt-1">{result.question}</p>
              </div>
            )}

            {/* Lines display */}
            <div className="bg-slate-800/30 border border-slate-700/50 rounded-xl p-6">
              <h3 className="text-amber-400 font-semibold mb-4 flex items-center gap-2">
                <span>☯</span> 爻象
              </h3>
              {result.lines && [...result.lines].reverse().map((line, i) => (
                <LineDisplay key={i} line={line} index={5 - i} />
              ))}
            </div>

            {/* Hexagrams */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <TakashimaHexagramDisplay hexagram={result.ben_gua} label="本卦" />
              {result.bian_gua && (
                <TakashimaHexagramDisplay hexagram={result.bian_gua} label="变卦" />
              )}
            </div>

            {/* Mutable lines info */}
            {result.mutable_lines && result.mutable_lines.length > 0 && (
              <div className="text-center text-sm text-amber-400">
                动爻：{result.mutable_lines.map(i => ['初', '二', '三', '四', '五', '上'][i] + '爻').join('、')}
              </div>
            )}

            {/* AI Reading */}
            {result && (
              <AIReading
                type="liuyao_v2"
                resultId={result.record_id}
                result={!result.record_id ? JSON.stringify(result) : undefined}
                question={result.question}
              />
            )}

            {/* Reset button */}
            <div className="text-center pt-4">
              <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                onClick={handleReset}
                className="px-8 py-3 bg-gradient-to-r from-amber-600 to-orange-600 hover:from-amber-700 hover:to-orange-700 text-white rounded-xl font-medium transition-all"
              >
                🔄 再次占卜
              </motion.button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}
