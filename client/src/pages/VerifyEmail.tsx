import { useState, useEffect, type FormEvent } from 'react'
import { useNavigate, useSearchParams, Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { emailApi } from '../services/email'
import { useAuthStore } from '../stores/authStore'

export default function VerifyEmail() {
  const [searchParams] = useSearchParams()
  const email = searchParams.get('email') || ''
  const [code, setCode] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [resendCooldown, setResendCooldown] = useState(0)
  const navigate = useNavigate()
  const clearAuthError = useAuthStore((state) => state.clearError)

  // Countdown timer for resend button
  useEffect(() => {
    if (resendCooldown <= 0) return
    const timer = setInterval(() => {
      setResendCooldown((prev) => {
        if (prev <= 1) {
          clearInterval(timer)
          return 0
        }
        return prev - 1
      })
    }, 1000)
    return () => clearInterval(timer)
  }, [resendCooldown])

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')

    if (!code.trim()) {
      setError('请输入验证码')
      return
    }
    if (code.trim().length !== 6) {
      setError('验证码为6位数字')
      return
    }

    setIsLoading(true)
    try {
      const { data: res } = await emailApi.verifyEmail(email, code.trim())
      if (res.code === 0) {
        clearAuthError()
        setSuccess(true)
        setTimeout(() => navigate('/login'), 2000)
      } else {
        setError(res.message || '验证失败')
      }
    } catch (err: unknown) {
      const message =
        (err as { response?: { data?: { message?: string } } })?.response?.data?.message || '验证失败，请稍后重试'
      setError(message)
    } finally {
      setIsLoading(false)
    }
  }

  const handleResend = async () => {
    if (resendCooldown > 0) return
    setError('')
    try {
      const { data: res } = await emailApi.resendVerification(email)
      if (res.code === 0) {
        setResendCooldown(60)
      } else {
        setError(res.message || '发送失败')
      }
    } catch (err: unknown) {
      const message =
        (err as { response?: { data?: { message?: string } } })?.response?.data?.message || '发送失败，请稍后重试'
      setError(message)
    }
  }

  if (!email) {
    return (
      <div className="flex-1 flex items-center justify-center px-4 py-12">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="w-full max-w-md text-center"
        >
          <div className="text-5xl mb-4">📧</div>
          <h1 className="text-2xl font-bold text-white mb-4">缺少邮箱信息</h1>
          <p className="text-slate-400 mb-6">请从注册页面跳转到此页面</p>
          <Link
            to="/register"
            className="inline-block px-6 py-2.5 bg-primary hover:bg-primary-dark text-white rounded-lg font-medium transition-colors"
          >
            去注册
          </Link>
        </motion.div>
      </div>
    )
  }

  return (
    <div className="flex-1 flex items-center justify-center px-4 py-12">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="w-full max-w-md"
      >
        <div className="text-center mb-8">
          <div className="text-5xl mb-4">📧</div>
          <h1 className="text-3xl font-bold text-white">验证邮箱</h1>
          <p className="text-slate-400 mt-2">
            验证码已发送至 <span className="text-primary">{email}</span>
          </p>
        </div>

        <form onSubmit={handleSubmit} className="bg-card rounded-2xl p-8 border border-slate-700/50 space-y-5">
          {error && (
            <div className="p-3 bg-red-500/10 border border-red-500/30 rounded-lg text-red-400 text-sm">
              {error}
            </div>
          )}

          {success && (
            <div className="p-3 bg-green-500/10 border border-green-500/30 rounded-lg text-green-400 text-sm">
              验证成功！正在跳转到登录页面...
            </div>
          )}

          {!success && (
            <>
              <div>
                <label htmlFor="code" className="block text-sm font-medium text-slate-300 mb-1.5">
                  验证码
                </label>
                <input
                  id="code"
                  type="text"
                  value={code}
                  onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                  placeholder="请输入6位验证码"
                  maxLength={6}
                  className="w-full px-4 py-2.5 bg-darker border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/50 transition-colors text-center text-2xl tracking-[0.5em]"
                />
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="w-full py-2.5 bg-primary hover:bg-primary-dark disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg font-medium transition-colors"
              >
                {isLoading ? '验证中...' : '验证'}
              </button>

              <div className="flex items-center justify-between text-sm">
                <button
                  type="button"
                  onClick={handleResend}
                  disabled={resendCooldown > 0}
                  className="text-primary hover:text-primary-light disabled:text-slate-500 disabled:cursor-not-allowed transition-colors"
                >
                  {resendCooldown > 0 ? `重新发送 (${resendCooldown}s)` : '重新发送验证码'}
                </button>
                <Link to="/login" className="text-slate-400 hover:text-slate-300 transition-colors">
                  返回登录
                </Link>
              </div>
            </>
          )}
        </form>
      </motion.div>
    </div>
  )
}
