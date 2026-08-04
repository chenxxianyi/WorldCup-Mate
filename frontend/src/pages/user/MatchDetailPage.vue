<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiGetGroupStandings } from '@/api/standings'
import TeamBadge from '@/components/theme/TeamBadge.vue'
import ThemeIcon from '@/components/theme/ThemeIcon.vue'
import { badgeColor, formatMatchday, matchToThemeMatch } from '@/data/themeAdapters'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import { useMatchStore } from '@/stores/useMatchStore'
import { useReminderStore } from '@/stores/useReminderStore'
import { normalizeStanding, type Standing } from '@/types/standing'

const route = useRoute()
const router = useRouter()
const theme = useLeagueThemeStore()
const auth = useAuthStore()
const matchStore = useMatchStore()
const favorites = useFavoriteStore()
const reminders = useReminderStore()
const groupStandings = ref<Standing[]>([])
const error = ref('')
const match = computed(() => matchStore.currentMatch)
const displayMatch = computed(() => match.value ? matchToThemeMatch(match.value) : null)

function standingBadge(row: Standing) {
  return [row.team_name, row.team_code || 'TBD', theme.current.name, badgeColor(row.team_code), row.flag] as const
}

async function loadMatch() {
  error.value = ''
  groupStandings.value = []
  const id = Number(route.params.id)
  if (!id) { error.value = '比赛不存在'; return }
  try {
    const current = await matchStore.fetchMatchDetail(id)
    if (current?.group_id) {
      const rows = await apiGetGroupStandings(current.group_id) as any[]
      groupStandings.value = rows.map(normalizeStanding)
    }
  } catch (reason) {
    error.value = reason instanceof Error ? reason.message : '比赛详情加载失败'
  }
}

function requireLogin() {
  if (auth.isLoggedIn) return true
  router.push({ path: '/login', query: { redirect: route.fullPath } })
  return false
}

async function toggleFavorite() {
  if (!match.value || !requireLogin()) return
  const was = favorites.isMatchFavorite(match.value.id)
  const saved = await favorites.toggleMatchFavorite(match.value.id)
  theme.showToast(saved ? (was ? '已取消收藏' : '比赛已加入收藏') : '收藏操作失败，请稍后重试')
}

async function toggleReminder() {
  if (!match.value || !requireLogin()) return
  const was = reminders.hasReminder(match.value.id)
  const saved = await reminders.toggleReminder(match.value.id, 15)
  theme.showToast(saved ? (was ? '开球提醒已关闭' : '开球提醒已开启') : '提醒操作失败，请稍后重试')
}

onMounted(loadMatch)
watch(() => route.params.id, loadMatch)
watch(() => theme.currentCode, () => router.push('/schedule'))
</script>

<template>
  <div class="page-view">
    <div class="back-row"><button class="back-button" type="button" @click="router.push('/schedule')"><ThemeIcon name="back" /> 返回赛程</button></div>
    <div v-if="matchStore.loading" class="page-state"><span class="state-spinner" />正在加载比赛</div>
    <article v-else-if="error || !displayMatch" class="card empty-compact"><span class="empty-art"><ThemeIcon name="calendar" /></span><span class="empty-copy"><h3>无法显示比赛</h3><p>{{ error || '比赛数据不存在。' }}</p></span></article>
    <template v-else>
      <section class="score-hero">
        <div class="match-card-top" style="position: relative"><span class="status-line"><i v-if="displayMatch.status === 'live'" class="live-dot" />{{ displayMatch.status === 'live' ? 'LIVE MATCH' : formatMatchday(displayMatch.source, theme.current.stage) }}</span><span class="next-kickoff">{{ displayMatch.date }} · {{ displayMatch.time }}</span></div>
        <div class="detail-score-line"><button class="detail-team" type="button" @click="router.push(`/teams/${displayMatch.homeTeamId}`)"><TeamBadge :team="displayMatch.home" size="large" /><h2>{{ displayMatch.home[0] }}</h2><p>{{ displayMatch.home[1] }}</p></button><div class="detail-score">{{ displayMatch.score }}<small>{{ displayMatch.status === 'live' ? `${match?.minute || ''} 比赛进行中` : displayMatch.venue }}</small></div><button class="detail-team" type="button" @click="router.push(`/teams/${displayMatch.awayTeamId}`)"><TeamBadge :team="displayMatch.away" size="large" /><h2>{{ displayMatch.away[0] }}</h2><p>{{ displayMatch.away[1] }}</p></button></div>
        <div class="detail-actions"><button class="secondary-button" :class="{ active: favorites.isMatchFavorite(displayMatch.id) }" type="button" @click="toggleFavorite"><ThemeIcon name="star" /> {{ favorites.isMatchFavorite(displayMatch.id) ? '已收藏' : '收藏比赛' }}</button><button class="primary-button" :class="{ active: reminders.hasReminder(displayMatch.id) }" type="button" @click="toggleReminder"><ThemeIcon name="bell" /> {{ reminders.hasReminder(displayMatch.id) ? '提醒已开启' : '开球提醒' }}</button></div>
      </section>

      <div class="detail-grid section">
        <article class="card detail-panel"><h3>比赛信息</h3><div class="info-list"><span><small>赛事阶段</small><strong>{{ formatMatchday(displayMatch.source, theme.current.stage) }}</strong></span><span><small>开球时间</small><strong>{{ displayMatch.date }} {{ displayMatch.time }}</strong></span><span><small>比赛城市</small><strong>{{ match?.city || '待定' }}</strong></span><span><small>比赛球场</small><strong>{{ match?.stadium || '待定' }}</strong></span><span><small>比赛编号</small><strong>#{{ match?.match_number || displayMatch.id }}</strong></span><span><small>当前状态</small><strong>{{ displayMatch.status === 'live' ? '直播中' : displayMatch.status === 'finished' ? '已结束' : '未开始' }}</strong></span></div></article>
        <article class="card timeline-card"><h3>数据说明</h3><div class="timeline-item"><span class="timeline-minute">API</span><i class="timeline-node" /><span class="timeline-copy"><strong>比赛基础数据已实时连接</strong><span>比分、状态、球队、开球时间和场地来自后端。</span></span></div><div class="timeline-item"><span class="timeline-minute">LIVE</span><i class="timeline-node" /><span class="timeline-copy"><strong>技术统计与事件流</strong><span>当前数据源未提供射门、控球和逐分钟事件，接口扩展后可直接接入。</span></span></div></article>
      </div>

      <section v-if="groupStandings.length" class="section"><div class="section-heading"><div><p class="eyebrow">GROUP TABLE</p><h2>{{ match?.group_name }} 积分</h2></div></div><div class="card standings-preview"><div v-for="(row, index) in groupStandings" :key="row.team_id" class="standing-row"><span class="standing-position" :class="{ 'zone-top': (row.rank || index + 1) <= 2 }">{{ row.rank || index + 1 }}</span><span class="standing-team"><TeamBadge :team="standingBadge(row)" size="small" /><span>{{ row.team_name }}</span></span><strong>{{ row.points }}</strong></div></div></section>
    </template>
  </div>
</template>
