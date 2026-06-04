<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import MatchInsightCard from '@/components/ai/MatchInsightCard.vue'
import ReminderControl from '@/components/common/ReminderControl.vue'
import TeamFlag from '@/components/common/TeamFlag.vue'
import { apiGetGroupStandings } from '@/api/standings'
import { useAIStore } from '@/stores/useAIStore'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useMatchStore } from '@/stores/useMatchStore'
import { useReminderStore } from '@/stores/useReminderStore'
import { useSettingStore } from '@/stores/useSettingStore'
import { normalizeStanding, type Standing } from '@/types/standing'
import type { Match } from '@/types/match'

const route = useRoute()
const ai = useAIStore()
const matchStore = useMatchStore()
const fav = useFavoriteStore()
const reminder = useReminderStore()
const auth = useAuthStore()
const settings = useSettingStore()

const match = ref<Match | null>(null)
const groupStandings = ref<Standing[]>([])

function formatLocalTime(utcString: string) {
  if (!utcString) return '时间待定'
  const d = new Date(utcString)
  if (Number.isNaN(d.getTime())) return '时间待定'
  return new Intl.DateTimeFormat('zh-CN', {
    timeZone: settings.timezone || Intl.DateTimeFormat().resolvedOptions().timeZone,
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hourCycle: 'h23',
  }).format(d)
}

function formatUTCTime(utcString: string) {
  if (!utcString) return '时间待定'
  const d = new Date(utcString)
  if (Number.isNaN(d.getTime())) return '时间待定'
  return d.toISOString().slice(0, 16).replace('T', ' ')
}

function statusText(status: Match['status']) {
  if (status === 'live') return '进行中'
  if (status === 'finished') return '已结束'
  if (status === 'postponed') return '延期'
  if (status === 'cancelled') return '取消'
  return '未开始'
}

function groupIdFromName(groupName?: string | null) {
  if (!groupName?.startsWith('Group ')) return 0
  return groupName.charCodeAt(6) - 64
}

function rowClass(status: string) {
  if (status === '晋级') return 'rank-ok'
  if (status === '待定') return 'rank-mid'
  return 'rank-out'
}

async function loadMatch() {
  const id = Number(route.params.id)
  if (!id) return

  ai.clearMatchInsight()
  match.value = await matchStore.fetchMatchDetail(id)

  if (auth.isLoggedIn) {
    fav.fetchFavoriteTeams()
    fav.fetchFavoriteMatches()
    reminder.fetchReminders()
  }

  const groupId = groupIdFromName(match.value?.group_name)
  if (!groupId) {
    groupStandings.value = []
    return
  }

  try {
    const res = await apiGetGroupStandings(groupId) as any[]
    groupStandings.value = res.map(normalizeStanding)
  } catch {
    groupStandings.value = []
  }
}

function generateInsight(forceRefresh = false) {
  if (!match.value) return
  ai.generateMatchInsight(match.value.id, forceRefresh).catch(() => {})
}

onMounted(loadMatch)
watch(() => route.params.id, loadMatch)
</script>

<template>
  <div v-if="match" class="match-detail-page">
    <article class="card detail-hero">
      <div class="match-top">
        <span class="tag gold">{{ match.group_name || '世界杯比赛' }}</span>
        <span
          class="tag"
          :class="{ live: match.status === 'live', green: match.status === 'finished' }"
        >
          <i v-if="match.status === 'live'" class="live-dot" style="background: #fff"></i>
          {{ statusText(match.status) }}
        </span>
      </div>

      <div class="detail-score">
        <div class="detail-team">
          <TeamFlag :value="match.home_flag" :alt="match.home_team_name" :fallback="match.home_team_code" size="lg" />
          <h2>{{ match.home_team_name || 'TBD' }}</h2>
          <p>{{ match.home_team_code }}</p>
        </div>

        <div class="big-vs">
          <template v-if="match.status === 'live' || match.status === 'finished'">
            {{ match.home_score ?? 0 }} - {{ match.away_score ?? 0 }}
          </template>
          <template v-else>VS</template>
        </div>

        <div class="detail-team">
          <TeamFlag :value="match.away_flag" :alt="match.away_team_name" :fallback="match.away_team_code" size="lg" />
          <h2>{{ match.away_team_name || 'TBD' }}</h2>
          <p>{{ match.away_team_code }}</p>
        </div>
      </div>

      <div class="info-grid">
        <div class="info-cell">
          <span>本地时间</span>
          <strong>{{ formatLocalTime(match.kickoff_time_utc) }}</strong>
        </div>
        <div class="info-cell">
          <span>UTC 时间</span>
          <strong>{{ formatUTCTime(match.kickoff_time_utc) }}</strong>
        </div>
        <div class="info-cell">
          <span>城市</span>
          <strong>{{ match.city || 'TBD' }}</strong>
        </div>
        <div class="info-cell">
          <span>球场</span>
          <strong>{{ match.stadium || 'TBD' }}</strong>
        </div>
      </div>

      <div class="actions">
        <button
          class="pill-btn"
          :class="{ active: fav.isMatchFavorite(match.id) }"
          type="button"
          @click="fav.toggleMatchFavorite(match.id)"
        >
          <span class="material-symbols-outlined">star</span>
          {{ fav.isMatchFavorite(match.id) ? '已收藏' : '收藏' }}
        </button>
        <ReminderControl :match-id="match.id" mode="pill" />
      </div>
    </article>

    <MatchInsightCard
      :insight="ai.currentMatchInsight"
      :loading="ai.matchInsightLoading"
      :error="ai.matchInsightError"
      @generate="generateInsight(false)"
      @refresh="generateInsight(true)"
    />

    <section v-if="groupStandings.length" class="section">
      <div class="section-head">
        <h2>所在小组积分</h2>
        <span>{{ match.group_name }}</span>
      </div>
      <div class="card table-card table-scroll">
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
            </tr>
          </thead>
          <tbody>
            <tr v-for="(s, i) in groupStandings" :key="s.team_id" :class="rowClass(s.status)">
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
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<style scoped>
.match-detail-page {
  display: grid;
  gap: 16px;
}

.detail-hero {
  padding: 18px;
  overflow: visible;
}

.match-top,
.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.detail-score {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 12px;
  align-items: center;
  margin: 22px 0;
}

.detail-team {
  min-width: 0;
  text-align: center;
}

.detail-team h2 {
  margin: 8px 0 0;
  overflow-wrap: anywhere;
  font-size: 18px;
  font-weight: 750;
}

.detail-team p {
  margin: 5px 0 0;
  color: var(--weak);
  font-size: 12px;
}

.big-vs {
  display: grid;
  place-items: center;
  min-width: 88px;
  min-height: 68px;
  border-radius: 18px;
  color: var(--primary);
  font-size: 24px;
  font-weight: 850;
  background: color-mix(in srgb, var(--primary) 10%, transparent);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.info-cell {
  padding: 13px;
  border-radius: var(--radius-lg);
  background: var(--card-soft);
}

.info-cell span,
.section-head span {
  display: block;
  color: var(--muted);
  font-size: 12px;
}

.info-cell strong {
  display: block;
  margin-top: 5px;
  overflow-wrap: anywhere;
  font-size: 14px;
}

.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 16px;
}

.actions .pill-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.actions .material-symbols-outlined {
  font-size: 18px;
}

.section {
  display: grid;
  gap: 12px;
}

.section-head h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
}

.table-card {
  overflow: hidden;
}

.table-scroll {
  overflow-x: auto;
}

.standing-table {
  width: 100%;
  min-width: 520px;
  border-collapse: collapse;
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
