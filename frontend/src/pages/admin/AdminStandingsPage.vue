<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { apiAdminListStandings, apiAdminRecalculateStandings, apiAdminRecalculateLeagueStanding } from '@/api/admin'
import { apiListCompetitions, apiGetCompetitionStandings } from '@/api/competitions'
import type { ApiStanding } from '@/types/standing'
import type { ApiLeagueStanding } from '@/types/standing'
import type { Competition } from '@/types/competition'

const standings = ref<ApiStanding[]>([])
const groupName = ref('')
const loading = ref(false)
const recalculating = ref(false)

// League mode: switch between group tables (World Cup) and league tables.
const mode = ref<'group' | 'league'>('group')
const leagueCode = ref('')
const leagueStandings = ref<ApiLeagueStanding[]>([])
const competitions = ref<Competition[]>([])
const leagueLoading = ref(false)

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
    const res = await apiAdminListStandings({ page: 1, page_size: 100 })
    standings.value = res || []
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

async function loadLeagueStandings() {
  if (!leagueCode.value) {
    leagueStandings.value = []
    return
  }
  leagueLoading.value = true
  try {
    const res = await apiGetCompetitionStandings(leagueCode.value, { type: 'total' })
    leagueStandings.value = res || []
  } catch {
    leagueStandings.value = []
  } finally {
    leagueLoading.value = false
  }
}

async function recalculateLeague() {
  if (!leagueCode.value) return
  const competition = competitions.value.find((c) => c.code === leagueCode.value)
  if (!competition) return
  recalculating.value = true
  try {
    await apiAdminRecalculateLeagueStanding({ competition_id: competition.id })
    await loadLeagueStandings()
  } finally {
    recalculating.value = false
  }
}

async function loadCompetitions() {
  try {
    competitions.value = await apiListCompetitions()
  } catch {
    competitions.value = []
  }
}

onMounted(() => {
  loadStandings()
  loadCompetitions()
})
watch(mode, () => {
  if (mode.value === 'league') loadLeagueStandings()
})
watch(leagueCode, loadLeagueStandings)
</script>

<template>
  <div class="admin-page">
    <div class="admin-head">
      <div>
        <h2>积分榜管理</h2>
        <span>{{ mode === 'group' ? filteredStandings.length + ' 条记录' : leagueStandings.length + ' 条记录' }}</span>
      </div>
      <div class="tools">
        <select v-model="mode" class="admin-select" aria-label="积分榜类型">
          <option value="group">小组积分榜（世界杯）</option>
          <option value="league">联赛积分榜</option>
        </select>
        <select v-if="mode === 'group'" v-model="groupName" class="admin-select">
          <option value="">全部小组</option>
          <option v-for="name in groupOptions" :key="name" :value="name">{{ name }}</option>
        </select>
        <template v-else>
          <select v-model="leagueCode" class="admin-select" aria-label="选择联赛">
            <option value="">选择联赛</option>
            <option v-for="c in competitions" :key="c.code" :value="c.code">{{ c.name }}</option>
          </select>
          <button class="pill-btn primary" :disabled="recalculating || !leagueCode" @click="recalculateLeague">
            {{ recalculating ? '重算中...' : '重算联赛积分' }}
          </button>
        </template>
        <button v-if="mode === 'group'" class="pill-btn primary" :disabled="recalculating" @click="recalculate">
          {{ recalculating ? '重算中...' : '重算积分' }}
        </button>
      </div>
    </div>

    <template v-if="mode === 'group'">
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
    </template>

    <template v-else>
      <div class="card table-card table-scroll">
        <table class="admin-table">
          <thead>
            <tr>
              <th>排名</th><th>球队</th><th>类型</th><th>场次</th><th>胜</th><th>平</th><th>负</th><th>净胜球</th><th>积分</th><th>分区</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in leagueStandings" :key="item.id">
              <td>{{ item.position }}</td>
              <td><b>{{ item.team?.name || '-' }}</b></td>
              <td>{{ item.type }}</td>
              <td>{{ item.played }}</td>
              <td>{{ item.won }}</td>
              <td>{{ item.drawn }}</td>
              <td>{{ item.lost }}</td>
              <td>{{ item.goal_difference }}</td>
              <td><b>{{ item.points }}</b></td>
              <td>{{ item.zone || '-' }}</td>
            </tr>
            <tr v-if="!leagueLoading && !leagueStandings.length">
              <td colspan="10" class="empty-row">请选择联赛查看积分榜</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
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
