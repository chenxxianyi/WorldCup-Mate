import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types/user'
import { normalizeUser } from '@/types/user'
import { apiLogin as apiLoginReq, apiRegister as apiRegisterReq, apiGetProfile, apiLogout as apiLogoutReq, apiUploadAvatar, apiChangePassword, apiUpdateProfile } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('wm-token') || sessionStorage.getItem('wm-token'))
  const user = ref<User | null>(null)
  const isLoggedIn = computed(() => !!token.value)

  function saveToken(value: string, persist: boolean) {
    localStorage.removeItem('wm-token')
    sessionStorage.removeItem('wm-token')
    ;(persist ? localStorage : sessionStorage).setItem('wm-token', value)
  }

  async function login(email: string, password: string, persist = true) {
    const res = await apiLoginReq({ email, password }) as any
    token.value = res.token
    saveToken(res.token, persist)
    user.value = normalizeUser(res.user)
    return user.value
  }

  async function register(username: string, email: string, password: string, persist = true) {
    const res = await apiRegisterReq({ username, email, password }) as any
    token.value = res.token
    saveToken(res.token, persist)
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
    sessionStorage.removeItem('wm-token')
    try { apiLogoutReq() } catch {}
  }

  async function uploadAvatar(file: File) {
    const res = await apiUploadAvatar(file) as any
    if (user.value) user.value.avatar = res.avatar
    return res.avatar
  }

  async function updateProfile(data: { timezone?: string; language?: string; notification_email?: string }) {
    const res = await apiUpdateProfile(data) as any
    user.value = normalizeUser(res)
    return user.value
  }

  async function changePassword(oldPassword: string, newPassword: string) {
    await apiChangePassword({ old_password: oldPassword, new_password: newPassword })
  }

  return { token, user, isLoggedIn, login, register, logout, fetchProfile, uploadAvatar, updateProfile, changePassword }
})
