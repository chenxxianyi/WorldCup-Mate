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

export function apiUpdateProfile(data: { avatar?: string; timezone?: string; language?: string }) {
  return request.put('/user/profile', data) as Promise<any>
}
