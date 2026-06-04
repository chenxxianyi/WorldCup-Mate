<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/useAuthStore'
import { useNotificationStore } from '@/stores/useNotificationStore'

const auth = useAuthStore()
const notification = useNotificationStore()

function timeText(value: string) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(() => {
  if (auth.isLoggedIn) notification.fetchNotifications()
})
</script>

<template>
  <section class="section">
    <div class="section-head">
      <h2>通知中心</h2>
      <button
        v-if="notification.unreadCount > 0"
        class="text-action"
        @click="notification.markAllRead"
      >
        全部已读
      </button>
    </div>

    <div v-if="!auth.isLoggedIn" class="card empty-card">登录后查看通知</div>
    <div v-else-if="notification.loading" class="card empty-card">通知加载中...</div>
    <div v-else-if="!notification.notifications.length" class="card empty-card">暂无通知</div>
    <div v-else class="card notification-card">
      <button
        v-for="item in notification.notifications"
        :key="item.id"
        class="notification-item"
        :class="{ unread: !item.is_read }"
        @click="!item.is_read && notification.markRead(item.id)"
      >
        <span class="dot"></span>
        <span class="notification-copy">
          <b>{{ item.title }}</b>
          <small>{{ item.content }}</small>
        </span>
        <time>{{ timeText(item.created_at) }}</time>
      </button>
    </div>
  </section>
</template>

<style scoped>
.section {
  margin-top: 18px;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.section-head h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
}

.text-action {
  border: 0;
  color: var(--primary);
  background: transparent;
  font-size: 13px;
  font-weight: 750;
  cursor: pointer;
}

.empty-card {
  padding: 20px;
  text-align: center;
  color: var(--muted);
}

.notification-card {
  overflow: hidden;
}

.notification-item {
  width: 100%;
  min-height: 64px;
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border: 0;
  border-bottom: 1px solid var(--line);
  color: var(--text);
  background: transparent;
  text-align: left;
  cursor: pointer;
}

.notification-item:last-child {
  border-bottom: 0;
}

.notification-item.unread {
  background: color-mix(in srgb, var(--primary) 6%, transparent);
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: transparent;
}

.notification-item.unread .dot {
  background: var(--primary);
}

.notification-copy {
  min-width: 0;
  display: grid;
  gap: 4px;
}

.notification-copy b,
.notification-copy small {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notification-copy b {
  font-size: 13px;
}

.notification-copy small,
.notification-item time {
  color: var(--muted);
  font-size: 12px;
}
</style>
