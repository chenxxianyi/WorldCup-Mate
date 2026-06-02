import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types/user'
import { normalizeUser } from '@/types/user'
import { apiLogin as apiLoginReq, apiRegister as apiRegisterReq, apiGetProfile, apiLogout as apiLogoutReq } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('wm-token'))
  const user = ref<User | null>(null)
  const isLoggedIn = computed(() => !!token.value)

  async function login(email: string, password: string) {
    const res = await apiLoginReq({ email, password }) as any
    token.value = res.token
    localStorage.setItem('wm-token', res.token)
    user.value = normalizeUser(res.user)
    return user.value
  }

  async function register(username: string, email: string, password: string) {
    const res = await apiRegisterReq({ username, email, password }) as any
    token.value = res.token
    localStorage.setItem('wm-token', res.token)
    user.value = normalizeUser(res.user)
    return user.value
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      const res = await apiGetProfile() as any
      user.value = normalizeUser(res)
    } catch {
      logout()
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('wm-token')
    try { apiLogoutReq() } catch {}
  }

  return { token, user, isLoggedIn, login, register, logout, fetchProfile }
})
