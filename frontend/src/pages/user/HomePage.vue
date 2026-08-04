<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import Countdown from '@/components/common/Countdown.vue'
import MatchCard from '@/components/common/MatchCard.vue'
import TeamFlag from '@/components/common/TeamFlag.vue'
import { apiGetGroupStandings } from '@/api/standings'
import { apiGetCompetitionStandings } from '@/api/competitions'
import { apiGetTodayMatches, apiGetTournamentProgress, apiGetUpcomingMatches, apiListMatches } from '@/api/matches'
import { useAuthStore } from '@/stores/useAuthStore'
import { useCompetitionStore } from '@/stores/useCompetitionStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useMatchStore } from '@/stores/useMatchStore'
import { useTeamStore } from '@/stores/useTeamStore'
import { normalizeMatch, type Match } from '@/types/match'
import { normalizeStanding, normalizeLeagueStanding, type Standing, type LeagueStanding } from '@/types/standing'
import { seasonLabel } from '@/types/competition'

const router = useRouter()
const matchStore = useMatchStore()
const teamStore = useTeamStore()
const fav = useFavoriteStore()
const auth = useAuthStore()
const comp = useCompetitionStore()

const groups = ['Group A', 'Group B', 'Group C', 'Group D', 'Group E', 'Group F', 'Group G', 'Group H', 'Group I', 'Group J', 'Group K', 'Group L']
const activeGroupIndex = ref(0)
const activeGroupName = computed(() => groups[activeGroupIndex.value])
const groupStandings = ref<Standing[]>([])
const groupSwipeStartX = ref(0)
const groupSwipeStartY = ref(0)
const groupSwipeActive = ref(false)
const groupSlideDirection = ref<'left' | 'right'>('left')
const groupTransitionName = computed(() =>
  groupSlideDirection.value === 'left' ? 'standings-next' : 'standings-prev',
)

interface TournamentProgress {
  stage_name: string
  total_matches: number
  completed: number
  live: number
  scheduled: number
  progress: number
}

const progress = ref<TournamentProgress>({
  stage_name: '小组赛阶段',
  total_matches: 0,
  completed: 0,
  live: 0,
  scheduled: 0,
  progress: 0,
})

// League-mode home data (World Cup mode keeps the legacy flow untouched).
const leagueTop = ref<LeagueStanding[]>([])
const leagueStats = ref({ total: 0, finished: 0 })

// Cross-competition "today's highlights": /matches/today is not bound to
// MatchQuery, so the interceptor's competitionId param is ignored by the
// backend — this endpoint naturally returns matches of every competition.
const focusMatches = ref<Match[]>([])

async function loadFocusMatches() {
  try {
    const res = await apiGetTodayMatches() as any[]
    const list = (res || []).map(normalizeMatch)
    const focused = list.filter((m) => m.is_featured || m.importance_level >= 2)
    const pool = focused.length ? focused : list
    const sorted = [...pool].sort(
      (a, b) => new Date(a.kickoff_time_utc).getTime() - new Date(b.kickoff_time_utc).getTime(),
    )
    focusMatches.value = sorted.slice(0, 6)
  } catch {
    focusMatches.value = []
  }
}

function todayParam() {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

async function loadLeagueNextMatch() {
  try {
    const res = await apiListMatches({ status: 'scheduled', page: 1, page_size: 1 }) as any
    const list = (res.list || res || []).map(normalizeMatch)
    nextMatch.value = list[0] || null
  } catch {
    nextMatch.value = null
  }
}

async function loadLeagueTodayMatches() {
  try {
    const res = await apiListMatches({ date: todayParam(), page: 1, page_size: 50 }) as any
    matchStore.todayMatches = (res.list || res || []).map(normalizeMatch)
  } catch {
    matchStore.todayMatches = []
  }
}

async function loadLeagueRecommended() {
  try {
    const res = await apiListMatches({ page: 1, page_size: 3 }) as any
    matchStore.recommendedMatches = (res.list || res || []).map(normalizeMatch)
  } catch {
    matchStore.recommendedMatches = []
  }
}

async function loadLeagueTop() {
  if (!comp.current) return
  try {
    const res = await apiGetCompetitionStandings(comp.currentCode, { type: 'total' }) as any[]
    leagueTop.value = (res || []).slice(0, 5).map(normalizeLeagueStanding)
  } catch {
    leagueTop.value = []
  }
}

async function loadLeagueStats() {
  try {
    const finishedRes = await apiListMatches({ status: 'finished', page: 1, page_size: 1 }) as any
    leagueStats.value.finished = Number(finishedRes.total || 0)
    const allRes = await apiListMatches({ page: 1, page_size: 1 }) as any
    leagueStats.value.total = Number(allRes.total || 0)
  } catch {
    leagueStats.value.total = 0
    leagueStats.value.finished = 0
  }
}

async function loadLeagueHomeData() {
  matchStore.todayMatches = []
  matchStore.recommendedMatches = []
  teamStore.fetchTeams({ page_size: 100, teamType: 'club' })
  loadLeagueNextMatch()
  loadLeagueTodayMatches()
  loadLeagueRecommended()
  loadLeagueTop()
  loadLeagueStats()

  if (auth.isLoggedIn) {
    await fav.fetchFavoriteTeams()
    await loadFollowedSchedule()
  }
}

const nextMatch = ref<Match | null>(null)
const followedSchedule = ref<Match[]>([])

const todayMatches = computed(() => matchStore.todayMatches)
const recommended = computed(() => matchStore.recommendedMatches)
const followedTeams = computed(() =>
  teamStore.teams.filter((team) => fav.isTeamFollowed(team.id)),
)

const nextMatchLocalTime = computed(() => {
  if (!nextMatch.value) return ''
  return nextMatch.value.local_kickoff_time || '时间待定'
})

const sortedFollowedSchedule = computed(() =>
  [...followedSchedule.value].sort((a, b) =>
    new Date(a.kickoff_time_utc).getTime() - new Date(b.kickoff_time_utc).getTime(),
  ),
)

const nextFollowedMatch = computed(() => sortedFollowedSchedule.value[0] || null)
const todayFollowedMatches = computed(() => {
  const todayKey = localDayKey()
  return sortedFollowedSchedule.value.filter((match) => match.local_kickoff_time.startsWith(todayKey))
})

function localDayKey(offset = 0) {
  const d = new Date()
  d.setDate(d.getDate() + offset)
  return `${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function swipeGroup(direction: 'left' | 'right') {
  if (direction === 'left' && activeGroupIndex.value < groups.length - 1) {
    groupSlideDirection.value = direction
    activeGroupIndex.value++
  } else if (direction === 'right' && activeGroupIndex.value > 0) {
    groupSlideDirection.value = direction
    activeGroupIndex.value--
  }
}

function startGroupSwipe(event: PointerEvent) {
  if (!event.isPrimary) return
  ;(event.currentTarget as HTMLElement).setPointerCapture?.(event.pointerId)
  groupSwipeStartX.value = event.clientX
  groupSwipeStartY.value = event.clientY
  groupSwipeActive.value = true
}

function finishGroupSwipe(event: PointerEvent) {
  if (!groupSwipeActive.value || !event.isPrimary) return
  ;(event.currentTarget as HTMLElement).releasePointerCapture?.(event.pointerId)
  groupSwipeActive.value = false

  const deltaX = event.clientX - groupSwipeStartX.value
  const deltaY = event.clientY - groupSwipeStartY.value
  const minDistance = 48
  const maxVerticalDrift = 64

  if (Math.abs(deltaX) < minDistance || Math.abs(deltaY) > maxVerticalDrift) return
  swipeGroup(deltaX < 0 ? 'left' : 'right')
}

function cancelGroupSwipe(event?: PointerEvent) {
  if (event) {
    ;(event.currentTarget as HTMLElement).releasePointerCapture?.(event.pointerId)
  }
  groupSwipeActive.value = false
}

async function loadGroupStandings() {
  try {
    const res = await apiGetGroupStandings(activeGroupIndex.value + 1) as any[]
    groupStandings.value = res.map(normalizeStanding)
  } catch {
    groupStandings.value = []
  }
}

async function loadNextMatch() {
  try {
    const res = await apiGetUpcomingMatches() as any[]
    nextMatch.value = res?.length ? normalizeMatch(res[0]) : null
  } catch {
    nextMatch.value = null
  }
}

async function loadFollowedSchedule() {
  followedSchedule.value = []
  if (!auth.isLoggedIn || fav.followedTeamIds.length === 0) return

  try {
    const res = await apiListMatches({ status: 'scheduled', page: 1, page_size: 100 }) as any
    const list = (res.list || res || []).map(normalizeMatch)
    const followedIds = new Set(fav.followedTeamIds)
    followedSchedule.value = list.filter(
      (match: Match) => followedIds.has(match.home_team_id) || followedIds.has(match.away_team_id),
    )
  } catch {
    followedSchedule.value = []
  }
}

async function loadHomeData() {
  loadFocusMatches() // cross-competition highlights, both modes
  await comp.fetchCompetitions()
  if (comp.isLeague) {
    await loadLeagueHomeData()
    return
  }
  matchStore.fetchTodayMatches()
  matchStore.fetchRecommendedMatches()
  teamStore.fetchTeams({ page_size: 100 })
  loadGroupStandings()
  loadNextMatch()

  try {
    const res = await apiGetTournamentProgress()
    if (res) progress.value = res
  } catch {}

  if (auth.isLoggedIn) {
    await fav.fetchFavoriteTeams()
    await loadFollowedSchedule()
  }
}

watch(activeGroupIndex, loadGroupStandings)
watch(
  () => comp.currentCode,
  () => {
    activeGroupIndex.value = 0
    loadHomeData()
  },
)

onMounted(loadHomeData)
</script>

<template>
  <div class="dashboard-grid">
    <div>
      <Countdown v-if="nextMatch" :targetTime="nextMatch.kickoff_time_utc">
        <span class="eyebrow"><i class="live-dot next-dot"></i> 下一场比赛</span>
        <div class="countdown-title">
          <div>
            <h2>{{ nextMatch.home_team_name }} vs {{ nextMatch.away_team_name }}</h2>
            <p>{{ nextMatchLocalTime }} 开球 · 本地时间</p>
          </div>
          <span v-if="nextMatch.is_featured" class="tag gold">推荐</span>
        </div>
      </Countdown>
      <article v-else class="card empty-card">暂无即将到来的比赛</article>

      <article class="card stage-card">
        <template v-if="comp.isLeague">
          <div class="stage-head">
            <span>{{ comp.current?.name || '联赛' }} · {{ seasonLabel(comp.current?.season || 0) }}</span>
            <span v-if="leagueStats.total > 0" style="color: var(--primary)">{{ leagueStats.total > 0 ? ((leagueStats.finished / leagueStats.total) * 100).toFixed(1) : '0.0' }}%</span>
          </div>
          <div class="stage-track">
            <div
              class="stage-progress"
              :style="{ width: leagueStats.total > 0 ? (leagueStats.finished / leagueStats.total) * 100 + '%' : '0%' }"
            ></div>
          </div>
          <div class="stage-meta">
            <span>已完成 {{ leagueStats.finished }} 场</span>
            <span>共 {{ leagueStats.total }} 场</span>
          </div>
        </template>
        <template v-else>
          <div class="stage-head">
            <span>{{ progress.stage_name }}</span>
            <span v-if="progress.live > 0" style="color: var(--hot)">进行中 {{ progress.progress.toFixed(1) }}%</span>
            <span v-else style="color: var(--primary)">{{ progress.progress.toFixed(1) }}%</span>
          </div>
          <div class="stage-track">
            <div class="stage-progress" :style="{ width: progress.progress + '%' }"></div>
          </div>
          <div class="stage-meta">
            <span>已完成 {{ progress.completed }} 场</span>
            <span>剩余 {{ progress.scheduled }} 场</span>
          </div>
        </template>
      </article>

      <section class="section">
        <div class="section-head">
          <h2>我的关注赛程</h2>
          <span v-if="auth.isLoggedIn">{{ followedTeams.length }} 支球队</span>
        </div>

        <div v-if="!auth.isLoggedIn" class="card follow-state">
          <span>登录后可以查看关注球队的下一场比赛</span>
          <button class="pill-btn primary" @click="router.push('/login')">去登录</button>
        </div>
        <div v-else-if="followedTeams.length === 0" class="card follow-state">
          <span>还没有关注球队</span>
          <button class="pill-btn primary" @click="router.push('/teams')">去关注</button>
        </div>
        <div v-else-if="!nextFollowedMatch" class="card follow-state">
          <span>关注球队暂无未开始比赛</span>
          <button class="pill-btn" @click="router.push('/schedule')">查看全部赛程</button>
        </div>
        <div v-else class="stack">
          <MatchCard :match="nextFollowedMatch" featured />
          <div v-if="todayFollowedMatches.length > 1" class="mini-list">
            <button
              v-for="match in todayFollowedMatches.slice(1, 4)"
              :key="match.id"
              class="mini-match"
              @click="router.push(`/matches/${match.id}`)"
            >
              <span>{{ match.home_team_name }} vs {{ match.away_team_name }}</span>
              <b>{{ match.local_kickoff_time.split(' ')[1] || 'TBD' }}</b>
            </button>
          </div>
        </div>
      </section>

      <section v-if="focusMatches.length" class="section">
        <div class="section-head">
          <h2>今日焦点</h2>
          <span>跨赛事精选</span>
        </div>
        <div class="match-strip">
          <MatchCard v-for="m in focusMatches" :key="`focus-${m.id}`" :match="m" featured />
        </div>
      </section>

      <section class="section">
        <div class="section-head">
          <h2>今日比赛</h2>
          <span>{{ todayMatches.length }} 场</span>
        </div>
        <div class="match-strip">
          <MatchCard v-for="m in todayMatches" :key="m.id" :match="m" />
        </div>
      </section>

      <section class="section">
        <div class="section-head">
          <h2>热门推荐</h2>
        </div>
        <div class="stack">
          <MatchCard v-for="m in recommended" :key="m.id" :match="m" featured />
        </div>
      </section>
    </div>

    <aside class="desktop-side">
      <section class="section" style="margin-top: 0">
        <div class="section-head">
          <h2>我的关注</h2>
          <span>{{ followedTeams.length }} 支球队</span>
        </div>
        <div v-if="followedTeams.length" class="stack">
          <article
            v-for="t in followedTeams"
            :key="t.id"
            class="card profile-card"
            @click="router.push(`/teams/${t.id}`)"
          >
            <TeamFlag :value="t.flag" :alt="t.name" :fallback="t.code" size="lg" />
            <div>
              <h2>{{ t.name }}</h2>
              <p>{{ t.group_name }}</p>
            </div>
          </article>
        </div>
        <div v-else class="card empty-card">暂无关注球队</div>
      </section>

      <section class="section">
        <div class="section-head">
          <h2>积分速览</h2>
          <div v-if="!comp.isLeague" class="group-nav">
            <button class="nav-btn" :disabled="activeGroupIndex === 0" @click="swipeGroup('right')">‹</button>
            <span class="group-label" @click="swipeGroup('left')">{{ activeGroupName }}</span>
            <button class="nav-btn" :disabled="activeGroupIndex === groups.length - 1" @click="swipeGroup('left')">›</button>
          </div>
          <span v-else class="muted-count">前 5 名</span>
        </div>

        <template v-if="comp.isLeague">
          <div class="card table-card">
            <table class="standing-table" style="min-width: 0">
              <tbody>
                <tr v-for="s in leagueTop" :key="s.team_id">
                  <td>{{ s.position }}</td>
                  <td class="team-cell">
                    <TeamFlag :value="s.flag" :alt="s.team_name" :fallback="s.team_code" size="sm" />
                    <span>{{ s.team_name }}</span>
                  </td>
                  <td><b>{{ s.points }}</b></td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>

        <div
          v-else
          class="card table-card standings-swipe-area"
          @pointerdown="startGroupSwipe"
          @pointerup="finishGroupSwipe"
          @pointercancel="cancelGroupSwipe"
          @pointerleave="cancelGroupSwipe"
        >
          <Transition :name="groupTransitionName" mode="out-in">
            <table :key="activeGroupName" class="standing-table" style="min-width: 0">
              <tbody>
                <tr
                  v-for="(s, i) in groupStandings"
                  :key="s.team_id"
                  :class="{
                    'rank-ok': s.status === '晋级',
                    'rank-mid': s.status === '待定',
                    'rank-out': s.status === '淘汰',
                  }"
                >
                  <td>{{ i + 1 }}</td>
                  <td class="team-cell">
                    <TeamFlag :value="s.flag" :alt="s.team_name" :fallback="s.team_code" size="sm" />
                    <span>{{ s.team_name }}</span>
                  </td>
                  <td><b>{{ s.points }}</b></td>
                </tr>
              </tbody>
            </table>
          </Transition>
        </div>
      </section>
    </aside>
  </div>
</template>

<style scoped>
.dashboard-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 16px;
}

.desktop-side {
  display: grid;
  gap: 16px;
}

.section {
  margin-top: 18px;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.section-head h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
}

.section-head span {
  color: var(--muted);
  font-size: 13px;
}

.match-strip {
  display: grid;
  grid-auto-flow: column;
  grid-auto-columns: minmax(280px, 84%);
  gap: 12px;
  overflow-x: auto;
  padding: 2px 2px 12px;
  scroll-snap-type: x mandatory;
}

.stack {
  display: grid;
  gap: 12px;
}

.eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 6px 10px;
  border-radius: 999px;
  color: #fff2c4;
  font-size: 12px;
  font-weight: 700;
  background: rgba(255, 255, 255, 0.1);
}

.countdown-title {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-top: 18px;
}

.countdown-title h2 {
  margin: 0;
  font-size: 21px;
  line-height: 1.2;
}

.countdown-title p {
  margin: 7px 0 0;
  color: rgba(255, 255, 255, 0.72);
  font-size: 13px;
}

.stage-card {
  display: grid;
  gap: 11px;
  margin-top: 12px;
  padding: 14px 15px;
}

.stage-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  font-size: 13px;
  font-weight: 750;
}

.stage-track {
  height: 7px;
  overflow: hidden;
  border-radius: 999px;
  background: var(--card-soft);
}

.stage-progress {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, var(--primary), var(--secondary));
  transition: width 0.6s ease;
}

.stage-meta {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  color: var(--muted);
  font-size: 11px;
}

.follow-state,
.empty-card {
  padding: 20px;
  display: grid;
  place-items: center;
  gap: 12px;
  text-align: center;
  color: var(--muted);
}

.mini-list {
  display: grid;
  gap: 8px;
}

.mini-match {
  min-height: 42px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 0 14px;
  border: 1px solid var(--line);
  border-radius: 12px;
  color: var(--text);
  background: var(--card);
  font-size: 13px;
  cursor: pointer;
}

.mini-match span {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.profile-card {
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 13px;
  cursor: pointer;
}

.profile-card h2 {
  margin: 0;
  font-size: 18px;
}

.profile-card p {
  margin: 5px 0 0;
  color: var(--muted);
  font-size: 13px;
}

.standing-table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  padding: 12px 10px;
  border-bottom: 1px solid var(--line);
  text-align: left;
  font-size: 13px;
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

.table-card {
  overflow: hidden;
}

.standings-swipe-area {
  min-height: 177px;
  touch-action: pan-y;
  user-select: none;
}

.standings-next-enter-active,
.standings-next-leave-active,
.standings-prev-enter-active,
.standings-prev-leave-active {
  transition: opacity 0.22s ease, transform 0.22s ease;
}

.standings-next-enter-from,
.standings-prev-leave-to {
  opacity: 0;
  transform: translateX(34px);
}

.standings-next-leave-to,
.standings-prev-enter-from {
  opacity: 0;
  transform: translateX(-34px);
}

.team-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.group-nav {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nav-btn {
  width: 28px;
  height: 28px;
  display: grid;
  place-items: center;
  border: 1px solid var(--line);
  border-radius: 999px;
  color: var(--text);
  background: var(--card);
  font-size: 16px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
}

.nav-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.nav-btn:not(:disabled):hover {
  color: #fff;
  border-color: var(--primary);
  background: var(--primary);
}

.group-label {
  min-width: 70px;
  text-align: center;
  color: var(--primary);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

.muted-count {
  color: var(--muted);
  font-size: 13px;
}

@media (min-width: 768px) {
  .dashboard-grid {
    grid-template-columns: minmax(0, 1.9fr) minmax(300px, 0.8fr);
    align-items: start;
  }

  .match-strip {
    grid-auto-flow: row;
    grid-auto-columns: initial;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    overflow: visible;
    padding-bottom: 0;
  }
}
</style>
