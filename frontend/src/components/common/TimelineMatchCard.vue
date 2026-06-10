<script setup lang="ts">
import { useRouter } from 'vue-router'
import TeamFlag from '@/components/common/TeamFlag.vue'

const props = defineProps<{
  match: any
  first?: boolean
  last?: boolean
}>()

const router = useRouter()

function timeText() {
  const local = props.match.local_kickoff_time || ''
  return local.split(' ')[1] || ''
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
    class="tl-card"
    :class="{
      'tl-live': match.status === 'live',
      'tl-finished': match.status === 'finished',
      'tl-scheduled': match.status === 'scheduled' || match.status === 'upcoming',
      'tl-first': first,
      'tl-last': last,
    }"
    role="link"
    tabindex="0"
    @click="goDetail"
    @keydown.enter.self="goDetail"
  >
    <!-- 时间标注 -->
    <div class="tl-time">
      <template v-if="match.status === 'live'">
        <span class="tl-live-pulse"></span>
        <span class="tl-time-text">{{ match.live_minute || match.minute || 0 }}'</span>
      </template>
      <template v-else-if="match.status === 'finished'">
        <span class="material-symbols-outlined tl-done">check_circle</span>
        <span class="tl-time-text done">FT</span>
      </template>
      <template v-else>
        <span class="tl-time-text">{{ timeText() || 'TBD' }}</span>
      </template>
    </div>

    <!-- 主内容 -->
    <div class="tl-body">
      <div class="tl-teams">
        <div class="tl-team tl-home">
          <button
            class="tl-flag-action"
            type="button"
            :aria-label="`查看${match.home_team_name}详情`"
            @click.stop="goTeam(match.home_team_id)"
          >
            <TeamFlag :value="match.home_flag" :alt="match.home_team_name" :fallback="match.home_team_code" size="sm" />
          </button>
          <span class="tl-team-name">{{ match.home_team_name }}</span>
          <span v-if="match.status === 'live' || match.status === 'finished'" class="tl-score">
            {{ match.home_score ?? '-' }}
          </span>
        </div>
        <div class="tl-vs">
          <template v-if="match.status === 'live' || match.status === 'finished'">-</template>
          <template v-else>vs</template>
        </div>
        <div class="tl-team tl-away">
          <span v-if="match.status === 'live' || match.status === 'finished'" class="tl-score">
            {{ match.away_score ?? '-' }}
          </span>
          <span class="tl-team-name">{{ match.away_team_name }}</span>
          <button
            class="tl-flag-action"
            type="button"
            :aria-label="`查看${match.away_team_name}详情`"
            @click.stop="goTeam(match.away_team_id)"
          >
            <TeamFlag :value="match.away_flag" :alt="match.away_team_name" :fallback="match.away_team_code" size="sm" />
          </button>
        </div>
      </div>
      <div class="tl-meta">
        <span v-if="match.stage_label || match.group_name">{{ match.group_name || match.stage }}</span>
        <span v-if="match.status === 'live'" class="tl-status live">直播中</span>
        <span v-else-if="match.status === 'finished'" class="tl-status finished">已结束</span>
        <span v-else class="tl-status scheduled">未开始</span>
        <span v-if="match.status === 'finished' && match.has_post_match_summary" class="tl-summary-link">赛后摘要</span>
      </div>
    </div>
  </article>
</template>

<style scoped>
.tl-card {
  display: grid;
  grid-template-columns: 54px minmax(0, 1fr);
  gap: 12px;
  min-height: 64px;
  padding: 6px 0 8px;
  border: 0;
  border-radius: 0;
  background: transparent;
  cursor: pointer;
  transition: background 180ms ease;
}

.tl-card:active {
  transform: scale(0.995);
  background: color-mix(in srgb, var(--primary) 4%, transparent);
}

.tl-card:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--primary) 55%, transparent);
  outline-offset: 3px;
}

.tl-live {
  border-radius: 12px;
  background: linear-gradient(90deg, color-mix(in srgb, var(--live) 8%, transparent), transparent 62%);
}

.tl-time {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 750;
  color: var(--muted);
}

.tl-time::before {
  content: '';
  position: absolute;
  top: -12px;
  bottom: -14px;
  left: 50%;
  width: 1px;
  background: color-mix(in srgb, var(--line) 82%, transparent);
  transform: translateX(-50%);
}

.tl-first .tl-time::before {
  top: 50%;
}

.tl-last .tl-time::before {
  bottom: 50%;
}

.tl-time::after {
  content: '';
  position: absolute;
  left: 50%;
  top: 50%;
  width: 7px;
  height: 7px;
  border: 2px solid var(--card);
  border-radius: 999px;
  background: var(--weak);
  box-shadow: 0 0 0 1px var(--line);
  transform: translate(-50%, -50%);
}

.tl-time > * {
  position: relative;
  z-index: 1;
}

.tl-time-text {
  min-width: 40px;
  display: inline-grid;
  place-items: center;
  padding: 2px 4px;
  border-radius: 999px;
  background: var(--bg);
  line-height: 1.2;
}

.tl-time-text.done {
  color: var(--success);
  font-size: 10px;
  font-weight: 850;
}

.tl-live .tl-time {
  color: var(--live);
}

.tl-live .tl-time::after {
  background: var(--live);
  box-shadow: 0 0 0 4px color-mix(in srgb, var(--live) 14%, transparent);
}

.tl-finished .tl-time::after {
  background: var(--success);
}

.tl-live-pulse {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: var(--live);
  animation: tl-pulse 1.5s infinite;
}

@keyframes tl-pulse {
  0% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(1.2); }
  100% { opacity: 1; transform: scale(1); }
}

.tl-done {
  font-size: 18px;
  color: var(--success);
  margin-bottom: -2px;
}

.tl-body {
  min-width: 0;
  display: grid;
  gap: 6px;
  padding: 10px 0 12px;
  border-bottom: 1px solid color-mix(in srgb, var(--line) 78%, transparent);
}

.tl-last .tl-body {
  border-bottom-color: transparent;
}

.tl-teams {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  align-items: center;
  gap: 8px;
}

.tl-team {
  display: flex;
  align-items: center;
  gap: 6px;
}

.tl-team.tl-away {
  justify-content: flex-end;
  text-align: right;
}

.tl-team-name {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
  font-weight: 650;
}

.tl-flag-action {
  width: 28px;
  height: 28px;
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

.tl-flag-action:active {
  transform: scale(0.94);
}

.tl-flag-action:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--primary) 55%, transparent);
  outline-offset: 3px;
}

.tl-score {
  font-weight: 850;
  font-size: 15px;
  color: var(--text);
}

.tl-live .tl-score {
  color: var(--live);
}

.tl-vs {
  color: var(--weak);
  font-size: 11px;
  font-weight: 700;
}

.tl-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  font-size: 11px;
  color: var(--weak);
}

.tl-status {
  font-weight: 800;
}

.tl-status.live {
  color: var(--live);
}

.tl-status.finished {
  color: var(--success);
}

.tl-status.scheduled {
  color: var(--muted);
}

.tl-summary-link {
  padding: 1px 7px;
  border-radius: 999px;
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 10%, transparent);
  font-size: 10px;
  font-weight: 800;
}
</style>
