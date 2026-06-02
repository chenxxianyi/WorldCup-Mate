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
      return Promise.reject(new Error(data.message || 'request failed'))
    }
    return data.data
  },
  (err) => {
    if (err.response?.status === 401) {
      const auth = useAuthStore()
      auth.logout()
    }
    return Promise.reject(err)
  }
)

export default request
