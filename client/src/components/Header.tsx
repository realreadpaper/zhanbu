import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { AnimatePresence, motion } from 'framer-motion'
import { useAuthStore } from '../stores/authStore'

export default function Header() {
  const { isAuthenticated, user, logout } = useAuthStore()
  const navigate = useNavigate()
  const [menuOpen, setMenuOpen] = useState(false)

  const handleLogout = () => {
    logout()
    setMenuOpen(false)
    navigate('/login')
  }

  const closeMenu = () => setMenuOpen(false)

  return (
    <>
      <header className="bg-darker/80 backdrop-blur-md border-b border-slate-800 sticky top-0 z-50 pt-[env(safe-area-inset-top)]">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <Link to="/" className="flex items-center gap-2 shrink-0" onClick={closeMenu}>
              <span className="text-3xl">🔮</span>
              <span className="text-xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                玄机占卜
              </span>
            </Link>

            {/* Desktop nav */}
            <div className="hidden sm:flex items-center gap-3">
              {isAuthenticated ? (
                <>
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
                    <span className="text-sm text-slate-300">
                      {user?.username}
                    </span>
                  </Link>
                  <button
                    onClick={handleLogout}
                    className="text-sm text-slate-400 hover:text-red-400 transition-colors"
                  >
                    退出
                  </button>
                </>
              ) : (
                <>
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
                </>
              )}
            </div>

            {/* Mobile hamburger */}
            <button
              onClick={() => setMenuOpen(!menuOpen)}
              className="sm:hidden w-10 h-10 flex items-center justify-center text-slate-400 hover:text-white transition-colors"
              aria-label="菜单"
            >
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                {menuOpen ? (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                ) : (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                )}
              </svg>
            </button>
          </div>
        </div>
      </header>

      {/* Mobile drawer overlay */}
      <AnimatePresence>
        {menuOpen && (
          <>
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={closeMenu}
              className="fixed inset-0 bg-black/50 z-40 sm:hidden"
            />
            <motion.div
              initial={{ x: '100%' }}
              animate={{ x: 0 }}
              exit={{ x: '100%' }}
              transition={{ type: 'spring', damping: 25, stiffness: 300 }}
              className="fixed top-0 right-0 bottom-0 w-64 bg-darker border-l border-slate-800 z-50 sm:hidden flex flex-col pt-[env(safe-area-inset-top)] pb-[env(safe-area-inset-bottom)]"
            >
              <div className="flex items-center justify-between p-4 border-b border-slate-800">
                <span className="text-lg font-semibold text-slate-200">菜单</span>
                <button
                  onClick={closeMenu}
                  className="w-8 h-8 flex items-center justify-center text-slate-400 hover:text-white"
                >
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              <nav className="flex-1 p-4 space-y-1">
                {isAuthenticated ? (
                  <>
                    <div className="flex items-center gap-3 px-3 py-3 mb-3">
                      <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white font-bold">
                        {user?.username?.[0]?.toUpperCase() || '?'}
                      </div>
                      <div>
                        <div className="text-sm font-medium text-slate-200">{user?.username}</div>
                        <div className="text-xs text-slate-500">{user?.email}</div>
                      </div>
                    </div>

                    <Link
                      to="/chat"
                      onClick={closeMenu}
                      className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-slate-300 hover:bg-slate-800/50 transition-colors"
                    >
                      <span>🔮</span>
                      <span>开始占卜</span>
                    </Link>
                    <Link
                      to="/history"
                      onClick={closeMenu}
                      className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-slate-300 hover:bg-slate-800/50 transition-colors"
                    >
                      <span>📋</span>
                      <span>历史记录</span>
                    </Link>
                    <Link
                      to="/profile"
                      onClick={closeMenu}
                      className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-slate-300 hover:bg-slate-800/50 transition-colors"
                    >
                      <span>👤</span>
                      <span>个人资料</span>
                    </Link>

                    <div className="border-t border-slate-800 my-3" />

                    <button
                      onClick={handleLogout}
                      className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-red-400 hover:bg-red-500/10 transition-colors w-full"
                    >
                      <span>🚪</span>
                      <span>退出登录</span>
                    </button>
                  </>
                ) : (
                  <>
                    <Link
                      to="/login"
                      onClick={closeMenu}
                      className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-slate-300 hover:bg-slate-800/50 transition-colors"
                    >
                      <span>🔑</span>
                      <span>登录</span>
                    </Link>
                    <Link
                      to="/register"
                      onClick={closeMenu}
                      className="flex items-center gap-3 px-3 py-2.5 rounded-lg bg-primary text-white hover:bg-primary-dark transition-colors"
                    >
                      <span>✨</span>
                      <span>注册</span>
                    </Link>
                  </>
                )}
              </nav>
            </motion.div>
          </>
        )}
      </AnimatePresence>
    </>
  )
}
