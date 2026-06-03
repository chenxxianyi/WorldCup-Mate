<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import SearchInput from '@/components/common/SearchInput.vue'
import ChipFilter from '@/components/common/ChipFilter.vue'
import MatchCard from '@/components/common/MatchCard.vue'
import { useMatchStore } from '@/stores/useMatchStore'

const matchStore = useMatchStore()
const search = ref('')
const activeFilter = ref('全部')

const filterOptions = [
  '全部', '今日', '明日', '未开始', '淘汰赛',
  'Group A', 'Group B', 'Group C', 'Group D', 'Group E', 'Group F',
  'Group G', 'Group H', 'Group I', 'Group J', 'Group K', 'Group L'
]

function loadMatches() {
  const params: Record<string, any> = {}
  if (activeFilter.value === '今日') {
    const now = new Date()
    params.date = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
  } else if (activeFilter.value === '明日') {
    const tmr = new Date()
    tmr.setDate(tmr.getDate() + 1)
    params.date = `${tmr.getFullYear()}-${String(tmr.getMonth() + 1).padStart(2, '0')}-${String(tmr.getDate()).padStart(2, '0')}`
  } else if (activeFilter.value.startsWith('Group')) {
    params.groupName = activeFilter.value
  } else if (activeFilter.value === '淘汰赛') {
    params.stage = 'round_of_16'
  } else if (activeFilter.value === '未开始') {
    params.status = 'scheduled'
  }
  if (search.value) params.keyword = search.value
  matchStore.fetchMatches(params)
}

onMounted(loadMatches)
watch(activeFilter, loadMatches)

const filteredMatches = computed(() => {
  let list = matchStore.matches

  if (activeFilter.value === '今日') {
    const now = new Date()
    const key = `${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
    list = list.filter((m) => (m.local_kickoff_time || '').startsWith(key))
  } else if (activeFilter.value === '明日') {
    const tmr = new Date()
    tmr.setDate(tmr.getDate() + 1)
    const key = `${String(tmr.getMonth() + 1).padStart(2, '0')}-${String(tmr.getDate()).padStart(2, '0')}`
    list = list.filter((m) => (m.local_kickoff_time || '').startsWith(key))
  } else if (activeFilter.value.startsWith('Group')) {
    list = list.filter((m) => m.group_name === activeFilter.value)
  } else if (activeFilter.value === '淘汰赛') {
    list = list.filter((m) => m.stage !== 'group' && m.stage !== 'group_stage')
  } else if (activeFilter.value === '未开始') {
    list = list.filter((m) => m.status === 'upcoming')
  }

  if (search.value) {
    const q = search.value.toLowerCase()
    list = list.filter(
      (m) =>
        m.home_team_name.toLowerCase().includes(q) ||
        m.away_team_name.toLowerCase().includes(q) ||
        m.city.toLowerCase().includes(q) ||
        m.stadium.toLowerCase().includes(q)
    )
  }

  return list
})

const groupedByDate = computed(() => {
  const groups: Record<string, typeof filteredMatches.value> = {}
  for (const m of filteredMatches.value) {
    let date = '其他日期'
    const timeStr = m.local_kickoff_time || ''
    if (timeStr.includes(' ')) {
      const md = timeStr.split(' ')[0]
      const parts = md.split('-')
      if (parts.length === 2) {
        date = `${Number(parts[0])}月${Number(parts[1])}日`
      }
    }
    if (!groups[date]) groups[date] = []
    groups[date].push(m)
  }
  return groups
})
</script>

<template>
  <div>
    <div class="section-head">
      <div>
        <h2>全部赛程</h2>
        <span>按日期、小组、阶段快速筛选</span>
      </div>
    </div>
    <SearchInput v-model="search" placeholder="搜索球队 / 城市 / 球场" />
    <div class="section">
      <ChipFilter v-model="activeFilter" :options="filterOptions" />
    </div>
    <template v-for="(matches, date) in groupedByDate" :key="date">
      <div class="date-group">{{ date }}</div>
      <div class="stack">
        <MatchCard v-for="m in matches" :key="m.id" :match="m" />
      </div>
    </template>
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
</style>
