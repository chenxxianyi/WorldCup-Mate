import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types/user'
import { normalizeUser } from '@/types/user'
import { apiLogin as apiLoginReq, apiRegister as apiRegisterReq, apiGetProfile, apiLogout as apiLogoutReq, apiUploadAvatar, apiChangePassword } from '@/api/auth'

const TOKEN_KEY = 'wm-token'
const getStoredToken = () => localStorage.getItem(TOKEN_KEY) || sessionStorage.getItem(TOKEN_KEY)

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getStoredToken())
  const user = ref<User | null>(null)
  const isLoggedIn = computed(() => !!token.value)

  function persistToken(nextToken: string, persistent: boolean) {
    token.value = nextToken
    if (persistent) {
      localStorage.setItem(TOKEN_KEY, nextToken)
      sessionStorage.removeItem(TOKEN_KEY)
    } else {
      sessionStorage.setItem(TOKEN_KEY, nextToken)
      localStorage.removeItem(TOKEN_KEY)
    }
  }

  async function login(email: string, password: string, autoLogin = false) {
    const res = await apiLoginReq({ email, password, remember_me: autoLogin }) as any
    persistToken(res.token, autoLogin)
    user.value = normalizeUser(res.user)
    return user.value
  }

  async function register(username: string, email: string, password: string) {
    const res = await apiRegisterReq({ username, email, password }) as any
    token.value = res.token
    sessionStorage.setItem(TOKEN_KEY, res.token)
    localStorage.removeItem(TOKEN_KEY)
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
    localStorage.removeItem(TOKEN_KEY)
    sessionStorage.removeItem(TOKEN_KEY)
    try { apiLogoutReq() } catch {}
  }

  async function uploadAvatar(file: File) {
    const res = await apiUploadAvatar(file) as any
    if (user.value) user.value.avatar = res.avatar
    return res.avatar
  }

  async function changePassword(oldPassword: string, newPassword: string) {
    await apiChangePassword({ old_password: oldPassword, new_password: newPassword })
  }

  return { token, user, isLoggedIn, login, register, logout, fetchProfile, uploadAvatar, changePassword }
})
