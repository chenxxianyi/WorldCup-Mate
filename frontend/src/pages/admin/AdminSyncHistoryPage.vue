<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { apiAdminSyncHistory } from '@/api/admin'
import type { SyncHistoryItem } from '@/api/admin'

const history = ref<SyncHistoryItem[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await apiAdminSyncHistory({ page: 1, page_size: 50 })
    history.value = res.list
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="admin-page">
    <div class="admin-head">
      <div>
        <h2>同步历史</h2>
        <span>{{ history.length }} 条记录</span>
      </div>
    </div>

    <div class="card table-card table-scroll">
      <table class="admin-table">
        <thead>
          <tr>
            <th>时间</th><th>提供方</th><th>资源</th><th>原因</th><th>状态</th><th>总数</th><th>新建</th><th>更新</th><th>跳过</th><th>错误</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in history" :key="row.id">
            <td>{{ row.started_at }}</td>
            <td>{{ row.provider }}</td>
            <td>{{ row.resource }}</td>
            <td>{{ row.reason }}</td>
            <td>{{ row.status }}</td>
            <td>{{ row.total }}</td>
            <td>{{ row.created }}</td>
            <td>{{ row.updated }}</td>
            <td>{{ row.skipped }}</td>
            <td style="max-width:200px;word-break:break-all">{{ row.error_message || '-' }}</td>
          </tr>
          <tr v-if="!loading && !history.length">
            <td colspan="10" class="empty-row">暂无同步历史</td>
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
