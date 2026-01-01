import api from '@/lib/api'
import { LoginResponse, User } from '@/types'

export interface RegisterData {
  email: string
  password: string
  fullName: string
  phone: string
}

export interface LoginData {
  email: string
  password: string
}

export const authService = {
  register: async (data: RegisterData): Promise<void> => {
    const response = await api.post('/auth/register', data)
    return response.data
  },

  login: async (data: LoginData): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/login', data)
    return response.data
  },

  getMe: async (): Promise<User> => {
    const response = await api.get<User>('/auth/me')
    return response.data
  },
}
