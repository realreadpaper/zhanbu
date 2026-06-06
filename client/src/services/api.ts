import axios from 'axios'

export const apiBaseURL = (import.meta.env.VITE_API_BASE_URL || '/api').replace(/\/$/, '')

export function apiURL(path: string) {
  return `${apiBaseURL}${path.startsWith('/') ? path : `/${path}`}`
}

const api = axios.create({
  baseURL: apiBaseURL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Navigation function that can be set from React Router context
let navigateFn: ((path: string) => void) | null = null

export function setNavigate(fn: (path: string) => void) {
  navigateFn = fn
}

// Redirect to login page - use React Router navigate if available, otherwise hard redirect
const redirectToLogin = () => {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
  if (navigateFn) {
    navigateFn('/login')
  } else {
    window.location.href = '/login'
  }
}

// Track if we're currently refreshing to avoid multiple refresh attempts
let isRefreshing = false
let failedQueue: Array<{
  resolve: (value: unknown) => void
  reject: (reason?: unknown) => void
}> = []

const processQueue = (error: unknown) => {
  const queue = [...failedQueue]
  failedQueue = []
  queue.forEach((prom) => {
    if (error) {
      prom.reject(error)
    } else {
      prom.resolve(undefined)
    }
  })
}

// Request interceptor - add token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor - handle 401 with token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    // If not a 401 or already retried, reject immediately
    if (error.response?.status !== 401 || originalRequest._retry) {
      return Promise.reject(error)
    }

    // Don't try to refresh if the failing request was itself a refresh
    if (originalRequest.url === '/auth/refresh') {
      redirectToLogin()
      return Promise.reject(error)
    }

    // If already refreshing, queue this request
    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        failedQueue.push({ resolve, reject })
      }).then(() => {
        return api(originalRequest)
      }).catch((err) => {
        return Promise.reject(err)
      })
    }

    originalRequest._retry = true
    isRefreshing = true

    const refreshToken = localStorage.getItem('refresh_token')

    if (!refreshToken) {
      isRefreshing = false
      redirectToLogin()
      return Promise.reject(error)
    }

    try {
      const { data } = await axios.post(apiURL('/auth/refresh'), {
        refresh_token: refreshToken,
      })

      if (data.code === 0 && data.data) {
        const { access_token, refresh_token } = data.data
        localStorage.setItem('access_token', access_token)
        if (refresh_token) {
          localStorage.setItem('refresh_token', refresh_token)
        }

        // Update the original request with new token
        originalRequest.headers.Authorization = `Bearer ${access_token}`

        // Reset flag BEFORE processing queue so that if any queued retry
        // hits another 401, it can start a fresh refresh cycle instead of
        // being added to a queue that nobody will drain.
        isRefreshing = false

        // Process queued requests - they will retry with the new token
        processQueue(null)

        // Retry the original request
        return api(originalRequest)
      } else {
        throw new Error('Refresh failed')
      }
    } catch (refreshError) {
      // Reset flag before draining the queue, same reasoning as above
      isRefreshing = false

      // Refresh failed - reject queued requests, clear tokens, redirect
      processQueue(refreshError)
      redirectToLogin()
      return Promise.reject(refreshError)
    }
  }
)

export default api
