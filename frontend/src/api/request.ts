import axios from 'axios'
import { useAuthStore } from '@/stores/useAuthStore'
import router from '@/router'

const request = axios.create({
  baseURL: '/api',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('wm-token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (res) => {
    const data = res.data
    if (data.code !== 0) {
      if (data.code === 401) {
        const auth = useAuthStore()
        auth.logout()
        router.push('/')
      }
      const message = data.code >= 500 ? '服务暂时不可用，请稍后重试' : data.message || '请求失败'
      return Promise.reject(new Error(message))
    }
    return data.data
  },
  (err) => {
    if (err.response?.status === 401) {
      const auth = useAuthStore()
      auth.logout()
    }
    const message = err.response?.status >= 500 || !err.response
      ? '服务暂时不可用，请稍后重试'
      : err.response?.data?.message || '请求失败'
    return Promise.reject(new Error(message))
  }
)

export default request
