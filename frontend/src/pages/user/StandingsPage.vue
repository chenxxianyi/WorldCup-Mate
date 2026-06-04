<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import ChipFilter from '@/components/common/ChipFilter.vue'
import GroupAnalysisCard from '@/components/ai/GroupAnalysisCard.vue'
import StandingTable from '@/components/common/StandingTable.vue'
import StatCard from '@/components/common/StatCard.vue'
import SyncStatusBadge from '@/components/common/SyncStatusBadge.vue'
import TeamFlag from '@/components/common/TeamFlag.vue'
import { apiGetBestThird, apiGetGroupStandings } from '@/api/standings'
import { useAIStore } from '@/stores/useAIStore'
import { normalizeStanding, type Standing } from '@/types/standing'

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

const ai = useAIStore()
const activeGroup = ref('Group A')
const currentStandings = ref<Standing[]>([])
const bestThird = ref<any[]>([])
const showBestThird = computed(() => activeGroup.value === BEST_THIRD)
const activeGroupId = computed(() => {
  if (!activeGroup.value.startsWith('Group ')) return 0
  return activeGroup.value.charCodeAt(6) - 64
})

async function loadStandings() {
  if (showBestThird.value) {
    try {
      const res = await apiGetBestThird() as any[]
      bestThird.value = res
    } catch {
      bestThird.value = []
    }
    return
  }

  try {
    const res = await apiGetGroupStandings(activeGroupId.value) as any[]
    currentStandings.value = res.map(normalizeStanding)
  } catch {
    currentStandings.value = []
  }
}

function generateGroupAnalysis(forceRefresh = false) {
  if (!activeGroupId.value) return
  ai.generateGroupAnalysis(activeGroupId.value, forceRefresh).catch(() => {})
}

onMounted(loadStandings)
watch(activeGroup, loadStandings)
</script>

<template>
  <div class="standings-page">
    <div class="section-head">
      <div>
        <h2>小组积分榜</h2>
        <span>Group A - Group L</span>
      </div>
      <SyncStatusBadge />
    </div>

    <ChipFilter v-model="activeGroup" :options="groupOptions" />

    <template v-if="!showBestThird">
      <StandingTable :standings="currentStandings" show-status />
      <GroupAnalysisCard
        :analysis="ai.groupAnalysisMap[activeGroupId] || null"
        :loading="ai.groupLoadingMap[activeGroupId]"
        :error="ai.groupErrorMap[activeGroupId]"
        @generate="generateGroupAnalysis(false)"
        @refresh="generateGroupAnalysis(true)"
      />
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
                <th>排名</th>
                <th>球队</th>
                <th>小组</th>
                <th>积分</th>
                <th>净胜球</th>
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
  </div>
</template>

<style scoped>
.standings-page {
  display: grid;
  gap: 14px;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
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
  margin-top: 4px;
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
</style>
