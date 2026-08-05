import axios, { AxiosError, type AxiosRequestConfig } from 'axios'
import { useAuthStore } from '@/stores/useAuthStore'
import router from '@/router'
import { ApiError } from '@/types/common'
import { getUserTimezone } from '@/utils/datetime'

const TOKEN_KEY = 'wm-token'
const REFRESH_KEY = 'wm-refresh-token'

const request = axios.create({
  baseURL: '/api',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY) || sessionStorage.getItem(TOKEN_KEY)
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  // DATA-05: today/tomorrow boundaries follow the user's timezone.
  if (config.url?.startsWith('/matches/today') || config.url?.startsWith('/matches/tomorrow')) {
    config.params = { ...(config.params || {}), tz: getUserTimezone() }
  }
  // Multi-competition: league calls use competitionId. Legacy World Cup
  // matches have a NULL competition_id, so match-list calls explicitly use
  // worldCup=true to prevent synced league fixtures leaking into that view.
  const compCode = localStorage.getItem('wm-competition') || 'WC'
  if (compCode === 'WC' && config.url?.startsWith('/matches') && !config.params?.competitionId) {
    config.params = { ...(config.params || {}), worldCup: true }
  } else if (compCode !== 'WC' && !config.params?.competitionId) {
    const compId = Number(localStorage.getItem('wm-competition-id'))
    if (Number.isFinite(compId) && compId > 0) {
      config.params = { ...(config.params || {}), competitionId: compId }
    }
  }
  return config
})

// ---------------------------------------------------------------------------
// SEC-04F: single-flight access-token refresh.
// Concurrent 401s share ONE refresh request; on success the original calls
// are replayed with the fresh token; on failure the session is cleared.
// ---------------------------------------------------------------------------

function readRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_KEY) || sessionStorage.getItem(REFRESH_KEY)
}

/** Persist the rotated pair under the same storage tier as before. */
function writeTokens(access: string, refresh: string) {
  const persist = !!localStorage.getItem(REFRESH_KEY)
  localStorage.removeItem(TOKEN_KEY)
  sessionStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(REFRESH_KEY)
  sessionStorage.removeItem(REFRESH_KEY)
  ;(persist ? localStorage : sessionStorage).setItem(TOKEN_KEY, access)
  ;(persist ? localStorage : sessionStorage).setItem(REFRESH_KEY, refresh)
  const auth = useAuthStore()
  auth.token = access
  auth.refreshToken = refresh
}

function forceLogout() {
  const auth = useAuthStore()
  auth.clearSession()
  if (router.currentRoute.value.path !== '/') router.push('/')
}

let refreshPromise: Promise<string> | null = null

function getFreshAccessToken(): Promise<string> {
  if (!refreshPromise) {
    refreshPromise = (async () => {
      const rt = readRefreshToken()
      if (!rt) throw new Error('no refresh token')
      // Plain axios on purpose: bypasses this instance's interceptors
      // (no recursion into the 401 flow).
      const res = await axios.post('/api/auth/refresh', { refresh_token: rt })
      const data = res.data
      if (data.code !== 0) throw new ApiError(data.code || 401, data.message || '刷新登录状态失败', data.request_id)
      writeTokens(data.data.token, data.data.refresh_token)
      return data.data.token as string
    })().finally(() => {
      refreshPromise = null
    })
  }
  return refreshPromise
}

async function refreshAndRetry(originalError: AxiosError): Promise<unknown> {
  const config = originalError.config as AxiosRequestConfig & { headers: Record<string, string>; _retried?: boolean }
  // Guard against a refresh loop (e.g. account disabled mid-flight):
  // replay exactly once, then force logout.
  if (config._retried) {
    forceLogout()
    return Promise.reject(new ApiError(401, '登录已过期，请重新登录'))
  }
  try {
    const newToken = await getFreshAccessToken()
    config._retried = true
    // AxiosHeaders (v1) has a typed `set`; fall back for plain objects.
    if (config.headers && typeof (config.headers as { set?: unknown }).set === 'function') {
      ;(config.headers as { set: (k: string, v: string) => void }).set('Authorization', `Bearer ${newToken}`)
    } else {
      ;(config as { headers?: unknown }).headers = { Authorization: `Bearer ${newToken}` }
    }
    return request(config)
  } catch {
    forceLogout()
    return Promise.reject(new ApiError(401, '登录已过期，请重新登录'))
  }
}

request.interceptors.response.use(
  (res) => {
    const data = res.data
    if (data.code !== 0) {
      const message = data.code >= 500 ? '服务暂时不可用，请稍后重试' : data.message || '请求失败'
      return Promise.reject(new ApiError(data.code || 500, message, data.request_id))
    }
    return data.data
  },
  (err) => {
    // LIVE-02: pass aborted requests through untouched — the caller's
    // isCancel() must see the original CanceledError (code ERR_CANCELED),
    // not a wrapped ApiError.
    if (axios.isCancel(err)) return Promise.reject(err)
    // API-01: backend now returns real HTTP status codes (4xx/5xx).
    const status = err.response?.status
    const url = err.config?.url || ''
    const isAuthCall = url.includes('/auth/login') || url.includes('/auth/register') || url.includes('/admin/login') || url.includes('/auth/refresh') || url.includes('/auth/logout')
    if (status === 401 && !isAuthCall) {
      // Login/register/refresh failures surface their own message; anything
      // else tries the refresh flow once.
      return refreshAndRetry(err)
    }
    let message = '请求失败'
    if (!err.response) {
      message = '网络连接失败，请稍后重试'
    } else if (status === 429) {
      message = '请求过于频繁，请稍后再试'
    } else if (status >= 500) {
      message = '服务暂时不可用，请稍后重试'
    } else {
      message = err.response?.data?.message || '请求失败'
    }
    return Promise.reject(new ApiError(status || 0, message, err.response?.data?.request_id))
  }
)

export default request
