import request from './request'

export function apiListNotifications(params?: Record<string, any>) {
  return request.get('/notifications', { params }) as Promise<any>
}

export function apiGetUnreadNotificationCount() {
  return request.get('/notifications/unread-count') as Promise<{ count: number }>
}

export function apiMarkNotificationRead(id: number) {
  return request.put(`/notifications/${id}/read`) as Promise<any>
}

export function apiMarkAllNotificationsRead() {
  return request.put('/notifications/read-all') as Promise<any>
}
