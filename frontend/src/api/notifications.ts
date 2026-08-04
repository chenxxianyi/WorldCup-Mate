import request from './request'
import type { ApiNotification } from '@/types/notification'
import type { PaginatedData } from '@/types/common'

export function apiListNotifications(params?: Record<string, any>) {
  return request.get('/notifications', { params }) as Promise<PaginatedData<ApiNotification>>
}

export function apiGetUnreadNotificationCount() {
  return request.get('/notifications/unread-count') as Promise<{ count: number }>
}

export function apiMarkNotificationRead(id: number) {
  return request.put(`/notifications/${id}/read`) as Promise<{ read: boolean }>
}

export function apiMarkAllNotificationsRead() {
  return request.put('/notifications/read-all') as Promise<{ read_all: boolean }>
}
