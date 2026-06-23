import { Routes, Route, useNavigate, Navigate, useLocation } from 'react-router-dom'
import { useEffect } from 'react'
import Header from './components/Header'
import Footer from './components/Footer'
import ErrorBoundary from './components/ErrorBoundary'
import ProtectedRoute from './components/ProtectedRoute'
import History from './pages/History'
import Profile from './pages/Profile'
import Login from './pages/Login'
import Register from './pages/Register'
import VerifyEmail from './pages/VerifyEmail'
import ChatPage from './pages/ChatPage'
import { useAuthStore } from './stores/authStore'
import { setNavigate } from './services/api'

export default function App() {
  const initAuth = useAuthStore((s) => s.initAuth)
  const navigate = useNavigate()
  const location = useLocation()
  const isChatRoute = location.pathname.startsWith('/chat')

  useEffect(() => {
    setNavigate(navigate)
  }, [navigate])

  useEffect(() => {
    initAuth()
  }, [initAuth])

  return (
    <ErrorBoundary>
      <Header />
      <main className="flex-1 min-h-0 flex flex-col">
        <Routes>
          <Route path="/" element={<Navigate to="/chat" replace />} />
          <Route path="/tarot" element={<Navigate to="/chat?type=tarot" replace />} />
          <Route path="/liuyao" element={<Navigate to="/chat?type=liuyao" replace />} />
          <Route path="/bazi" element={<Navigate to="/chat?type=bazi" replace />} />
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
          <Route
            path="/chat"
            element={
              <ProtectedRoute>
                <ChatPage key={location.search} />
              </ProtectedRoute>
            }
          />
          <Route
            path="/chat/:id"
            element={
              <ProtectedRoute>
                <ChatPage />
              </ProtectedRoute>
            }
          />
        </Routes>
      </main>
      {!isChatRoute && <Footer />}
    </ErrorBoundary>
  )
}
