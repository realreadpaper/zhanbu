export type ChatDivinationType = 'tarot' | 'liuyao' | 'liuyao_v2' | 'bazi' | 'horoscope' | 'meihua'

export interface DivinationPersona {
  type: ChatDivinationType
  icon: string
  name: string
  title: string
  subtitle: string
  welcomeTitle: string
  welcomeDescription: string
  welcomeHint: string
  ritualTitle: string
  ritualSubtitle: string
  ritualSteps: string[]
  accentClass: string
  softClass: string
}

const PERSONAS: Record<ChatDivinationType, DivinationPersona> = {
  tarot: {
    type: 'tarot',
    icon: '🎴',
    name: '星牌解语师',
    title: '星牌解语师',
    subtitle: '从牌面象征中解读你的当下与可能性',
    welcomeTitle: '塔罗牌占卜',
    welcomeDescription: '以牌面、正逆位与象征关系观察问题的核心脉络。',
    welcomeHint: '输入你的问题，我会先抽牌，再为你解读牌意。',
    ritualTitle: '正在洗牌与抽取牌面',
    ritualSubtitle: '让问题沉入牌阵，等待象征浮现。',
    ritualSteps: ['凝神洗牌', '抽取牌面', '翻开象征', '整理解读线索'],
    accentClass: 'from-fuchsia-500 via-violet-500 to-indigo-400',
    softClass: 'border-fuchsia-400/30 bg-fuchsia-500/10 text-fuchsia-200',
  },
  liuyao: {
    type: 'liuyao',
    icon: '☯️',
    name: '玄机卦师',
    title: '玄机卦师',
    subtitle: '以卦象、动爻与变卦推演事情走势',
    welcomeTitle: '六爻占卜',
    welcomeDescription: '以高岛易断与六爻卦象为据，观察事情的趋向与应对。',
    welcomeHint: '输入你想问的事情，我会先起卦，再为你推演。',
    ritualTitle: '正在凝神起卦',
    ritualSubtitle: '六爻自下而上成象，动静之间显露变化。',
    ritualSteps: ['凝神定问', '起卦成爻', '审视动爻', '推演变卦'],
    accentClass: 'from-amber-400 via-orange-500 to-rose-500',
    softClass: 'border-amber-400/30 bg-amber-500/10 text-amber-200',
  },
  liuyao_v2: {
    type: 'liuyao_v2',
    icon: '☯️',
    name: '玄机卦师',
    title: '玄机卦师',
    subtitle: '以高岛易断为据，为你推演卦象变化',
    welcomeTitle: '高岛易断',
    welcomeDescription: '结合本卦、变卦、动爻与高岛原文，判断问题的机势。',
    welcomeHint: '输入你想问的事情，我会先起卦，再结合高岛易断解读。',
    ritualTitle: '正在推演高岛卦象',
    ritualSubtitle: '卦象、动爻与原典线索正在汇合。',
    ritualSteps: ['定问起卦', '生成六爻', '检索高岛原文', '归纳判断依据'],
    accentClass: 'from-amber-400 via-orange-500 to-rose-500',
    softClass: 'border-amber-400/30 bg-amber-500/10 text-amber-200',
  },
  bazi: {
    type: 'bazi',
    icon: '📋',
    name: '四柱明鉴师',
    title: '四柱明鉴师',
    subtitle: '从四柱五行中观察命理结构与阶段趋势',
    welcomeTitle: '八字排盘',
    welcomeDescription: '以出生年月日时排出四柱，分析五行、日主与命理结构。',
    welcomeHint: '请输入出生年月日时，例如：1990-05-12 08:30 女，看看事业。',
    ritualTitle: '正在排布四柱命盘',
    ritualSubtitle: '年月日时依次落位，五行气势逐步展开。',
    ritualSteps: ['校验生辰', '四柱落位', '分析五行', '整理命理重点'],
    accentClass: 'from-cyan-400 via-sky-500 to-blue-500',
    softClass: 'border-cyan-400/30 bg-cyan-500/10 text-cyan-200',
  },
  horoscope: {
    type: 'horoscope',
    icon: '⭐',
    name: '星轨知微师',
    title: '星轨知微师',
    subtitle: '沿星象轨迹观察近期状态与能量变化',
    welcomeTitle: '星座运势',
    welcomeDescription: '以星座能量观察近期状态，覆盖关系、事业、财务与节奏。',
    welcomeHint: '请告诉我你的星座，例如：天蝎座今天事业运如何。',
    ritualTitle: '正在观测星轨能量',
    ritualSubtitle: '星点连线，近期状态正在浮现。',
    ritualSteps: ['定位星座', '点亮星轨', '读取能量', '生成趋势提示'],
    accentClass: 'from-yellow-300 via-emerald-300 to-sky-400',
    softClass: 'border-emerald-300/30 bg-emerald-400/10 text-emerald-100',
  },
  meihua: {
    type: 'meihua',
    icon: '🌸',
    name: '观梅心易师',
    title: '观梅心易师',
    subtitle: '以数起卦，观本卦、互卦、变卦与体用生克',
    welcomeTitle: '梅花易数',
    welcomeDescription: '以邵雍一脉的时间起卦与数字起卦为核心，结合本卦、互卦、变卦与体用生克进行推演。',
    welcomeHint: '输入你想问的事，我会以当前时间起卦。也可输入：数字 12 34，问事业机会。',
    ritualTitle: '正在以数起卦',
    ritualSubtitle: '梅花瓣落，卦象渐成。',
    ritualSteps: ['取数成象', '分出上下卦', '定动爻', '推演体用生克'],
    accentClass: 'from-pink-400 via-rose-500 to-red-400',
    softClass: 'border-pink-400/30 bg-pink-500/10 text-pink-200',
  },
}

export function normalizeChatType(type?: string): ChatDivinationType {
  if (type === 'liuyao') return 'liuyao_v2'
  if (type === 'tarot' || type === 'liuyao_v2' || type === 'bazi' || type === 'horoscope' || type === 'meihua') {
    return type
  }
  return 'tarot'
}

export function getDivinationPersona(type?: string): DivinationPersona {
  return PERSONAS[normalizeChatType(type)]
}
