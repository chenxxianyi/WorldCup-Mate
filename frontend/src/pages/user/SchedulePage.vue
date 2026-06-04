<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import SearchInput from '@/components/common/SearchInput.vue'
import ChipFilter from '@/components/common/ChipFilter.vue'
import MatchCard from '@/components/common/MatchCard.vue'
import { apiListMatches } from '@/api/matches'
import { normalizeMatch, type Match } from '@/types/match'

const INITIAL_DAYS = 3
const DAYS_PER_LOAD = 2
const PAGE_SIZE = 8
const LOAD_MORE_MIN_MS = 760

const search = ref('')
const activeFilter = ref('全部')
const matches = ref<Match[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const errorMessage = ref('')
const page = ref(1)
const total = ref(0)
const reachedEnd = ref(false)
const visibleDateLimit = ref(INITIAL_DAYS)
const loadMoreTarget = ref<HTMLElement | null>(null)
const showBackTop = ref(false)
const loadHintVisible = ref(false)

let observer: IntersectionObserver | null = null
let requestToken = 0
let searchTimer: number | undefined

const filterOptions = [
  '全部',
  '今日',
  '明日',
  '未开始',
  '淘汰赛',
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
    pageSize: PAGE_SIZE,
  }

  if (activeFilter.value === '今日') {
    params.date = dayParam()
  } else if (activeFilter.value === '明日') {
    params.date = dayParam(1)
  } else if (activeFilter.value.startsWith('Group')) {
    params.groupName = activeFilter.value
  } else if (activeFilter.value === '淘汰赛') {
    params.stage = 'knockout'
  } else if (activeFilter.value === '未开始') {
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

const loadedDateKeys = computed(() => {
  const keys: string[] = []
  for (const match of matches.value) {
    const key = matchDateKey(match)
    if (!keys.includes(key)) keys.push(key)
  }
  return keys
})

const visibleMatches = computed(() => {
  const allowedKeys = new Set(loadedDateKeys.value.slice(0, visibleDateLimit.value))
  return matches.value.filter((match) => allowedKeys.has(matchDateKey(match)))
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
const hasMore = computed(() => hasBufferedDays.value || !reachedEnd.value)
const loadText = computed(() => {
  if (loadingMore.value) return '加载中...'
  if (hasMore.value) {
    return loadHintVisible.value ? '继续下滑，加载更多赛程' : '下滑查看更多赛程'
  }
  return matches.value.length ? '已加载全部赛程' : ''
})

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
    await fillInitialWindow(token)
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

function updateBackTopVisibility() {
  showBackTop.value = window.scrollY > 420
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
  updateBackTopVisibility()
  maybeLoadAfterHint()
}

function backToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(async () => {
  observer = new IntersectionObserver(
    ([entry]) => {
      if (entry?.isIntersecting) maybeLoadAfterHint()
    },
    { rootMargin: '0px 0px 80px' }
  )
  window.addEventListener('scroll', onWindowScroll, { passive: true })
  updateBackTopVisibility()
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
        <span>默认展示最近三天，下滑继续加载</span>
      </div>
    </div>

    <SearchInput v-model="search" placeholder="搜索球队 / 城市 / 球场" />

    <div class="section">
      <ChipFilter v-model="activeFilter" :options="filterOptions" />
    </div>

    <div v-if="loading" class="state-text">赛程加载中...</div>
    <div v-else-if="errorMessage" class="state-text error">{{ errorMessage }}</div>
    <div v-else-if="!groupedByDate.length" class="state-text">暂无符合条件的比赛</div>

    <template v-for="group in groupedByDate" :key="group.key">
      <div class="date-group">{{ group.title }}</div>
      <div class="stack">
        <MatchCard v-for="m in group.matches" :key="m.id" :match="m" />
      </div>
    </template>

    <div ref="loadMoreTarget" class="load-more" :class="{ active: hasMore && !loading }">
      {{ loadText }}
    </div>

    <Transition name="back-top">
      <button
        v-if="showBackTop"
        class="back-top-btn"
        type="button"
        aria-label="返回顶部"
        title="返回顶部"
        @click="backToTop"
      >
        <span class="material-symbols-outlined">keyboard_arrow_up</span>
      </button>
    </Transition>
  </div>
</template>

<style scoped>
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
  margin-top: 18px;
}

.state-text.error {
  color: var(--primary);
}

.load-more {
  margin: 16px auto 88px;
  opacity: 0.72;
  transition: opacity 0.18s ease;
}

.load-more.active {
  opacity: 1;
}

.back-top-btn {
  position: fixed;
  right: max(18px, calc((100vw - 1280px) / 2 + 18px));
  bottom: calc(var(--nav-h) + 22px);
  z-index: 40;
  width: 46px;
  height: 46px;
  display: grid;
  place-items: center;
  border: 1px solid color-mix(in srgb, var(--primary) 24%, transparent);
  border-radius: 50%;
  color: #fff;
  background: var(--primary);
  box-shadow: 0 14px 30px color-mix(in srgb, var(--primary) 28%, transparent);
}

.back-top-btn .material-symbols-outlined {
  font-size: 28px;
  line-height: 1;
}

.back-top-enter-active,
.back-top-leave-active {
  transition: opacity 0.18s ease, transform 0.18s ease;
}

.back-top-enter-from,
.back-top-leave-to {
  opacity: 0;
  transform: translateY(10px) scale(0.94);
}

@media (min-width: 768px) {
  .back-top-btn {
    bottom: 28px;
  }
}
</style>
