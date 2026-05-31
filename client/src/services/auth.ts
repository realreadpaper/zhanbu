import api from './api'

export interface RegisterParams {
  username: string
  email: string
  password: string
}

export interface LoginParams {
  email: string
  password: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: UserProfile
}

export interface UserProfile {
  id: number
  username: string
  email: string
  email_verified: boolean
}

export interface RegisterResponse {
  user: UserProfile
  need_verify: boolean
}

export interface ApiResponse<T> {
  code: number
  data: T
  message: string
}

export const authApi = {
  register: (params: RegisterParams) =>
    api.post<ApiResponse<RegisterResponse>>('/auth/register', params),

  login: (params: LoginParams) =>
    api.post<ApiResponse<LoginResponse>>('/auth/login', params),

  getProfile: () =>
    api.get<ApiResponse<UserProfile>>('/auth/profile'),
}
