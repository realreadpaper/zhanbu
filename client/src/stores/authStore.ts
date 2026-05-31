import { create } from 'zustand'
import { authApi, type UserProfile, type LoginParams, type RegisterParams, type RegisterResponse } from '../services/auth'

interface AuthState {
  user: UserProfile | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
  errorCode: number | null

  login: (params: LoginParams) => Promise<void>
  register: (params: RegisterParams) => Promise<RegisterResponse>
  logout: () => void
  fetchProfile: () => Promise<void>
  clearError: () => void
  initAuth: () => Promise<void>
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: !!localStorage.getItem('access_token'),
  isLoading: false,
  error: null,
  errorCode: null,

  login: async (params) => {
    set({ isLoading: true, error: null, errorCode: null })
    try {
      const { data: res } = await authApi.login(params)
      if (res.code === 0 && res.data) {
        localStorage.setItem('access_token', res.data.access_token)
        localStorage.setItem('refresh_token', res.data.refresh_token)
        set({
          user: res.data.user,
          isAuthenticated: true,
          isLoading: false,
          errorCode: null,
        })
      } else {
        set({ error: res.message || '登录失败', errorCode: res.code, isLoading: false })
      }
    } catch (err: unknown) {
      const axiosData = (err as { response?: { data?: { code?: number; message?: string } } })?.response?.data
      const message = axiosData?.message || '登录失败，请稍后重试'
      const code = axiosData?.code || null
      set({ error: message, errorCode: code, isLoading: false })
      throw err
    }
  },

  register: async (params) => {
    set({ isLoading: true, error: null })
    try {
      const { data: res } = await authApi.register(params)
      if (res.code === 0 && res.data) {
        set({ isLoading: false })
        return res.data
      } else {
        set({ error: res.message || '注册失败', isLoading: false })
        throw new Error(res.message)
      }
    } catch (err: unknown) {
      const message =
        (err as { response?: { data?: { message?: string } } })?.response?.data?.message || '注册失败，请稍后重试'
      set({ error: message, isLoading: false })
      throw err
    }
  },

  logout: () => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    set({ user: null, isAuthenticated: false, error: null })
  },

  fetchProfile: async () => {
    try {
      const { data: res } = await authApi.getProfile()
      if (res.code === 0 && res.data) {
        set({ user: res.data, isAuthenticated: true })
      }
    } catch {
      // Token invalid
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      set({ user: null, isAuthenticated: false })
    }
  },

  clearError: () => set({ error: null, errorCode: null }),

  initAuth: async () => {
    const token = localStorage.getItem('access_token')
    if (token) {
      try {
        const { data: res } = await authApi.getProfile()
        if (res.code === 0 && res.data) {
          set({ user: res.data, isAuthenticated: true })
        } else {
          set({ isAuthenticated: false })
        }
      } catch {
        set({ isAuthenticated: false })
      }
    }
  },
}))
