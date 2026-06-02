import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { useAuthStore } from '../stores/authStore'

const divinationTypes = [
  {
    path: '/liuyao',
    icon: '☰',
    title: '周易六爻',
    description: '传承千年的易经智慧，通过六爻卦象预测事物的发展趋势。',
    gradient: 'from-emerald-500 to-teal-500',
  },
  {
    path: '/bazi',
    icon: '📅',
    title: '八字排盘',
    description: '根据出生时辰推算八字，解读命运密码，规划人生蓝图。',
    gradient: 'from-blue-500 to-cyan-500',
  },
  {
    path: '/tarot',
    icon: '🔮',
    title: '塔罗牌占卜',
    description: '通过塔罗牌的神秘图案，探索内心深处的答案，指引人生方向。',
    gradient: 'from-purple-500 to-pink-500',
  },
]

export default function Home() {
  const { isAuthenticated, user } = useAuthStore()

  return (
    <div className="flex-1">
      {/* Hero Section */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-b from-primary/5 to-transparent"></div>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 relative">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="text-center"
          >
            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold mb-6">
              <span className="bg-gradient-to-r from-primary via-secondary to-accent bg-clip-text text-transparent">
                探索命运的奥秘
              </span>
            </h1>
            <p className="text-lg sm:text-xl text-slate-400 max-w-2xl mx-auto mb-8">
              {isAuthenticated
                ? `欢迎回来，${user?.username}！选择你感兴趣的占卜方式，开启探索之旅。`
                : '融合东西方神秘学智慧，为你提供专业的在线占卜服务。登录后即可开始体验。'}
            </p>
            {!isAuthenticated && (
              <div className="flex justify-center gap-4">
                <Link
                  to="/register"
                  className="px-8 py-3 bg-primary hover:bg-primary-dark text-white rounded-lg font-medium transition-colors"
                >
                  免费注册
                </Link>
                <Link
                  to="/login"
                  className="px-8 py-3 border border-slate-600 hover:border-slate-500 text-slate-300 hover:text-white rounded-lg font-medium transition-colors"
                >
                  立即登录
                </Link>
              </div>
            )}
          </motion.div>
        </div>
      </section>

      {/* Divination Types */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <motion.h2
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
          className="text-2xl sm:text-3xl font-bold text-center text-white mb-12"
        >
          选择你的占卜方式
        </motion.h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {divinationTypes.map((type, index) => (
            <motion.div
              key={type.path}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
            >
              <Link
                to={type.path}
                className="block group h-full"
              >
                <div className="h-full bg-card hover:bg-card-hover rounded-2xl p-6 border border-slate-700/50 hover:border-primary/30 transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-primary/5">
                  <div className={`w-16 h-16 rounded-2xl bg-gradient-to-br ${type.gradient} flex items-center justify-center text-3xl mb-4 group-hover:scale-110 transition-transform`}>
                    {type.icon}
                  </div>
                  <h3 className="text-xl font-bold text-white mb-2">{type.title}</h3>
                  <p className="text-slate-400 text-sm leading-relaxed">{type.description}</p>
                </div>
              </Link>
            </motion.div>
          ))}
        </div>
      </section>
    </div>
  )
}
