import axios from 'axios'
import { clearAuthStorage } from '@/utils/logout'

declare module 'axios' {
  interface AxiosRequestConfig {
    skipUnauthorizedRedirect?: boolean
  }
}

const request = axios.create({
  baseURL: '/api',
  timeout: 15000,
})

let loginRedirect: Promise<void> | null = null

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('wm-token') || sessionStorage.getItem('wm-token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (res) => {
    const data = res.data
    if (data.code !== 0) {
      if (data.code === 401 && !res.config.skipUnauthorizedRedirect) {
        void redirectToLogin()
      }
      const message = data.code >= 500 ? 'Service temporarily unavailable, please try again later' : data.message || 'Request failed'
      return Promise.reject(new Error(message))
    }
    return data.data
  },
  (err) => {
    if (err.response?.status === 401 && !err.config?.skipUnauthorizedRedirect) {
      void redirectToLogin()
    }
    const message = err.response?.status >= 500 || !err.response
      ? 'Service temporarily unavailable, please try again later'
      : err.response?.data?.message || 'Request failed'
    return Promise.reject(new Error(message))
  },
)

async function redirectToLogin() {
  if (loginRedirect) return loginRedirect

  loginRedirect = runLoginRedirect().finally(() => {
    loginRedirect = null
  })

  return loginRedirect
}

async function runLoginRedirect() {
  const { default: router } = await import('@/router')
  const currentRoute = router.currentRoute.value

  clearAuthStorage()
  if (currentRoute.name !== 'login') {
    await router.replace({ name: 'login', query: { redirect: currentRoute.fullPath } }).catch(() => {})
  }
}

export default request
