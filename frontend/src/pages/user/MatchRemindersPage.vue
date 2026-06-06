<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import TeamFlag from '@/components/common/TeamFlag.vue'
import { useAuthStore } from '@/stores/useAuthStore'
import { useReminderStore } from '@/stores/useReminderStore'
import type { Match } from '@/types/match'

const router = useRouter()
const auth = useAuthStore()
const reminder = useReminderStore()

const loading = ref(false)
const error = ref('')

const sortedReminders = computed(() =>
  [...reminder.reminders].sort((a, b) => {
    const at = a.match?.kickoff_time_utc ? new Date(a.match.kickoff_time_utc).getTime() : 0
    const bt = b.match?.kickoff_time_utc ? new Date(b.match.kickoff_time_utc).getTime() : 0
    return at - bt
  }),
)

async function load() {
  if (!auth.isLoggedIn) return
  loading.value = true
  error.value = ''
  try {
    await reminder.fetchReminders()
  } catch {
    error.value = '比赛提醒加载失败'
  } finally {
    loading.value = false
  }
}

function channelText(channel: string) {
  if (channel === 'email') return '邮件通知'
  return '站内通知'
}

function kickoffText(match: Match) {
  return match.local_kickoff_time || '时间待定'
}

function goMatch(matchId: number) {
  router.push(`/matches/${matchId}`)
}

onMounted(load)
</script>

<template>
  <div class="detail-page">
    <div class="detail-head">
      <button class="icon-back" type="button" title="返回" @click="router.back()">
        <span class="material-symbols-outlined">arrow_back</span>
      </button>
      <div>
        <h2>比赛提醒</h2>
        <p>{{ reminder.count }} 个提醒</p>
      </div>
    </div>

    <div v-if="loading" class="card state-card">正在加载比赛提醒...</div>
    <div v-else-if="!auth.isLoggedIn" class="card state-card">
      <span>登录后可以查看比赛提醒</span>
      <button class="pill-btn primary" @click="router.push('/login')">去登录</button>
    </div>
    <div v-else-if="error" class="card state-card error">{{ error }}</div>
    <div v-else-if="!sortedReminders.length" class="card state-card">
      <span>还没有设置比赛提醒</span>
      <button class="pill-btn" @click="router.push('/schedule')">去赛程设置</button>
    </div>

    <div v-else class="reminder-list">
      <article
        v-for="item in sortedReminders"
        :key="item.id"
        class="card reminder-item"
        tabindex="0"
        @click="goMatch(item.match_id)"
        @keydown.enter="goMatch(item.match_id)"
      >
        <div class="reminder-meta">
          <span class="material-symbols-outlined">notifications_active</span>
          <strong>提前 {{ item.remind_before_minutes }} 分钟</strong>
          <small>{{ channelText(item.channel) }}</small>
        </div>
        <div v-if="item.match" class="match-summary">
          <div class="team-side">
            <TeamFlag :value="item.match.home_flag" :alt="item.match.home_team_name" :fallback="item.match.home_team_code" size="sm" />
            <span>{{ item.match.home_team_name }}</span>
          </div>
          <div class="match-mid">
            <strong v-if="item.match.status === 'finished' || item.match.status === 'live'">
              {{ item.match.home_score ?? '-' }}-{{ item.match.away_score ?? '-' }}
            </strong>
            <strong v-else>VS</strong>
            <small>{{ kickoffText(item.match) }}</small>
          </div>
          <div class="team-side away">
            <span>{{ item.match.away_team_name }}</span>
            <TeamFlag :value="item.match.away_flag" :alt="item.match.away_team_name" :fallback="item.match.away_team_code" size="sm" />
          </div>
        </div>
        <div v-else class="missing-match">
          <span>比赛 #{{ item.match_id }}</span>
          <small>点击查看比赛详情</small>
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.detail-page {
  display: grid;
  gap: 14px;
}

.detail-head {
  display: flex;
  align-items: center;
  gap: 10px;
}

.icon-back {
  width: 36px;
  height: 36px;
  display: grid;
  place-items: center;
  border: 1px solid var(--line);
  border-radius: 999px;
  color: var(--text);
  background: var(--card);
}

.icon-back .material-symbols-outlined {
  font-size: 20px;
}

.detail-head h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 800;
}

.detail-head p {
  margin: 4px 0 0;
  color: var(--muted);
  font-size: 13px;
}

.reminder-list {
  display: grid;
  gap: 12px;
}

.reminder-item {
  display: grid;
  gap: 10px;
  padding: 12px;
  cursor: pointer;
  transition: border-color 160ms ease-out, background 160ms ease-out, transform 160ms ease-out;
}

.reminder-item:hover {
  border-color: color-mix(in srgb, var(--primary) 28%, var(--line));
  background: color-mix(in srgb, var(--primary) 4%, var(--card));
}

.reminder-item:active {
  transform: scale(0.99);
}

.reminder-item:focus-visible {
  outline: 2px solid var(--primary);
  outline-offset: 2px;
}

.reminder-meta {
  display: flex;
  align-items: center;
  gap: 7px;
  min-width: 0;
  color: var(--text);
}

.reminder-meta .material-symbols-outlined {
  color: var(--primary);
  font-size: 19px;
}

.reminder-meta strong {
  font-size: 13px;
  font-weight: 850;
}

.reminder-meta small {
  margin-left: auto;
  color: var(--muted);
  font-size: 12px;
}

.match-summary {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  align-items: center;
  gap: 10px;
  padding: 12px;
  border-radius: var(--radius-md);
  background: var(--card-soft);
}

.team-side {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.team-side.away {
  justify-content: flex-end;
  text-align: right;
}

.team-side span {
  min-width: 0;
  overflow: hidden;
  font-size: 13px;
  font-weight: 800;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.match-mid {
  min-width: 54px;
  display: grid;
  justify-items: center;
  gap: 4px;
}

.match-mid strong {
  font-size: 16px;
  font-weight: 900;
}

.match-mid small {
  color: var(--muted);
  font-size: 11px;
  white-space: nowrap;
}

.missing-match {
  display: grid;
  gap: 4px;
  padding: 18px;
  border-radius: var(--radius-md);
  background: var(--card-soft);
}

.missing-match span {
  font-weight: 800;
}

.missing-match small {
  color: var(--muted);
  font-size: 12px;
}

.state-card {
  min-height: 120px;
  display: grid;
  place-items: center;
  gap: 12px;
  padding: 20px;
  color: var(--muted);
  text-align: center;
}

.state-card.error {
  color: var(--primary);
}
</style>
