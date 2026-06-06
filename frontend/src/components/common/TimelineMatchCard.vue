<script setup lang="ts">
import { useRouter } from 'vue-router'
import TeamFlag from '@/components/common/TeamFlag.vue'

const props = defineProps<{
  match: any
}>()

const router = useRouter()

function timeText() {
  const local = props.match.local_kickoff_time || ''
  return local.split(' ')[1] || ''
}

function matchDateKey() {
  const local = props.match.local_kickoff_time || ''
  return local.split(' ')[0] || ''
}

function goDetail() {
  router.push(`/matches/${props.match.id}`)
}
</script>

<template>
  <article
    class="tl-card"
    :class="{
      'tl-live': match.status === 'live',
      'tl-finished': match.status === 'finished',
      'tl-scheduled': match.status === 'scheduled' || match.status === 'upcoming',
    }"
    @click="goDetail"
  >
    <!-- 时间标注 -->
    <div class="tl-time">
      <template v-if="match.status === 'live'">
        <span class="tl-live-pulse"></span>
        {{ match.live_minute || match.minute || 0 }}'
      </template>
      <template v-else-if="match.status === 'finished'">
        <span class="material-symbols-outlined tl-done">check_circle</span>
      </template>
      <template v-else>
        {{ timeText() || 'TBD' }}
      </template>
    </div>

    <!-- 主内容 -->
    <div class="tl-body">
      <div class="tl-teams">
        <div class="tl-team tl-home">
          <TeamFlag :value="match.home_flag" :alt="match.home_team_name" :fallback="match.home_team_code" size="sm" />
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
          <TeamFlag :value="match.away_flag" :alt="match.away_team_name" :fallback="match.away_team_code" size="sm" />
        </div>
      </div>
      <div class="tl-meta">
        <span v-if="match.stage_label || match.group_name">{{ match.group_name || match.stage }}</span>
        <span v-if="match.status === 'finished' && match.has_post_match_summary" class="tl-summary-link">赛后摘要</span>
      </div>
    </div>
  </article>
</template>

<style scoped>
.tl-card {
  display: grid;
  grid-template-columns: 48px minmax(0, 1fr);
  gap: 10px;
  min-height: 64px;
  padding: 11px 14px;
  border: 1px solid var(--line);
  border-radius: var(--radius-lg);
  background: var(--card);
  cursor: pointer;
  transition: border-color 180ms ease, background 180ms ease;
}

.tl-card:active {
  transform: scale(0.99);
}

.tl-live {
  border-color: var(--live);
  background: color-mix(in srgb, var(--live) 5%, var(--card));
}

.tl-time {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 750;
  color: var(--muted);
}

.tl-live .tl-time {
  color: var(--live);
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
}

.tl-body {
  min-width: 0;
  display: grid;
  gap: 6px;
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
  align-items: center;
  gap: 8px;
  font-size: 11px;
  color: var(--weak);
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
