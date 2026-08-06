<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiListMatches } from '@/api/matches'
import { apiCompetitionOverview } from '@/api/competitions'
import ThemeIcon from '@/components/theme/ThemeIcon.vue'
import ThemeMatchCard from '@/components/theme/ThemeMatchCard.vue'
import { localDateLabel, matchToThemeMatch, type ThemeMatch } from '@/data/themeAdapters'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import { useLivePolling, pollingStatusFrom } from '@/composables/useLivePolling'
import { useGeneration } from '@/composables/useRequestGuard'
import { isCancel } from '@/types/common'
import { normalizeMatch } from '@/types/match'

const route = useRoute()
const router = useRouter()
const theme = useLeagueThemeStore()
const auth = useAuthStore()
const favorites = useFavoriteStore()

const filter = ref<string>('全部')
const matches = ref<ThemeMatch[]>([])
const loading = ref(false)
const error = ref('')
const hasMore = ref(false)
const page = ref(1)
const PAGE_SIZE = 40

// DATA-09: matchday ranges are derived from the backend overview, not
// hard-coded. Once loaded, `roundFilters` holds the available rounds.
const roundFilters = ref<string[]>([])
const overviewLoading = ref(false)

const isLeague = computed(() => theme.current.format === 'league')
const filters = computed(() => {
  if (theme.current.slug === 'wc') return ['全部', '小组赛', '淘汰赛', '收藏']
  const rounds = roundFilters.value.length ? roundFilters.value : ['第 1 轮']
  return ['全部', ...rounds, '直播', '收藏']
})

const groupedMatches = computed(() => {
  const groups = new Map<string, ThemeMatch[]>()
  matches.value.forEach((match) => groups.set(match.kickoffKey, [...(groups.get(match.kickoffKey) || []), match]))
  return [...groups.entries()].map(([key, items]) => ({ key, label: localDateLabel(key), matches: items }))
})

// ---------- UI state derives from the URL query (DATA-09) ----------
const qsValue = () => route.query.filter as string | undefined
const roundFromQuery = () => (route.query.matchday ? Number(route.query.matchday) : null)

watch(() => route.query, () => {
  filter.value = qsValue() || '全部'
}, { immediate: true })

// Toggle a filter; push it to the URL query so refresh/share keeps it.
function selectFilter(value: string) {
  const query: Record<string, any> = { ...route.query }
  if (value === '全部') delete query.filter
  else query.filter = value
  const round = value.match(/第\s*(\d+)\s*轮/)
  if (round) query.matchday = Number(round[1])
  else delete query.matchday
  router.push({ query }).catch(() => {})
}

function queryForFilter() {
  const params: Record<string, unknown> = { page: page.value, page_size: PAGE_SIZE }
  if (filter.value === '小组赛') params.stage = 'group'
  if (filter.value === '淘汰赛') params.stage = 'knockout'
  if (filter.value === '直播') params.status = 'live'
  const round = filter.value.match(/第\s*(\d+)\s*轮/)
  if (round) params.matchday = Number(round[1])
  const league = theme.current
  if (league.id) {
    params.competitionId = league.id
    if (league.season) params.season = league.season
  }
  return params
}

// LIVE-02: abort in-flight list requests when a newer load starts or the
// component unmounts (belt & braces; generation guard keeps correctness).
let listController: AbortController | null = null
function freshSignal(): AbortSignal {
  listController?.abort()
  listController = new AbortController()
  return listController.signal
}
onBeforeUnmount(() => listController?.abort())
onBeforeUnmount(() => gen.bump()) // LIVE-02: drop in-flight flows writing refs of this page

// PERF-04: incremental loading. First load pulls one page; subsequent
// load calls stack pages instead of re-fetching everything.
async function loadMatches(quiet = false, append = false) {
  const g = gen.next()
  if (!quiet && !append) loading.value = true
  error.value = ''
  try {
    if (filter.value === '收藏' && !auth.isLoggedIn) {
      matches.value = []
      hasMore.value = false
      return
    }
    if (filter.value === '收藏') await favorites.fetchFavoriteMatches()
    const params = queryForFilter()
    const signal = freshSignal()
    const response = await apiListMatches(params, { signal })
    if (!gen.isCurrent(g)) return // stale: a newer load/switch won

    const rows = response.list.map(normalizeMatch)
    let final = rows
    if (filter.value === '收藏') final = rows.filter((match: any) => favorites.isMatchFavorite(match.id))

    if (append) {
      matches.value = [...matches.value, ...final.map(matchToThemeMatch)]
    } else {
      matches.value = final.map(matchToThemeMatch)
    }
    // Determine whether more pages remain.
    const total = Number(response.total ?? rows.length)
    const pageSize = Number(response.page_size || PAGE_SIZE)
    hasMore.value = page.value * pageSize < total
  } catch (reason) {
    if (isCancel(reason)) return
    if (!quiet) {
      // On an append failure keep what we already have; only a first page
      // failure clears the list.
      if (!append) matches.value = []
      error.value = reason instanceof Error ? reason.message : '赛程加载失败'
    }
  } finally {
    if (!quiet && !append) loading.value = false
    polling.schedule() // LIVE-01: (re)arm the adaptive timer
  }
}

async function loadMore() {
  page.value += 1
  await loadMatches(true, true)
}

function resetToTop() {
  router.push({ query: { ...route.query, filter: route.query.filter } }).catch(() => {})
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// LIVE-01: adaptive polling of the schedule.
const gen = useGeneration()
const polling = useLivePolling(
  () => pollingStatusFrom(matches.value),
  () => loadMatches(true),
  () => matches.value.length > 0,
)

// Load overview (available rounds) for leagues on mount & on switch.
async function loadOverview() {
  const competition = theme.current
  if (!competition.code) return
  if (competition.slug === 'wc') return
  overviewLoading.value = true
  try {
    const overview = await apiCompetitionOverview(competition.code)
    const maxMatchday = Number(overview.matchday) || 0
    roundFilters.value = Array.from({ length: maxMatchday }, (_, i) => `第 ${i + 1} 轮`)
  } catch {
    roundFilters.value = []
  } finally {
    overviewLoading.value = false
  }
}

onMounted(() => {
  loadOverview()
  loadMatches()
})
// DATA-09: switching to a league loads its current season/round range.
watch(() => theme.currentCode, () => {
  gen.bump()
  page.value = 1
  roundFilters.value = []
  loadOverview()
  loadMatches()
})
watch(filter, () => {
  page.value = 1
  loadMatches()
})
</script>

<template>
  <div class="page-view">
    <header class="page-heading">
      <div>
        <p class="eyebrow">
          {{ theme.current.en }}
        </p><h1>赛程中心</h1>
      </div><p class="muted">
        当地时间 · 自动转换为你的时区
      </p>
    </header>
    <section class="card filter-card">
      <div class="filter-row">
        <button
          v-for="item in filters"
          :key="item"
          class="chip"
          :class="{ active: item === filter }"
          type="button"
          @click="selectFilter(item)"
        >
          {{ item }}
        </button>
      </div>
    </section>

    <div
      v-if="loading"
      class="page-state"
    >
      <span class="state-spinner" />正在加载赛程
    </div>
    <template v-else-if="groupedMatches.length">
      <section
        v-for="group in groupedMatches"
        :key="group.key"
        class="schedule-day"
      >
        <div class="day-marker">
          <span>{{ group.label.weekday }}</span><strong>{{ group.label.day }}</strong><span>{{ group.label.month }}</span>
        </div>
        <div class="match-list">
          <ThemeMatchCard
            v-for="match in group.matches"
            :key="match.id"
            :match="match"
          />
        </div>
      </section>
      <!-- PERF-04: load-more + back-to-top -->
      <div class="load-more-row">
        <button
          v-if="hasMore"
          class="load-more-button"
          type="button"
          :disabled="loading"
          @click="loadMore"
        >
          {{ loading ? '加载中…' : '加载更多' }}
        </button>
        <button
          v-if="groupedMatches.length > 10"
          class="load-more-button ghost"
          type="button"
          @click="resetToTop"
        >
          返回顶部
        </button>
      </div>
    </template>
    <article
      v-else
      class="card empty-compact schedule-empty"
    >
      <span class="empty-art"><ThemeIcon :name="filter === '收藏' ? 'star' : 'calendar'" /></span><span class="empty-copy"><h3>{{ filter === '收藏' && !auth.isLoggedIn ? '登录后查看收藏赛程' : '暂无符合条件的比赛' }}</h3><p>{{ error || '可以切换筛选条件或等待后台同步比赛。' }}</p></span><button
        v-if="filter === '收藏' && !auth.isLoggedIn"
        class="primary-button"
        type="button"
        @click="router.push({ path: '/login', query: { redirect: '/schedule' } })"
      >
        登录
      </button>
    </article>
  </div>
</template>

<style scoped>
.load-more-row {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 20px;
}

.load-more-button {
  min-height: 40px;
  padding: 0 22px;
  border: 1px solid var(--line);
  border-radius: 20px;
  background: var(--card);
  color: var(--text);
  cursor: pointer;
  font: inherit;
}

.load-more-button:disabled {
  opacity: 0.6;
  cursor: default;
}

.load-more-button.ghost {
  background: transparent;
}
</style>
