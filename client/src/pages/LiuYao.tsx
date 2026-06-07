import { useState, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { throwHexagrams } from '../services/liuyao'
import { useAI } from '../hooks/useAI'
import AIReading from '../components/common/AIReading'
import type { LiuYaoResult, LineResult } from '../types/liuyao'

const CASTING_RITUAL_MS = 2800

const wait = (ms: number) => new Promise((resolve) => window.setTimeout(resolve, ms))

function CastingRitual() {
  const coinFaces = ['乾', '坤', '震']
  const lineLabels = ['初', '二', '三', '四', '五', '上']

  return (
    <motion.div
      key="casting"
      initial={{ opacity: 0, scale: 0.96, y: 18 }}
      animate={{ opacity: 1, scale: 1, y: 0 }}
      exit={{ opacity: 0, scale: 0.98, y: -18 }}
      transition={{ duration: 0.28 }}
      className="mx-auto max-w-xl py-8 text-center"
    >
      <div className="relative mx-auto mb-8 flex h-48 w-48 items-center justify-center">
        <motion.div
          animate={{ rotate: 360 }}
          transition={{ repeat: Infinity, duration: 4.8, ease: 'linear' }}
          className="absolute inset-0 rounded-full border border-dashed border-emerald-400/40"
        />
        <motion.div
          animate={{ rotate: -360 }}
          transition={{ repeat: Infinity, duration: 7, ease: 'linear' }}
          className="absolute inset-5 rounded-full border border-teal-300/20"
        />
        <motion.div
          animate={{ scale: [0.88, 1.16, 0.88], opacity: [0.25, 0.85, 0.25] }}
          transition={{ repeat: Infinity, duration: 1.6, ease: 'easeInOut' }}
          className="absolute h-28 w-28 rounded-full bg-emerald-400/15 blur-md"
        />
        {coinFaces.map((face, index) => (
          <motion.div
            key={face}
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
            className="relative mx-1 flex h-14 w-14 items-center justify-center rounded-full border border-amber-300/60 bg-gradient-to-br from-amber-200 to-amber-700 text-base font-bold text-slate-950 shadow-lg shadow-amber-900/30"
          >
            {face}
          </motion.div>
        ))}
      </div>

      <div className="mx-auto max-w-sm space-y-2">
        {lineLabels.map((label, index) => (
          <motion.div
            key={label}
            initial={{ opacity: 0.25 }}
            animate={{ opacity: [0.25, 1, 0.25] }}
            transition={{ repeat: Infinity, duration: 1.2, delay: index * 0.12 }}
            className="grid grid-cols-[2.5rem_1fr] items-center gap-3 text-sm"
          >
            <span className="text-emerald-300">{label}爻</span>
            <span className="h-2 rounded-full bg-gradient-to-r from-emerald-500/20 via-cyan-300/70 to-emerald-500/20" />
          </motion.div>
        ))}
      </div>

      <motion.p
        animate={{ opacity: [0.45, 1, 0.45] }}
        transition={{ repeat: Infinity, duration: 1.4, ease: 'easeInOut' }}
        className="mt-6 text-sm text-emerald-300"
      >
        凝神起念，铜钱旋落，六爻成象...
      </motion.p>
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

function HexagramDisplay({ hexagram, label }: { hexagram: { name: string; name_short: string; judgment: string; description: string }; label: string }) {
  return (
    <div className="bg-slate-800/50 border border-emerald-500/20 rounded-xl p-5">
      <h3 className="text-emerald-400 text-sm font-medium mb-2">{label}</h3>
      <div className="flex items-baseline gap-3 mb-2">
        <span className="text-3xl font-bold text-white">{hexagram.name_short}</span>
        <span className="text-lg text-slate-300">{hexagram.name}</span>
      </div>
      <p className="text-slate-400 text-sm leading-relaxed">{hexagram.judgment}</p>
      {hexagram.description && (
        <p className="text-slate-500 text-xs mt-2">{hexagram.description}</p>
      )}
    </div>
  )
}

export default function LiuYao() {
  const [question, setQuestion] = useState('')
  const [result, setResult] = useState<LiuYaoResult | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const ai = useAI({
    type: 'liuyao',
    resultId: result?.record_id,
    result: result && !result.record_id ? JSON.stringify(result) : undefined,
    question: result?.question,
  })

  const handleThrow = useCallback(async () => {
    setIsLoading(true)
    setError(null)

    try {
      const [data] = await Promise.all([
        throwHexagrams(question || undefined),
        wait(CASTING_RITUAL_MS),
      ])
      setResult(data)
    } catch (err) {
      setError('投掷失败，请稍后重试')
      console.error('LiuYao throw error:', err)
    } finally {
      setIsLoading(false)
    }
  }, [question])

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
          <span className="bg-gradient-to-r from-emerald-400 via-teal-400 to-cyan-400 bg-clip-text text-transparent">
            六爻占卜
          </span>
        </h1>
        <p className="text-slate-400">心诚则灵，静心冥想你的问题</p>
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
        {isLoading ? (
          <CastingRitual />
        ) : !result ? (
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
              <label className="block text-emerald-400 text-sm font-medium mb-3">
                ☯ 请输入你的问题（可选）
              </label>
              <input
                type="text"
                value={question}
                onChange={(e) => setQuestion(e.target.value)}
                placeholder="例如：此事能否成功？近期运势如何？"
                className="w-full px-5 py-4 bg-slate-800/80 border border-slate-700 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all"
              />
            </div>

            {/* Throw button */}
            <div className="text-center pt-4">
              <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                onClick={handleThrow}
                disabled={isLoading}
                className="px-10 py-4 bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 text-white rounded-xl font-bold text-lg shadow-lg shadow-emerald-500/20 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
              >
                {isLoading ? (
                  <span className="flex items-center gap-2">
                    <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                    </svg>
                    掷铜钱中...
                  </span>
                ) : (
                  '🪙 掷铜钱'
                )}
              </motion.button>
              <p className="text-slate-500 text-sm mt-3">模拟三枚铜钱投掷六次，生成卦象</p>
            </div>
          </motion.div>
        ) : (
          /* Result Phase */
          <motion.div
            key="result"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="space-y-8"
          >
            {/* Question display */}
            {result.question && (
              <div className="text-center">
                <p className="text-slate-400 text-sm">所问之事</p>
                <p className="text-white text-lg mt-1">{result.question}</p>
              </div>
            )}

            {/* Lines display */}
            <div className="bg-slate-800/30 border border-slate-700/50 rounded-xl p-6">
              <h3 className="text-emerald-400 font-semibold mb-4 flex items-center gap-2">
                <span>☯</span> 爻象
              </h3>
              {result.lines && [...result.lines].reverse().map((line, i) => (
                <LineDisplay key={i} line={line} index={5 - i} />
              ))}
            </div>

            {/* Hexagrams */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <HexagramDisplay hexagram={result.ben_gua} label="本卦" />
              {result.bian_gua && (
                <HexagramDisplay hexagram={result.bian_gua} label="变卦" />
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
                type="liuyao"
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
                className="px-8 py-3 bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 text-white rounded-xl font-medium transition-all"
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
