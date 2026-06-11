<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import Countdown from '@/components/common/Countdown.vue'
import MatchTicketRow from '@/components/common/MatchTicketRow.vue'
import PointerGlow from '@/components/common/PointerGlow.vue'
import TimelineDayGroup from '@/components/common/TimelineDayGroup.vue'
import TeamFlag from '@/components/common/TeamFlag.vue'
import { apiGetGroupStandings } from '@/api/standings'
import { apiGetTournamentProgress, apiGetUpcomingMatches, apiGetTimeline } from '@/api/matches'
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

const stageTotalMatches = computed(() =>
  progress.value.total_matches ||
  progress.value.completed + progress.value.live + progress.value.scheduled,
)

const stageStatusCopy = computed(() => {
  if (progress.value.live > 0) return `${progress.value.live} 场比赛正在进行，关注的球队开赛会同步更新。`
  if (progress.value.completed === 0) return '小组赛尚未开赛，首场比赛结束后这里会展示阶段走势。'
  if (progress.value.scheduled > 0) return `还有 ${progress.value.scheduled} 场等待开球，赛程进度会持续推进。`
  return '本阶段比赛已经完成，可以继续查看淘汰赛路径。'
})

const nextMatch = ref<Match | null>(null)
type TimelineDisplayMatch = Match & { has_post_match_summary?: boolean }
const timelineRaw = ref<TimelineDisplayMatch[]>([])
const followedMatchRail = ref<HTMLElement | null>(null)
const activeFollowedMatchIndex = ref(0)
const followedSlideDirection = ref<'from-left' | 'from-right'>('from-right')

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

const sortedFavoriteMatches = computed(() =>
  [...fav.favoriteMatches].sort((a, b) =>
    new Date(a.kickoff_time_utc).getTime() - new Date(b.kickoff_time_utc).getTime(),
  ),
)

const followedMatchCards = computed(() => sortedFavoriteMatches.value)

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

function updateFollowedMatchIndex() {
  const rail = followedMatchRail.value
  if (!rail) return
  const width = rail.clientWidth || 1
  const lastIndex = Math.max(0, followedMatchCards.value.length - 1)
  const nextIndex = Math.min(lastIndex, Math.max(0, Math.round(rail.scrollLeft / width)))
  if (nextIndex === activeFollowedMatchIndex.value) return
  followedSlideDirection.value = nextIndex > activeFollowedMatchIndex.value ? 'from-right' : 'from-left'
  activeFollowedMatchIndex.value = nextIndex
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
    await fav.fetchFavoriteMatches()
  }
}

watch(activeGroupIndex, loadGroupStandings)
watch(
  () => auth.isLoggedIn,
  async (loggedIn) => {
    if (!loggedIn) {
      fav.clearFavorites()
      return
    }
    await fav.fetchFavoriteTeams()
    await fav.fetchFavoriteMatches()
  },
)

watch(
  () => followedMatchCards.value.map((match) => match.id).join(','),
  () => {
    activeFollowedMatchIndex.value = 0
    followedMatchRail.value?.scrollTo({ left: 0 })
  },
)

onMounted(loadHomeData)

onBeforeUnmount(() => {
  if (standingsMistTimer) window.clearTimeout(standingsMistTimer)
})
</script>

<template>
  <div class="dashboard-grid">
    <PointerGlow class="home-pointer-glow" />
    <div class="home-main">
      <!-- Countdown -->
      <Countdown v-if="nextMatch" class="home-countdown" :targetTime="nextMatch.kickoff_time_utc">
        <span class="eyebrow"><i class="live-dot next-dot"></i> 下一场比赛</span>
        <div class="countdown-title">
          <div>
            <h2>{{ nextMatch.home_team_name }} vs {{ nextMatch.away_team_name }}</h2>
            <p>{{ nextMatchLocalTime }} 开球 · 北京时间</p>
          </div>
          <span v-if="nextMatch.is_featured" class="tag gold">推荐</span>
        </div>
      </Countdown>
      <article v-else class="empty-state home-countdown">暂无即将到来的比赛</article>

      <!-- Stage progress -->
      <section class="stage-strip">
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
        <div class="stage-stats-grid" aria-label="赛事阶段统计">
          <div class="stage-stat">
            <span>总场次</span>
            <strong>{{ stageTotalMatches }}</strong>
          </div>
          <div class="stage-stat">
            <span>已完成</span>
            <strong>{{ progress.completed }}</strong>
          </div>
          <div class="stage-stat">
            <span>进行中</span>
            <strong>{{ progress.live }}</strong>
          </div>
          <div class="stage-stat">
            <span>未开始</span>
            <strong>{{ progress.scheduled }}</strong>
          </div>
        </div>
        <div class="stage-context">
          <span class="material-symbols-outlined" aria-hidden="true">event_available</span>
          <div>
            <b>{{ progress.stage_name }}</b>
            <p>{{ stageStatusCopy }}</p>
          </div>
        </div>
      </section>

      <!-- My followed items -->
      <section v-if="auth.isLoggedIn" class="section follow-section">
        <div class="section-head">
          <h2>我的关注</h2>
          <span>{{ followedTeams.length }} 支球队</span>
        </div>
        <div class="follow-block">
          <div class="follow-subhead">关注的比赛</div>
          <div v-if="followedMatchCards.length" class="followed-match-carousel">
            <div
              ref="followedMatchRail"
              class="followed-match-rail"
              @scroll.passive="updateFollowedMatchIndex"
            >
              <div
                v-for="match in followedMatchCards"
                :key="match.id"
                class="followed-match-slide"
                :class="[
                  { active: match.id === followedMatchCards[activeFollowedMatchIndex]?.id },
                  followedSlideDirection,
                ]"
              >
                <MatchTicketRow :match="match" />
              </div>
            </div>
          </div>
          <div v-else class="empty-state compact">暂无关注比赛</div>
        </div>
        <div class="follow-block">
          <div class="follow-subhead">关注的球队</div>
          <div v-if="followedTeams.length" class="followed-team-rail">
            <article
              v-for="t in followedTeams"
              :key="t.id"
              class="team-chip"
              @click="router.push(`/teams/${t.id}`)"
            >
              <TeamFlag :value="t.flag" :alt="t.name" :fallback="t.code" size="md" />
              <span>{{ t.name }}</span>
            </article>
          </div>
          <div v-else class="empty-state compact">暂无关注球队</div>
        </div>
      </section>

      <!-- ⏱ Timeline -->
      <section class="section timeline-section timeline-mobile">
        <div class="section-head">
          <h2>比赛时间线</h2>
          <span>近 1 周</span>
        </div>
        <TimelineDayGroup
          v-for="group in groupedTimeline"
          :key="group.key"
          :title="group.title"
          :matches="group.matches"
          :is-today="group.isToday"
          :is-yesterday="group.isYesterday"
          :is-tomorrow="group.isTomorrow"
        />
        <div v-if="!timelineRaw.length" class="empty-state">暂无比赛数据</div>
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
          class="standings-panel standings-swipe-area"
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

      <section class="section timeline-section timeline-desktop timeline-main">
        <div class="section-head">
          <h2>比赛时间线</h2>
          <span>近 1 周</span>
        </div>
        <TimelineDayGroup
          v-for="group in groupedTimeline"
          :key="group.key"
          :title="group.title"
          :matches="group.matches"
          :is-today="group.isToday"
          :is-yesterday="group.isYesterday"
          :is-tomorrow="group.isTomorrow"
        />
        <div v-if="!timelineRaw.length" class="empty-state">暂无比赛数据</div>
      </section>
    </aside>
  </div>
</template>

<style scoped>
.dashboard-grid {
  position: relative;
  isolation: isolate;
  width: 100%;
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 22px;
}

.dashboard-grid > :not(.home-pointer-glow) {
  position: relative;
  z-index: 1;
}

.desktop-side {
  display: grid;
  gap: 24px;
}

.home-main {
  min-width: 0;
}

.section {
  margin-top: 26px;
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
  padding-top: 18px;
  border-top: 1px solid color-mix(in srgb, var(--line) 82%, transparent);
}

.follow-section .section-head {
  margin-bottom: 0;
}

.follow-block {
  display: grid;
  gap: 10px;
  padding-top: 12px;
  border-top: 1px solid color-mix(in srgb, var(--line) 68%, transparent);
}

.follow-block:first-of-type {
  border-top: 0;
  padding-top: 0;
}

.follow-subhead {
  color: var(--muted);
  font-size: 13px;
  font-weight: 750;
}

.followed-match-carousel {
  display: grid;
  gap: 8px;
  min-width: 0;
  overflow: hidden;
}

.followed-match-rail {
  display: grid;
  grid-auto-flow: column;
  grid-auto-columns: 100%;
  min-width: 0;
  overflow-x: auto;
  overflow-y: hidden;
  scroll-snap-type: x mandatory;
  overscroll-behavior-x: contain;
  scroll-behavior: smooth;
  touch-action: pan-x pan-y;
  -webkit-overflow-scrolling: touch;
}

.followed-match-slide {
  min-width: 0;
  scroll-snap-align: start;
  scroll-snap-stop: always;
  opacity: 0.88;
  transform: translateX(0);
  transition: opacity 260ms ease-out;
}

.followed-match-slide.active {
  opacity: 1;
  animation: followed-ticket-in-right 420ms cubic-bezier(0.22, 0.72, 0.2, 1) both;
}

.followed-match-slide.active.from-left {
  animation-name: followed-ticket-in-left;
}

.followed-match-slide.active.from-right {
  animation-name: followed-ticket-in-right;
}

.followed-match-slide :deep(.match-ticket-row) {
  padding-left: 0;
  border-left: 0;
}

.followed-match-slide.active :deep(.match-ticket-row) {
  padding-left: 12px;
  border-left: 3px solid color-mix(in srgb, var(--primary) 48%, transparent);
}

.followed-match-slide :deep(.match-ticket-row::before) {
  left: 0;
}

@keyframes followed-ticket-in-right {
  from {
    opacity: 0.72;
    transform: translateX(18px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes followed-ticket-in-left {
  from {
    opacity: 0.72;
    transform: translateX(-18px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.followed-team-rail {
  display: flex;
  gap: 10px;
  overflow-x: auto;
  padding: 2px 0 4px;
  scroll-snap-type: x proximity;
}

.team-chip {
  min-width: 118px;
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  gap: 9px;
  padding: 7px 11px;
  border: 1px solid color-mix(in srgb, var(--line) 76%, transparent);
  border-radius: 999px;
  color: var(--text);
  background: transparent;
  cursor: pointer;
  scroll-snap-align: start;
  transition: border-color 160ms ease-out, background 160ms ease-out, transform 160ms ease-out;
}

.team-chip:active {
  transform: scale(0.98);
}

.team-chip span {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 13px;
  font-weight: 750;
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

.stage-strip {
  display: grid;
  gap: 14px;
  margin-top: 18px;
  padding: 18px 0;
  border-top: 1px solid color-mix(in srgb, var(--line) 82%, transparent);
  border-bottom: 1px solid color-mix(in srgb, var(--line) 82%, transparent);
  background: transparent;
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
  height: 8px;
  overflow: hidden;
  border-radius: 999px;
  background: color-mix(in srgb, var(--line) 52%, transparent);
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

.stage-stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.stage-stat {
  min-width: 0;
  min-height: 74px;
  display: grid;
  align-content: center;
  gap: 7px;
  padding: 4px 10px;
  border-left: 1px solid color-mix(in srgb, var(--line) 72%, transparent);
  background: transparent;
}

.stage-stat:first-child {
  padding-left: 0;
  border-left: 0;
}

.stage-stat span {
  overflow: hidden;
  color: var(--muted);
  font-size: 12px;
  font-weight: 700;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stage-stat strong {
  color: var(--text);
  font-size: 22px;
  font-weight: 850;
  line-height: 1;
}

.stage-context {
  min-height: 86px;
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px 0 0;
  border-top: 1px solid color-mix(in srgb, var(--line) 72%, transparent);
  background: transparent;
}

.stage-context .material-symbols-outlined {
  width: 36px;
  height: 36px;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: 999px;
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 9%, transparent);
  font-size: 20px;
}

.stage-context b {
  display: block;
  color: var(--text);
  font-size: 14px;
  font-weight: 800;
}

.stage-context p {
  margin: 6px 0 0;
  color: var(--muted);
  font-size: 13px;
  line-height: 1.55;
}

.timeline-desktop {
  display: none;
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

.standings-panel {
  overflow: hidden;
  border-top: 1px solid color-mix(in srgb, var(--line) 76%, transparent);
  border-bottom: 1px solid color-mix(in srgb, var(--line) 76%, transparent);
  background: transparent;
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
  .followed-match-slide,
  .followed-match-slide.active {
    animation: none;
    transition: none;
    opacity: 1;
    transform: none;
  }

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
  border: 1px solid color-mix(in srgb, var(--line) 76%, transparent);
  border-radius: 999px;
  color: var(--text);
  background: transparent;
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

@media (max-width: 520px) {
  .stage-strip {
    padding: 15px;
  }

  .stage-stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .stage-stat {
    min-height: 68px;
  }

  .stage-context {
    min-height: 0;
  }
}

@media (min-width: 768px) {
  .followed-team-rail {
    flex-wrap: wrap;
    overflow: visible;
  }
}

@media (min-width: 1024px) {
  .dashboard-grid {
    max-width: 920px;
    margin: 0 auto;
    grid-template-columns: minmax(0, 580px) minmax(280px, 300px);
    align-items: start;
    justify-content: center;
    gap: 28px;
  }

  .home-main {
    display: contents;
  }

  .home-countdown {
    grid-column: 1;
    grid-row: 1;
  }

  .home-countdown {
    padding: 18px;
  }

  .home-countdown :deep(.countdown-time) {
    gap: 8px;
    margin-top: 20px;
  }

  .home-countdown :deep(.flip-stack) {
    height: 92px;
    border-radius: 12px;
  }

  .home-countdown :deep(.flip-value) {
    font-size: 48px;
  }

  .home-countdown :deep(.digits-medium .flip-value) {
    font-size: 40px;
  }

  .home-countdown :deep(.digits-tight .flip-value) {
    font-size: 34px;
  }

  .stage-strip {
    grid-column: 1;
    grid-row: 2;
    align-self: stretch;
    grid-template-rows: auto auto auto auto;
    padding: 16px 0;
  }

  .stage-stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 0;
    border-top: 1px solid color-mix(in srgb, var(--line) 72%, transparent);
  }

  .stage-stat {
    min-height: 64px;
    padding: 12px 0;
    border-left: 0;
    border-top: 1px solid color-mix(in srgb, var(--line) 64%, transparent);
  }

  .stage-stat:nth-child(-n + 2) {
    border-top: 0;
  }

  .stage-stat:nth-child(2n) {
    padding-left: 18px;
    border-left: 1px solid color-mix(in srgb, var(--line) 72%, transparent);
  }

  .stage-context {
    min-height: 0;
    padding-top: 12px;
  }

  .timeline-mobile {
    grid-column: 1;
    grid-row: 3;
  }

  .desktop-side {
    display: contents;
  }

  .standings-section {
    grid-column: 2;
    grid-row: 1;
  }

  .follow-section {
    grid-column: 2;
    grid-row: 2;
    margin-top: 0;
    padding-top: 2px;
  }

  .follow-section .section-head {
    align-items: flex-start;
  }

  .follow-section .section-head h2 {
    font-size: 18px;
  }

  .followed-match-rail {
    grid-auto-columns: 100%;
  }

  .followed-match-slide.active :deep(.match-ticket-row) {
    padding-left: 10px;
  }

  .followed-match-slide :deep(.match-top) {
    gap: 6px;
  }

  .followed-match-slide :deep(.tag) {
    max-width: 148px;
    padding-inline: 9px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .followed-match-slide :deep(.teams-line) {
    gap: 8px;
    margin: 12px 0 10px;
  }

  .followed-match-slide :deep(.team-side) {
    min-height: 48px;
    gap: 7px;
  }

  .followed-match-slide :deep(.team-flag-action) {
    width: 32px;
    height: 32px;
  }

  .followed-match-slide :deep(.team-name) {
    font-size: 14px;
  }

  .followed-match-slide :deep(.score) {
    min-width: 50px;
    font-size: 24px;
  }

  .followed-match-slide :deep(.vs) {
    min-width: 46px;
    height: 34px;
  }

  .followed-match-slide :deep(.match-bottom) {
    gap: 8px;
  }

  .followed-match-slide :deep(.where) {
    font-size: 12px;
  }

  .timeline-mobile {
    display: none;
  }

  .timeline-desktop {
    display: block;
  }

  .timeline-main {
    grid-column: 1;
    grid-row: 3;
    margin-top: 18px;
  }

  .timeline-desktop .section-head {
    margin-bottom: 10px;
  }

  .timeline-desktop :deep(.timeline-day) {
    margin-top: 12px;
  }
}

@media (min-width: 1200px) {
  .dashboard-grid {
    max-width: 960px;
    grid-template-columns: minmax(0, 620px) minmax(280px, 300px);
    gap: 32px;
  }
}
</style>
