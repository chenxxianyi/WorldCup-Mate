<script setup lang="ts">
import MatchTicketRow from '@/components/common/MatchTicketRow.vue'
import type { Match } from '@/types/match'

defineProps<{
  match: Match
}>()

function kickoffDateText(match: Match) {
  const localTime = match.local_kickoff_time || ''
  const [date] = localTime.split(' ')
  if (!date) return ''
  const [, month, day] = date.split('-')
  return month && day ? `${month}月${day}日` : date
}
</script>

<template>
  <article class="featured-match-panel">
    <div v-if="kickoffDateText(match)" class="date-ribbon">
      {{ kickoffDateText(match) }}
    </div>
    <div class="hot-ribbon">热门比赛</div>
    <MatchTicketRow :match="match" />
  </article>
</template>

<style scoped>
.featured-match-panel {
  position: relative;
  margin-top: 14px;
  padding: 34px 15px 15px;
  border: 1px solid color-mix(in srgb, var(--hot) 30%, var(--line));
  border-radius: var(--radius-lg);
  background:
    linear-gradient(180deg, color-mix(in srgb, var(--hot) 4%, var(--card)), var(--card));
  box-shadow: var(--shadow);
}

.featured-match-panel :deep(.match-ticket-row) {
  padding: 0;
  border-left: 0;
}

.featured-match-panel :deep(.match-ticket-row::before) {
  content: none;
}

.featured-match-panel :deep(.teams-line) {
  margin: 18px 0 16px;
}

.date-ribbon {
  position: absolute;
  top: -14px;
  left: 0;
  min-height: 28px;
  display: inline-flex;
  align-items: center;
  padding: 0 18px;
  border: 1px solid color-mix(in srgb, var(--green-2) 16%, var(--line));
  border-bottom: 0;
  border-radius: 14px 14px 0 0;
  color: #13543a;
  font-size: 12px;
  font-weight: 850;
  background: color-mix(in srgb, var(--green-2) 14%, #fff);
  box-shadow: 0 -4px 12px rgba(0, 104, 71, 0.06);
}

.hot-ribbon {
  position: absolute;
  top: 0;
  right: 0;
  min-height: 34px;
  display: inline-flex;
  align-items: center;
  padding: 0 12px;
  border-radius: 0 0 0 12px;
  color: #fff;
  font-size: 10px;
  font-weight: 800;
  background: var(--hot);
}
</style>
