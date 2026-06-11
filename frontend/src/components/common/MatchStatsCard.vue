<script setup lang="ts">
import { computed } from 'vue'
import type { Match } from '@/types/match'

const props = defineProps<{
  match: Match
}>()

interface StatItem {
  label: string
  home: number | null | undefined
  away: number | null | undefined
  showBar?: boolean
  barHome?: number
  barAway?: number
}

const possessionHome = computed(() => props.match.home_possession ?? null)
const possessionAway = computed(() => props.match.away_possession ?? null)

const stats = computed<StatItem[]>(() => {
  const m = props.match
  const items: StatItem[] = []

  if (possessionHome.value != null && possessionAway.value != null) {
    const total = possessionHome.value + possessionAway.value
    items.push({
      label: '控球率',
      home: possessionHome.value,
      away: possessionAway.value,
      showBar: true,
      barHome: total > 0 ? Math.round((possessionHome.value / total) * 100) : 50,
      barAway: total > 0 ? Math.round((possessionAway.value / total) * 100) : 50,
    })
  }
  if (m.home_shots != null || m.away_shots != null) {
    items.push({ label: '射门', home: m.home_shots, away: m.away_shots })
  }
  if (m.home_shots_on_target != null || m.away_shots_on_target != null) {
    items.push({ label: '射正', home: m.home_shots_on_target, away: m.away_shots_on_target })
  }
  if (m.home_corners != null || m.away_corners != null) {
    items.push({ label: '角球', home: m.home_corners, away: m.away_corners })
  }
  if (m.home_offsides != null || m.away_offsides != null) {
    items.push({ label: '越位', home: m.home_offsides, away: m.away_offsides })
  }
  if (m.home_fouls != null || m.away_fouls != null) {
    items.push({ label: '犯规', home: m.home_fouls, away: m.away_fouls })
  }
  if (m.home_yellow_cards != null || m.away_yellow_cards != null) {
    items.push({ label: '黄牌', home: m.home_yellow_cards, away: m.away_yellow_cards })
  }
  if (m.home_red_cards != null || m.away_red_cards != null) {
    items.push({ label: '红牌', home: m.home_red_cards, away: m.away_red_cards })
  }

  return items
})
</script>

<template>
  <article v-if="stats.length" class="stats-card">
    <div class="stats-head">
      <span class="material-symbols-outlined">bar_chart</span>
      <span>比赛数据</span>
    </div>

    <div class="stats-table">
      <div class="stats-row stats-header">
        <div class="stats-cell home">{{ match.home_team_name }}</div>
        <div class="stats-cell label"></div>
        <div class="stats-cell away">{{ match.away_team_name }}</div>
      </div>

      <div v-for="item in stats" :key="item.label" class="stats-row">
        <div class="stats-cell home">
          <template v-if="item.showBar">
            <div class="bar-wrap bar-home">
              <div class="bar-fill" :style="{ width: item.barHome + '%' }"></div>
            </div>
            <span>{{ item.home ?? '-' }}%</span>
          </template>
          <template v-else>
            <span :class="{ highlight: item.home != null && item.away != null && item.home > item.away }">
              {{ item.home ?? '-' }}
            </span>
          </template>
        </div>
        <div class="stats-cell label">{{ item.label }}</div>
        <div class="stats-cell away">
          <template v-if="item.showBar">
            <span>{{ item.away ?? '-' }}%</span>
            <div class="bar-wrap bar-away">
              <div class="bar-fill" :style="{ width: item.barAway + '%' }"></div>
            </div>
          </template>
          <template v-else>
            <span :class="{ highlight: item.away != null && item.home != null && item.away > item.home }">
              {{ item.away ?? '-' }}
            </span>
          </template>
        </div>
      </div>
    </div>
  </article>
</template>

<style scoped>
.stats-card {
  display: grid;
  gap: 14px;
  padding-top: 20px;
  border-top: 1px solid color-mix(in srgb, var(--line) 78%, transparent);
}

.stats-head {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 750;
  color: var(--muted);
}

.stats-head .material-symbols-outlined {
  font-size: 18px;
}

.stats-table {
  display: grid;
  border-top: 1px solid color-mix(in srgb, var(--line) 76%, transparent);
}

.stats-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  align-items: center;
  gap: 10px;
  padding: 11px 0;
  border-bottom: 1px solid color-mix(in srgb, var(--line) 70%, transparent);
  background: transparent;
}

.stats-header {
  font-size: 12px;
  font-weight: 750;
  color: var(--muted);
  padding-bottom: 7px;
}

.stats-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 650;
  font-variant-numeric: tabular-nums;
}

.stats-cell.home {
  justify-content: flex-start;
}

.stats-cell.away {
  justify-content: flex-end;
}

.stats-cell.label {
  color: var(--weak);
  font-size: 11px;
  font-weight: 700;
  justify-content: center;
  min-width: 36px;
  white-space: nowrap;
}

.highlight {
  font-weight: 850;
  color: var(--text);
}

.bar-wrap {
  width: 48px;
  height: 6px;
  border-radius: 999px;
  background: var(--line);
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: inherit;
  background: var(--primary);
}

.bar-wrap.bar-away .bar-fill {
  float: right;
}
</style>
