import { Link, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'

const navItems = [
  { path: '/tarot', label: '塔罗占卜' },
  { path: '/horoscope', label: '星座运势' },
  { path: '/liuyao', label: '六爻占卜' },
  { path: '/bazi', label: '八字命理' },
]

export default function Header() {
  const { isAuthenticated, user, logout } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <header className="bg-darker/80 backdrop-blur-md border-b border-slate-800 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link to="/" className="flex items-center gap-2 shrink-0">
            <span className="text-3xl">🔮</span>
            <span className="text-xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
              玄机占卜
            </span>
          </Link>

          {/* Navigation */}
          <nav className="hidden md:flex items-center gap-1">
            {navItems.map((item) => (
              <Link
                key={item.path}
                to={item.path}
                className="px-4 py-2 text-sm text-slate-300 hover:text-white rounded-lg hover:bg-white/5 transition-colors"
              >
                {item.label}
              </Link>
            ))}
          </nav>

          {/* Right side */}
          <div className="flex items-center gap-3">
            {isAuthenticated ? (
              <div className="flex items-center gap-3">
                <Link
                  to="/history"
                  className="text-sm text-slate-400 hover:text-white transition-colors"
                >
                  历史记录
                </Link>
                <Link
                  to="/profile"
                  className="flex items-center gap-2 px-3 py-1.5 rounded-lg hover:bg-white/5 transition-colors"
                >
                  <div className="w-8 h-8 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white text-sm font-bold">
                    {user?.username?.[0]?.toUpperCase() || '?'}
                  </div>
                  <span className="text-sm text-slate-300 hidden sm:inline">
                    {user?.username}
                  </span>
                </Link>
                <button
                  onClick={handleLogout}
                  className="text-sm text-slate-400 hover:text-red-400 transition-colors"
                >
                  退出
                </button>
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <Link
                  to="/login"
                  className="px-4 py-1.5 text-sm text-slate-300 hover:text-white transition-colors"
                >
                  登录
                </Link>
                <Link
                  to="/register"
                  className="px-4 py-1.5 text-sm bg-primary hover:bg-primary-dark text-white rounded-lg transition-colors"
                >
                  注册
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  )
}
