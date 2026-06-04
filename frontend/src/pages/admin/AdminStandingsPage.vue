<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { apiAdminListStandings, apiAdminRecalculateStandings } from '@/api/admin'

const standings = ref<any[]>([])
const groupName = ref('')
const loading = ref(false)
const recalculating = ref(false)

const filteredStandings = computed(() => {
  if (!groupName.value) return standings.value
  return standings.value.filter((item) => item.group?.name === groupName.value)
})

const groupOptions = computed(() =>
  Array.from(new Set(standings.value.map((item) => item.group?.name).filter(Boolean))),
)

async function loadStandings() {
  loading.value = true
  try {
    const res = await apiAdminListStandings({ page: 1, page_size: 100 }) as any
    standings.value = res.list || res || []
  } finally {
    loading.value = false
  }
}

async function recalculate() {
  recalculating.value = true
  try {
    await apiAdminRecalculateStandings()
    await loadStandings()
  } finally {
    recalculating.value = false
  }
}

onMounted(loadStandings)
</script>

<template>
  <div class="admin-page">
    <div class="admin-head">
      <div>
        <h2>积分榜管理</h2>
        <span>{{ filteredStandings.length }} 条记录</span>
      </div>
      <div class="tools">
        <select v-model="groupName" class="admin-select">
          <option value="">全部小组</option>
          <option v-for="name in groupOptions" :key="name" :value="name">{{ name }}</option>
        </select>
        <button class="pill-btn primary" :disabled="recalculating" @click="recalculate">
          {{ recalculating ? '重算中...' : '重算积分' }}
        </button>
      </div>
    </div>

    <div class="card table-card table-scroll">
      <table class="admin-table">
        <thead>
          <tr>
            <th>小组</th><th>球队</th><th>场次</th><th>胜</th><th>平</th><th>负</th><th>净胜球</th><th>积分</th><th>状态</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in filteredStandings" :key="item.id">
            <td>{{ item.group?.name || '-' }}</td>
            <td><b>{{ item.team?.name || '-' }}</b></td>
            <td>{{ item.played }}</td>
            <td>{{ item.won }}</td>
            <td>{{ item.drawn }}</td>
            <td>{{ item.lost }}</td>
            <td>{{ item.goal_difference }}</td>
            <td><b>{{ item.points }}</b></td>
            <td>{{ item.qualification_status || '-' }}</td>
          </tr>
          <tr v-if="!loading && !filteredStandings.length">
            <td colspan="9" class="empty-row">暂无积分数据</td>
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

.admin-select {
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
