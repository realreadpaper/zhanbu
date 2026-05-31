import api from './api'
import type { ApiResponse } from './auth'

export const emailApi = {
  verifyEmail: (email: string, code: string) =>
    api.post<ApiResponse<null>>('/auth/verify-email', { email, code }),

  resendVerification: (email: string) =>
    api.post<ApiResponse<null>>('/auth/resend-verification', { email }),
}
