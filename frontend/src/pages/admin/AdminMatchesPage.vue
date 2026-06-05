<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { apiAdminListMatches, apiAdminSyncMatchLineups } from '@/api/admin'
import { normalizeMatch, type Match } from '@/types/match'

const search = ref('')
const status = ref('')
const matches = ref<Match[]>([])
const loading = ref(false)
const syncingLineupId = ref<number | null>(null)
const syncMessage = ref('')

const filteredMatches = computed(() => {
  const q = search.value.trim().toLowerCase()
  return matches.value.filter((match) => {
    const statusOk = !status.value || match.status === status.value
    const textOk = !q || [
      match.home_team_name,
      match.away_team_name,
      match.home_team_code,
      match.away_team_code,
      match.group_name,
      match.city,
      match.stadium,
    ].filter(Boolean).some((value) => String(value).toLowerCase().includes(q))
    return statusOk && textOk
  })
})

async function loadMatches() {
  loading.value = true
  try {
    const res = await apiAdminListMatches({ page: 1, page_size: 100 }) as any
    matches.value = (res.list || res || []).map(normalizeMatch)
  } finally {
    loading.value = false
  }
}

async function syncLineup(matchId: number) {
  syncingLineupId.value = matchId
  syncMessage.value = ''
  try {
    const result = await apiAdminSyncMatchLineups(matchId)
    syncMessage.value = `阵容同步完成：更新 ${result.updated || 0}，可用 ${result.available || 0}`
  } catch (err: any) {
    syncMessage.value = err?.message || '阵容同步失败'
  } finally {
    syncingLineupId.value = null
  }
}

onMounted(loadMatches)
</script>

<template>
  <div class="admin-page">
    <div class="admin-head">
      <div>
        <h2>比赛管理</h2>
        <span>{{ filteredMatches.length }} 场比赛</span>
      </div>
      <div class="tools">
        <input v-model="search" class="admin-search" placeholder="搜索球队 / 城市 / 球场" />
        <select v-model="status" class="admin-select">
          <option value="">全部状态</option>
          <option value="scheduled">未开始</option>
          <option value="live">进行中</option>
          <option value="finished">已结束</option>
        </select>
      </div>
    </div>

    <div v-if="syncMessage" class="card result-card">{{ syncMessage }}</div>

    <div class="card table-card table-scroll">
      <table class="admin-table">
        <thead>
          <tr>
            <th>比赛</th><th>小组</th><th>时间</th><th>城市</th><th>球场</th><th>状态</th><th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="match in filteredMatches" :key="match.id">
            <td><b>{{ match.home_team_name }} vs {{ match.away_team_name }}</b></td>
            <td>{{ match.group_name || '-' }}</td>
            <td>{{ match.local_kickoff_time || '-' }}</td>
            <td>{{ match.city || '-' }}</td>
            <td>{{ match.stadium || '-' }}</td>
            <td>{{ match.status }}</td>
            <td>
              <button class="mini-btn" :disabled="syncingLineupId === match.id" @click="syncLineup(match.id)">
                {{ syncingLineupId === match.id ? '同步中' : '同步阵容' }}
              </button>
            </td>
          </tr>
          <tr v-if="!loading && !filteredMatches.length">
            <td colspan="7" class="empty-row">暂无比赛数据</td>
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

.admin-head,
.tools {
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

.admin-search,
.admin-select {
  min-height: 40px;
  padding: 0 14px;
  border: 1px solid var(--line);
  border-radius: 8px;
  color: var(--text);
  background: var(--card);
}

.admin-search {
  width: min(320px, 100%);
}

.table-card {
  overflow: hidden;
}

.result-card {
  padding: 12px 14px;
  color: var(--muted);
  font-size: 13px;
}

.mini-btn {
  min-height: 30px;
  padding: 0 10px;
  border: 1px solid var(--line);
  border-radius: 8px;
  color: var(--text);
  background: var(--card-soft);
  white-space: nowrap;
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
