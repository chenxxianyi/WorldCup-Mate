<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
  apiAdminListReminders,
  apiAdminRetryReminder,
} from '@/api/admin'
import type { AdminReminderItem } from '@/api/admin'

const reminders = ref<AdminReminderItem[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await apiAdminListReminders({ page: 1, page_size: 50 })
    reminders.value = res.list
  } finally {
    loading.value = false
  }
}

async function retry(id: number) {
  try {
    await apiAdminRetryReminder(id)
    ElMessage.success('已重新入队')
    await load()
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '重试失败')
  }
}

onMounted(load)
</script>

<template>
  <div class="admin-page">
    <div class="admin-head">
      <div>
        <h2>提醒管理</h2>
        <span>{{ reminders.length }} 条提醒</span>
      </div>
    </div>

    <div class="card table-card table-scroll">
      <table class="admin-table">
        <thead>
          <tr>
            <th>ID</th><th>用户</th><th>比赛</th><th>提醒时间</th><th>渠道</th><th>状态</th><th>重试</th><th>下次重试</th><th>错误</th><th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in reminders" :key="r.id">
            <td>{{ r.id }}</td>
            <td>{{ r.user_id }}</td>
            <td>{{ r.match ? `${r.match.home_team?.name || ''} vs ${r.match.away_team?.name || ''}` : `match #${r.match_id}` }}</td>
            <td>{{ r.remind_at }}</td>
            <td>{{ r.channel }}</td>
            <td>{{ r.status }}</td>
            <td>{{ r.retry_count }}</td>
            <td>{{ r.next_retry_at || '-' }}</td>
            <td style="max-width:200px;word-break:break-all">{{ r.last_error || '-' }}</td>
            <td>
              <el-button
                size="small"
                :disabled="r.status === 'sent'"
                @click="retry(r.id)"
              >
                重试
              </el-button>
            </td>
          </tr>
          <tr v-if="!loading && !reminders.length">
            <td colspan="10" class="empty-row">暂无提醒</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.admin-page { display: grid; gap: 14px; }
.admin-head { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.admin-head h2 { margin: 0; font-size: 20px; }
.admin-head span { color: var(--muted); font-size: 13px; }
.card { padding: 14px; border-radius: 12px; background: var(--card); border: 1px solid var(--line); }
.table-card { overflow: hidden; }
.table-scroll { overflow-x: auto; }
.admin-table { width: 100%; min-width: 900px; border-collapse: collapse; }
th, td { padding: 10px 8px; border-bottom: 1px solid var(--line); text-align: left; font-size: 12px; }
th { color: var(--muted); background: var(--card-soft); }
.empty-row { color: var(--muted); text-align: center; padding: 24px; }
</style>
