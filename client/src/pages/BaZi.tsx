import { useState, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { DayPicker } from 'react-day-picker'
import { zhCN } from 'react-day-picker/locale'
import dayjs from 'dayjs'
import { calculateBaZi } from '../services/bazi'
import { useAI } from '../hooks/useAI'
import AIReading from '../components/common/AIReading'
import type { BaZiResult, BaZiPillar, FiveElementAnalysis } from '../types/bazi'
import 'react-day-picker/style.css'

const HOURS = Array.from({ length: 24 }, (_, i) => String(i).padStart(2, '0'))
const MINUTES = Array.from({ length: 60 }, (_, i) => String(i).padStart(2, '0'))

function getCurrentBirthDefaults() {
  const now = dayjs()
  return {
    date: now.format('YYYY-MM-DD'),
    hour: now.format('HH'),
    minute: now.format('mm'),
  }
}

function PillarCard({ pillar, label }: { pillar: BaZiPillar; label: string }) {
  return (
    <div className="bg-slate-800/50 border border-blue-500/20 rounded-xl p-4 text-center">
      <p className="text-blue-400 text-xs font-medium mb-2">{label}</p>
      <div className="flex flex-col items-center gap-1">
        <span className="text-3xl font-bold text-white">{pillar.tian_gan}</span>
        <span className="text-3xl font-bold text-amber-300">{pillar.di_zhi}</span>
      </div>
      <div className="mt-3 space-y-1">
        <p className="text-slate-400 text-xs">五行：{pillar.wu_xing}</p>
        <p className="text-slate-500 text-xs">纳音：{pillar.na_yin}</p>
        {pillar.hidden_gan && pillar.hidden_gan.length > 0 && (
          <p className="text-slate-500 text-xs">藏干：{pillar.hidden_gan.join(' ')}</p>
        )}
      </div>
    </div>
  )
}

function FiveElementBar({ analysis }: { analysis: FiveElementAnalysis }) {
  const elements = [
    { name: '金', count: analysis.metal, color: 'bg-yellow-400', label: '金' },
    { name: '木', count: analysis.wood, color: 'bg-green-400', label: '木' },
    { name: '水', count: analysis.water, color: 'bg-blue-400', label: '水' },
    { name: '火', count: analysis.fire, color: 'bg-red-400', label: '火' },
    { name: '土', count: analysis.earth, color: 'bg-amber-600', label: '土' },
  ]

  const maxCount = Math.max(...elements.map(e => e.count), 1)

  return (
    <div className="bg-slate-800/30 border border-slate-700/50 rounded-xl p-5">
      <h3 className="text-blue-400 font-semibold mb-4">五行分布</h3>
      <div className="space-y-3">
        {elements.map((el) => (
          <div key={el.name} className="flex items-center gap-3">
            <span className="text-sm text-slate-300 w-6">{el.label}</span>
            <div className="flex-1 h-4 bg-slate-700/50 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                animate={{ width: `${(el.count / maxCount) * 100}%` }}
                transition={{ duration: 0.5, delay: 0.1 }}
                className={`h-full ${el.color} rounded-full`}
              />
            </div>
            <span className="text-slate-400 text-sm w-6 text-right">{el.count}</span>
          </div>
        ))}
      </div>
      <div className="mt-4 grid grid-cols-2 gap-3 text-sm">
        <div>
          <span className="text-slate-500">日主：</span>
          <span className="text-white">{analysis.day_master}</span>
        </div>
        <div>
          <span className="text-slate-500">强弱：</span>
          <span className="text-white">{analysis.strength}</span>
        </div>
        <div>
          <span className="text-slate-500">用神：</span>
          <span className="text-emerald-400">{analysis.yong_shen}</span>
        </div>
        <div>
          <span className="text-slate-500">忌神：</span>
          <span className="text-red-400">{analysis.ji_shen}</span>
        </div>
      </div>
    </div>
  )
}

export default function BaZi() {
  const [birthInputs, setBirthInputs] = useState(getCurrentBirthDefaults)
  const [showDatePicker, setShowDatePicker] = useState(false)
  const [gender, setGender] = useState('')
  const [result, setResult] = useState<BaZiResult | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const birthDate = birthInputs.date
  const birthHour = birthInputs.hour
  const birthMinute = birthInputs.minute
  const birthTime = birthHour && birthMinute ? `${birthHour}:${birthMinute}` : ''
  const selectedBirthDate = birthDate ? dayjs(birthDate).toDate() : undefined

  const ai = useAI({
    type: 'bazi',
    resultId: result?.record_id,
    result: result && !result.record_id ? JSON.stringify(result) : undefined,
    question: result ? `八字排盘 ${result.birth.solar}` : undefined,
  })

  const handleBirthDateSelect = useCallback((date: Date | undefined) => {
    if (!date) return
    setBirthInputs((inputs) => ({
      ...inputs,
      date: dayjs(date).format('YYYY-MM-DD'),
    }))
    setShowDatePicker(false)
    setError(null)
  }, [])

  const handleCalculate = useCallback(async () => {
    if (!birthDate) {
      setError('请选择出生日期')
      return
    }
    if (!birthTime) {
      setError('请选择出生时间')
      return
    }

    setIsLoading(true)
    setError(null)

    try {
      const data = await calculateBaZi(birthDate, birthTime, gender)
      setResult(data)
    } catch (err) {
      setError('排盘失败，请检查输入信息')
      console.error('BaZi calculate error:', err)
    } finally {
      setIsLoading(false)
    }
  }, [birthDate, birthTime, gender])

  const handleReset = useCallback(() => {
    setResult(null)
    setBirthInputs(getCurrentBirthDefaults())
    setShowDatePicker(false)
    setGender('')
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
          <span className="bg-gradient-to-r from-blue-400 via-cyan-400 to-teal-400 bg-clip-text text-transparent">
            八字排盘
          </span>
        </h1>
        <p className="text-slate-400">输入出生时间，解读命运密码</p>
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
            className="space-y-6"
          >
            {/* Birth date */}
            <div>
              <label className="block text-blue-400 text-sm font-medium mb-2">
                📅 出生日期
              </label>
              <button
                type="button"
                onClick={() => setShowDatePicker((visible) => !visible)}
                aria-expanded={showDatePicker}
                aria-controls="birth-date-picker"
                className="w-full px-5 py-4 bg-slate-800/80 border border-slate-700 rounded-xl text-white focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all flex items-center justify-between text-left"
              >
                <span>{birthDate ? dayjs(birthDate).format('YYYY年M月D日') : '请选择出生日期'}</span>
                <span className="text-slate-400" aria-hidden="true">📅</span>
              </button>
              <AnimatePresence>
                {showDatePicker && (
                  <motion.div
                    id="birth-date-picker"
                    initial={{ opacity: 0, y: -8 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -8 }}
                    className="mt-3 rounded-xl border border-slate-700 bg-slate-800/95 p-4 shadow-xl shadow-black/20 backdrop-blur-sm"
                  >
                    <DayPicker
                      mode="single"
                      selected={selectedBirthDate}
                      onSelect={handleBirthDateSelect}
                      locale={zhCN}
                      captionLayout="dropdown"
                      reverseYears
                      startMonth={new Date(1900, 0)}
                      endMonth={new Date()}
                      defaultMonth={selectedBirthDate ?? new Date(1990, 0)}
                      disabled={{ after: new Date() }}
                      className="text-slate-300"
                      classNames={{
                        root: 'w-full',
                        months: 'flex justify-center',
                        month: 'flex flex-col gap-3',
                        month_caption: 'flex justify-center',
                        caption_label: 'text-sm font-medium text-slate-300',
                        dropdowns: 'flex items-center justify-center gap-2',
                        dropdown: 'rounded-lg border border-slate-600 bg-slate-900 px-2 py-1 text-sm text-slate-100 outline-none focus:border-blue-500',
                        nav: 'flex items-center justify-between gap-2',
                        button_previous: 'rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-700 hover:text-white',
                        button_next: 'rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-700 hover:text-white',
                        weekdays: 'grid grid-cols-7',
                        weekday: 'py-2 text-center text-xs font-medium text-slate-500',
                        week: 'grid grid-cols-7',
                        day: 'p-0.5 text-center',
                        day_button: 'h-9 w-9 rounded-lg text-sm transition-all hover:bg-slate-700/70 focus:outline-none focus:ring-1 focus:ring-blue-500',
                        selected: 'bg-blue-500/20 text-blue-300',
                        today: 'font-bold text-cyan-300',
                        disabled: 'text-slate-600 opacity-50',
                        outside: 'text-slate-600',
                      }}
                    />
                  </motion.div>
                )}
              </AnimatePresence>
            </div>

            {/* Birth time */}
            <div>
              <label className="block text-blue-400 text-sm font-medium mb-2">
                🕐 出生时间
              </label>
              <div className="grid grid-cols-[1fr_auto_1fr] items-center gap-3">
                <select
                  id="birth-hour-select"
                  value={birthHour}
                  onChange={(e) => {
                    setBirthInputs((inputs) => ({
                      ...inputs,
                      hour: e.target.value,
                    }))
                    setError(null)
                  }}
                  className="w-full appearance-none rounded-xl border border-slate-700 bg-slate-800/80 px-5 py-4 text-white focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all"
                >
                  <option value="">时</option>
                  {HOURS.map((hour) => (
                    <option key={hour} value={hour}>
                      {hour} 时
                    </option>
                  ))}
                </select>
                <span className="text-slate-500">:</span>
                <select
                  id="birth-minute-select"
                  value={birthMinute}
                  onChange={(e) => {
                    setBirthInputs((inputs) => ({
                      ...inputs,
                      minute: e.target.value,
                    }))
                    setError(null)
                  }}
                  className="w-full appearance-none rounded-xl border border-slate-700 bg-slate-800/80 px-5 py-4 text-white focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all"
                >
                  <option value="">分</option>
                  {MINUTES.map((minute) => (
                    <option key={minute} value={minute}>
                      {minute} 分
                    </option>
                  ))}
                </select>
              </div>
            </div>

            {/* Gender */}
            <div>
              <label className="block text-blue-400 text-sm font-medium mb-2">
                ⚧ 性别（可选）
              </label>
              <div className="flex gap-3">
                {[
                  { value: 'male', label: '男', icon: '♂' },
                  { value: 'female', label: '女', icon: '♀' },
                ].map((g) => (
                  <button
                    key={g.value}
                    onClick={() => setGender(g.value)}
                    className={`flex-1 py-3 rounded-xl border-2 transition-all text-sm font-medium ${
                      gender === g.value
                        ? 'border-blue-500 bg-blue-500/10 text-blue-400'
                        : 'border-slate-700 bg-slate-800/50 text-slate-400 hover:border-slate-600'
                    }`}
                  >
                    {g.icon} {g.label}
                  </button>
                ))}
              </div>
            </div>

            {/* Calculate button */}
            <div className="text-center pt-4">
              <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                onClick={handleCalculate}
                disabled={isLoading}
                className="px-10 py-4 bg-gradient-to-r from-blue-600 to-cyan-600 hover:from-blue-700 hover:to-cyan-700 text-white rounded-xl font-bold text-lg shadow-lg shadow-blue-500/20 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
              >
                {isLoading ? (
                  <span className="flex items-center gap-2">
                    <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                    </svg>
                    排盘中...
                  </span>
                ) : (
                  '📜 开始排盘'
                )}
              </motion.button>
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
            {/* Birth info */}
            <div className="text-center">
              <p className="text-slate-400 text-sm">出生信息</p>
              <p className="text-white text-lg mt-1">{result.birth.solar}</p>
              <p className="text-slate-500 text-sm mt-1">{result.birth.lunar}</p>
            </div>

            {/* Four Pillars */}
            <div>
              <h3 className="text-blue-400 font-semibold mb-4 text-center">四柱八字</h3>
              <div className="grid grid-cols-4 gap-3">
                <PillarCard pillar={result.pillars.year} label="年柱" />
                <PillarCard pillar={result.pillars.month} label="月柱" />
                <PillarCard pillar={result.pillars.day} label="日柱" />
                <PillarCard pillar={result.pillars.hour} label="时柱" />
              </div>
            </div>

            {/* Five Elements */}
            <FiveElementBar analysis={result.five_elements} />

            {/* Ten Gods */}
            {result.ten_gods.length > 0 && (
              <div className="bg-slate-800/30 border border-slate-700/50 rounded-xl p-5">
                <h3 className="text-blue-400 font-semibold mb-4">十神</h3>
                <div className="grid grid-cols-2 gap-2">
                  {result.ten_gods.map((god, i) => (
                    <div key={i} className="flex items-center gap-2 text-sm">
                      <span className="text-slate-500 w-16">{god.position}</span>
                      <span className="text-white">{god.tian_gan}</span>
                      <span className="text-amber-400">({god.god})</span>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* AI Reading */}
            {result && (
              <AIReading
                type="bazi"
                resultId={result.record_id}
                result={!result.record_id ? JSON.stringify(result) : undefined}
                question={`八字排盘 ${result.birth.solar}`}
              />
            )}

            {/* Reset button */}
            <div className="text-center pt-4">
              <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                onClick={handleReset}
                className="px-8 py-3 bg-gradient-to-r from-blue-600 to-cyan-600 hover:from-blue-700 hover:to-cyan-700 text-white rounded-xl font-medium transition-all"
              >
                🔄 重新排盘
              </motion.button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}
