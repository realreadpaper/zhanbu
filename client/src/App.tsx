import { Routes, Route, useNavigate } from 'react-router-dom'
import { useEffect } from 'react'
import Header from './components/Header'
import Footer from './components/Footer'
import ErrorBoundary from './components/ErrorBoundary'
import ProtectedRoute from './components/ProtectedRoute'
import Home from './pages/Home'
import Tarot from './pages/Tarot'
import Horoscope from './pages/Horoscope'
import LiuYao from './pages/LiuYao'
import BaZi from './pages/BaZi'
import History from './pages/History'
import Profile from './pages/Profile'
import Login from './pages/Login'
import Register from './pages/Register'
import VerifyEmail from './pages/VerifyEmail'
import { useAuthStore } from './stores/authStore'
import { setNavigate } from './services/api'

export default function App() {
  const initAuth = useAuthStore((s) => s.initAuth)
  const navigate = useNavigate()

  useEffect(() => {
    setNavigate(navigate)
  }, [navigate])

  useEffect(() => {
    initAuth()
  }, [initAuth])

  return (
    <ErrorBoundary>
      <Header />
      <main className="flex-1 flex flex-col">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/tarot" element={<Tarot />} />
          <Route path="/horoscope" element={<Horoscope />} />
          <Route path="/liuyao" element={<LiuYao />} />
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
