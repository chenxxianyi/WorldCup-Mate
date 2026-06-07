import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types/user'
import { normalizeUser } from '@/types/user'
import { apiLogin as apiLoginReq, apiRegister as apiRegisterReq, apiGetProfile, apiUploadAvatar, apiChangePassword } from '@/api/auth'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useReminderStore } from '@/stores/useReminderStore'
import { useNotificationStore } from '@/stores/useNotificationStore'

const TOKEN_KEY = 'wm-token'
const getStoredToken = () => localStorage.getItem(TOKEN_KEY) || sessionStorage.getItem(TOKEN_KEY)

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getStoredToken())
  const user = ref<User | null>(null)
  const isCheckingProfile = ref(!!token.value)
  const hasToken = computed(() => !!token.value)
  const isLoggedIn = computed(() => !!token.value && !!user.value)

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
    isCheckingProfile.value = false
    return user.value
  }

  async function register(username: string, email: string, password: string) {
    const res = await apiRegisterReq({ username, email, password }) as any
    token.value = res.token
    sessionStorage.setItem(TOKEN_KEY, res.token)
    localStorage.removeItem(TOKEN_KEY)
    user.value = normalizeUser(res.user)
    isCheckingProfile.value = false
    return user.value
  }

  async function fetchProfile() {
    if (!token.value) {
      isCheckingProfile.value = false
      return null
    }
    isCheckingProfile.value = true
    try {
      const res = await apiGetProfile() as any
      user.value = normalizeUser(res)
      return user.value
    } catch {
      clearSession()
      return null
    } finally {
      isCheckingProfile.value = false
    }
  }

  function clearSession() {
    token.value = null
    user.value = null
    isCheckingProfile.value = false
    localStorage.removeItem(TOKEN_KEY)
    sessionStorage.removeItem(TOKEN_KEY)
    clearUserScopedStores()
  }

  function logout() {
    clearSession()
  }

  function clearUserScopedStores() {
    useFavoriteStore().clearFavorites()
    useReminderStore().clearReminders()
    useNotificationStore().clearNotifications()
  }

  async function uploadAvatar(file: File) {
    const res = await apiUploadAvatar(file) as any
    if (user.value) user.value.avatar = res.avatar
    return res.avatar
  }

  async function changePassword(oldPassword: string, newPassword: string) {
    await apiChangePassword({ old_password: oldPassword, new_password: newPassword })
  }

  return { token, user, hasToken, isCheckingProfile, isLoggedIn, login, register, logout, fetchProfile, uploadAvatar, changePassword }
})
