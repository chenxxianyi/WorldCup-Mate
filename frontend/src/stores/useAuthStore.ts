import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import type { User } from '@/types/user'
import { normalizeUser } from '@/types/user'
import { setUserTimezone } from '@/utils/datetime'
import { apiLogin as apiLoginReq, apiRegister as apiRegisterReq, apiAdminLogin as apiAdminLoginReq, apiGetProfile, apiUploadAvatar, apiChangePassword, apiUpdateProfile } from '@/api/auth'

const TOKEN_KEY = 'wm-token'
const REFRESH_KEY = 'wm-refresh-token'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY) || sessionStorage.getItem(TOKEN_KEY))
  const refreshToken = ref<string | null>(localStorage.getItem(REFRESH_KEY) || sessionStorage.getItem(REFRESH_KEY))
  const user = ref<User | null>(null)
  const isLoggedIn = computed(() => !!token.value)

  /** Persist both tokens under the same storage tier (SEC-04E). */
  function saveSession(access: string, refresh: string, persist: boolean) {
    localStorage.removeItem(TOKEN_KEY)
    sessionStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_KEY)
    sessionStorage.removeItem(REFRESH_KEY)
    ;(persist ? localStorage : sessionStorage).setItem(TOKEN_KEY, access)
    ;(persist ? localStorage : sessionStorage).setItem(REFRESH_KEY, refresh)
  }

  async function login(email: string, password: string, persist = true) {
    const res = await apiLoginReq({ email, password })
    token.value = res.token
    refreshToken.value = res.refresh_token
    saveSession(res.token, res.refresh_token, persist)
    user.value = normalizeUser(res.user)
    setUserTimezone(res.user.timezone)
    return user.value
  }

  async function register(username: string, email: string, password: string, persist = true) {
    const res = await apiRegisterReq({ username, email, password })
    token.value = res.token
    refreshToken.value = res.refresh_token
    saveSession(res.token, res.refresh_token, persist)
    user.value = normalizeUser(res.user)
    setUserTimezone(res.user.timezone)
    return user.value
  }

  // Admin sign-in (ADM-07): dedicated endpoint; backend enforces admin role.
  async function adminLogin(email: string, password: string, persist = true) {
    const res = await apiAdminLoginReq({ email, password })
    token.value = res.token
    refreshToken.value = res.refresh_token
    saveSession(res.token, res.refresh_token, persist)
    user.value = normalizeUser(res.user)
    setUserTimezone(res.user.timezone)
    return user.value
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      const res = await apiGetProfile()
      user.value = normalizeUser(res)
      setUserTimezone(res.timezone)
    } catch {
      logout()
    }
  }

  function clearSession() {
    token.value = null
    refreshToken.value = null
    user.value = null
    localStorage.removeItem(TOKEN_KEY)
    sessionStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_KEY)
    sessionStorage.removeItem(REFRESH_KEY)
    // DATA-05: logout resets the timezone on shared devices.
    setUserTimezone(null)
  }

  function logout() {
    // SEC-04C: revoke server-side sessions first, then clear local state.
    // The revocation request goes out on plain axios with the *current*
    // access token: it must not flow through this instance's 401
    // refresh-and-retry interceptor, or an expired access token would
    // trigger a refresh that writes tokens back after clearSession()
    // (session resurrection bug).
    const access = token.value
    clearSession()
    try {
      if (access) {
        axios.post('/api/auth/logout', null, {
          headers: { Authorization: `Bearer ${access}` },
        })
      }
    } catch { /* server-side revocation is best-effort; logout is idempotent */ }
  }

  async function uploadAvatar(file: File) {
    const res = await apiUploadAvatar(file)
    if (user.value) user.value.avatar = res.avatar
    return res.avatar
  }

  async function updateProfile(data: { timezone?: string; language?: string; notification_email?: string }) {
    const res = await apiUpdateProfile(data)
    user.value = normalizeUser(res)
    return user.value
  }

  async function changePassword(oldPassword: string, newPassword: string) {
    await apiChangePassword({ old_password: oldPassword, new_password: newPassword })
    // SEC-04D: backend revokes all sessions on password change; the client
    // must discard its (now invalid) refresh token as well.
    clearSession()
  }

  return { token, refreshToken, user, isLoggedIn, login, register, adminLogin, logout, fetchProfile, uploadAvatar, updateProfile, changePassword, clearSession }
})
