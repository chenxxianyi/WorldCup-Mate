<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
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
const groupAStandings = ref<Standing[]>([])

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
  try {
    const res = await apiGetGroupStandings(1) as any[]
    groupAStandings.value = res.map(normalizeStanding)
  } catch {}
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
          <span>importance_level ≥ 1</span>
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
          <span>Group A</span>
        </div>
        <div class="card table-card">
          <table class="standing-table" style="min-width: 0">
            <tbody>
              <tr
                v-for="(s, i) in groupAStandings"
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

.team-cell {
  display: flex;
  align-items: center;
  gap: 8px;
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
