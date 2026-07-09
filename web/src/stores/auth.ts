import { defineStore } from 'pinia'

import { getUserInfo, login } from '@/services/auth'
import { clearToken, getToken, setToken } from '@/services/storage'
import type { LoginPayload, UserInfo } from '@/types/api'

interface AuthState {
  token: string
  user: UserInfo | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: getToken(),
    user: null,
  }),
  getters: {
    isLogin: (state) => Boolean(state.token),
  },
  actions: {
    async login(payload: LoginPayload) {
      const token = await login(payload)
      this.token = token
      setToken(token)
      await this.loadUser()
    },
    async loadUser() {
      if (!this.token) {
        return
      }
      this.user = await getUserInfo()
    },
    logout() {
      this.token = ''
      this.user = null
      clearToken()
    },
  },
})
