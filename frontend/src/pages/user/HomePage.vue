<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import Countdown from '@/components/common/Countdown.vue'
import MatchCard from '@/components/common/MatchCard.vue'
import { useMatchStore } from '@/stores/useMatchStore'
import { useTeamStore } from '@/stores/useTeamStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useAuthStore } from '@/stores/useAuthStore'
import { apiGetGroupStandings } from '@/api/standings'
import { apiGetTournamentProgress, apiGetUpcomingMatches } from '@/api/matches'
import { normalizeStanding, type Standing } from '@/types/standing'
import TeamFlag from '@/components/common/TeamFlag.vue'

const matchStore = useMatchStore()
const teamStore = useTeamStore()
const fav = useFavoriteStore()
const auth = useAuthStore()

const groups = ['Group A', 'Group B', 'Group C', 'Group D', 'Group E', 'Group F', 'Group G', 'Group H', 'Group I', 'Group J', 'Group K', 'Group L']
const activeGroupIndex = ref(0)
const activeGroupName = computed(() => groups[activeGroupIndex.value])
const groupStandings = ref<Standing[]>([])
const groupSwipeStartX = ref(0)
const groupSwipeStartY = ref(0)
const groupSwipeActive = ref(false)
const groupSlideDirection = ref<'left' | 'right'>('left')
const groupTransitionName = computed(() =>
  groupSlideDirection.value === 'left' ? 'standings-next' : 'standings-prev'
)

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

  if (Math.abs(deltaX) < minDistance || Math.abs(deltaY) > maxVerticalDrift) {
    return
  }

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

watch(activeGroupIndex, loadGroupStandings)

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
  progress: 0
})

interface NextMatch {
  id: number
  home_team: { name: string; flag_url: string }
  away_team: { name: string; flag_url: string }
  kickoff_time_utc: string
  stadium: { name: string }
  city: { name: string }
  recommend_tag: string
}

const nextMatch = ref<NextMatch | null>(null)

const nextMatchLocalTime = computed(() => {
  if (!nextMatch.value) return ''
  const date = new Date(nextMatch.value.kickoff_time_utc)
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${month}月${day}日 ${hours}:${minutes}`
})

const todayMatches = computed(() => matchStore.todayMatches)
const recommended = computed(() => matchStore.recommendedMatches)
const followedTeams = computed(() =>
  teamStore.teams.filter((t) => fav.isTeamFollowed(t.id))
)

onMounted(async () => {
  matchStore.fetchTodayMatches()
  matchStore.fetchRecommendedMatches()
  teamStore.fetchTeams()
  if (auth.isLoggedIn) {
    fav.fetchFavoriteTeams()
  }
  loadGroupStandings()
  try {
    const res = await apiGetTournamentProgress()
    if (res) {
      progress.value = res
    }
  } catch {}
  try {
    const res = await apiGetUpcomingMatches() as any[]
    if (res && res.length > 0) {
      nextMatch.value = res[0]
    }
  } catch {}
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
            <h2>{{ nextMatch.home_team.name }} vs {{ nextMatch.away_team.name }}</h2>
            <p>{{ nextMatchLocalTime }} 开球 · 本地时间</p>
          </div>
          <span v-if="nextMatch.recommend_tag" class="tag gold">{{ nextMatch.recommend_tag }}</span>
        </div>
      </Countdown>
      <article v-else class="card" style="padding: 20px; text-align: center; color: var(--muted)">
        暂无即将到来的比赛
      </article>

      <!-- Stage Progress -->
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

      <!-- Today Matches -->
      <section class="section">
        <div class="section-head">
          <h2>今日比赛</h2>
          <span>{{ todayMatches.length }} 场</span>
        </div>
        <div class="match-strip">
          <MatchCard
            v-for="m in todayMatches"
            :key="m.id"
            :match="m"
          />
        </div>
      </section>

      <!-- Recommended -->
      <section class="section">
        <div class="section-head">
          <h2>热门推荐</h2>
        </div>
        <div class="stack">
          <MatchCard
            v-for="m in recommended"
            :key="m.id"
            :match="m"
            featured
          />
        </div>
      </section>
    </div>

    <!-- Desktop Sidebar -->
    <aside class="desktop-side">
      <section class="section" style="margin-top: 0">
        <div class="section-head">
          <h2>我的关注</h2>
          <span>{{ followedTeams.length }} 支球队</span>
        </div>
        <div class="stack">
          <article v-for="t in followedTeams" :key="t.id" class="card profile-card">
            <TeamFlag :value="t.flag" :alt="t.name" :fallback="t.code" size="lg" />
            <div>
              <h2>{{ t.name }}</h2>
              <p>{{ t.group_name }}</p>
            </div>
          </article>
        </div>
      </section>

      <section class="section">
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

.stage-head span:last-child {
  color: var(--primary);
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

.profile-card {
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 13px;
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
  background: var(--card);
  color: var(--text);
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
  background: var(--primary);
  color: #fff;
  border-color: var(--primary);
}

.group-label {
  min-width: 70px;
  text-align: center;
  font-size: 13px;
  font-weight: 600;
  color: var(--primary);
  cursor: pointer;
}

@media (min-width: 768px) {
  .dashboard-grid {
    grid-template-columns: minmax(0, 1.9fr) minmax(300px, 0.8fr);
    align-items: start;
  }

  .match-strip {
    grid-auto-columns: minmax(320px, 1fr);
    grid-template-columns: repeat(2, minmax(0, 1fr));
    grid-auto-flow: row;
    overflow: visible;
    padding-bottom: 0;
  }
}
</style>
