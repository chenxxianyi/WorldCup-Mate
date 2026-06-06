import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  apiGetUnreadNotificationCount,
  apiListNotifications,
  apiMarkAllNotificationsRead,
  apiMarkNotificationRead,
} from '@/api/notifications'

export interface NotificationItem {
  id: number
  title: string
  content: string
  type: string
  target_type?: string
  target_id?: number
  is_read: boolean
  created_at: string
}

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<NotificationItem[]>([])
  const unreadCount = ref(0)
  const loading = ref(false)

  async function fetchUnreadCount() {
    try {
      const res = await apiGetUnreadNotificationCount() as any
      unreadCount.value = Number(res.count || 0)
    } catch {
      unreadCount.value = 0
    }
  }

  async function fetchNotifications() {
    loading.value = true
    try {
      const res = await apiListNotifications({ page: 1, page_size: 20 }) as any
      notifications.value = res.list || res || []
      unreadCount.value = notifications.value.filter((item) => !item.is_read).length
    } finally {
      loading.value = false
    }
  }

  async function markRead(id: number) {
    await apiMarkNotificationRead(id)
    const item = notifications.value.find((n) => n.id === id)
    if (item && !item.is_read) {
      item.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    }
  }

  async function markAllRead() {
    await apiMarkAllNotificationsRead()
    notifications.value = notifications.value.map((item) => ({ ...item, is_read: true }))
    unreadCount.value = 0
  }

  return {
    notifications,
    unreadCount,
    loading,
    fetchUnreadCount,
    fetchNotifications,
    markRead,
    markAllRead,
  }
})
