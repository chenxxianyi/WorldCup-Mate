<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import SearchInput from '@/components/common/SearchInput.vue'
import ChipFilter from '@/components/common/ChipFilter.vue'
import TeamCard from '@/components/common/TeamCard.vue'
import { useTeamStore } from '@/stores/useTeamStore'

const teamStore = useTeamStore()
const search = ref('')
const activeFilter = ref('全部')

const filterOptions = [
  '全部',
  '亚洲',
  '欧洲',
  '南美洲',
  '北美洲',
  '非洲',
  '大洋洲',
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

function isPlaceholderTeam(team: (typeof teamStore.teams)[number]) {
  const code = team.code.toUpperCase()
  return code === 'TBD' || code.startsWith('TBD') || team.name === 'TBD' || team.name_en === 'TBD'
}

const filteredTeams = computed(() => {
  let list = teamStore.teams.filter((team) => !isPlaceholderTeam(team))

  if (activeFilter.value !== '全部') {
    if (activeFilter.value.startsWith('Group')) {
      list = list.filter((t) => t.group_name === activeFilter.value)
    } else {
      list = list.filter((t) => t.continent === activeFilter.value)
    }
  }

  if (search.value) {
    const q = search.value.toLowerCase()
    list = list.filter(
      (t) =>
        t.name.toLowerCase().includes(q) ||
        t.name_en.toLowerCase().includes(q) ||
        t.code.toLowerCase().includes(q)
    )
  }

  return list
})

onMounted(() => {
  teamStore.fetchTeams({ page_size: 100 })
})
</script>

<template>
  <div>
    <div class="section-head">
      <div>
        <h2>参赛球队</h2>
        <span>48 支球队 · 按大洲与小组筛选</span>
      </div>
    </div>
    <SearchInput v-model="search" placeholder="搜索球队" />
    <div class="section">
      <ChipFilter v-model="activeFilter" :options="filterOptions" />
    </div>
    <div class="team-grid">
      <TeamCard v-for="t in filteredTeams" :key="t.id" :team="t" />
    </div>
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

.team-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 12px;
}

@media (min-width: 768px) {
  .team-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
}

@media (min-width: 1120px) {
  .team-grid {
    grid-template-columns: repeat(6, minmax(0, 1fr));
  }
}
</style>
