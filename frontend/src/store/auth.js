import { defineStore } from 'pinia'
import { login as loginApi, logout as logoutApi, profile as profileApi } from '../api/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    profile: JSON.parse(localStorage.getItem('profile') || 'null')
  }),
  getters: {
    permissions: (state) => state.profile?.permissions || [],
    roles: (state) => state.profile?.roles || [],
    isLoggedIn: (state) => Boolean(state.token)
  },
  actions: {
    async login(form) {
      const result = await loginApi(form)
      this.token = result.token
      this.profile = result.user
      localStorage.setItem('token', result.token)
      localStorage.setItem('profile', JSON.stringify(result.user))
    },
    async loadProfile() {
      const result = await profileApi()
      this.profile = result
      localStorage.setItem('profile', JSON.stringify(result))
      return result
    },
    async logout() {
      try {
        await logoutApi()
      } finally {
        this.token = ''
        this.profile = null
        localStorage.removeItem('token')
        localStorage.removeItem('profile')
      }
    },
    hasPermission(permission) {
      if (!permission) return true
      if (this.roles.includes('admin')) return true
      return this.permissions.includes(permission) || this.permissions.includes('*:*:*')
    }
  }
})
