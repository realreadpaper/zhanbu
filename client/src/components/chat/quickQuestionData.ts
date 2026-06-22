export interface QuickQuestion {
  id: string
  icon: string
  label?: string
  text: string
}

const DEFAULT_QUESTIONS: Record<string, QuickQuestion[]> = {
  tarot: [
    { id: 'detail', icon: '🃏', label: '牌意展开', text: '请结合牌面象征、正逆位和我的问题，进一步展开这张牌的核心含义。' },
    { id: 'advice', icon: '💡', label: '行动建议', text: '基于这次塔罗结果，我接下来最适合采取什么行动？请给出具体建议。' },
    { id: 'warning', icon: '⚠️', label: '需要避开', text: '这次结果里有哪些需要警惕或避免的地方？请直接指出风险点。' },
    { id: 'improve', icon: '🔧', label: '如何改善', text: '如果我想让事情往更好的方向发展，现在可以从哪些方面改善？' },
  ],
  liuyao: [
    { id: 'detail', icon: '☯️', label: '卦象细解', text: '请结合本卦、卦辞和爻象，详细解释这卦对我所问事情的含义。' },
    { id: 'change', icon: '🔄', label: '变化趋势', text: '如果有动爻或变卦，请重点分析事情接下来会如何变化。' },
    { id: 'advice', icon: '💡', label: '取舍建议', text: '基于这个卦象，我现在应该主动推进、观望，还是调整方向？请说明理由。' },
    { id: 'time', icon: '⏰', label: '时间节点', text: '请分析这件事可能在哪些时间节点出现变化或结果。' },
  ],
  liuyao_v2: [
    { id: 'detail', icon: '☯️', label: '高岛细解', text: '请结合高岛易断原文、卦象和我的问题，深入解释这卦的判断依据。' },
    { id: 'change', icon: '🔄', label: '变动分析', text: '请重点分析动爻、变卦以及它们代表的变化趋势。' },
    { id: 'advice', icon: '💡', label: '应对策略', text: '基于高岛易断的判断，我现在应该如何应对？请给出可执行建议。' },
    { id: 'time', icon: '⏰', label: '应期判断', text: '请根据卦象分析这件事可能何时有明显进展、阻力或结果。' },
  ],
  bazi: [
    { id: 'detail', icon: '📋', label: '命盘总览', text: '请从五行、日主强弱和格局角度，给我做一个清晰的命盘总览。' },
    { id: 'career', icon: '💼', label: '事业方向', text: '结合我的八字，哪些事业方向更适合我？当前阶段应该注意什么？' },
    { id: 'wealth', icon: '💰', label: '财运节奏', text: '请分析我的财运特点、适合的求财方式，以及近期需要避开的风险。' },
    { id: 'relationship', icon: '💕', label: '关系婚恋', text: '请结合八字分析我的亲密关系/婚恋倾向，以及相处中的关键建议。' },
  ],
  horoscope: [
    { id: 'love', icon: '💕', label: '感情提醒', text: '请结合当前星座运势，具体分析我在感情关系里需要注意什么。' },
    { id: 'career', icon: '💼', label: '事业重点', text: '请分析我近期在工作和事业上的机会、阻力和行动重点。' },
    { id: 'wealth', icon: '💰', label: '财务建议', text: '请结合星座运势，给我近期消费、投资或财务安排上的建议。' },
    { id: 'energy', icon: '✨', label: '状态调整', text: '请分析我近期的能量状态，并给出适合的调整方式。' },
  ],
  default: [
    { id: 'detail', icon: '✨', label: '继续细讲', text: '请把刚才的解读进一步展开，重点讲清楚原因和判断依据。' },
    { id: 'advice', icon: '💡', label: '下一步', text: '基于当前结果，我接下来最应该做什么？请给出具体行动建议。' },
    { id: 'warning', icon: '⚠️', label: '风险点', text: '请指出当前结果中最需要注意的风险、误区或不宜做的事。' },
    { id: 'improve', icon: '🔧', label: '改善方式', text: '如果我想改善当前局面，可以从哪些具体方面入手？' },
  ],
}

export function getQuickQuestions(type: string): QuickQuestion[] {
  return DEFAULT_QUESTIONS[type] || DEFAULT_QUESTIONS.default
}
