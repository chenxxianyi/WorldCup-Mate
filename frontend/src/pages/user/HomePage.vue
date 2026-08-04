<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { apiGetCompetitionStandings } from '@/api/competitions'
import { apiListMatches, apiGetMatchesByTeam } from '@/api/matches'
import { apiGetAllStandings } from '@/api/standings'
import TeamBadge from '@/components/theme/TeamBadge.vue'
import ThemeIcon from '@/components/theme/ThemeIcon.vue'
import ThemeMatchCard from '@/components/theme/ThemeMatchCard.vue'
import { badgeColor, matchToThemeMatch, type ThemeMatch } from '@/data/themeAdapters'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import { normalizeMatch, type Match } from '@/types/match'
import { normalizeLeagueStanding, normalizeStanding } from '@/types/standing'

interface PreviewStanding { teamId: number; name: string; code: string; flag: string; points: number; position: number }

const router = useRouter()
const theme = useLeagueThemeStore()
const auth = useAuthStore()
const favorites = useFavoriteStore()
const loading = ref(true)
const error = ref('')
const matches = ref<Match[]>([])
const standings = ref<PreviewStanding[]>([])
const totalMatches = ref(0)
const finishedMatches = ref(0)
const followedSchedule = ref<Match[]>([])

const themeMatches = computed(() => matches.value.map(matchToThemeMatch))
const nextMatch = computed(() => themeMatches.value.find((item) => item.source.status !== 'finished') || themeMatches.value[0] || null)
const now = new Date()
const todayKey = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
const todayMatches = computed(() => themeMatches.value.filter((item) => item.kickoffKey === todayKey))
const matchdayMatches = computed(() => (todayMatches.value.length ? todayMatches.value : themeMatches.value).slice(0, 4))
const followedThemeMatches = computed(() => followedSchedule.value.map(matchToThemeMatch).slice(0, 4))
const progress = computed(() => totalMatches.value ? Math.round((finishedMatches.value / totalMatches.value) * 100) : 0)

function standingBadge(item: PreviewStanding) {
  return [item.name, item.code || 'TBD', theme.current.name, badgeColor(item.code), item.flag] as const
}

async function loadStandings() {
  try {
    if (theme.current.slug === 'wc') {
      const rows = await apiGetAllStandings()
      standings.value = rows.slice(0, 5).map((row, index) => {
        const item = normalizeStanding(row)
        return { teamId: item.team_id, name: item.team_name, code: item.team_code, flag: item.flag, points: item.points, position: item.rank || index + 1 }
      })
    } else {
      const rows = await apiGetCompetitionStandings(theme.currentCode, { type: 'total' })
      standings.value = rows.slice(0, 5).map((row) => {
        const item = normalizeLeagueStanding(row)
        return { teamId: item.team_id, name: item.team_name, code: item.team_code, flag: item.flag, points: item.points, position: item.position }
      })
    }
  } catch {
    standings.value = []
  }
}

async function loadFollowedSchedule() {
  followedSchedule.value = []
  if (!auth.isLoggedIn || !favorites.followedTeamIds.length) return
  const lists = await Promise.all(favorites.followedTeamIds.slice(0, 5).map(async (teamId) => {
    try {
      const rows = await apiGetMatchesByTeam(teamId)
      return rows.map(normalizeMatch)
    } catch { return [] }
  }))
  const competitionId = theme.competition.current?.id
  const unique = new Map<number, Match>()
  lists.flat()
    .filter((match) => !competitionId || !match.competition_id || match.competition_id === competitionId)
    .sort((a, b) => new Date(a.kickoff_time_utc).getTime() - new Date(b.kickoff_time_utc).getTime())
    .forEach((match) => unique.set(match.id, match))
  followedSchedule.value = [...unique.values()]
}

async function loadHome() {
  loading.value = true
  error.value = ''
  try {
    await theme.initialize()
    const [response, finishedResponse] = await Promise.all([
      apiListMatches({ page: 1, page_size: 50 }),
      apiListMatches({ page: 1, page_size: 1, status: 'finished' }),
    ])
    const list = response.list.map(normalizeMatch)
    matches.value = list
    totalMatches.value = Number(response.total ?? list.length)
    const finishedList = finishedResponse.list
    finishedMatches.value = Number(finishedResponse.total ?? finishedList.length)
    if (auth.isLoggedIn) await favorites.fetchFavoriteTeams()
    await Promise.all([loadStandings(), loadFollowedSchedule()])
  } catch (reason) {
    matches.value = []
    error.value = reason instanceof Error ? reason.message : '赛事数据加载失败'
  } finally {
    loading.value = false
  }
}

function setNextReminder() {
  if (!nextMatch.value) return
  if (!auth.isLoggedIn) {
    router.push({ path: '/login', query: { redirect: '/' } })
    return
  }
  theme.showToast('可在比赛卡片上开启开球提醒')
}

onMounted(loadHome)
watch(() => theme.currentCode, loadHome)
</script>

<template>
  <div class="page-view">
    <section class="competition-masthead" aria-labelledby="hero-title">
      <div class="hero-copy">
        <div>
          <div class="hero-league-line"><span class="hero-mark">{{ theme.current.mark }}</span><span><b>{{ theme.current.name }}</b><br><small>{{ theme.current.en }} · {{ theme.competition.current ? `${theme.competition.current.season} 赛季` : theme.current.season }}</small></span></div>
          <h1 id="hero-title">{{ theme.current.tagline }}</h1><p class="hero-description">{{ theme.current.description }}</p>
        </div>
        <div class="hero-progress"><div class="hero-progress-meta"><span>{{ theme.current.stage }}</span><span>{{ finishedMatches }} / {{ totalMatches }} 场 · {{ progress }}%</span></div><div class="progress-track"><div class="progress-fill" :style="{ width: `${progress}%` }" /></div></div>
      </div>
      <div v-if="nextMatch" class="next-match-panel">
        <div class="next-match-head"><span class="status-line"><i class="live-dot" />下一场焦点</span><span class="next-kickoff">{{ nextMatch.date }} · {{ nextMatch.time }}</span></div>
        <div class="team-versus"><span class="versus-team"><TeamBadge :team="nextMatch.home" /><strong>{{ nextMatch.home[0] }}</strong></span><span class="versus-center"><strong>{{ nextMatch.score }}</strong><small>{{ theme.current.stage }}</small></span><span class="versus-team"><TeamBadge :team="nextMatch.away" /><strong>{{ nextMatch.away[0] }}</strong></span></div>
        <button class="primary-button full-button" type="button" @click="router.push(`/matches/${nextMatch.id}`)">查看比赛中心 <ThemeIcon name="arrow" /></button>
      </div>
      <div v-else class="next-match-panel empty-panel"><ThemeIcon name="calendar" /><strong>{{ loading ? '正在加载比赛…' : '暂无即将开始的比赛' }}</strong><span v-if="error">{{ error }}</span></div>
    </section>

    <section class="quick-grid" aria-label="赛事速览">
      <article class="card quick-card"><span class="quick-card-top"><span>赛季进度</span><ThemeIcon name="trophy" /></span><strong>{{ progress }}%</strong></article>
      <article class="card quick-card"><span class="quick-card-top"><span>关注球队</span><ThemeIcon name="star" /></span><strong>{{ favorites.followedTeamIds.length }} 支</strong></article>
      <article class="card quick-card"><span class="quick-card-top"><span>今日比赛</span><ThemeIcon name="calendar" /></span><strong>{{ todayMatches.length }} 场</strong></article>
    </section>

    <div class="content-grid">
      <div>
        <section class="section"><div class="section-heading"><div><p class="eyebrow">MATCHDAY</p><h2>{{ todayMatches.length ? '今日比赛' : '近期比赛' }}</h2></div><button class="text-link" type="button" @click="router.push('/schedule')">完整赛程 <ThemeIcon name="arrow" /></button></div>
          <div v-if="matchdayMatches.length" class="match-strip"><ThemeMatchCard v-for="match in matchdayMatches" :key="match.id" :match="match" compact /></div>
          <article v-else class="card empty-compact"><span class="empty-art"><ThemeIcon name="calendar" /></span><span class="empty-copy"><h3>暂无比赛数据</h3><p>{{ error || '后台同步比赛后将在这里显示。' }}</p></span></article>
        </section>

        <section class="section"><div class="section-heading"><div><p class="eyebrow">FOLLOWING</p><h2>我的关注赛程</h2></div><span>{{ favorites.followedTeamIds.length }} 支球队</span></div>
          <div v-if="followedThemeMatches.length" class="match-list"><ThemeMatchCard v-for="match in followedThemeMatches" :key="match.id" :match="match" /></div>
          <article v-else class="card empty-compact"><span class="empty-art"><ThemeIcon name="star" /></span><span class="empty-copy"><h3>{{ auth.isLoggedIn ? '还没有关注球队的赛程' : '登录后查看关注赛程' }}</h3><p>关注球队后，相关比赛会自动聚合到这里。</p></span><button class="primary-button" type="button" @click="router.push(auth.isLoggedIn ? '/teams' : '/login')">{{ auth.isLoggedIn ? '去关注' : '登录' }}</button></article>
        </section>

        <section class="section"><div class="section-heading"><div><p class="eyebrow">PERSONAL</p><h2>比赛提醒</h2></div></div><article class="card empty-compact"><span class="empty-art"><ThemeIcon name="bell" /></span><span class="empty-copy"><h3>不错过真正重要的比赛</h3><p>开球前 15 分钟，通过站内消息和邮件提醒你。</p></span><button class="primary-button" type="button" @click="setNextReminder">设置提醒</button></article></section>
      </div>

      <aside>
        <section class="section"><div class="section-heading"><div><p class="eyebrow">TABLE</p><h2>积分速览</h2></div><button class="text-link" type="button" @click="router.push('/standings')">查看全部</button></div>
          <div v-if="standings.length" class="card standings-preview"><div v-for="item in standings" :key="item.teamId" class="standing-row"><span class="standing-position" :class="{ 'zone-top': item.position <= 4 }">{{ item.position }}</span><span class="standing-team"><TeamBadge :team="standingBadge(item)" size="small" /><span>{{ item.name }}</span></span><strong>{{ item.points }}</strong></div></div>
          <article v-else class="card empty-mini">暂无积分数据</article>
        </section>
        <section v-if="standings[0]" class="section"><div class="section-heading"><div><p class="eyebrow">FORM</p><h2>领头羊近况</h2></div></div><article class="card quick-card"><span class="standing-team"><TeamBadge :team="standingBadge(standings[0])" size="small" /><span>{{ standings[0].name }}</span></span><span class="form-dots"><i class="form-dot">胜</i><i class="form-dot">胜</i><i class="form-dot draw">平</i><i class="form-dot">胜</i><i class="form-dot">胜</i></span></article></section>
      </aside>
    </div>
  </div>
</template>
