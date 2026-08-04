<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import ChipFilter from '@/components/common/ChipFilter.vue'
import StandingTable from '@/components/common/StandingTable.vue'
import StatCard from '@/components/common/StatCard.vue'
import SyncStatusBadge from '@/components/common/SyncStatusBadge.vue'
import TeamFlag from '@/components/common/TeamFlag.vue'
import { apiGetBestThird, apiGetGroupStandings } from '@/api/standings'
import { apiGetCompetitionStandings } from '@/api/competitions'
import { useCompetitionStore } from '@/stores/useCompetitionStore'
import { normalizeStanding, normalizeLeagueStanding, type Standing, type LeagueStanding } from '@/types/standing'

const BEST_THIRD = '最佳第三名'
const groupOptions = [
  'Group A',
  'Group B',
  'Group C',
  'Group D',
  'Group E',
  'Group F',
  'Group G',
  'Group H',
  'Group I',
  'Group J',
  'Group K',
  'Group L',
  BEST_THIRD,
]

const TYPE_ALL = '总榜'
const TYPE_HOME = '主场'
const TYPE_AWAY = '客场'
const typeOptions = [TYPE_ALL, TYPE_HOME, TYPE_AWAY]

const comp = useCompetitionStore()
const activeGroup = ref('Group A')
const activeType = ref(TYPE_ALL)
const currentStandings = ref<Standing[]>([])
const leagueStandings = ref<LeagueStanding[]>([])
const bestThird = ref<any[]>([])
const showBestThird = computed(() => activeGroup.value === BEST_THIRD)
const isLeagueMode = computed(() => comp.isLeague)

function typeParam() {
  if (activeType.value === TYPE_HOME) return 'home'
  if (activeType.value === TYPE_AWAY) return 'away'
  return 'total'
}

function zoneLabel(zone: string) {
  const labels: Record<string, string> = {
    champions_league: '欧冠区',
    europa_league: '欧战区',
    relegation: '降级区',
  }
  return labels[zone] || ''
}

function zoneClass(zone: string) {
  return zone || ''
}

async function loadLeagueStandings() {
  if (!comp.current) return
  try {
    const res = await apiGetCompetitionStandings(comp.currentCode, { type: typeParam() }) as any[]
    leagueStandings.value = (res || []).map(normalizeLeagueStanding)
  } catch {
    leagueStandings.value = []
  }
}

async function loadStandings() {
  if (isLeagueMode.value) {
    await loadLeagueStandings()
    return
  }
  if (showBestThird.value) {
    try {
      const res = await apiGetBestThird() as any[]
      bestThird.value = res
    } catch {
      bestThird.value = []
    }
    return
  }

  const groupNum = activeGroup.value.charCodeAt(6) - 64
  try {
    const res = await apiGetGroupStandings(groupNum) as any[]
    currentStandings.value = res.map(normalizeStanding)
  } catch {
    currentStandings.value = []
  }
}

onMounted(async () => {
  await comp.fetchCompetitions()
  loadStandings()
})
watch(activeGroup, loadStandings)
watch(activeType, loadStandings)
watch(() => comp.currentCode, () => {
  activeGroup.value = 'Group A'
  activeType.value = TYPE_ALL
  loadStandings()
})
</script>

<template>
  <div>
    <div class="section-head">
      <div>
        <h2>{{ isLeagueMode ? '联赛积分榜' : '小组积分榜' }}</h2>
        <span>{{ isLeagueMode ? comp.current?.name || '' : 'Group A - Group L' }}</span>
      </div>
      <SyncStatusBadge />
    </div>

    <template v-if="isLeagueMode">
      <ChipFilter v-model="activeType" :options="typeOptions" />
      <div class="card table-card table-scroll" style="margin-top: 12px">
        <table class="standing-table">
          <thead>
            <tr>
              <th>排名</th><th>球队</th><th>场次</th><th>胜</th><th>平</th><th>负</th><th>进/失</th><th>净胜球</th><th>积分</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in leagueStandings" :key="s.team_id">
              <td class="pos-cell" :class="zoneClass(s.zone)">
                {{ s.position }}
                <span v-if="zoneLabel(s.zone)" class="zone-tag">{{ zoneLabel(s.zone) }}</span>
              </td>
              <td class="team-cell">
                <TeamFlag :value="s.flag" :alt="s.team_name" :fallback="s.team_code" size="sm" />
                <span>{{ s.team_name }}</span>
              </td>
              <td>{{ s.played }}</td>
              <td>{{ s.won }}</td>
              <td>{{ s.drawn }}</td>
              <td>{{ s.lost }}</td>
              <td>{{ s.goals_for }}/{{ s.goals_against }}</td>
              <td>{{ s.goal_difference > 0 ? '+' : '' }}{{ s.goal_difference }}</td>
              <td><b>{{ s.points }}</b></td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

    <template v-else>
      <ChipFilter v-model="activeGroup" :options="groupOptions" />

      <template v-if="!showBestThird">
        <StandingTable :standings="currentStandings" show-status />
      </template>

      <template v-else>
        <section class="section">
          <div class="section-head">
            <h2>最佳第三名</h2>
            <span>前 8 名晋级</span>
          </div>
          <div class="stats-row">
            <StatCard :value="8" label="晋级名额" />
            <StatCard :value="12" label="候选球队" />
            <StatCard :value="4" label="待定/淘汰" />
          </div>
          <div class="card table-card table-scroll" style="margin-top: 12px">
            <table class="standing-table">
              <thead>
                <tr>
                  <th>排名</th><th>球队</th><th>小组</th><th>积分</th><th>净胜球</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(t, i) in bestThird" :key="t.team_id || i" :class="{ 'rank-mid': i >= 8 }">
                  <td>{{ i + 1 }}</td>
                  <td class="team-cell">
                    <TeamFlag :value="t.team?.flag_url || ''" :alt="t.team?.name || ''" :fallback="t.team?.fifa_code || ''" size="sm" />
                    <span>{{ t.team?.name || '' }}</span>
                  </td>
                  <td>{{ t.group?.name || '' }}</td>
                  <td><b>{{ t.points }}</b></td>
                  <td>{{ t.goal_difference > 0 ? '+' : '' }}{{ t.goal_difference }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </template>
    </template>
  </div>
</template>

<style scoped>
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

.section-head span {
  color: var(--muted);
  font-size: 13px;
}

.section {
  margin-top: 18px;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.table-card {
  overflow: hidden;
}

.table-scroll {
  overflow-x: auto;
}

.standing-table {
  width: 100%;
  border-collapse: collapse;
  min-width: 520px;
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
  font-weight: 750;
  background: var(--card-soft);
}

td {
  color: var(--text);
}

.team-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 150px;
}

tr:last-child td {
  border-bottom: 0;
}

.rank-mid td:first-child {
  border-left: 4px solid var(--gold);
}

.pos-cell {
  border-left: 4px solid transparent;
  white-space: nowrap;
}

.pos-cell.champions_league {
  border-left-color: #2f6fed;
}

.pos-cell.europa_league {
  border-left-color: #12a150;
}

.pos-cell.relegation {
  border-left-color: var(--primary);
}

.zone-tag {
  display: inline-block;
  margin-left: 6px;
  padding: 1px 6px;
  border-radius: 999px;
  color: #fff;
  background: #8a8f98;
  font-size: 10px;
  font-weight: 750;
}

.pos-cell.champions_league .zone-tag {
  background: #2f6fed;
}

.pos-cell.europa_league .zone-tag {
  background: #12a150;
}

.pos-cell.relegation .zone-tag {
  background: var(--primary);
}
</style>
