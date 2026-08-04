<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import TeamBadge from '@/components/theme/TeamBadge.vue'
import ThemeIcon from '@/components/theme/ThemeIcon.vue'
import { teamToThemeTeam } from '@/data/themeAdapters'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import { useTeamStore } from '@/stores/useTeamStore'
import type { Team } from '@/types/team'

const route = useRoute()
const router = useRouter()
const theme = useLeagueThemeStore()
const auth = useAuthStore()
const favorites = useFavoriteStore()
const teamStore = useTeamStore()
const query = ref('')
const error = ref('')
const teams = computed(() => {
  const keyword = query.value.trim().toLowerCase()
  if (!keyword) return teamStore.teams
  return teamStore.teams.filter((team) => [team.name, team.name_en, team.code, team.country, team.group_name].join(' ').toLowerCase().includes(keyword))
})

async function loadTeams() {
  error.value = ''
  try {
    await theme.initialize()
    const params: Record<string, unknown> = { page_size: 100, teamType: theme.current.slug === 'wc' ? 'national' : 'club' }
    if (theme.current.slug !== 'wc' && theme.competition.current?.country) params.country = theme.competition.current.country
    await teamStore.fetchTeams(params)
    if (auth.isLoggedIn) await favorites.fetchFavoriteTeams()
  } catch (reason) {
    error.value = reason instanceof Error ? reason.message : '球队数据加载失败'
  }
}

async function toggleFollow(team: Team) {
  if (!auth.isLoggedIn) {
    theme.showToast('登录后才能关注球队')
    router.push({ path: '/login', query: { redirect: route.fullPath } })
    return
  }
  const wasFollowed = favorites.isTeamFollowed(team.id)
  const saved = await favorites.toggleTeamFollow(team.id)
  theme.showToast(saved ? (wasFollowed ? '已取消关注' : `已关注${team.name}`) : '关注操作失败，请稍后重试')
}

onMounted(loadTeams)
watch(() => theme.currentCode, () => { query.value = ''; loadTeams() })
</script>

<template>
  <div class="page-view">
    <header class="page-heading"><div><p class="eyebrow">{{ theme.current.en }}</p><h1>{{ theme.current.slug === 'wc' ? '国家队' : '俱乐部' }}</h1></div><p class="muted">{{ teamStore.teams.length }} 支球队 · 数据来自赛事中心</p></header>
    <label class="search-box"><ThemeIcon name="search" /><input v-model="query" type="search" placeholder="搜索球队名称、国家或缩写" aria-label="搜索球队" /></label>
    <section class="section">
      <div v-if="teamStore.loading" class="page-state"><span class="state-spinner" />正在加载球队</div>
      <div v-else-if="teams.length" class="teams-grid">
        <article v-for="team in teams" :key="team.id" class="team-card">
          <button class="team-card-open" type="button" @click="router.push(`/teams/${team.id}`)">
            <TeamBadge :team="teamToThemeTeam(team)" /><h3>{{ team.name }}</h3><p>{{ team.country || team.continent }} · {{ team.code }}</p>
            <span class="team-card-footer"><span>{{ team.group_name || team.venue || '查看球队资料' }}</span><ThemeIcon name="arrow" /></span>
          </button>
          <button class="team-follow" :class="{ active: favorites.isTeamFollowed(team.id) }" type="button" :aria-label="favorites.isTeamFollowed(team.id) ? `取消关注${team.name}` : `关注${team.name}`" @click="toggleFollow(team)"><ThemeIcon name="star" /></button>
        </article>
      </div>
      <article v-else class="card empty-compact"><span class="empty-art"><ThemeIcon name="search" /></span><span class="empty-copy"><h3>{{ query ? '没有找到球队' : '暂无球队数据' }}</h3><p>{{ error || '请调整搜索条件，或先在后台同步当前赛事球队。' }}</p></span></article>
    </section>
  </div>
</template>
