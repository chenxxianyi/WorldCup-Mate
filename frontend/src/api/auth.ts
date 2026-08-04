import request from './request'
import type { ApiUser } from '@/types/user'

export interface AuthSession {
  token: string
  refresh_token: string
  user: ApiUser
}

export function apiRegister(data: { username: string; email: string; password: string }) {
  return request.post('/auth/register', data) as Promise<AuthSession>
}

export function apiLogin(data: { email: string; password: string }) {
  return request.post('/auth/login', data) as Promise<AuthSession>
}

export function apiAdminLogin(data: { email: string; password: string }) {
  return request.post('/admin/login', data) as Promise<AuthSession>
}

export function apiLogout() {
  return request.post('/auth/logout') as Promise<null>
}

export function apiGetProfile() {
  return request.get('/user/profile') as Promise<ApiUser>
}

export function apiUpdateProfile(data: { timezone?: string; language?: string; notification_email?: string }) {
  return request.put('/user/profile', data) as Promise<ApiUser>
}

export function apiUploadAvatar(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/user/avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }) as Promise<{ avatar: string }>
}

export function apiChangePassword(data: { old_password: string; new_password: string }) {
  return request.put('/user/password', data) as Promise<null>
}
