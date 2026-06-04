<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { apiAdminSyncMatches } from '@/api/admin'
import { apiGetSyncStatus } from '@/api/sync'
import SyncStatusBadge from '@/components/common/SyncStatusBadge.vue'

const states = ref<any[]>([])
const syncing = ref(false)
const result = ref<any | null>(null)
const error = ref('')

async function loadStates() {
  try {
    states.value = await apiGetSyncStatus() as any[]
  } catch {
    states.value = []
  }
}

async function syncMatches() {
  syncing.value = true
  result.value = null
  error.value = ''
  try {
    result.value = await apiAdminSyncMatches()
    await loadStates()
  } catch (err: any) {
    error.value = err?.message || '同步失败'
  } finally {
    syncing.value = false
  }
}

function timeText(value: string) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleString('zh-CN')
}

onMounted(loadStates)
</script>

<template>
  <div class="admin-page">
    <div class="admin-head">
      <div>
        <h2>同步管理</h2>
        <span>赛事数据同步状态和手动同步</span>
      </div>
      <button class="pill-btn primary" :disabled="syncing" @click="syncMatches">
        {{ syncing ? '同步中...' : '手动同步比赛' }}
      </button>
    </div>

    <SyncStatusBadge mode="card" />

    <div v-if="result" class="card result-card">
      <b>本次同步完成</b>
      <span>总数 {{ result.total }}，创建 {{ result.created }}，更新 {{ result.updated }}，跳过 {{ result.skipped }}</span>
    </div>
    <div v-if="error" class="card result-card error">{{ error }}</div>

    <div class="card table-card table-scroll">
      <table class="admin-table">
        <thead>
          <tr>
            <th>Provider</th><th>Resource</th><th>Status</th><th>Last Synced</th><th>Next Sync</th><th>Error</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in states" :key="`${item.provider}-${item.resource}`">
            <td>{{ item.provider }}</td>
            <td>{{ item.resource }}</td>
            <td>{{ item.status }}</td>
            <td>{{ timeText(item.last_synced_at) }}</td>
            <td>{{ timeText(item.next_sync_at) }}</td>
            <td>{{ item.last_error || '-' }}</td>
          </tr>
          <tr v-if="!states.length">
            <td colspan="6" class="empty-row">暂无同步记录</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.admin-page {
  display: grid;
  gap: 14px;
}

.admin-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.admin-head h2 {
  margin: 0;
  font-size: 20px;
}

.admin-head span {
  color: var(--muted);
  font-size: 13px;
}

.result-card {
  padding: 14px;
  display: grid;
  gap: 4px;
}

.result-card span {
  color: var(--muted);
  font-size: 13px;
}

.result-card.error {
  color: var(--primary);
}

.table-card {
  overflow: hidden;
}

.table-scroll {
  overflow-x: auto;
}

.admin-table {
  width: 100%;
  min-width: 760px;
  border-collapse: collapse;
}

th,
td {
  padding: 12px 10px;
  border-bottom: 1px solid var(--line);
  text-align: left;
  font-size: 13px;
}

th {
  color: var(--muted);
  background: var(--card-soft);
}

.empty-row {
  color: var(--muted);
  text-align: center;
}
</style>
