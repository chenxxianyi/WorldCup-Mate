<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import ChipFilter from '@/components/common/ChipFilter.vue'
import StandingTable from '@/components/common/StandingTable.vue'
import StatCard from '@/components/common/StatCard.vue'
import { apiGetGroupStandings, apiGetAllStandings, apiGetBestThird } from '@/api/standings'
import { normalizeStanding, type Standing } from '@/types/standing'
import TeamFlag from '@/components/common/TeamFlag.vue'

const groupOptions = ['Group A', 'Group B', 'Group C', 'Group D', '最佳第三名']
const activeGroup = ref('Group A')
const currentStandings = ref<Standing[]>([])
const bestThird = ref<any[]>([])
const showBestThird = computed(() => activeGroup.value === '最佳第三名')

async function loadStandings() {
  if (activeGroup.value === '最佳第三名') {
    try {
      const res = await apiGetBestThird() as any[]
      bestThird.value = res
    } catch { bestThird.value = [] }
  } else {
    const groupNum = activeGroup.value.charCodeAt(6) - 64
    try {
      const res = await apiGetGroupStandings(groupNum) as any[]
      currentStandings.value = res.map(normalizeStanding)
    } catch { currentStandings.value = [] }
  }
}

onMounted(loadStandings)
watch(activeGroup, loadStandings)
</script>

<template>
  <div>
    <div class="section-head">
      <div>
        <h2>小组积分榜</h2>
        <span>Group A - Group L</span>
      </div>
    </div>
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

th, td {
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

.rank-ok td:first-child {
  border-left: 4px solid var(--success);
}

.rank-mid td:first-child {
  border-left: 4px solid var(--gold);
}

.rank-out {
  color: var(--weak);
  opacity: 0.66;
}
</style>
