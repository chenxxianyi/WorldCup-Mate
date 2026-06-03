import request from './request'

export function apiRegister(data: { username: string; email: string; password: string }) {
  return request.post('/auth/register', data) as Promise<{ token: string; user: any }>
}

export function apiLogin(data: { email: string; password: string }) {
  return request.post('/auth/login', data) as Promise<{ token: string; user: any }>
}

export function apiLogout() {
  return request.post('/auth/logout')
}

export function apiGetProfile() {
  return request.get('/user/profile') as Promise<any>
}

export function apiUpdateProfile(data: { avatar?: string; timezone?: string; language?: string; notification_email?: string }) {
  return request.put('/user/profile', data) as Promise<any>
}

export function apiUploadAvatar(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/user/avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }) as Promise<{ avatar: string }>
}

export function apiChangePassword(data: { old_password: string; new_password: string }) {
  return request.put('/user/password', data) as Promise<any>
}
