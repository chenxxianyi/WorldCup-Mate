<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiGetCompetitionStandings } from '@/api/competitions'
import { apiGetGroupStandings } from '@/api/standings'
import { apiGetTeamMatches } from '@/api/teams'
import TeamBadge from '@/components/theme/TeamBadge.vue'
import ThemeIcon from '@/components/theme/ThemeIcon.vue'
import ThemeMatchCard from '@/components/theme/ThemeMatchCard.vue'
import { matchToThemeMatch, teamToThemeTeam } from '@/data/themeAdapters'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import { useTeamStore } from '@/stores/useTeamStore'
import { normalizeMatch, type Match } from '@/types/match'
import { normalizeLeagueStanding, normalizeStanding, type LeagueStanding, type Standing } from '@/types/standing'
import type { Team } from '@/types/team'

const route = useRoute()
const router = useRouter()
const theme = useLeagueThemeStore()
const auth = useAuthStore()
const favorites = useFavoriteStore()
const teamStore = useTeamStore()
const team = ref<Team | null>(null)
const matches = ref<Match[]>([])
const groupStanding = ref<Standing | null>(null)
const leagueStanding = ref<LeagueStanding | null>(null)
const loading = ref(false)
const error = ref('')

const displayTeam = computed(() => team.value ? teamToThemeTeam(team.value) : null)
const sortedMatches = computed(() => [...matches.value].sort((a, b) => new Date(a.kickoff_time_utc).getTime() - new Date(b.kickoff_time_utc).getTime()))
const nextMatch = computed(() => sortedMatches.value.find((item) => item.status !== 'finished') || sortedMatches.value[0] || null)
const displayMatches = computed(() => sortedMatches.value.map(matchToThemeMatch))
const position = computed(() => leagueStanding.value?.position || groupStanding.value?.rank || null)
const points = computed(() => leagueStanding.value?.points ?? groupStanding.value?.points ?? null)
const goalsFor = computed(() => leagueStanding.value?.goals_for ?? groupStanding.value?.goals_for ?? null)

async function loadTeam() {
  const id = Number(route.params.id)
  if (!id) { error.value = '球队不存在'; return }
  loading.value = true
  error.value = ''
  team.value = null
  matches.value = []
  groupStanding.value = null
  leagueStanding.value = null
  try {
    const detail = await teamStore.fetchTeamDetail(id)
    if (!detail) throw new Error('球队不存在')
    team.value = detail
    const matchRows = await apiGetTeamMatches(id) as any[]
    matches.value = (matchRows || []).map(normalizeMatch)
    if (theme.current.slug === 'wc' && detail.group_id) {
      const rows = await apiGetGroupStandings(detail.group_id) as any[]
      const standings = rows.map(normalizeStanding)
      groupStanding.value = standings.find((item) => item.team_id === id) || null
    } else {
      const rows = await apiGetCompetitionStandings(theme.currentCode, { type: 'total' }) as any[]
      leagueStanding.value = rows.map(normalizeLeagueStanding).find((item) => item.team_id === id) || null
    }
    if (auth.isLoggedIn) await favorites.fetchFavoriteTeams()
  } catch (reason) {
    error.value = reason instanceof Error ? reason.message : '球队详情加载失败'
  } finally {
    loading.value = false
  }
}

async function toggleFollow() {
  if (!team.value) return
  if (!auth.isLoggedIn) {
    router.push({ path: '/login', query: { redirect: route.fullPath } })
    return
  }
  const was = favorites.isTeamFollowed(team.value.id)
  const saved = await favorites.toggleTeamFollow(team.value.id)
  theme.showToast(saved ? (was ? '已取消关注' : `已关注${team.value.name}`) : '关注操作失败，请稍后重试')
}

onMounted(loadTeam)
watch(() => route.params.id, loadTeam)
watch(() => theme.currentCode, () => router.push('/teams'))
</script>

<template>
  <div class="page-view">
    <div class="back-row"><button class="back-button" type="button" @click="router.push('/teams')"><ThemeIcon name="back" /> 返回球队</button></div>
    <div v-if="loading" class="page-state"><span class="state-spinner" />正在加载球队详情</div>
    <article v-else-if="error || !team || !displayTeam" class="card empty-compact"><span class="empty-art"><ThemeIcon name="shield" /></span><span class="empty-copy"><h3>无法显示球队</h3><p>{{ error || '球队数据不存在。' }}</p></span></article>
    <template v-else>
      <section class="team-hero"><div class="team-hero-content"><TeamBadge :team="displayTeam" size="large" /><div class="team-hero-copy"><p class="eyebrow" style="color: var(--competition-accent)">{{ theme.current.en }}</p><h1>{{ team.name }}</h1><p>{{ team.name_en }} · {{ team.code }}</p><div class="team-detail-meta"><span class="hero-meta-pill">当前排名 {{ position || '—' }}</span><span class="hero-meta-pill">主场 {{ team.venue || '待定' }}</span><span class="hero-meta-pill">{{ team.group_name || team.country || team.continent }}</span></div></div><button class="team-hero-follow" :class="{ active: favorites.isTeamFollowed(team.id) }" type="button" @click="toggleFollow"><ThemeIcon name="star" />{{ favorites.isTeamFollowed(team.id) ? '已关注' : '关注' }}</button></div></section>
      <div class="detail-grid section">
        <div class="stack">
          <article class="card detail-panel"><h3>赛季表现</h3><div class="quick-grid" style="margin-top: 0"><span class="quick-card"><small class="muted">积分</small><strong>{{ points ?? '—' }}</strong></span><span class="quick-card"><small class="muted">进球</small><strong>{{ goalsFor ?? '—' }}</strong></span><span class="quick-card"><small class="muted">比赛</small><strong>{{ matches.length }}</strong></span></div></article>
          <article class="card detail-panel"><h3>球队资料</h3><div class="info-list"><span><small>所在赛事</small><strong>{{ theme.current.name }}</strong></span><span><small>国家 / 地区</small><strong>{{ team.country || team.continent }}</strong></span><span><small>球队代码</small><strong>{{ team.code || '—' }}</strong></span><span><small>主场</small><strong>{{ team.venue || '待定' }}</strong></span></div></article>
        </div>
        <section><div class="section-heading"><div><p class="eyebrow">NEXT MATCH</p><h2>下一场比赛</h2></div></div><ThemeMatchCard v-if="nextMatch" :match="matchToThemeMatch(nextMatch)" /><article v-else class="card empty-mini">暂无球队赛程</article></section>
      </div>
      <section v-if="displayMatches.length > 1" class="section"><div class="section-heading"><div><p class="eyebrow">FIXTURES</p><h2>球队赛程</h2></div><span>{{ displayMatches.length }} 场</span></div><div class="match-list"><ThemeMatchCard v-for="match in displayMatches" :key="match.id" :match="match" /></div></section>
    </template>
  </div>
</template>
