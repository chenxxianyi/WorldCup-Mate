<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import type { Match } from '@/types/match'
import ReminderControl from '@/components/common/ReminderControl.vue'
import TeamFlag from '@/components/common/TeamFlag.vue'

const props = defineProps<{
  match: Match
  featured?: boolean
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
  const [date, time] = localTime.split(' ')
  if (!date || !time) return localTime
  return time
}

function kickoffDateText(match: Match) {
  const localTime = match.local_kickoff_time || ''
  const [date] = localTime.split(' ')
  if (!date) return ''
  const [, month, day] = date.split('-')
  return month && day ? `${month}月${day}日` : date
}

function goDetail() {
  router.push(`/matches/${props.match.id}`)
}
</script>

<template>
  <article class="card match-card" :class="{ featured: featured || match.is_featured }">
    <div v-if="(featured || match.is_featured) && kickoffDateText(match)" class="date-ribbon">
      {{ kickoffDateText(match) }}
    </div>
    <div v-if="featured || match.is_featured" class="hot-ribbon">热门比赛</div>
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

    <div class="teams-line" @click="goDetail">
      <div class="team-side">
        <TeamFlag :value="match.home_flag" :alt="match.home_team_name" :fallback="match.home_team_code" />
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
        <TeamFlag :value="match.away_flag" :alt="match.away_team_name" :fallback="match.away_team_code" />
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
          @click.stop="fav.toggleMatchFavorite(match.id)"
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
.match-card {
  position: relative;
  overflow: visible;
  padding: 15px;
}

.match-card.featured {
  border-color: color-mix(in srgb, var(--hot) 30%, var(--line));
  margin-top: 14px;
  padding-top: 34px;
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
  margin: 18px 0 16px;
  cursor: pointer;
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
  height: 42px;
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
  padding-top: 12px;
  border-top: 1px solid var(--line);
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
