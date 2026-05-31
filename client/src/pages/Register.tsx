import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { useAuthStore } from '../stores/authStore'

export default function Register() {
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [validationError, setValidationError] = useState('')
  const { register, isLoading, error, clearError } = useAuthStore()
  const navigate = useNavigate()

  const validate = (): boolean => {
    if (!username.trim()) {
      setValidationError('请输入用户名')
      return false
    }
    if (username.length < 2 || username.length > 20) {
      setValidationError('用户名长度为2-20个字符')
      return false
    }
    if (!email.trim()) {
      setValidationError('请输入邮箱')
      return false
    }
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      setValidationError('请输入有效的邮箱地址')
      return false
    }
    if (!password) {
      setValidationError('请输入密码')
      return false
    }
    if (password.length < 6) {
      setValidationError('密码至少6位')
      return false
    }
    if (password !== confirmPassword) {
      setValidationError('两次输入的密码不一致')
      return false
    }
    setValidationError('')
    return true
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    clearError()
    if (!validate()) return

    try {
      const result = await register({ username, email, password })
      if (result?.need_verify) {
        navigate(`/verify-email?email=${encodeURIComponent(email)}`)
      } else {
        navigate('/login')
      }
    } catch {
      // Error is handled by store
    }
  }

  const displayError = validationError || error

  return (
    <div className="flex-1 flex items-center justify-center px-4 py-12">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="w-full max-w-md"
      >
        <div className="text-center mb-8">
          <div className="text-5xl mb-4">🔮</div>
          <h1 className="text-3xl font-bold text-white">创建账号</h1>
          <p className="text-slate-400 mt-2">开启你的玄机占卜之旅</p>
        </div>

        <form onSubmit={handleSubmit} className="bg-card rounded-2xl p-8 border border-slate-700/50 space-y-5">
          {displayError && (
            <div className="p-3 bg-red-500/10 border border-red-500/30 rounded-lg text-red-400 text-sm">
              {displayError}
            </div>
          )}

          <div>
            <label htmlFor="username" className="block text-sm font-medium text-slate-300 mb-1.5">
              用户名
            </label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="你的昵称"
              className="w-full px-4 py-2.5 bg-darker border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/50 transition-colors"
            />
          </div>

          <div>
            <label htmlFor="email" className="block text-sm font-medium text-slate-300 mb-1.5">
              邮箱
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="your@email.com"
              className="w-full px-4 py-2.5 bg-darker border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/50 transition-colors"
            />
          </div>

          <div>
            <label htmlFor="password" className="block text-sm font-medium text-slate-300 mb-1.5">
              密码
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="至少6位"
              className="w-full px-4 py-2.5 bg-darker border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/50 transition-colors"
            />
          </div>

          <div>
            <label htmlFor="confirmPassword" className="block text-sm font-medium text-slate-300 mb-1.5">
              确认密码
            </label>
            <input
              id="confirmPassword"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="再次输入密码"
              className="w-full px-4 py-2.5 bg-darker border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/50 transition-colors"
            />
          </div>

          <button
            type="submit"
            disabled={isLoading}
            className="w-full py-2.5 bg-primary hover:bg-primary-dark disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg font-medium transition-colors"
          >
            {isLoading ? '注册中...' : '注册'}
          </button>

          <p className="text-center text-sm text-slate-400">
            已有账号？{' '}
            <Link to="/login" className="text-primary hover:text-primary-light">
              立即登录
            </Link>
          </p>
        </form>
      </motion.div>
    </div>
  )
}
