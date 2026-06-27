import { motion } from 'framer-motion'
import type { DivinationRecord } from '../../services/chat'
import { getDivinationPersona } from './divinationPersona'

interface DivinationResultCardProps {
  record: DivinationRecord
}

export default function DivinationResultCard({ record }: DivinationResultCardProps) {
  const persona = getDivinationPersona(record.type)
  const profileName = record.prompt_profile_name || persona.name
  const summary = buildSummary(record)

  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      className="flex gap-2 sm:gap-3"
    >
      <div className={`flex h-8 w-8 shrink-0 items-center justify-center rounded-lg border text-sm sm:h-9 sm:w-9 ${persona.softClass}`}>
        {persona.icon}
      </div>
      <div className="min-w-0 flex-1">
        <div className="mb-1 px-1 text-xs text-slate-500">占卜结果</div>
        <div className="rounded-2xl rounded-tl-md border border-violet-500/30 bg-violet-500/10 px-4 py-3">
          <div className="mb-2 flex items-center justify-between gap-3">
            <div>
              <div className="text-sm font-semibold text-slate-100">{summary.title}</div>
              <div className="mt-0.5 text-xs text-slate-500">{persona.subtitle}</div>
            </div>
            <span className="rounded-full border border-violet-400/30 px-2 py-1 text-[11px] text-violet-200">
              {summary.badge}
            </span>
          </div>
          <div className="grid gap-2 sm:grid-cols-2">
            {summary.items.map((item) => (
              <div key={item.label} className="rounded-lg border border-slate-700/50 bg-slate-900/45 px-3 py-2">
                <div className="text-[11px] text-slate-500">{item.label}</div>
                <div className="mt-1 text-sm text-slate-200">{item.value}</div>
              </div>
            ))}
          </div>
          <div className="mt-3 text-xs text-slate-500">下面由 {profileName} 基于此结果继续解读。</div>
        </div>
      </div>
    </motion.div>
  )
}

function buildSummary(record: DivinationRecord) {
  const data = safeJSON(record.result)

  if (record.type === 'liuyao' || record.type === 'liuyao_v2') {
    const benGua = data?.ben_gua
    const bianGua = data?.bian_gua
    const mutableLines = Array.isArray(data?.mutable_lines) ? data.mutable_lines : []
    const lines = Array.isArray(data?.lines) ? data.lines : []
    const benGuaMap = asMap(benGua)
    return {
      title: asString(benGuaMap?.full_name) || asString(benGuaMap?.name) || '卦象已成',
      badge: record.type === 'liuyao_v2' ? '高岛易断' : '六爻',
      items: [
        { label: '本卦', value: formatHexagram(benGua) },
        { label: '变卦', value: bianGua ? formatHexagram(bianGua) : '无变卦' },
        { label: '动爻', value: mutableLines.length ? mutableLines.map((line: number) => `第${line + 1}爻`).join('、') : '无动爻' },
        { label: '六爻', value: lines.map((line: { symbol?: string }) => line.symbol).filter(Boolean).join(' ') || '已生成' },
      ],
    }
  }

  if (record.type === 'tarot') {
    const cards = Array.isArray(data?.cards) ? data.cards : []
    const first = asMap(cards[0])
    const firstCard = asMap(first?.card)
    return {
      title: asString(firstCard?.name) || '牌面已抽取',
      badge: '塔罗牌',
      items: [
        { label: '牌阵', value: asString(data?.spread) || '单牌' },
        { label: '牌面', value: asString(firstCard?.name) || '已抽牌' },
        { label: '位置', value: asString(first?.position_name) || '核心信息' },
        { label: '方向', value: first?.orientation === 'reversed' ? '逆位' : '正位' },
      ],
    }
  }

  if (record.type === 'bazi') {
    const birth = asMap(data?.birth)
    const pillars = asMap(data?.pillars)
    return {
      title: asString(birth?.solar) || '四柱已排定',
      badge: '八字',
      items: [
        { label: '年柱', value: formatPillar(pillars?.year) },
        { label: '月柱', value: formatPillar(pillars?.month) },
        { label: '日柱', value: formatPillar(pillars?.day) },
        { label: '时柱', value: formatPillar(pillars?.hour) },
      ],
    }
  }

  if (record.type === 'horoscope') {
    return {
      title: asString(data?.zodiac_cn) || '星座运势已生成',
      badge: '星座',
      items: [
        { label: '综合', value: score(data?.overall) },
        { label: '爱情', value: score(data?.love) },
        { label: '事业', value: score(data?.career) },
        { label: '财运', value: score(data?.wealth) },
      ],
    }
  }

  if (record.type === 'meihua') {
    const benGua = asMap(data?.ben_gua)
    const huGua = asMap(data?.hu_gua)
    const bianGua = asMap(data?.bian_gua)
    const tiYong = asMap(data?.ti_yong)
    const ti = asMap(tiYong?.ti)
    const yong = asMap(tiYong?.yong)
    const movingLine = data?.moving_line
    const method = asString(data?.method)
    const sourceValues = asMap(data?.source_values)

    return {
      title: asString(benGua?.name) || '卦象已成',
      badge: '梅花易数',
      items: [
        { label: '起卦方式', value: formatMeihuaMethod(method, sourceValues) },
        { label: '本卦', value: formatMeihuaHexagram(benGua) },
        { label: '互卦', value: formatMeihuaHexagram(huGua) },
        { label: '变卦', value: formatMeihuaHexagram(bianGua) },
        {
          label: '体用',
          value: tiYong && ti && yong
            ? `${asString(ti.name)}为体 · ${asString(yong.name)}为用 · ${asString(tiYong.relation)}`
            : '待分析',
        },
        {
          label: '动爻',
          value: typeof movingLine === 'number' ? formatMovingLineName(movingLine) : '无动爻',
        },
      ],
    }
  }

  return {
    title: '占卜结果已生成',
    badge: '占卜',
    items: [{ label: '问题', value: record.question || '已记录' }],
  }
}

type ResultMap = Record<string, unknown>

function safeJSON(value: string): ResultMap | null {
  try {
    return JSON.parse(value)
  } catch {
    return null
  }
}

function asMap(value: unknown): ResultMap | null {
  return value && typeof value === 'object' ? value as ResultMap : null
}

function asString(value: unknown) {
  return typeof value === 'string' ? value : ''
}

function formatHexagram(hexagram: unknown) {
  const h = asMap(hexagram)
  if (!h) return '无'
  const fullName = asString(h.full_name)
  const nameShort = asString(h.name_short)
  const name = asString(h.name)
  if (fullName) return fullName
  if (nameShort && name) return `${nameShort}（${name}）`
  return name || '已生成'
}

function formatPillar(pillar: unknown) {
  const p = asMap(pillar)
  if (!p) return '待分析'
  const tianGan = asString(p.tian_gan)
  const diZhi = asString(p.di_zhi)
  const wuXing = asString(p.wu_xing)
  return `${tianGan}${diZhi}${wuXing ? ` · ${wuXing}` : ''}`
}

function score(value: unknown) {
  return typeof value === 'number' ? `${value}/5` : '待分析'
}

function formatMeihuaHexagram(hex: ResultMap | null) {
  if (!hex) return '无'
  const name = asString(hex.name)
  const upper = asMap(hex.upper)
  const lower = asMap(hex.lower)
  if (name && upper && lower) {
    return `${name}（${asString(upper.name)}上${asString(lower.name)}下）`
  }
  return name || '已生成'
}

function formatMovingLineName(line: number) {
  const names = ['', '初爻', '二爻', '三爻', '四爻', '五爻', '上爻']
  return (line >= 1 && line <= 6 ? names[line] : `第${line}爻`) + '动'
}

function formatMeihuaMethod(method: string, sourceValues: ResultMap | null) {
  if (method === 'number') {
    const numbers = Array.isArray(sourceValues?.numbers)
      ? sourceValues.numbers.join(' ')
      : ''
    return numbers ? `数字起卦（${numbers}）` : '数字起卦'
  }
  // 时间起卦：显示农历信息
  if (sourceValues) {
    const lunarDisplay = asString(sourceValues.lunar_display)
    if (lunarDisplay) return lunarDisplay

    const yearBranch = asString(sourceValues.year_branch)
    const month = sourceValues.lunar_month
    const day = sourceValues.lunar_day
    const hourBranch = asString(sourceValues.hour_branch)
    if (yearBranch && month && day && hourBranch) {
      return `${yearBranch}年${month}月${day}日${hourBranch}时`
    }
  }
  return '时间起卦'
}
