<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { apiGetCompetitionStandings } from '@/api/competitions'
import { apiGetAllStandings, apiGetBestThird } from '@/api/standings'
import TeamBadge from '@/components/theme/TeamBadge.vue'
import ThemeIcon from '@/components/theme/ThemeIcon.vue'
import { badgeColor } from '@/data/themeAdapters'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import { normalizeLeagueStanding, normalizeStanding } from '@/types/standing'

interface TableStanding {
  teamId: number; name: string; code: string; flag: string; position: number; played: number; won: number; drawn: number; lost: number; goalDifference: number; points: number; zone: string
}

const theme = useLeagueThemeStore()
const mode = ref('总榜')
const rows = ref<TableStanding[]>([])
const loading = ref(false)
const error = ref('')
const worldCupGroups = Array.from({ length: 12 }, (_, index) => `小组 ${String.fromCharCode(65 + index)}`)
const modes = computed(() => theme.current.slug === 'wc' ? [...worldCupGroups, '最佳第三'] : ['总榜', '主场', '客场'])

function badge(row: TableStanding) {
  return [row.name, row.code || 'TBD', theme.current.name, badgeColor(row.code), row.flag] as const
}

function zoneClass(row: TableStanding) {
  const value = row.zone.toLowerCase()
  if (value.includes('releg') || value.includes('淘汰')) return 'zone-danger'
  if (value.includes('europa') || value.includes('possible')) return 'zone-europa'
  if (value.includes('champion') || value.includes('qual') || value.includes('晋级')) return 'zone-top'
  return ''
}

async function loadStandings() {
  loading.value = true
  error.value = ''
  try {
    if (theme.current.slug === 'wc') {
      const response = mode.value === '最佳第三' ? await apiGetBestThird() : await apiGetAllStandings()
      const groupName = mode.value.replace('小组 ', 'Group ')
      const filtered = mode.value === '最佳第三' ? response : response.filter((item) => item.group?.name === groupName)
      rows.value = filtered.map((item, index) => {
        const row = normalizeStanding(item)
        return { teamId: row.team_id, name: row.team_name, code: row.team_code, flag: row.flag, position: row.rank || index + 1, played: row.played, won: row.won, drawn: row.drawn, lost: row.lost, goalDifference: row.goal_difference, points: row.points, zone: row.status }
      })
    } else {
      const type = ({ 总榜: 'total', 主场: 'home', 客场: 'away' } as Record<string, string>)[mode.value] || 'total'
      const response = await apiGetCompetitionStandings(theme.currentCode, { type })
      rows.value = response.map((item) => {
        const row = normalizeLeagueStanding(item)
        return { teamId: row.team_id, name: row.team_name, code: row.team_code, flag: row.flag, position: row.position, played: row.played, won: row.won, drawn: row.drawn, lost: row.lost, goalDifference: row.goal_difference, points: row.points, zone: row.zone }
      })
    }
  } catch (reason) {
    rows.value = []
    error.value = reason instanceof Error ? reason.message : '积分榜加载失败'
  } finally {
    loading.value = false
  }
}

function selectMode(value: string) {
  mode.value = value
  loadStandings()
}

onMounted(() => { mode.value = modes.value[0]; loadStandings() })
watch(() => theme.currentCode, () => { mode.value = modes.value[0]; loadStandings() })
</script>

<template>
  <div class="page-view">
    <header class="page-heading"><div><p class="eyebrow">{{ theme.current.en }} · {{ theme.competition.current?.season || theme.current.season }}</p><h1>积分榜</h1></div><div class="standings-tabs"><button v-for="item in modes" :key="item" class="chip" :class="{ active: item === mode }" type="button" @click="selectMode(item)">{{ item }}</button></div></header>
    <div v-if="loading" class="page-state"><span class="state-spinner" />正在加载积分榜</div>
    <section v-else-if="rows.length" class="card table-card">
      <table class="standings-table"><thead><tr><th>排名</th><th>球队</th><th>赛</th><th>胜</th><th>平</th><th>负</th><th>净胜</th><th>积分</th></tr></thead><tbody>
        <tr v-for="row in rows" :key="row.teamId"><td><span class="standing-position" :class="zoneClass(row)">{{ row.position }}</span></td><td><span class="table-team"><TeamBadge :team="badge(row)" size="small" /><span>{{ row.name }}</span></span></td><td>{{ row.played }}</td><td>{{ row.won }}</td><td>{{ row.drawn }}</td><td>{{ row.lost }}</td><td>{{ row.goalDifference > 0 ? '+' : '' }}{{ row.goalDifference }}</td><td><strong>{{ row.points }}</strong></td></tr>
      </tbody></table>
      <div class="zone-legend"><span class="zone-item"><i class="zone-swatch" />晋级 / 欧冠区</span><span class="zone-item"><i class="zone-swatch europa" />次级欧战区</span><span class="zone-item"><i class="zone-swatch danger" />淘汰 / 降级区</span></div>
    </section>
    <article v-else class="card empty-compact"><span class="empty-art"><ThemeIcon name="trophy" /></span><span class="empty-copy"><h3>暂无积分数据</h3><p>{{ error || '完成比赛同步或积分重算后将在这里展示。' }}</p></span></article>
  </div>
</template>
