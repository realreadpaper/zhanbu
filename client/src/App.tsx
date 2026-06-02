import { Routes, Route, useNavigate } from 'react-router-dom'
import { useEffect, useState } from 'react'
import Header from './components/Header'
import Footer from './components/Footer'
import ErrorBoundary from './components/ErrorBoundary'
import ProtectedRoute from './components/ProtectedRoute'
import Home from './pages/Home'
import Tarot from './pages/Tarot'
import LiuYao from './pages/LiuYao'
import LiuYaoV2 from './pages/LiuYaoV2'
import BaZi from './pages/BaZi'
import History from './pages/History'
import Profile from './pages/Profile'
import Login from './pages/Login'
import Register from './pages/Register'
import VerifyEmail from './pages/VerifyEmail'
import { useAuthStore } from './stores/authStore'
import { setNavigate } from './services/api'
import { fetchLiuYaoConfig } from './services/liuyao'

export default function App() {
  const initAuth = useAuthStore((s) => s.initAuth)
  const navigate = useNavigate()
  const [liuyaoVersion, setLiuYaoVersion] = useState<string>('v1')

  useEffect(() => {
    setNavigate(navigate)
  }, [navigate])

  useEffect(() => {
    initAuth()
  }, [initAuth])

  useEffect(() => {
    fetchLiuYaoConfig()
      .then((config) => {
        setLiuYaoVersion(config.version)
      })
      .catch(() => {
        // 默认使用v1
        setLiuYaoVersion('v1')
      })
  }, [])

  // 根据配置渲染对应的六爻组件
  const LiuYaoPage = liuyaoVersion === 'v2' ? LiuYaoV2 : LiuYao

  return (
    <ErrorBoundary>
      <Header />
      <main className="flex-1 flex flex-col">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/tarot" element={<Tarot />} />
          <Route path="/liuyao" element={<LiuYaoPage />} />
          <Route path="/bazi" element={<BaZi />} />
          <Route
            path="/history"
            element={
              <ProtectedRoute>
                <History />
              </ProtectedRoute>
            }
          />
          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <Profile />
              </ProtectedRoute>
            }
          />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/verify-email" element={<VerifyEmail />} />
        </Routes>
      </main>
      <Footer />
    </ErrorBoundary>
  )
}
