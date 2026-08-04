<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { apiGetSyncStatus } from '@/api/sync'
import type { SyncState } from '@/types/sync'

const props = withDefaults(defineProps<{
  mode?: 'line' | 'card'
}>(), {
  mode: 'line',
})

const loading = ref(false)
const state = ref<SyncState | null>(null)

const statusText = computed(() => {
  if (loading.value) return '同步状态读取中'
  if (!state.value) return '演示数据或未同步'
  if (state.value.status === 'success') return '已同步'
  if (state.value.status === 'running') return '同步中'
  if (state.value.status === 'failed') return '同步失败'
  return state.value.status || '未知状态'
})

const timeText = computed(() => {
  const value = state.value?.last_synced_at
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
})

async function loadStatus() {
  loading.value = true
  try {
    const res = await apiGetSyncStatus()
    state.value = (res || []).find((item) => item.resource === 'matches') || res?.[0] || null
  } catch {
    state.value = null
  } finally {
    loading.value = false
  }
}

onMounted(loadStatus)
</script>

<template>
  <div :class="mode === 'card' ? 'sync-card' : 'sync-line'">
    <span class="sync-dot" :class="state?.status || 'idle'"></span>
    <div class="sync-copy">
      <b>{{ statusText }}</b>
      <span v-if="timeText">更新时间 {{ timeText }}</span>
      <span v-else>暂无同步时间</span>
      <small v-if="mode === 'card' && state?.last_error">{{ state.last_error }}</small>
    </div>
  </div>
</template>

<style scoped>
.sync-line,
.sync-card {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--muted);
}

.sync-line {
  font-size: 12px;
}

.sync-card {
  width: 100%;
  padding: 14px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--card);
}

.sync-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: var(--weak);
  flex: 0 0 auto;
}

.sync-dot.success {
  background: var(--success);
}

.sync-dot.running {
  background: var(--gold);
}

.sync-dot.failed {
  background: var(--primary);
}

.sync-copy {
  min-width: 0;
  display: grid;
  gap: 3px;
}

.sync-copy b {
  color: var(--text);
  font-size: 13px;
}

.sync-copy span,
.sync-copy small {
  color: var(--muted);
  font-size: 12px;
}

.sync-copy small {
  overflow-wrap: anywhere;
}
</style>
