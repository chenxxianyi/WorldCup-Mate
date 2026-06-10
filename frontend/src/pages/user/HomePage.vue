<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import Countdown from '@/components/common/Countdown.vue'
import MatchCard from '@/components/common/MatchCard.vue'
import TimelineMatchCard from '@/components/common/TimelineMatchCard.vue'
import TeamFlag from '@/components/common/TeamFlag.vue'
import { apiGetGroupStandings } from '@/api/standings'
import { apiGetTournamentProgress, apiGetUpcomingMatches, apiGetTimeline, apiListMatches } from '@/api/matches'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useMatchStore } from '@/stores/useMatchStore'
import { useTeamStore } from '@/stores/useTeamStore'
import { normalizeMatch, type Match } from '@/types/match'

const router = useRouter()
const matchStore = useMatchStore()
const teamStore = useTeamStore()
const fav = useFavoriteStore()
const auth = useAuthStore()

const groups = ['Group A', 'Group B', 'Group C', 'Group D', 'Group E', 'Group F', 'Group G', 'Group H', 'Group I', 'Group J', 'Group K', 'Group L']
const activeGroupIndex = ref(0)
const activeGroupName = computed(() => groups[activeGroupIndex.value])
const groupStandings = ref<any[]>([])
const groupSwipeStartX = ref(0)
const groupSwipeStartY = ref(0)
const groupSwipeActive = ref(false)
const groupSlideDirection = ref<'left' | 'right'>('left')
const groupTransitionName = computed(() =>
  groupSlideDirection.value === 'left' ? 'standings-next' : 'standings-prev',
)
const standingsMistVisible = ref(false)
const standingsMistKey = ref(0)
let standingsMistTimer: number | undefined

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

const nextMatch = ref<Match | null>(null)
type TimelineDisplayMatch = Match & { has_post_match_summary?: boolean }
const timelineRaw = ref<TimelineDisplayMatch[]>([])
const followedSchedule = ref<Match[]>([])

const todayMatches = computed(() => matchStore.todayMatches)
const followedTeams = computed(() =>
  teamStore.teams.filter((team) => fav.isTeamFollowed(team.id)),
)

const nextMatchLocalTime = computed(() => {
  if (!nextMatch.value) return ''
  return nextMatch.value.local_kickoff_time || '时间待定'
})

// ── Timeline grouping ──

interface TimelineGroup {
  key: string
  title: string
  isToday: boolean
  isYesterday: boolean
  isTomorrow: boolean
  matches: any[]
}

const groupedTimeline = computed(() => {
  const groupsMap = new Map<string, any[]>()
  const todayKey = localDayKey(0)
  const yesterdayKey = localDayKey(-1)
  const tomorrowKey = localDayKey(1)

  for (const m of timelineRaw.value) {
    const local = m.local_kickoff_time || ''
    const key = local.split(' ')[0] || 'other'
    if (!groupsMap.has(key)) groupsMap.set(key, [])
    groupsMap.get(key)!.push(m)
  }

  const result: TimelineGroup[] = []
  for (const [key, matches] of groupsMap) {
    result.push({
      key,
      title: dateTitle(key),
      isToday: key === todayKey,
      isYesterday: key === yesterdayKey,
      isTomorrow: key === tomorrowKey,
      matches,
    })
  }
  result.sort((a, b) => a.key.localeCompare(b.key))
  return result
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

function dateTitle(key: string) {
  const [month, day] = key.split('-')
  if (!month || !day) return '时间待定'
  return `${Number(month)}月${Number(day)}日`
}

function swipeGroup(direction: 'left' | 'right') {
  if (direction === 'left' && activeGroupIndex.value < groups.length - 1) {
    groupSlideDirection.value = direction
    activeGroupIndex.value++
    triggerStandingsMist()
  } else if (direction === 'right' && activeGroupIndex.value > 0) {
    groupSlideDirection.value = direction
    activeGroupIndex.value--
    triggerStandingsMist()
  }
}

function triggerStandingsMist() {
  standingsMistKey.value++
  standingsMistVisible.value = true
  if (standingsMistTimer) window.clearTimeout(standingsMistTimer)
  standingsMistTimer = window.setTimeout(() => {
    standingsMistVisible.value = false
  }, 760)
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
  if (Math.abs(deltaX) < 48 || Math.abs(deltaY) > 64) return
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
    groupStandings.value = res.map((s: any) => ({
      team_id: s.team_id,
      team_name: s.team?.name || '',
      team_code: s.team?.fifa_code || '',
      flag: s.team?.flag_url || '',
      points: s.points,
    }))
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
  matchStore.fetchTodayMatches()
  matchStore.fetchRecommendedMatches()
  teamStore.fetchTeams({ page_size: 100 })
  loadGroupStandings()
  loadNextMatch()

  // Load timeline
  try {
    const res = await apiGetTimeline({ days_back: 3, days_ahead: 3 }) as any[]
    timelineRaw.value = (res || []).map((m) => ({
      ...normalizeMatch(m),
      has_post_match_summary: m.has_post_match_summary,
    }))
  } catch {
    timelineRaw.value = []
  }

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
  () => auth.isLoggedIn,
  async (loggedIn) => {
    if (!loggedIn) {
      followedSchedule.value = []
      return
    }
    await fav.fetchFavoriteTeams()
    await loadFollowedSchedule()
  },
)

onMounted(loadHomeData)

onBeforeUnmount(() => {
  if (standingsMistTimer) window.clearTimeout(standingsMistTimer)
})
</script>

<template>
  <div class="dashboard-grid">
    <div>
      <!-- Countdown -->
      <Countdown v-if="nextMatch" :targetTime="nextMatch.kickoff_time_utc">
        <span class="eyebrow"><i class="live-dot next-dot"></i> 下一场比赛</span>
        <div class="countdown-title">
          <div>
            <h2>{{ nextMatch.home_team_name }} vs {{ nextMatch.away_team_name }}</h2>
            <p>{{ nextMatchLocalTime }} 开球 · 北京时间</p>
          </div>
          <span v-if="nextMatch.is_featured" class="tag gold">推荐</span>
        </div>
      </Countdown>
      <article v-else class="card empty-card">暂无即将到来的比赛</article>

      <!-- Stage progress -->
      <article class="card stage-card">
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
      </article>

      <!-- My followed items -->
      <section v-if="auth.isLoggedIn" class="section follow-section">
        <div class="section-head">
          <h2>我的关注</h2>
          <span>{{ followedTeams.length }} 支球队</span>
        </div>
        <div class="follow-block">
          <div class="follow-subhead">关注的比赛</div>
          <div v-if="nextFollowedMatch" class="stack">
            <MatchCard :match="nextFollowedMatch" featured />
          </div>
          <div v-else class="card empty-card">关注球队暂无未开始比赛</div>
        </div>
        <div class="follow-block">
          <div class="follow-subhead">关注的球队</div>
          <div v-if="followedTeams.length" class="followed-team-grid">
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
        </div>
      </section>

      <!-- ⏱ Timeline -->
      <section class="section timeline-section timeline-mobile">
        <div class="section-head">
          <h2>比赛时间线</h2>
          <span>近 1 周</span>
        </div>
        <template v-for="group in groupedTimeline" :key="group.key">
          <div class="tl-date-head">
            <span>{{ group.title }}</span>
            <span v-if="group.isToday" class="tag live">今天</span>
            <span v-else-if="group.isYesterday" class="tag">昨天</span>
            <span v-else-if="group.isTomorrow" class="tag blue">明天</span>
          </div>
          <div class="tl-list">
            <TimelineMatchCard
              v-for="m in group.matches"
              :key="m.id"
              :match="m"
            />
          </div>
        </template>
        <div v-if="!timelineRaw.length" class="card empty-card">暂无比赛数据</div>
      </section>
    </div>

    <!-- Desktop sidebar -->
    <aside class="desktop-side">
      <section class="section standings-section">
        <div class="section-head">
          <h2>积分速览</h2>
          <div class="group-nav">
            <button class="nav-btn" :disabled="activeGroupIndex === 0" @click="swipeGroup('right')">‹</button>
            <span class="group-label" @click="swipeGroup('left')">{{ activeGroupName }}</span>
            <button class="nav-btn" :disabled="activeGroupIndex === groups.length - 1" @click="swipeGroup('left')">›</button>
          </div>
        </div>
        <div
          class="card table-card standings-swipe-area"
          @pointerdown="startGroupSwipe"
          @pointerup="finishGroupSwipe"
          @pointercancel="cancelGroupSwipe"
          @pointerleave="cancelGroupSwipe"
        >
          <div
            v-if="standingsMistVisible"
            :key="standingsMistKey"
            class="standings-mist"
            :class="groupSlideDirection"
            aria-hidden="true"
          ></div>
          <Transition :name="groupTransitionName" mode="out-in">
            <table :key="activeGroupName" class="standing-table" style="min-width: 0">
              <tbody>
                <tr v-for="(s, i) in groupStandings" :key="s.team_id">
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

      <section class="section timeline-section timeline-desktop">
        <div class="section-head">
          <h2>比赛时间线</h2>
          <span>近 1 周</span>
        </div>
        <template v-for="group in groupedTimeline" :key="group.key">
          <div class="tl-date-head">
            <span>{{ group.title }}</span>
            <span v-if="group.isToday" class="tag live">今天</span>
            <span v-else-if="group.isYesterday" class="tag">昨天</span>
            <span v-else-if="group.isTomorrow" class="tag blue">明天</span>
          </div>
          <div class="tl-list">
            <TimelineMatchCard
              v-for="m in group.matches"
              :key="m.id"
              :match="m"
            />
          </div>
        </template>
        <div v-if="!timelineRaw.length" class="card empty-card">暂无比赛数据</div>
      </section>
    </aside>
  </div>
</template>

<style scoped>
.dashboard-grid {
  width: 100%;
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

.standings-section {
  margin-top: 0;
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

.stack {
  display: grid;
  gap: 12px;
}

.follow-section {
  display: grid;
  gap: 14px;
}

.follow-section .section-head {
  margin-bottom: 0;
}

.follow-block {
  display: grid;
  gap: 10px;
}

.follow-subhead {
  color: var(--muted);
  font-size: 13px;
  font-weight: 750;
}

.followed-team-grid {
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

.empty-card {
  padding: 20px;
  display: grid;
  place-items: center;
  gap: 12px;
  text-align: center;
  color: var(--muted);
}

/* ── Timeline ── */
.tl-date-head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 16px 0 8px;
  font-size: 14px;
  font-weight: 750;
  color: var(--muted);
}

.tl-list {
  display: grid;
  gap: 8px;
}

.timeline-desktop {
  display: none;
}

/* ── Sidebar ── */
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

th, td {
  padding: 12px 10px;
  border-bottom: 1px solid var(--line);
  text-align: left;
  font-size: 13px;
}

tr:last-child td {
  border-bottom: 0;
}

.team-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.table-card {
  overflow: hidden;
}

.standings-swipe-area {
  position: relative;
  min-height: 177px;
  touch-action: pan-y;
  user-select: none;
}

.standings-mist {
  position: absolute;
  inset: -14px -22px;
  z-index: 1;
  overflow: hidden;
  pointer-events: none;
  border-radius: inherit;
  filter: saturate(1.08);
}

.standings-mist::before,
.standings-mist::after {
  content: "";
  position: absolute;
  top: 4px;
  bottom: 4px;
  width: 82%;
  opacity: 0;
  filter: blur(16px);
  transform: translateX(0) scaleX(0.68);
  background:
    radial-gradient(ellipse at 42% 48%, rgba(255, 255, 255, 0.94) 0 18%, rgba(255, 255, 255, 0.56) 38%, transparent 72%),
    radial-gradient(ellipse at 72% 28%, rgba(255, 255, 255, 0.5) 0 18%, transparent 54%),
    linear-gradient(90deg, rgba(14, 165, 233, 0.24), rgba(255, 255, 255, 0.82), rgba(245, 158, 11, 0.2));
  mix-blend-mode: screen;
}

.standings-mist::after {
  top: 16%;
  bottom: 10%;
  width: 64%;
  filter: blur(24px);
  background:
    radial-gradient(ellipse at center, rgba(255, 255, 255, 0.86) 0 24%, rgba(255, 255, 255, 0.36) 52%, transparent 80%),
    linear-gradient(90deg, rgba(99, 102, 241, 0.18), rgba(255, 255, 255, 0.52), rgba(14, 165, 233, 0.16));
}

.standings-mist {
  box-shadow:
    inset 0 0 40px rgba(255, 255, 255, 0.7),
    0 0 34px rgba(14, 165, 233, 0.08);
}

.standings-mist.left::before {
  right: -28%;
  animation: standings-mist-left 760ms ease-out both;
}

.standings-mist.left::after {
  right: 2%;
  animation: standings-mist-left-soft 760ms ease-out both;
}

.standings-mist.right::before {
  left: -28%;
  animation: standings-mist-right 760ms ease-out both;
}

.standings-mist.right::after {
  left: 2%;
  animation: standings-mist-right-soft 760ms ease-out both;
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

@keyframes standings-mist-left {
  0% {
    opacity: 0;
    transform: translateX(52px) scaleX(0.48);
  }
  18% {
    opacity: 0.9;
  }
  46% {
    opacity: 0.56;
  }
  100% {
    opacity: 0;
    transform: translateX(-118px) scaleX(1.28);
  }
}

@keyframes standings-mist-left-soft {
  0% {
    opacity: 0;
    transform: translateX(34px) scaleX(0.58);
  }
  20% {
    opacity: 0.68;
  }
  52% {
    opacity: 0.34;
  }
  100% {
    opacity: 0;
    transform: translateX(-76px) scaleX(1.36);
  }
}

@keyframes standings-mist-right {
  0% {
    opacity: 0;
    transform: translateX(-52px) scaleX(0.48);
  }
  18% {
    opacity: 0.9;
  }
  46% {
    opacity: 0.56;
  }
  100% {
    opacity: 0;
    transform: translateX(118px) scaleX(1.28);
  }
}

@keyframes standings-mist-right-soft {
  0% {
    opacity: 0;
    transform: translateX(-34px) scaleX(0.58);
  }
  20% {
    opacity: 0.68;
  }
  52% {
    opacity: 0.34;
  }
  100% {
    opacity: 0;
    transform: translateX(76px) scaleX(1.36);
  }
}

@media (prefers-reduced-motion: reduce) {
  .standings-mist {
    display: none;
  }

  .standings-next-enter-active,
  .standings-next-leave-active,
  .standings-prev-enter-active,
  .standings-prev-leave-active {
    transition: none;
  }
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

@media (min-width: 768px) {
  .followed-team-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (min-width: 1024px) {
  .dashboard-grid {
    max-width: 1080px;
    margin: 0 auto;
    grid-template-columns: minmax(0, 720px) minmax(300px, 320px);
    align-items: start;
    justify-content: center;
    gap: 20px;
  }

  .timeline-mobile {
    display: none;
  }

  .timeline-desktop {
    display: block;
  }

  .timeline-desktop .section-head {
    margin-bottom: 10px;
  }

  .timeline-desktop .tl-date-head {
    margin: 13px 0 7px;
  }

  .timeline-desktop .tl-list {
    gap: 7px;
  }
}

@media (min-width: 1200px) {
  .dashboard-grid {
    max-width: 1120px;
    grid-template-columns: minmax(0, 760px) minmax(300px, 320px);
    gap: 24px;
  }
}
</style>
