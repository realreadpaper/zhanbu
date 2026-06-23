import { motion } from 'framer-motion'
import type { DivinationPersona } from './divinationPersona'

interface DivinationRitualProps {
  persona: DivinationPersona
}

export default function DivinationRitual({ persona }: DivinationRitualProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -8 }}
      className="mx-auto w-full max-w-xl rounded-xl border border-slate-700/60 bg-slate-900/80 px-5 py-5 shadow-2xl shadow-black/20"
    >
      <div className="flex items-start gap-4">
        <div className={`relative flex h-14 w-14 shrink-0 items-center justify-center rounded-xl border ${persona.softClass}`}>
          <motion.div
            className={`absolute inset-0 rounded-xl bg-gradient-to-br ${persona.accentClass} opacity-20 blur-md`}
            animate={{ opacity: [0.16, 0.34, 0.16], scale: [0.95, 1.08, 0.95] }}
            transition={{ duration: 2.4, repeat: Infinity, ease: 'easeInOut' }}
          />
          <span className="relative text-2xl">{persona.icon}</span>
        </div>

        <div className="min-w-0 flex-1">
          <div className="text-sm font-semibold text-slate-100">{persona.ritualTitle}</div>
          <div className="mt-1 text-xs leading-relaxed text-slate-500">{persona.ritualSubtitle}</div>

          <RitualVisual persona={persona} />

          <div className="mt-4 grid grid-cols-2 gap-2 sm:grid-cols-4">
            {persona.ritualSteps.map((step, index) => (
              <motion.div
                key={step}
                initial={{ opacity: 0.35 }}
                animate={{ opacity: [0.35, 1, 0.35] }}
                transition={{
                  duration: 1.8,
                  repeat: Infinity,
                  delay: index * 0.35,
                  ease: 'easeInOut',
                }}
                className="rounded-lg border border-slate-700/50 bg-slate-800/40 px-2.5 py-2 text-center text-[11px] text-slate-300"
              >
                {step}
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </motion.div>
  )
}

function RitualVisual({ persona }: DivinationRitualProps) {
  if (persona.type === 'tarot') return <TarotVisual persona={persona} />
  if (persona.type === 'bazi') return <BaziVisual persona={persona} />
  if (persona.type === 'horoscope') return <HoroscopeVisual persona={persona} />
  return <LiuYaoVisual persona={persona} />
}

function TarotVisual({ persona }: DivinationRitualProps) {
  return (
    <div className="mt-4 flex h-20 items-center justify-center gap-3">
      {[0, 1, 2].map((i) => (
        <motion.div
          key={i}
          className={`h-16 w-11 rounded-md border border-violet-300/30 bg-gradient-to-br ${persona.accentClass} p-px shadow-lg shadow-violet-950/40`}
          animate={{
            y: [0, i === 1 ? -10 : -4, 0],
            rotate: [-6 + i * 6, 2 - i * 2, -6 + i * 6],
          }}
          transition={{ duration: 1.8, repeat: Infinity, delay: i * 0.18, ease: 'easeInOut' }}
        >
          <div className="flex h-full w-full items-center justify-center rounded-[5px] bg-slate-950/80 text-sm text-violet-100">
            ✦
          </div>
        </motion.div>
      ))}
    </div>
  )
}

function LiuYaoVisual({ persona }: DivinationRitualProps) {
  return (
    <div className="mt-4 flex h-20 items-center justify-center">
      <div className="flex flex-col-reverse gap-1.5">
        {[0, 1, 2, 3, 4, 5].map((i) => (
          <motion.div
            key={i}
            className="flex items-center justify-center gap-1"
            initial={{ opacity: 0.2, width: 40 }}
            animate={{ opacity: [0.2, 1, 0.55], width: [42, 112, 112] }}
            transition={{ duration: 2.2, repeat: Infinity, delay: i * 0.18, ease: 'easeInOut' }}
          >
            <span className={`h-1.5 flex-1 rounded-full bg-gradient-to-r ${persona.accentClass}`} />
            {i % 3 === 0 && <span className="h-1.5 w-4 rounded-full bg-slate-900" />}
            <span className={`h-1.5 flex-1 rounded-full bg-gradient-to-r ${persona.accentClass}`} />
          </motion.div>
        ))}
      </div>
    </div>
  )
}

function BaziVisual({ persona }: DivinationRitualProps) {
  return (
    <div className="mt-4 grid h-20 grid-cols-4 gap-2">
      {['年', '月', '日', '时'].map((label, i) => (
        <motion.div
          key={label}
          className="flex flex-col items-center justify-center rounded-lg border border-cyan-300/20 bg-slate-800/50"
          animate={{ y: [8, 0, 0], opacity: [0.25, 1, 0.65] }}
          transition={{ duration: 1.9, repeat: Infinity, delay: i * 0.2, ease: 'easeInOut' }}
        >
          <div className={`mb-1 h-1.5 w-8 rounded-full bg-gradient-to-r ${persona.accentClass}`} />
          <span className="text-sm font-semibold text-cyan-100">{label}柱</span>
          <span className="mt-1 text-[10px] text-slate-500">五行</span>
        </motion.div>
      ))}
    </div>
  )
}

function HoroscopeVisual({ persona }: DivinationRitualProps) {
  return (
    <div className="mt-4 flex h-20 items-center justify-center">
      <motion.div
        className="relative h-20 w-40"
        animate={{ rotate: [0, 2, -2, 0] }}
        transition={{ duration: 4, repeat: Infinity, ease: 'easeInOut' }}
      >
        {[12, 36, 68, 92, 124].map((left, i) => (
          <motion.span
            key={left}
            className={`absolute h-2 w-2 rounded-full bg-gradient-to-br ${persona.accentClass}`}
            style={{ left, top: i % 2 === 0 ? 18 : 48 }}
            animate={{ scale: [0.8, 1.35, 0.8], opacity: [0.45, 1, 0.45] }}
            transition={{ duration: 1.8, repeat: Infinity, delay: i * 0.22 }}
          />
        ))}
        <svg className="absolute inset-0 h-full w-full" viewBox="0 0 160 80" aria-hidden="true">
          <motion.path
            d="M16 22 L40 52 L72 22 L96 52 L128 22"
            fill="none"
            stroke="rgba(125, 211, 252, 0.45)"
            strokeWidth="1"
            strokeDasharray="180"
            initial={{ strokeDashoffset: 180 }}
            animate={{ strokeDashoffset: [180, 0, 0] }}
            transition={{ duration: 2.4, repeat: Infinity, ease: 'easeInOut' }}
          />
        </svg>
      </motion.div>
    </div>
  )
}
