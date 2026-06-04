<script setup lang="ts">
import type { Standing } from '@/types/standing'
import TeamFlag from '@/components/common/TeamFlag.vue'

defineProps<{
  standings: Standing[]
  showStatus?: boolean
}>()

function rowClass(status: string) {
  if (status === '晋级') return 'rank-ok'
  if (status === '待定') return 'rank-mid'
  return 'rank-out'
}
</script>

<template>
  <div class="card table-card">
    <div class="table-scroll">
      <table class="standing-table">
        <thead>
          <tr>
            <th>排名</th>
            <th>球队</th>
            <th>场次</th>
            <th>胜</th>
            <th>平</th>
            <th>负</th>
            <th>净胜球</th>
            <th>积分</th>
            <th v-if="showStatus">状态</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(s, i) in standings" :key="s.team_id" :class="rowClass(s.status)">
            <td>{{ i + 1 }}</td>
            <td class="team-cell">
              <TeamFlag :value="s.flag" :alt="s.team_name" :fallback="s.team_code" size="sm" />
              <span>{{ s.team_name }}</span>
            </td>
            <td>{{ s.played }}</td>
            <td>{{ s.won }}</td>
            <td>{{ s.drawn }}</td>
            <td>{{ s.lost }}</td>
            <td>{{ s.goal_difference > 0 ? '+' : '' }}{{ s.goal_difference }}</td>
            <td><b>{{ s.points }}</b></td>
            <td v-if="showStatus">
              <span
                class="tag"
                :class="{
                  green: s.status === '晋级',
                  gold: s.status === '待定',
                }"
              >
                {{ s.status }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
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
