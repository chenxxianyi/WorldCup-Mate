<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import type { Match } from '@/types/match'
import ReminderControl from '@/components/common/ReminderControl.vue'
import TeamFlag from '@/components/common/TeamFlag.vue'

const props = defineProps<{
  match: Match
}>()

const router = useRouter()
const fav = useFavoriteStore()

function stageLabel(stage: Match['stage']) {
  const labels: Record<Match['stage'], string> = {
    group: '小组赛',
    group_stage: '小组赛',
    round_of_32: '32 强赛',
    round_of_16: '16 强赛',
    quarter_final: '四分之一决赛',
    semi_final: '半决赛',
    third_place: '季军赛',
    final: '决赛',
  }
  return labels[stage] || '世界杯比赛'
}

function stageTagText(match: Match) {
  const label = stageLabel(match.stage)
  if (match.group_name) return `${label} · ${match.group_name}`
  return label
}

function statusText(match: Match) {
  if (match.status === 'live') return `${match.minute || 0}' LIVE`
  if (match.status === 'finished') return '已结束'
  return '未开始'
}

function kickoffText(match: Match) {
  const localTime = match.local_kickoff_time || ''
  if (!localTime) return '时间待定'
  const [, time] = localTime.split(' ')
  return time || localTime
}

function goDetail() {
  router.push(`/matches/${props.match.id}`)
}

function goTeam(teamId: number) {
  if (!teamId) return
  router.push(`/teams/${teamId}`)
}
</script>

<template>
  <article
    class="match-ticket-row"
    role="link"
    tabindex="0"
    @click="goDetail"
    @keydown.enter.self="goDetail"
  >
    <div class="match-top">
      <span class="tag blue">{{ stageTagText(match) }}</span>
      <span
        class="tag"
        :class="{ live: match.status === 'live', green: match.status === 'finished' }"
      >
        <i v-if="match.status === 'live'" class="live-dot" style="background: #fff"></i>
        {{ statusText(match) }}
      </span>
    </div>

    <div class="teams-line">
      <div class="team-side">
        <button
          class="team-flag-action"
          type="button"
          :aria-label="`查看${match.home_team_name}详情`"
          @click.stop="goTeam(match.home_team_id)"
        >
          <TeamFlag :value="match.home_flag" :alt="match.home_team_name" :fallback="match.home_team_code" />
        </button>
        <div class="team-copy">
          <span class="team-name">{{ match.home_team_name }}</span>
          <span class="team-meta">{{ match.home_team_code }}</span>
        </div>
      </div>
      <div v-if="match.status === 'live' || match.status === 'finished'" class="score">
        {{ match.home_score }}-{{ match.away_score }}
      </div>
      <div v-else class="vs">VS</div>
      <div class="team-side away">
        <div class="team-copy">
          <span class="team-name">{{ match.away_team_name }}</span>
          <span class="team-meta">{{ match.away_team_code }}</span>
        </div>
        <button
          class="team-flag-action"
          type="button"
          :aria-label="`查看${match.away_team_name}详情`"
          @click.stop="goTeam(match.away_team_id)"
        >
          <TeamFlag :value="match.away_flag" :alt="match.away_team_name" :fallback="match.away_team_code" />
        </button>
      </div>
    </div>

    <div class="match-bottom">
      <div class="where">
        <template v-if="match.status === 'live'">
          进行中 · {{ match.minute || 0 }}'<br />{{ match.stadium || 'TBD' }}
        </template>
        <template v-else>
          {{ kickoffText(match) }} · {{ match.city || 'TBD' }}<br />{{ match.stadium || 'TBD' }}
        </template>
      </div>
      <div class="actions">
        <ReminderControl :match-id="match.id" />
        <button
          class="icon-action star"
          :class="{ active: fav.isMatchFavorite(match.id), ghost: !fav.isMatchFavorite(match.id) }"
          title="收藏"
          @click.stop="fav.toggleMatchFavorite(match.id, match)"
        >
          <span
            class="material-symbols-outlined"
            :style="fav.isMatchFavorite(match.id) ? 'font-variation-settings: \'FILL\' 1' : ''"
          >star</span>
        </button>
      </div>
    </div>
  </article>
</template>

<style scoped>
.match-ticket-row {
  position: relative;
  overflow: visible;
  padding: 13px 0 13px 12px;
  border-left: 2px solid color-mix(in srgb, var(--primary) 36%, transparent);
  border-radius: 0;
  background: transparent;
  cursor: pointer;
  transition: background 160ms ease-out, border-color 160ms ease-out;
}

.match-ticket-row::before {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 1px;
  background: var(--line);
}

.match-ticket-row:active {
  background: color-mix(in srgb, var(--primary) 4%, transparent);
}

.match-ticket-row:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--primary) 55%, transparent);
  outline-offset: 3px;
}

.match-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.teams-line {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  align-items: center;
  gap: 14px;
  margin: 14px 0 12px;
}

.team-side {
  min-width: 0;
  min-height: 54px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.team-side.away {
  justify-content: flex-end;
  text-align: right;
}

.team-flag-action {
  width: 38px;
  height: 38px;
  display: inline-grid;
  place-items: center;
  flex: 0 0 auto;
  padding: 0;
  border: 0;
  border-radius: 999px;
  color: inherit;
  background: transparent;
  cursor: pointer;
}

.team-flag-action:active {
  transform: scale(0.94);
}

.team-flag-action:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--primary) 55%, transparent);
  outline-offset: 3px;
}

.team-copy {
  min-width: 0;
  display: grid;
  gap: 5px;
}

.team-name {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 750;
  font-size: 16px;
}

.team-meta {
  display: block;
  color: var(--weak);
  font-size: 12px;
}

.score {
  min-width: 70px;
  text-align: center;
  font-weight: 850;
  font-size: 30px;
  line-height: 1;
}

.vs {
  min-width: 64px;
  height: 38px;
  display: grid;
  place-items: center;
  border-radius: 999px;
  color: var(--muted);
  font-size: 13px;
  font-weight: 800;
  background: var(--card-soft);
}

.match-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding-top: 10px;
  border-top: 1px solid color-mix(in srgb, var(--line) 74%, transparent);
}

.where {
  min-width: 0;
  color: var(--muted);
  font-size: 12px;
  line-height: 1.5;
}

.actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.icon-action {
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  border: 0;
  border-radius: 999px;
  color: var(--weak);
  background: transparent;
  cursor: pointer;
}

.icon-action.active {
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 10%, transparent);
}

.icon-action .material-symbols-outlined {
  font-size: 22px;
}
</style>
