import { motion, AnimatePresence } from 'framer-motion'
import { useHoroscope } from '../hooks/useHoroscope'
import { ZODIAC_SIGNS, ELEMENT_COLORS } from '../types/horoscope'
import type { HoroscopeResult } from '../types/horoscope'
import AIReading from '../components/common/AIReading'

// ─── Zodiac Wheel ────────────────────────────────────────────────────────────

function ZodiacWheel({
  selected,
  onSelect,
}: {
  selected: string | null
  onSelect: (name: string) => void
}) {
  const radius = 140
  const centerX = 170
  const centerY = 170

  return (
    <div className="flex justify-center">
      <svg width="340" height="340" viewBox="0 0 340 340">
        {/* Outer ring */}
        <circle
          cx={centerX}
          cy={centerY}
          r={radius + 20}
          fill="none"
          stroke="rgba(139,92,246,0.15)"
          strokeWidth="1"
        />
        <circle
          cx={centerX}
          cy={centerY}
          r={radius - 20}
          fill="none"
          stroke="rgba(139,92,246,0.1)"
          strokeWidth="1"
        />

        {ZODIAC_SIGNS.map((sign, i) => {
          const angle = (i * 30 - 90) * (Math.PI / 180)
          const x = centerX + radius * Math.cos(angle)
          const y = centerY + radius * Math.sin(angle)
          const isSelected = selected === sign.name
          const color = ELEMENT_COLORS[sign.element]

          return (
            <g
              key={sign.name}
              onClick={() => onSelect(sign.name)}
              className="cursor-pointer"
            >
              {/* Selection highlight */}
              {isSelected && (
                <circle
                  cx={x}
                  cy={y}
                  r={30}
                  fill={color}
                  opacity={0.15}
                />
              )}
              {/* Symbol circle */}
              <circle
                cx={x}
                cy={y}
                r={24}
                fill={isSelected ? color : 'rgba(37,37,66,0.8)'}
                stroke={isSelected ? color : 'rgba(139,92,246,0.3)'}
                strokeWidth={isSelected ? 2 : 1}
                className="transition-all duration-200"
              />
              {/* Zodiac symbol */}
              <text
                x={x}
                y={y}
                textAnchor="middle"
                dominantBaseline="central"
                fill={isSelected ? '#fff' : '#e2e8f0'}
                fontSize="16"
                className="pointer-events-none select-none"
              >
                {sign.symbol}
              </text>
              {/* Name label */}
              <text
                x={x}
                y={y + 34}
                textAnchor="middle"
                fill={isSelected ? color : '#94a3b8'}
                fontSize="11"
                className="pointer-events-none select-none"
              >
                {sign.name_cn}
              </text>
            </g>
          )
        })}

        {/* Center decoration */}
        <circle
          cx={centerX}
          cy={centerY}
          r={40}
          fill="rgba(37,37,66,0.6)"
          stroke="rgba(139,92,246,0.2)"
          strokeWidth="1"
        />
        <text
          x={centerX}
          y={centerY - 6}
          textAnchor="middle"
          fill="#a78bfa"
          fontSize="12"
        >
          ✦
        </text>
        <text
          x={centerX}
          y={centerY + 12}
          textAnchor="middle"
          fill="#94a3b8"
          fontSize="10"
        >
          星座运势
        </text>
      </svg>
    </div>
  )
}

// ─── Radar Chart ─────────────────────────────────────────────────────────────

function RadarChart({ result }: { result: HoroscopeResult }) {
  const size = 240
  const cx = size / 2
  const cy = size / 2
  const maxR = 90

  const dims = [
    { label: '综合', value: result.overall, key: 'overall' },
    { label: '爱情', value: result.love, key: 'love' },
    { label: '事业', value: result.career, key: 'career' },
    { label: '财运', value: result.wealth, key: 'wealth' },
    { label: '健康', value: result.health, key: 'health' },
  ]

  const n = dims.length
  const angleStep = (2 * Math.PI) / n

  // Grid levels (1-5)
  const gridLevels = [1, 2, 3, 4, 5]

  function getPoint(index: number, value: number) {
    const angle = angleStep * index - Math.PI / 2
    const r = (value / 5) * maxR
    return {
      x: cx + r * Math.cos(angle),
      y: cy + r * Math.sin(angle),
    }
  }

  // Data polygon
  const dataPoints = dims.map((d, i) => getPoint(i, d.value))
  const polygonPath =
    dataPoints.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ') + ' Z'

  // Axis endpoints
  const axisEnds = dims.map((_, i) => getPoint(i, 5))

  return (
    <div className="flex justify-center">
      <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
        {/* Grid polygons */}
        {gridLevels.map((level) => {
          const pts = dims.map((_, i) => getPoint(i, level))
          const path =
            pts.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ') + ' Z'
          return (
            <path
              key={level}
              d={path}
              fill="none"
              stroke="rgba(139,92,246,0.15)"
              strokeWidth="1"
            />
          )
        })}

        {/* Axis lines */}
        {axisEnds.map((end, i) => (
          <line
            key={i}
            x1={cx}
            y1={cy}
            x2={end.x}
            y2={end.y}
            stroke="rgba(139,92,246,0.2)"
            strokeWidth="1"
          />
        ))}

        {/* Data polygon */}
        <path
          d={polygonPath}
          fill="rgba(139,92,246,0.2)"
          stroke="#8b5cf6"
          strokeWidth="2"
        />

        {/* Data points */}
        {dataPoints.map((p, i) => (
          <circle
            key={i}
            cx={p.x}
            cy={p.y}
            r={4}
            fill="#8b5cf6"
            stroke="#fff"
            strokeWidth="1.5"
          />
        ))}

        {/* Labels */}
        {dims.map((d, i) => {
          const labelPoint = getPoint(i, 5.8)
          return (
            <text
              key={d.key}
              x={labelPoint.x}
              y={labelPoint.y}
              textAnchor="middle"
              dominantBaseline="central"
              fill="#94a3b8"
              fontSize="12"
            >
              {d.label}
            </text>
          )
        })}

        {/* Score labels near data points */}
        {dims.map((d, i) => {
          const scorePoint = getPoint(i, d.value + 0.6)
          return (
            <text
              key={`score-${d.key}`}
              x={scorePoint.x}
              y={scorePoint.y}
              textAnchor="middle"
              dominantBaseline="central"
              fill="#f59e0b"
              fontSize="11"
              fontWeight="bold"
            >
              {d.value}
            </text>
          )
        })}
      </svg>
    </div>
  )
}

// ─── Stars ───────────────────────────────────────────────────────────────────

function Stars({ count }: { count: number }) {
  return (
    <span className="inline-flex gap-0.5">
      {[1, 2, 3, 4, 5].map((i) => (
        <span
          key={i}
          className={i <= count ? 'text-amber-400' : 'text-slate-600'}
        >
          ★
        </span>
      ))}
    </span>
  )
}

// ─── Dimension Card ──────────────────────────────────────────────────────────

function DimensionCard({
  icon,
  label,
  score,
  text,
}: {
  icon: string
  label: string
  score: number
  text: string
}) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="bg-card rounded-xl p-4 border border-slate-700/50"
    >
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2">
          <span className="text-lg">{icon}</span>
          <span className="text-white font-medium">{label}</span>
        </div>
        <Stars count={score} />
      </div>
      <p className="text-slate-400 text-sm leading-relaxed">{text}</p>
    </motion.div>
  )
}

// ─── Lucky Elements ──────────────────────────────────────────────────────────

function LuckyElements({
  number,
  color,
}: {
  number: number
  color: string
}) {
  return (
    <div className="flex justify-center gap-6">
      <motion.div
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        transition={{ type: 'spring', stiffness: 200, delay: 0.2 }}
        className="flex flex-col items-center"
      >
        <div className="w-16 h-16 rounded-full bg-gradient-to-br from-amber-500 to-orange-500 flex items-center justify-center text-white text-2xl font-bold shadow-lg shadow-amber-500/20">
          {number}
        </div>
        <span className="text-slate-400 text-xs mt-2">幸运数字</span>
      </motion.div>
      <motion.div
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        transition={{ type: 'spring', stiffness: 200, delay: 0.3 }}
        className="flex flex-col items-center"
      >
        <div className="w-16 h-16 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center text-white text-lg font-medium shadow-lg shadow-purple-500/20">
          {color}
        </div>
        <span className="text-slate-400 text-xs mt-2">幸运颜色</span>
      </motion.div>
    </div>
  )
}

// ─── Period Tabs ─────────────────────────────────────────────────────────────

const PERIODS: { key: 'daily' | 'weekly' | 'monthly'; label: string }[] = [
  { key: 'daily', label: '今日' },
  { key: 'weekly', label: '本周' },
  { key: 'monthly', label: '本月' },
]

function PeriodTabs({
  selected,
  onChange,
}: {
  selected: 'daily' | 'weekly' | 'monthly'
  onChange: (p: 'daily' | 'weekly' | 'monthly') => void
}) {
  return (
    <div className="flex justify-center gap-2">
      {PERIODS.map((p) => (
        <button
          key={p.key}
          onClick={() => onChange(p.key)}
          className={`px-5 py-2 rounded-lg text-sm font-medium transition-all ${
            selected === p.key
              ? 'bg-gradient-to-r from-amber-500 to-orange-500 text-white shadow-lg shadow-amber-500/20'
              : 'bg-slate-800 text-slate-400 hover:text-white hover:bg-slate-700'
          }`}
        >
          {p.label}
        </button>
      ))}
    </div>
  )
}

// ─── Result Card ─────────────────────────────────────────────────────────────

function ResultCard({ result }: { result: HoroscopeResult }) {
  const resultJson = result.record_id ? undefined : JSON.stringify(result)

  return (
    <motion.div
      key={`${result.zodiac}-${result.period}-${result.date}`}
      initial={{ opacity: 0, y: 30 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -30 }}
      transition={{ duration: 0.4 }}
      className="space-y-6"
    >
      {/* Zodiac header */}
      <div className="text-center">
        <div className="text-5xl mb-2">
          {ZODIAC_SIGNS.find((z) => z.name === result.zodiac)?.symbol}
        </div>
        <h2 className="text-2xl font-bold text-white">{result.zodiac_cn}</h2>
        <p className="text-slate-400 text-sm mt-1">{result.date}</p>
      </div>

      {/* Summary */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.1 }}
        className="bg-gradient-to-r from-purple-500/10 to-pink-500/10 border border-purple-500/20 rounded-xl p-5 text-center"
      >
        <p className="text-slate-300 leading-relaxed">{result.summary}</p>
      </motion.div>

      {/* Radar chart */}
      <div className="bg-card rounded-xl p-4 border border-slate-700/50">
        <h3 className="text-amber-400 text-sm font-medium mb-2 text-center">
          📊 运势雷达图
        </h3>
        <RadarChart result={result} />
      </div>

      {/* Lucky elements */}
      <div className="bg-card rounded-xl p-5 border border-slate-700/50">
        <h3 className="text-amber-400 text-sm font-medium mb-4 text-center">
          🍀 幸运元素
        </h3>
        <LuckyElements number={result.lucky_number} color={result.lucky_color} />
      </div>

      {/* Dimension details */}
      <div className="space-y-3">
        <DimensionCard icon="💕" label="爱情" score={result.love} text={result.detail.love} />
        <DimensionCard icon="💼" label="事业" score={result.career} text={result.detail.career} />
        <DimensionCard icon="💰" label="财运" score={result.wealth} text={result.detail.wealth} />
        <DimensionCard icon="🏥" label="健康" score={result.health} text={result.detail.health} />
      </div>

      {/* AI Reading */}
      <AIReading
        type="horoscope"
        resultId={result.record_id}
        result={resultJson}
      />
    </motion.div>
  )
}

// ─── Main Page ───────────────────────────────────────────────────────────────

export default function Horoscope() {
  const {
    selectedZodiac,
    period,
    result,
    isLoading,
    error,
    selectZodiac,
    setPeriod,
  } = useHoroscope()

  return (
    <div className="flex-1 px-4 py-12 max-w-2xl mx-auto w-full">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center mb-8"
      >
        <h1 className="text-4xl font-bold mb-3">
          <span className="bg-gradient-to-r from-amber-400 via-orange-400 to-pink-400 bg-clip-text text-transparent">
            星座运势
          </span>
        </h1>
        <p className="text-slate-400">选择你的星座，探索星象为你带来的指引</p>
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

      {/* Zodiac Wheel */}
      <motion.div
        initial={{ opacity: 0, scale: 0.9 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.5 }}
        className="mb-6"
      >
        <ZodiacWheel selected={selectedZodiac} onSelect={selectZodiac} />
      </motion.div>

      {/* Period Tabs (only shown after zodiac selected) */}
      <AnimatePresence>
        {selectedZodiac && (
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0 }}
            className="mb-6"
          >
            <PeriodTabs selected={period} onChange={setPeriod} />
          </motion.div>
        )}
      </AnimatePresence>

      {/* Loading */}
      <AnimatePresence>
        {isLoading && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="flex justify-center py-12"
          >
            <div className="flex items-center gap-3 text-slate-400">
              <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
              </svg>
              <span>解读星象中...</span>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Result */}
      <AnimatePresence mode="wait">
        {!isLoading && result && <ResultCard result={result} />}
      </AnimatePresence>

      {/* Empty state */}
      {!selectedZodiac && !isLoading && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="text-center py-8"
        >
          <p className="text-slate-500">👆 点击上方星座图标查看运势</p>
        </motion.div>
      )}
    </div>
  )
}
