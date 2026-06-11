<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MatchInsightCard from '@/components/ai/MatchInsightCard.vue'
import PostMatchSummaryCard from '@/components/ai/PostMatchSummaryCard.vue'
import MatchStatsCard from '@/components/common/MatchStatsCard.vue'
import MatchLineups from '@/components/common/MatchLineups.vue'
import PointerGlow from '@/components/common/PointerGlow.vue'
import ReminderControl from '@/components/common/ReminderControl.vue'
import TeamFlag from '@/components/common/TeamFlag.vue'
import { apiGetMatchLineups } from '@/api/lineups'
import { apiGetGroupStandings } from '@/api/standings'
import { apiGetPostMatchSummary, apiGeneratePostMatchSummary } from '@/api/postMatchSummary'
import { useAIStore } from '@/stores/useAIStore'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useMatchStore } from '@/stores/useMatchStore'
import { useReminderStore } from '@/stores/useReminderStore'
import { useSettingStore } from '@/stores/useSettingStore'
import { normalizeMatchLineups, type MatchLineups as MatchLineupsData } from '@/types/lineup'
import { normalizeStanding, type Standing } from '@/types/standing'
import type { Match } from '@/types/match'
import type { PostMatchSummary } from '@/types/postMatchSummary'
import { hasPostMatchSummary } from '@/types/postMatchSummary'

const route = useRoute()
const router = useRouter()
const ai = useAIStore()
const matchStore = useMatchStore()
const fav = useFavoriteStore()
const reminder = useReminderStore()
const auth = useAuthStore()
const settings = useSettingStore()

const match = ref<Match | null>(null)
const groupStandings = ref<Standing[]>([])
const lineups = ref<MatchLineupsData | null>(null)
const lineupsLoading = ref(false)
const lineupsError = ref('')
let lineupRequestSeq = 0

const postSummary = ref<PostMatchSummary | null>(null)
const postSummaryLoading = ref(false)
const postSummaryError = ref('')

const canGeneratePostSummary = computed(() =>
  match.value?.status === 'finished' &&
  match.value?.home_score != null &&
  match.value?.away_score != null
)

const hasStats = computed(() =>
  match.value?.home_possession != null ||
  match.value?.home_shots != null ||
  match.value?.home_shots_on_target != null
)

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
  lineups.value = null
  lineupsError.value = ''
  lineupsLoading.value = false
  match.value = await matchStore.fetchMatchDetail(id)
  loadLineups(id)
  loadPostMatchSummary(id)

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

async function loadLineups(matchId: number) {
  const seq = ++lineupRequestSeq
  lineupsLoading.value = true
  lineupsError.value = ''
  lineups.value = null

  try {
    const res = await apiGetMatchLineups(matchId)
    if (seq !== lineupRequestSeq) return
    lineups.value = normalizeMatchLineups(res)
  } catch {
    if (seq !== lineupRequestSeq) return
    lineupsError.value = '首发阵容加载失败'
  } finally {
    if (seq === lineupRequestSeq) {
      lineupsLoading.value = false
    }
  }
}

async function loadPostMatchSummary(matchId: number) {
  postSummary.value = null
  postSummaryError.value = ''

  try {
    const res = await apiGetPostMatchSummary(matchId)
    if (hasPostMatchSummary(res)) {
      postSummary.value = res
    }
  } catch {
    postSummary.value = null
  }
}

async function generatePostSummary(forceRefresh = false) {
  if (!match.value || !canGeneratePostSummary.value) return
  postSummaryLoading.value = true
  postSummaryError.value = ''

  try {
    postSummary.value = await apiGeneratePostMatchSummary(match.value.id, forceRefresh)
  } catch {
    postSummaryError.value = '赛后摘要生成失败'
  } finally {
    postSummaryLoading.value = false
  }
}

function generateInsight(forceRefresh = false) {
  if (!match.value) return
  ai.generateMatchInsight(match.value.id, forceRefresh).catch(() => {})
}

function goBack() {
  if (window.history.length > 1) {
    router.back()
    return
  }
  router.push('/schedule')
}

onMounted(loadMatch)
watch(() => route.params.id, loadMatch)
</script>

<template>
  <div v-if="match" class="match-detail-page">
    <PointerGlow class="detail-pointer-glow" />
    <div class="detail-toolbar">
      <button class="back-action" type="button" title="返回" aria-label="返回上一页" @click="goBack">
        <span class="material-symbols-outlined">arrow_back</span>
      </button>
      <div>
        <h1>比赛详情</h1>
        <p>{{ match.home_team_name || 'TBD' }} vs {{ match.away_team_name || 'TBD' }}</p>
      </div>
    </div>

    <article class="detail-hero">
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
          <span>北京时间</span>
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
          @click="fav.toggleMatchFavorite(match.id, match)"
        >
          <span class="material-symbols-outlined">star</span>
          {{ fav.isMatchFavorite(match.id) ? '已收藏' : '收藏' }}
        </button>
        <ReminderControl :match-id="match.id" mode="pill" />
      </div>
    </article>

    <MatchStatsCard
      v-if="hasStats"
      :match="match"
    />

    <PostMatchSummaryCard
      v-if="match.status === 'finished'"
      :summary="postSummary"
      :loading="postSummaryLoading"
      :error="postSummaryError"
      :can-generate="canGeneratePostSummary"
      @generate="generatePostSummary(false)"
      @refresh="generatePostSummary(true)"
    />

    <MatchInsightCard
      :insight="ai.currentMatchInsight"
      :loading="ai.matchInsightLoading"
      :error="ai.matchInsightError"
      @generate="generateInsight(false)"
      @refresh="generateInsight(true)"
    />

    <MatchLineups
      :lineups="lineups"
      :loading="lineupsLoading"
      :error="lineupsError"
    />

    <section v-if="groupStandings.length" class="section">
      <div class="section-head">
        <h2>所在小组积分</h2>
        <span>{{ match.group_name }}</span>
      </div>
      <div class="table-card table-scroll">
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
  position: relative;
  isolation: isolate;
  display: grid;
  gap: 24px;
}

.match-detail-page > :not(.detail-pointer-glow) {
  position: relative;
  z-index: 1;
}

.detail-toolbar {
  display: grid;
  gap: 8px;
  justify-items: start;
  min-width: 0;
}

.back-action {
  width: 28px;
  height: 28px;
  display: inline-grid;
  place-items: center;
  justify-self: start;
  margin-left: -12px;
  margin-bottom: -2px;
  border: 0;
  border-radius: 999px;
  color: var(--text);
  background: transparent;
  transition: color 160ms ease-out, background 160ms ease-out, transform 160ms ease-out;
}

.back-action:active {
  transform: scale(0.97);
}

.back-action:hover {
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 7%, transparent);
}

.back-action:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--blue) 55%, transparent);
  outline-offset: 2px;
}

.back-action .material-symbols-outlined {
  font-size: 20px;
}

.detail-toolbar > div {
  min-width: 0;
  width: 100%;
}

.detail-toolbar h1 {
  margin: 0;
  font-size: 18px;
  font-weight: 800;
  letter-spacing: 0;
}

.detail-toolbar p {
  min-width: 0;
  overflow: hidden;
  margin: 4px 0 0;
  color: var(--muted);
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.detail-hero {
  padding: 22px;
  border: 1px solid color-mix(in srgb, var(--line) 72%, transparent);
  border-radius: var(--radius-xl);
  background:
    radial-gradient(circle at 12% 0%, color-mix(in srgb, var(--primary) 10%, transparent), transparent 34%),
    linear-gradient(180deg, color-mix(in srgb, var(--card) 88%, transparent), transparent);
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
  border-radius: var(--radius-lg);
  color: var(--primary);
  font-size: 24px;
  font-weight: 850;
  background: color-mix(in srgb, var(--primary) 10%, transparent);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 18px;
  padding: 12px 0;
  border-top: 1px solid color-mix(in srgb, var(--line) 74%, transparent);
  border-bottom: 1px solid color-mix(in srgb, var(--line) 74%, transparent);
}

.info-cell {
  padding: 10px 0;
  border-top: 1px solid color-mix(in srgb, var(--line) 58%, transparent);
  background: transparent;
}

.info-cell:nth-child(-n + 2) {
  border-top: 0;
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
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin-top: 16px;
}

.actions .pill-btn {
  width: 100%;
  min-height: 46px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 0 14px;
  font-size: 14px;
  font-weight: 750;
}

.actions .material-symbols-outlined {
  font-size: 18px;
}

.section {
  display: grid;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid color-mix(in srgb, var(--line) 78%, transparent);
}

.section-head h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
}

.table-card {
  overflow: hidden;
  border-top: 1px solid color-mix(in srgb, var(--line) 78%, transparent);
  border-bottom: 1px solid color-mix(in srgb, var(--line) 78%, transparent);
}

.table-scroll {
  overflow-x: hidden;
}

.standing-table {
  width: 100%;
  table-layout: fixed;
  border-collapse: collapse;
}

th,
td {
  padding: 12px 8px;
  border-bottom: 1px solid var(--line);
  text-align: left;
  font-size: 13px;
}

th:first-child,
td:first-child {
  width: 54px;
}

th:nth-child(2),
td:nth-child(2) {
  width: auto;
}

th:nth-child(n + 3),
td:nth-child(n + 3) {
  width: 42px;
  text-align: center;
}

th {
  color: var(--muted);
  font-weight: 750;
  background: transparent;
  white-space: nowrap;
  word-break: keep-all;
}

td {
  color: var(--text);
}

.team-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.team-cell span {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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

@media (max-width: 720px) {
  th,
  td {
    padding: 12px 4px;
    font-size: 12px;
  }

  th {
    font-size: 11px;
  }

  th:first-child,
  td:first-child {
    width: 42px;
  }

  th:nth-child(n + 3),
  td:nth-child(n + 3) {
    width: 30px;
  }

  th:nth-child(7),
  td:nth-child(7) {
    width: 48px;
  }

  th:nth-child(8),
  td:nth-child(8) {
    width: 40px;
  }

  .team-cell {
    gap: 6px;
  }
}
</style>
