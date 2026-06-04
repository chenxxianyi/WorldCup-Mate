<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiGetTeamMatches } from '@/api/teams'
import { apiGetGroupStandings } from '@/api/standings'
import MatchCard from '@/components/common/MatchCard.vue'
import StandingTable from '@/components/common/StandingTable.vue'
import TeamFlag from '@/components/common/TeamFlag.vue'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useTeamStore } from '@/stores/useTeamStore'
import { normalizeMatch, type Match } from '@/types/match'
import { normalizeStanding, type Standing } from '@/types/standing'
import type { Team } from '@/types/team'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const fav = useFavoriteStore()
const teamStore = useTeamStore()

const team = ref<Team | null>(null)
const matches = ref<Match[]>([])
const standings = ref<Standing[]>([])
const loading = ref(false)
const error = ref('')

const sortedMatches = computed(() =>
  [...matches.value].sort((a, b) => {
    const ta = new Date(a.kickoff_time_utc).getTime()
    const tb = new Date(b.kickoff_time_utc).getTime()
    return ta - tb
  }),
)

const nextMatch = computed(() => {
  const now = Date.now()
  return sortedMatches.value.find((match) => {
    if (match.status === 'finished' || !match.kickoff_time_utc) return false
    return new Date(match.kickoff_time_utc).getTime() >= now
  }) || sortedMatches.value[0] || null
})

const teamStanding = computed(() =>
  team.value ? standings.value.find((item) => item.team_id === team.value?.id) || null : null,
)

async function loadTeamDetail() {
  const id = Number(route.params.id)
  if (!id) {
    error.value = '球队不存在'
    return
  }

  loading.value = true
  error.value = ''
  team.value = null
  matches.value = []
  standings.value = []

  try {
    const detail = await teamStore.fetchTeamDetail(id)
    if (!detail) {
      error.value = '球队不存在'
      return
    }

    team.value = detail

    const matchRes = await apiGetTeamMatches(id) as any[]
    matches.value = (matchRes || []).map(normalizeMatch)

    if (detail.group_id) {
      const standingRes = await apiGetGroupStandings(detail.group_id) as any[]
      standings.value = (standingRes || []).map(normalizeStanding)
    }

    if (auth.isLoggedIn) {
      await fav.fetchFavoriteTeams()
    }
  } catch {
    error.value = '球队详情加载失败'
  } finally {
    loading.value = false
  }
}

function toggleFollow() {
  if (!team.value) return
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  fav.toggleTeamFollow(team.value.id)
}

onMounted(loadTeamDetail)

watch(() => route.params.id, loadTeamDetail)
</script>

<template>
  <div class="team-detail">
    <div v-if="loading" class="state-text">加载中...</div>
    <div v-else-if="error" class="state-text">{{ error }}</div>

    <template v-else-if="team">
      <article class="card team-hero">
        <button class="back-btn" title="返回" @click="router.back()">
          <span class="material-symbols-outlined">arrow_back</span>
        </button>

        <div class="hero-main">
          <TeamFlag :value="team.flag" :alt="team.name" :fallback="team.code" size="lg" />
          <div class="team-copy">
            <div class="hero-tags">
              <span class="tag blue">{{ team.group_name || '未分组' }}</span>
              <span class="tag">{{ team.continent }}</span>
            </div>
            <h1>{{ team.name }}</h1>
            <p>{{ team.name_en }} · {{ team.code }}</p>
          </div>
        </div>

        <div class="hero-actions">
          <button
            class="pill-btn"
            :class="{ active: fav.isTeamFollowed(team.id) }"
            @click="toggleFollow"
          >
            <span
              class="material-symbols-outlined"
              :style="fav.isTeamFollowed(team.id) ? 'font-variation-settings: \'FILL\' 1' : ''"
            >star</span>
            {{ fav.isTeamFollowed(team.id) ? '已关注' : '关注球队' }}
          </button>
          <button class="pill-btn primary" @click="router.push('/schedule')">
            <span class="material-symbols-outlined">calendar_month</span>
            查看赛程
          </button>
        </div>

        <div class="stat-grid">
          <div class="stat-cell">
            <span>比赛</span>
            <strong>{{ matches.length }}</strong>
          </div>
          <div class="stat-cell">
            <span>小组排名</span>
            <strong>{{ teamStanding ? standings.indexOf(teamStanding) + 1 : '-' }}</strong>
          </div>
          <div class="stat-cell">
            <span>积分</span>
            <strong>{{ teamStanding ? teamStanding.points : '-' }}</strong>
          </div>
        </div>
      </article>

      <section v-if="nextMatch" class="section">
        <div class="section-head">
          <h2>下一场比赛</h2>
          <span>{{ nextMatch.local_kickoff_time || '时间待定' }}</span>
        </div>
        <MatchCard :match="nextMatch" featured />
      </section>

      <section class="section">
        <div class="section-head">
          <h2>球队赛程</h2>
          <span>{{ matches.length }} 场比赛</span>
        </div>
        <div v-if="sortedMatches.length" class="match-list">
          <MatchCard v-for="match in sortedMatches" :key="match.id" :match="match" />
        </div>
        <div v-else class="card empty-card">暂无球队赛程</div>
      </section>

      <section v-if="standings.length" class="section">
        <div class="section-head">
          <h2>所在小组积分</h2>
          <span>{{ team.group_name }}</span>
        </div>
        <StandingTable :standings="standings" show-status />
      </section>
    </template>
  </div>
</template>

<style scoped>
.team-detail {
  display: grid;
  gap: 18px;
}

.team-hero {
  position: relative;
  padding: 18px;
}

.back-btn {
  position: absolute;
  top: 14px;
  right: 14px;
  width: 36px;
  height: 36px;
  display: grid;
  place-items: center;
  border: 1px solid var(--line);
  border-radius: 999px;
  color: var(--muted);
  background: var(--card);
  cursor: pointer;
}

.back-btn .material-symbols-outlined {
  font-size: 20px;
}

.hero-main {
  display: flex;
  align-items: center;
  gap: 14px;
  padding-right: 42px;
}

.team-copy {
  min-width: 0;
}

.hero-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;
}

.team-copy h1 {
  margin: 0;
  overflow-wrap: anywhere;
  font-size: 24px;
  line-height: 1.15;
  font-weight: 850;
}

.team-copy p {
  margin: 6px 0 0;
  color: var(--muted);
  font-size: 13px;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 18px;
}

.hero-actions .material-symbols-outlined {
  font-size: 18px;
  vertical-align: -4px;
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  margin-top: 16px;
}

.stat-cell {
  padding: 12px;
  border-radius: 14px;
  background: var(--card-soft);
}

.stat-cell span {
  display: block;
  color: var(--muted);
  font-size: 12px;
}

.stat-cell strong {
  display: block;
  margin-top: 4px;
  font-size: 18px;
  font-weight: 850;
}

.section {
  display: grid;
  gap: 12px;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.section-head h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 750;
}

.section-head span {
  color: var(--muted);
  font-size: 13px;
}

.match-list {
  display: grid;
  gap: 12px;
}

.empty-card,
.state-text {
  padding: 24px;
  text-align: center;
  color: var(--muted);
}

@media (max-width: 380px) {
  .stat-grid {
    grid-template-columns: 1fr;
  }
}
</style>
