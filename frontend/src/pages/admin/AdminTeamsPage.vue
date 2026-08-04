<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { apiAdminListTeams } from '@/api/admin'
import type { ApiTeam } from '@/types/team'

const search = ref('')
const teamType = ref('')
const teams = ref<ApiTeam[]>([])
const loading = ref(false)

const filteredTeams = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return teams.value
  return teams.value.filter((team) =>
    [team.name, team.name_en, team.fifa_code, team.continent, team.group?.name]
      .filter(Boolean)
      .some((value) => String(value).toLowerCase().includes(q)),
  )
})

async function loadTeams() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: 1, page_size: 100 }
    if (teamType.value) params.teamType = teamType.value
    const res = await apiAdminListTeams(params)
    teams.value = res.list
  } finally {
    loading.value = false
  }
}

onMounted(loadTeams)
watch(teamType, loadTeams)
watch(search, () => {})
</script>

<template>
  <div class="admin-page">
    <div class="admin-head">
      <div>
        <h2>球队管理</h2>
        <span>{{ filteredTeams.length }} 支球队</span>
      </div>
      <div class="tools">
        <select v-model="teamType" class="admin-select" aria-label="球队类型">
          <option value="">全部类型</option>
          <option value="national">国家队</option>
          <option value="club">俱乐部</option>
        </select>
        <input v-model="search" class="admin-search" placeholder="搜索球队 / 小组 / 大洲" />
      </div>
    </div>

    <div class="card table-card table-scroll">
      <table class="admin-table">
        <thead>
          <tr>
            <th>球队</th><th>代码</th><th>英文名</th><th>大洲</th><th>小组</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="team in filteredTeams" :key="team.id">
            <td><b>{{ team.name }}</b></td>
            <td>{{ team.fifa_code }}</td>
            <td>{{ team.name_en }}</td>
            <td>{{ team.continent }}</td>
            <td>{{ team.group?.name || '-' }}</td>
          </tr>
          <tr v-if="!loading && !filteredTeams.length">
            <td colspan="5" class="empty-row">暂无球队数据</td>
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

.tools {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  flex-wrap: wrap;
}

.admin-head h2 {
  margin: 0;
  font-size: 20px;
}

.admin-head span {
  color: var(--muted);
  font-size: 13px;
}

.admin-search {
  width: min(320px, 100%);
  min-height: 40px;
  padding: 0 14px;
  border: 1px solid var(--line);
  border-radius: 8px;
  color: var(--text);
  background: var(--card);
}

.table-card {
  overflow: hidden;
}

.table-scroll {
  overflow-x: auto;
}

.admin-table {
  width: 100%;
  min-width: 640px;
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
