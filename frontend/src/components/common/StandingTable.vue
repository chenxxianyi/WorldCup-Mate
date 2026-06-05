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
  overflow-x: visible;
}

.standing-table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
  min-width: 0;
}

th,
td {
  padding: 12px 6px;
  border-bottom: 1px solid var(--line);
  text-align: center;
  font-size: 12px;
  white-space: nowrap;
}

th {
  color: var(--muted);
  font-weight: 750;
  background: var(--card-soft);
}

td {
  color: var(--text);
}

th:nth-child(1),
td:nth-child(1) {
  width: 42px;
}

th:nth-child(2),
td:nth-child(2) {
  width: auto;
  text-align: left;
}

th:nth-child(3),
td:nth-child(3),
th:nth-child(4),
td:nth-child(4),
th:nth-child(5),
td:nth-child(5),
th:nth-child(6),
td:nth-child(6),
th:nth-child(8),
td:nth-child(8) {
  width: 38px;
}

th:nth-child(7),
td:nth-child(7) {
  width: 52px;
}

th:nth-child(9),
td:nth-child(9) {
  width: 58px;
}

.team-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}

.team-cell span {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

td .tag {
  max-width: 100%;
  padding: 4px 7px;
  font-size: 11px;
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

@media (max-width: 430px) {
  th,
  td {
    padding: 10px 4px;
    font-size: 11px;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 34px;
  }

  th:nth-child(3),
  td:nth-child(3),
  th:nth-child(4),
  td:nth-child(4),
  th:nth-child(5),
  td:nth-child(5),
  th:nth-child(6),
  td:nth-child(6),
  th:nth-child(8),
  td:nth-child(8) {
    width: 30px;
  }

  th:nth-child(7),
  td:nth-child(7) {
    width: 42px;
  }

  th:nth-child(9),
  td:nth-child(9) {
    width: 46px;
  }

  .team-cell {
    gap: 4px;
  }

  td .tag {
    padding: 3px 5px;
    font-size: 10px;
  }
}
</style>
