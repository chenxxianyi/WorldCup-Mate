<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import SearchInput from '@/components/common/SearchInput.vue'
import ChipFilter from '@/components/common/ChipFilter.vue'
import MatchCard from '@/components/common/MatchCard.vue'
import SyncStatusBadge from '@/components/common/SyncStatusBadge.vue'
import { apiListMatches } from '@/api/matches'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { normalizeMatch, type Match } from '@/types/match'

const INITIAL_DAYS = 3
const DAYS_PER_LOAD = 2
const PAGE_SIZE = 8
const LOAD_MORE_MIN_MS = 760

const FILTER_ALL = '全部'
const FILTER_TODAY = '今日'
const FILTER_TOMORROW = '明日'
const FILTER_SCHEDULED = '未开始'
const FILTER_KNOCKOUT = '淘汰赛'
const FILTER_FOLLOWED = '只看关注'

const router = useRouter()
const auth = useAuthStore()
const fav = useFavoriteStore()

const search = ref('')
const activeFilter = ref(FILTER_ALL)
const matches = ref<Match[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const errorMessage = ref('')
const page = ref(1)
const total = ref(0)
const reachedEnd = ref(false)
const visibleDateLimit = ref(INITIAL_DAYS)
const loadMoreTarget = ref<HTMLElement | null>(null)
const loadHintVisible = ref(false)

let observer: IntersectionObserver | null = null
let requestToken = 0
let searchTimer: number | undefined

const filterOptions = [
  FILTER_ALL,
  FILTER_TODAY,
  FILTER_TOMORROW,
  FILTER_SCHEDULED,
  FILTER_KNOCKOUT,
  FILTER_FOLLOWED,
  'Group A',
  'Group B',
  'Group C',
  'Group D',
  'Group E',
  'Group F',
  'Group G',
  'Group H',
  'Group I',
  'Group J',
  'Group K',
  'Group L',
]

const followedOnly = computed(() => activeFilter.value === FILTER_FOLLOWED)

const clientFilteredMatches = computed(() => {
  if (!followedOnly.value) return matches.value
  const followedIds = new Set(fav.followedTeamIds)
  return matches.value.filter(
    (match) => followedIds.has(match.home_team_id) || followedIds.has(match.away_team_id),
  )
})

const loadedDateKeys = computed(() => {
  const keys: string[] = []
  for (const match of clientFilteredMatches.value) {
    const key = matchDateKey(match)
    if (!keys.includes(key)) keys.push(key)
  }
  return keys
})

const visibleMatches = computed(() => {
  const allowedKeys = new Set(loadedDateKeys.value.slice(0, visibleDateLimit.value))
  return clientFilteredMatches.value.filter((match) => allowedKeys.has(matchDateKey(match)))
})

const groupedByDate = computed(() => {
  const groups: Array<{ key: string; title: string; matches: Match[] }> = []
  for (const match of visibleMatches.value) {
    const key = matchDateKey(match)
    let group = groups.find((item) => item.key === key)
    if (!group) {
      group = { key, title: dateTitle(key), matches: [] }
      groups.push(group)
    }
    group.matches.push(match)
  }
  return groups
})

const hasBufferedDays = computed(() => loadedDateKeys.value.length > visibleDateLimit.value)
const hasMore = computed(() => {
  if (followedOnly.value && !auth.isLoggedIn) return false
  if (followedOnly.value && auth.isLoggedIn && fav.followedTeamIds.length === 0) return false
  return hasBufferedDays.value || !reachedEnd.value
})

const emptyText = computed(() => {
  if (followedOnly.value && !auth.isLoggedIn) return '登录后可以查看关注球队的赛程'
  if (followedOnly.value && fav.followedTeamIds.length === 0) return '还没有关注球队'
  return '暂无符合条件的比赛'
})

const loadText = computed(() => {
  if (loadingMore.value) return '加载中...'
  if (hasMore.value) {
    return loadHintVisible.value ? '继续下滑，加载更多赛程' : '下滑查看更多赛程'
  }
  return clientFilteredMatches.value.length ? '已加载全部赛程' : ''
})

function dayParam(offset = 0) {
  const d = new Date()
  d.setDate(d.getDate() + offset)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function wait(ms: number) {
  return new Promise((resolve) => window.setTimeout(resolve, ms))
}

function requestParams() {
  const params: Record<string, any> = {
    page: page.value,
    page_size: PAGE_SIZE,
  }

  if (activeFilter.value === FILTER_TODAY) {
    params.date = dayParam()
  } else if (activeFilter.value === FILTER_TOMORROW) {
    params.date = dayParam(1)
  } else if (activeFilter.value.startsWith('Group')) {
    params.groupName = activeFilter.value
  } else if (activeFilter.value === FILTER_KNOCKOUT) {
    params.stage = 'knockout'
  } else if (activeFilter.value === FILTER_SCHEDULED) {
    params.status = 'scheduled'
  }

  const keyword = search.value.trim()
  if (keyword) params.keyword = keyword
  return params
}

function matchDateKey(match: Match) {
  return (match.local_kickoff_time || '').split(' ')[0] || '其他日期'
}

function dateTitle(key: string) {
  const [month, day] = key.split('-')
  if (!month || !day) return key
  return `${Number(month)}月${Number(day)}日`
}

async function prepareFollowedFilter() {
  if (!followedOnly.value) return true
  if (!auth.isLoggedIn) {
    reachedEnd.value = true
    return false
  }
  await fav.fetchFavoriteTeams()
  if (fav.followedTeamIds.length === 0) {
    reachedEnd.value = true
    return false
  }
  return true
}

async function fetchNextPage(token: number) {
  const res = (await apiListMatches(requestParams())) as any
  if (token !== requestToken) return false

  const list = (res.list || res || []).map(normalizeMatch)
  const knownIds = new Set(matches.value.map((match) => match.id))
  matches.value.push(...list.filter((match: Match) => !knownIds.has(match.id)))

  total.value = Number(res.total || matches.value.length)
  const currentPage = Number(res.page || page.value)
  const pageSize = Number(res.page_size || PAGE_SIZE)
  page.value = currentPage + 1
  reachedEnd.value = matches.value.length >= total.value || list.length < pageSize
  return list.length > 0
}

async function fillInitialWindow(token: number) {
  while (
    token === requestToken &&
    loadedDateKeys.value.length < INITIAL_DAYS &&
    !reachedEnd.value
  ) {
    const gotRows = await fetchNextPage(token)
    if (!gotRows) break
  }
}

async function resetAndLoad() {
  const token = ++requestToken
  loading.value = true
  loadingMore.value = false
  errorMessage.value = ''
  page.value = 1
  total.value = 0
  reachedEnd.value = false
  visibleDateLimit.value = INITIAL_DAYS
  loadHintVisible.value = false
  matches.value = []

  try {
    const canLoad = await prepareFollowedFilter()
    if (canLoad) await fillInitialWindow(token)
  } catch (err: any) {
    if (token === requestToken) errorMessage.value = err?.message || '赛程加载失败'
  } finally {
    if (token === requestToken) {
      loading.value = false
      await nextTick()
      observeLoadMoreTarget()
    }
  }
}

async function loadMore() {
  if (loading.value || loadingMore.value || !hasMore.value) return
  loadHintVisible.value = false

  const token = requestToken
  const nextDateLimit = visibleDateLimit.value + DAYS_PER_LOAD
  loadingMore.value = true
  errorMessage.value = ''

  try {
    const minimumLoading = wait(LOAD_MORE_MIN_MS)

    while (
      token === requestToken &&
      loadedDateKeys.value.length < nextDateLimit &&
      !reachedEnd.value
    ) {
      const gotRows = await fetchNextPage(token)
      if (!gotRows) break
    }

    await minimumLoading
    if (token === requestToken) {
      visibleDateLimit.value = nextDateLimit
    }
  } catch (err: any) {
    if (token === requestToken) errorMessage.value = err?.message || '加载更多失败'
  } finally {
    if (token === requestToken) {
      loadingMore.value = false
      await nextTick()
      observeLoadMoreTarget()
    }
  }
}

function observeLoadMoreTarget() {
  if (!observer || !loadMoreTarget.value) return
  observer.disconnect()
  observer.observe(loadMoreTarget.value)
}

function maybeLoadAfterHint() {
  if (!loadMoreTarget.value || loading.value || loadingMore.value || !hasMore.value) return

  const rect = loadMoreTarget.value.getBoundingClientRect()
  if (rect.top > window.innerHeight - 24) return

  if (!loadHintVisible.value) {
    loadHintVisible.value = true
    return
  }

  if (rect.top <= window.innerHeight - 112) {
    loadMore()
  }
}

function onWindowScroll() {
  maybeLoadAfterHint()
}

function goToTeams() {
  router.push('/teams')
}

function goToLogin() {
  router.push('/login')
}

onMounted(async () => {
  observer = new IntersectionObserver(
    ([entry]) => {
      if (entry?.isIntersecting) maybeLoadAfterHint()
    },
    { rootMargin: '0px 0px 80px' },
  )
  window.addEventListener('scroll', onWindowScroll, { passive: true })
  await resetAndLoad()
})

onBeforeUnmount(() => {
  observer?.disconnect()
  window.removeEventListener('scroll', onWindowScroll)
  if (searchTimer) window.clearTimeout(searchTimer)
})

watch(activeFilter, resetAndLoad)
watch(search, () => {
  if (searchTimer) window.clearTimeout(searchTimer)
  searchTimer = window.setTimeout(resetAndLoad, 320)
})
</script>

<template>
  <div>
    <div class="section-head">
      <div>
        <h2>全部赛程</h2>
      </div>
      <button class="pill-btn bracket-entry" @click="router.push('/bracket')">
        <span class="material-symbols-outlined">account_tree</span>
        对阵图
      </button>
      <div class="schedule-actions">
        <SyncStatusBadge />
      </div>
    </div>

    <SearchInput v-model="search" placeholder="搜索球队 / 城市 / 球场" />

    <div class="section">
      <ChipFilter v-model="activeFilter" :options="filterOptions" />
    </div>

    <div v-if="loading" class="state-text">赛程加载中...</div>
    <div v-else-if="errorMessage" class="state-text error">{{ errorMessage }}</div>
    <div v-else-if="!groupedByDate.length" class="state-text">
      <span>{{ emptyText }}</span>
      <button v-if="followedOnly && !auth.isLoggedIn" class="inline-action" @click="goToLogin">去登录</button>
      <button v-else-if="followedOnly" class="inline-action" @click="goToTeams">去关注球队</button>
    </div>

    <template v-for="group in groupedByDate" :key="group.key">
      <div class="date-group">{{ group.title }}</div>
      <div class="stack">
        <MatchCard v-for="m in group.matches" :key="m.id" :match="m" />
      </div>
    </template>

    <div ref="loadMoreTarget" class="load-more" :class="{ active: hasMore && !loading }">
      {{ loadText }}
    </div>

  </div>
</template>

<style scoped>
.section-head {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 10px;
  margin-bottom: 12px;
}

.section-head > div:first-of-type {
  min-width: 0;
  flex: 0 0 auto;
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

.schedule-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}

.bracket-entry {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  flex: 0 0 auto;
  font-size: 12px;
  min-height: 32px;
  padding: 0 10px;
}

.bracket-entry .material-symbols-outlined {
  font-size: 16px;
}

.section {
  margin-top: 18px;
}

.date-group {
  margin: 18px 0 10px;
  color: var(--muted);
  font-size: 13px;
  font-weight: 750;
}

.stack {
  display: grid;
  gap: 12px;
}

.state-text,
.load-more {
  min-height: 44px;
  display: grid;
  place-items: center;
  color: var(--muted);
  font-size: 13px;
  font-weight: 650;
}

.state-text {
  gap: 10px;
  margin-top: 18px;
}

.state-text.error {
  color: var(--primary);
}

.inline-action {
  min-height: 34px;
  padding: 0 14px;
  border: 0;
  border-radius: 999px;
  color: #fff;
  background: var(--primary);
  font-size: 13px;
  font-weight: 750;
  cursor: pointer;
}

.load-more {
  margin: 16px auto 88px;
  opacity: 0.72;
  transition: opacity 0.18s ease;
}

.load-more.active {
  opacity: 1;
}

</style>
