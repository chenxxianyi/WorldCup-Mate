<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { apiListMatches } from '@/api/matches'
import ThemeIcon from '@/components/theme/ThemeIcon.vue'
import ThemeMatchCard from '@/components/theme/ThemeMatchCard.vue'
import { localDateLabel, matchToThemeMatch, type ThemeMatch } from '@/data/themeAdapters'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import { normalizeMatch } from '@/types/match'

const router = useRouter()
const theme = useLeagueThemeStore()
const auth = useAuthStore()
const favorites = useFavoriteStore()
const filter = ref('全部')
const matches = ref<ThemeMatch[]>([])
const loading = ref(false)
const error = ref('')
const filters = computed(() => theme.current.slug === 'wc' ? ['全部', '小组赛', '淘汰赛', '收藏'] : ['全部', '第 11 轮', '第 12 轮', '直播', '收藏'])
const groupedMatches = computed(() => {
  const groups = new Map<string, ThemeMatch[]>()
  matches.value.forEach((match) => groups.set(match.kickoffKey, [...(groups.get(match.kickoffKey) || []), match]))
  return [...groups.entries()].map(([key, items]) => ({ key, label: localDateLabel(key), matches: items }))
})

function queryForFilter() {
  const params: Record<string, unknown> = { page: 1, page_size: 100 }
  if (filter.value === '小组赛') params.stage = 'group'
  if (filter.value === '淘汰赛') params.stage = 'knockout'
  if (filter.value === '直播') params.status = 'live'
  const round = filter.value.match(/第\s*(\d+)\s*轮/)
  if (round) params.matchday = Number(round[1])
  return params
}

async function loadMatches() {
  loading.value = true
  error.value = ''
  try {
    if (filter.value === '收藏' && !auth.isLoggedIn) {
      matches.value = []
      return
    }
    if (filter.value === '收藏') await favorites.fetchFavoriteMatches()
    const params = queryForFilter()
    const response = await apiListMatches(params) as any
    let rawRows = response.list || response || []
    const total = Number(response.total ?? rawRows.length)
    if (rawRows.length && rawRows.length < total) {
      const pageSize = Number(response.page_size || params.page_size || 100)
      const pageCount = Math.ceil(total / pageSize)
      const remaining = await Promise.all(Array.from({ length: pageCount - 1 }, (_, index) =>
        apiListMatches({ ...params, page: index + 2, page_size: pageSize }) as Promise<any>,
      ))
      rawRows = rawRows.concat(...remaining.map((page) => page.list || page || []))
    }
    let rows = rawRows.map(normalizeMatch)
    if (filter.value === '收藏') rows = rows.filter((match: any) => favorites.isMatchFavorite(match.id))
    matches.value = rows.map(matchToThemeMatch)
  } catch (reason) {
    matches.value = []
    error.value = reason instanceof Error ? reason.message : '赛程加载失败'
  } finally {
    loading.value = false
  }
}

async function selectFilter(value: string) {
  filter.value = value
  await loadMatches()
}

onMounted(loadMatches)
watch(() => theme.currentCode, () => { filter.value = '全部'; loadMatches() })
</script>

<template>
  <div class="page-view">
    <header class="page-heading"><div><p class="eyebrow">{{ theme.current.en }}</p><h1>赛程中心</h1></div><p class="muted">当地时间 · 自动转换为你的时区</p></header>
    <section class="card filter-card"><div class="filter-row"><button v-for="item in filters" :key="item" class="chip" :class="{ active: item === filter }" type="button" @click="selectFilter(item)">{{ item }}</button></div></section>

    <div v-if="loading" class="page-state"><span class="state-spinner" />正在加载赛程</div>
    <template v-else-if="groupedMatches.length">
      <section v-for="group in groupedMatches" :key="group.key" class="schedule-day">
        <div class="day-marker"><span>{{ group.label.weekday }}</span><strong>{{ group.label.day }}</strong><span>{{ group.label.month }}</span></div>
        <div class="match-list"><ThemeMatchCard v-for="match in group.matches" :key="match.id" :match="match" /></div>
      </section>
    </template>
    <article v-else class="card empty-compact schedule-empty"><span class="empty-art"><ThemeIcon :name="filter === '收藏' ? 'star' : 'calendar'" /></span><span class="empty-copy"><h3>{{ filter === '收藏' && !auth.isLoggedIn ? '登录后查看收藏赛程' : '暂无符合条件的比赛' }}</h3><p>{{ error || '可以切换筛选条件或等待后台同步比赛。' }}</p></span><button v-if="filter === '收藏' && !auth.isLoggedIn" class="primary-button" type="button" @click="router.push({ path: '/login', query: { redirect: '/schedule' } })">登录</button></article>
  </div>
</template>
