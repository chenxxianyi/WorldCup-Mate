<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiGetMatchesByStage } from '@/api/matches'
import type { Match, Stage } from '@/types/match'
import { normalizeMatch } from '@/types/match'
import TeamFlag from '@/components/common/TeamFlag.vue'

const router = useRouter()

interface BracketRound {
  key: Stage
  name: string
  matches: Match[]
  loading: boolean
}

interface BracketPath {
  key: string
  index: number
  matches: Match[]
  nodeTitle: string
  nodeName: string
  slotLabel: string
}

const rounds = ref<BracketRound[]>([
  { key: 'round_of_32', name: '32强', matches: [], loading: false },
  { key: 'round_of_16', name: '16强', matches: [], loading: false },
  { key: 'quarter_final', name: '8强', matches: [], loading: false },
  { key: 'semi_final', name: '半决赛', matches: [], loading: false },
  { key: 'third_place', name: '季军赛', matches: [], loading: false },
  { key: 'final', name: '决赛', matches: [], loading: false },
])

const nextRoundName: Partial<Record<Stage, string>> = {
  round_of_32: '16强',
  round_of_16: '8强',
  quarter_final: '4强',
  semi_final: '决赛',
  third_place: '季军',
  final: '冠军',
}

const activeRoundIndex = ref(0)
const activeRound = computed(() => rounds.value[activeRoundIndex.value] ?? rounds.value[0])
const totalRounds = computed(() => rounds.value.filter((r) => r.matches.length > 0).length)
const totalMatches = computed(() => rounds.value.reduce((sum, round) => sum + round.matches.length, 0))

const activeMatches = computed(() => (
  [...activeRound.value.matches].sort((a, b) => a.match_number - b.match_number)
))

const activeFinishedCount = computed(() =>
  activeMatches.value.filter((match) => match.status === 'finished' || winner(match)).length,
)

const activeProgress = computed(() => {
  if (!activeMatches.value.length) return 0
  return Math.round((activeFinishedCount.value / activeMatches.value.length) * 100)
})

const headerMeta = computed(() => {
  if (totalMatches.value <= 0) return '等待官方淘汰赛赛程'
  return `${totalRounds.value} 轮已同步 · ${totalMatches.value} 场比赛`
})

const activePaths = computed<BracketPath[]>(() => {
  const round = activeRound.value
  const groupSize = round.key === 'final' || round.key === 'third_place' ? 1 : 2
  const nodeName = nextRoundName[round.key] || '下一轮'
  const nodeTitle = round.key === 'final' ? '冠军' : round.key === 'third_place' ? '季军' : '晋级'
  const paths: BracketPath[] = []

  for (let i = 0; i < activeMatches.value.length; i += groupSize) {
    const matches = activeMatches.value.slice(i, i + groupSize)
    const index = paths.length + 1
    const slotLabel = round.key === 'final' || round.key === 'third_place'
      ? `${nodeName}归属`
      : `${nodeName}席位 ${index}`

    paths.push({
      key: `${round.key}-${index}`,
      index,
      matches,
      nodeTitle,
      nodeName,
      slotLabel,
    })
  }

  return paths
})

async function loadRound(round: BracketRound) {
  if (round.matches.length > 0 || round.loading) return
  round.loading = true
  try {
    const res = await apiGetMatchesByStage(round.key) as any[]
    round.matches = (res || []).map(normalizeMatch)
  } catch {
    round.matches = []
  } finally {
    round.loading = false
  }
}

async function loadAll() {
  for (const round of rounds.value) {
    await loadRound(round)
  }

  const idx = rounds.value.findIndex((r) => r.matches.length > 0)
  if (idx >= 0) activeRoundIndex.value = idx
}

function goMatch(matchId?: number) {
  if (matchId) router.push(`/matches/${matchId}`)
}

function goTeam(teamId?: number) {
  if (teamId) router.push(`/teams/${teamId}`)
}

function goBack() {
  const canGoBack = Boolean(window.history.state?.back)
  if (canGoBack) {
    router.back()
    return
  }
  router.push('/schedule').catch(() => {})
}

function winner(m: Match): 'home' | 'away' | null {
  if (m.home_score == null || m.away_score == null) return null
  if (m.home_score > m.away_score) return 'home'
  if (m.away_score > m.home_score) return 'away'
  return null
}

function teamName(name?: string) {
  const value = (name || '').trim()
  return value && value !== 'TBD' ? value : '待定'
}

function kickoffText(match: Match) {
  return match.local_kickoff_time || '时间待定'
}

function statusText(match: Match) {
  if (match.status === 'live') return `进行中 ${match.minute || 0}'`
  if (match.status === 'finished') return '已结束'
  if (match.status === 'postponed') return '延期'
  if (match.status === 'cancelled') return '取消'
  return '未开始'
}

function venueText(match: Match) {
  const city = match.city && match.city !== 'TBD' ? match.city : ''
  const stadium = match.stadium && match.stadium !== 'TBD' ? match.stadium : ''
  if (city && stadium) return `${city} · ${stadium}`
  return city || stadium || '场地待定'
}

onMounted(loadAll)
</script>

<template>
  <div class="bracket-page">
    <section class="bracket-hero">
      <div class="hero-toolbar">
        <button class="back-action" type="button" title="返回" aria-label="返回上一页" @click="goBack">
          <span class="material-symbols-outlined" aria-hidden="true">arrow_back</span>
        </button>
        <span class="hero-chip">
          <span class="material-symbols-outlined" aria-hidden="true">account_tree</span>
          淘汰赛
        </span>
        <span v-if="totalRounds > 0" class="hero-rounds">{{ totalRounds }} 轮</span>
      </div>

      <div class="page-head">
        <div>
          <h2>淘汰赛对阵图</h2>
          <p>{{ headerMeta }}</p>
        </div>
      </div>

      <div class="bracket-overview">
        <div class="overview-title">
          <span>当前轮次</span>
          <strong>{{ activeRound.name }}</strong>
        </div>
        <div class="overview-stat">
          <strong>{{ activeRound.matches.length || '-' }}</strong>
          <span>场比赛</span>
        </div>
        <div class="overview-stat accent">
          <strong>{{ activePaths.length || '-' }}</strong>
          <span>晋级位</span>
        </div>
      </div>

      <div class="round-progress" aria-hidden="true">
        <span :style="{ width: `${activeProgress}%` }"></span>
      </div>
    </section>

    <div class="round-tabs" role="tablist" aria-label="淘汰赛轮次">
      <button
        v-for="(round, i) in rounds"
        :key="round.key"
        class="round-tab"
        :class="{ active: i === activeRoundIndex, has: round.matches.length > 0 }"
        type="button"
        role="tab"
        :aria-selected="i === activeRoundIndex"
        @click="activeRoundIndex = i; loadRound(round)"
      >
        <span>{{ round.name }}</span>
        <small v-if="round.loading" class="tab-loading">...</small>
        <small v-else-if="round.matches.length" class="tab-count">{{ round.matches.length }}</small>
      </button>
    </div>

    <div class="round-summary">
      <span>{{ activeFinishedCount }} / {{ activeRound.matches.length || 0 }} 场已决出结果</span>
      <strong>{{ activeProgress }}%</strong>
    </div>

    <div class="bracket-grid">
      <div v-if="activeRound.loading" class="state-text">加载中...</div>

      <template v-else-if="activePaths.length">
        <section
          v-for="path in activePaths"
          :key="path.key"
          class="bracket-path"
          :class="{ single: path.matches.length === 1 }"
        >
          <div class="path-title">
            <span>
              <span class="path-number">{{ path.index }}</span>
              路径 {{ path.index }}
            </span>
            <small>{{ path.slotLabel }}</small>
          </div>

          <div class="path-body">
            <div class="match-stack">
              <article
                v-for="m in path.matches"
                :key="m.id"
                class="bracket-game"
              >
                <button class="game-slot" type="button" @click="goMatch(m.id)">
                  <span class="match-meta">
                    <span>
                      <span class="material-symbols-outlined" aria-hidden="true">schedule</span>
                      {{ kickoffText(m) }}
                    </span>
                    <strong>{{ statusText(m) }}</strong>
                  </span>

                  <span
                    class="team-line"
                    :class="{ win: winner(m) === 'home' }"
                    @click.stop="goTeam(m.home_team_id)"
                  >
                    <TeamFlag :value="m.home_flag" :alt="m.home_team_name" :fallback="m.home_team_code" size="sm" />
                    <span class="team-name">{{ teamName(m.home_team_name) }}</span>
                    <span class="game-score">{{ m.home_score ?? '-' }}</span>
                  </span>

                  <span
                    class="team-line"
                    :class="{ win: winner(m) === 'away' }"
                    @click.stop="goTeam(m.away_team_id)"
                  >
                    <TeamFlag :value="m.away_flag" :alt="m.away_team_name" :fallback="m.away_team_code" size="sm" />
                    <span class="team-name">{{ teamName(m.away_team_name) }}</span>
                    <span class="game-score">{{ m.away_score ?? '-' }}</span>
                  </span>

                  <span class="match-venue">
                    <span class="material-symbols-outlined" aria-hidden="true">stadium</span>
                    {{ venueText(m) }}
                  </span>
                </button>
              </article>
            </div>

            <div class="connector" aria-hidden="true">
              <span></span>
            </div>

            <div class="advance-node">
              <span>{{ path.nodeTitle }}</span>
              <strong>{{ path.nodeName }}</strong>
            </div>
          </div>
        </section>
      </template>

      <div v-else class="state-text">该轮比赛尚未开始</div>
    </div>
  </div>
</template>

<style scoped>
.bracket-page {
  --bracket-line: color-mix(in srgb, var(--primary) 30%, var(--line));
  display: grid;
  gap: 14px;
}

.bracket-hero {
  position: relative;
  display: grid;
  gap: 14px;
  overflow: hidden;
  padding: 14px;
  border: 1px solid color-mix(in srgb, var(--line) 76%, transparent);
  border-radius: var(--radius-xl);
  background:
    radial-gradient(100% 120% at 0% 0%, color-mix(in srgb, var(--primary) 10%, transparent), transparent 48%),
    radial-gradient(84% 110% at 88% 4%, color-mix(in srgb, var(--accent) 18%, transparent), transparent 46%),
    linear-gradient(180deg, color-mix(in srgb, var(--card) 94%, transparent), color-mix(in srgb, var(--card-soft) 48%, transparent)),
    var(--card);
  box-shadow:
    0 16px 38px rgba(15, 23, 42, 0.06),
    inset 0 1px 0 color-mix(in srgb, #fff 72%, transparent);
}

.bracket-hero::after {
  content: "";
  position: absolute;
  inset: auto 14px 10px;
  height: 1px;
  background: linear-gradient(90deg, transparent, color-mix(in srgb, var(--primary) 22%, transparent), transparent);
  pointer-events: none;
}

.hero-toolbar {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 8px;
}

.back-action {
  width: 44px;
  height: 44px;
  display: inline-grid;
  flex: 0 0 auto;
  place-items: center;
  margin-left: -6px;
  border: 1px solid color-mix(in srgb, var(--line) 82%, transparent);
  border-radius: 999px;
  color: var(--text);
  background: color-mix(in srgb, var(--card) 80%, transparent);
  cursor: pointer;
  transition:
    color 160ms ease-out,
    background 160ms ease-out,
    border-color 160ms ease-out,
    transform 160ms ease-out;
}

.back-action:hover {
  color: var(--primary);
  border-color: color-mix(in srgb, var(--primary) 24%, var(--line));
  background: color-mix(in srgb, var(--primary) 7%, var(--card));
}

.back-action:active {
  transform: scale(0.97);
}

.back-action:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--blue) 55%, transparent);
  outline-offset: 2px;
}

.back-action .material-symbols-outlined {
  font-size: 24px;
}

.hero-chip,
.hero-rounds {
  min-height: 34px;
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 800;
  white-space: nowrap;
}

.hero-chip {
  gap: 5px;
  padding: 0 11px;
  color: var(--primary);
  border: 1px solid color-mix(in srgb, var(--primary) 16%, var(--line));
  background: color-mix(in srgb, var(--primary) 8%, var(--card));
}

.hero-chip .material-symbols-outlined {
  font-size: 16px;
}

.hero-rounds {
  margin-left: auto;
  padding: 0 10px;
  color: var(--muted);
  background: color-mix(in srgb, var(--card) 70%, transparent);
}

.page-head {
  position: relative;
  z-index: 1;
  min-width: 0;
}

.page-head > div {
  min-width: 0;
}

.page-head h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 900;
  line-height: 1.15;
  letter-spacing: 0;
}

.page-head p {
  margin: 8px 0 0;
  color: var(--muted);
  font-size: 13px;
  font-weight: 650;
}

.round-tabs {
  display: flex;
  gap: 9px;
  overflow-x: auto;
  padding: 0 1px 5px;
  scrollbar-width: none;
  scroll-snap-type: x proximity;
}

.round-tabs::-webkit-scrollbar {
  display: none;
}

.round-tab {
  flex-shrink: 0;
  min-height: 40px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 0 15px;
  border: 1px solid color-mix(in srgb, var(--line) 88%, transparent);
  border-radius: 999px;
  color: var(--text);
  background: color-mix(in srgb, var(--card) 92%, transparent);
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  scroll-snap-align: start;
  transition:
    transform 150ms ease-out,
    color 150ms ease-out,
    background 150ms ease-out,
    border-color 150ms ease-out,
    box-shadow 150ms ease-out;
}

.round-tab:active {
  transform: scale(0.96);
}

.round-tab.active {
  color: #fff;
  border-color: var(--primary);
  background: linear-gradient(135deg, var(--primary), color-mix(in srgb, var(--primary-dark) 82%, var(--secondary)));
  box-shadow: 0 10px 22px color-mix(in srgb, var(--primary) 24%, transparent);
}

.round-tab.has:not(.active) {
  border-color: color-mix(in srgb, var(--primary) 24%, var(--line));
  background: color-mix(in srgb, var(--primary) 6%, var(--card));
}

.tab-count,
.tab-loading {
  min-width: 18px;
  height: 18px;
  display: inline-grid;
  place-items: center;
  border-radius: 999px;
  color: inherit;
  background: color-mix(in srgb, currentColor 10%, transparent);
  font-size: 10px;
  font-weight: 900;
}

.bracket-overview {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto auto;
  align-items: center;
  gap: 9px;
  padding: 10px;
  border: 1px solid color-mix(in srgb, var(--line) 78%, transparent);
  border-radius: var(--radius-md);
  background: color-mix(in srgb, var(--card) 68%, transparent);
  backdrop-filter: blur(10px);
}

.overview-title,
.overview-stat {
  min-width: 0;
  display: grid;
  gap: 2px;
}

.overview-title span,
.overview-stat span {
  color: var(--muted);
  font-size: 11px;
  font-weight: 700;
  white-space: nowrap;
}

.overview-title strong {
  overflow: hidden;
  color: var(--text);
  font-size: 18px;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.overview-stat {
  min-width: 62px;
  justify-items: center;
  padding: 8px 9px;
  border-radius: var(--radius-sm);
  background: color-mix(in srgb, var(--card-soft) 72%, var(--card));
}

.overview-stat strong {
  color: var(--text);
  font-size: 18px;
  font-weight: 900;
  line-height: 1;
}

.overview-stat.accent {
  background: color-mix(in srgb, var(--accent) 32%, var(--card));
}

.round-progress {
  position: relative;
  z-index: 1;
  height: 5px;
  overflow: hidden;
  border-radius: 999px;
  background: color-mix(in srgb, var(--line) 70%, transparent);
}

.round-progress span {
  display: block;
  height: 100%;
  min-width: 8px;
  border-radius: inherit;
  background: linear-gradient(90deg, var(--primary), var(--accent));
  transition: width 220ms ease-out;
}

.round-summary {
  min-height: 38px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 0 12px;
  border: 1px solid color-mix(in srgb, var(--line) 82%, transparent);
  border-radius: var(--radius-md);
  color: var(--muted);
  background: color-mix(in srgb, var(--card) 76%, transparent);
  font-size: 12px;
  font-weight: 750;
}

.round-summary strong {
  flex: 0 0 auto;
  color: var(--primary);
  font-size: 13px;
  font-weight: 900;
}

.bracket-grid {
  display: grid;
  gap: 14px;
}

.bracket-path {
  position: relative;
  display: grid;
  gap: 10px;
  padding: 12px;
  border: 1px solid color-mix(in srgb, var(--line) 82%, transparent);
  border-radius: var(--radius-lg);
  background:
    linear-gradient(180deg, color-mix(in srgb, var(--card) 82%, transparent), transparent),
    color-mix(in srgb, var(--card-soft) 42%, transparent);
}

.path-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  min-width: 0;
}

.path-title > span {
  min-width: 0;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  color: var(--text);
  font-size: 13px;
  font-weight: 900;
}

.path-number {
  width: 22px;
  height: 22px;
  display: inline-grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: 999px;
  color: #fff;
  background: var(--primary);
  font-size: 11px;
  line-height: 1;
}

.path-title small {
  min-width: 0;
  overflow: hidden;
  padding: 5px 9px;
  border-radius: 999px;
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 9%, transparent);
  font-size: 11px;
  font-weight: 850;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.path-body {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 28px 66px;
  align-items: center;
  min-width: 0;
}

.match-stack {
  display: grid;
  gap: 8px;
  min-width: 0;
}

.bracket-game {
  min-width: 0;
}

.game-slot {
  width: 100%;
  min-width: 0;
  display: grid;
  padding: 0;
  overflow: hidden;
  border: 1px solid color-mix(in srgb, var(--line) 88%, transparent);
  border-radius: var(--radius-md);
  color: var(--text);
  background: color-mix(in srgb, var(--card) 94%, transparent);
  text-align: left;
  transition:
    border-color 160ms ease-out,
    background 160ms ease-out,
    transform 160ms ease-out;
}

.game-slot:hover {
  border-color: color-mix(in srgb, var(--primary) 48%, var(--line));
  background: color-mix(in srgb, var(--primary) 3%, var(--card));
}

.game-slot:active {
  transform: scale(0.995);
}

.game-slot:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--primary) 55%, transparent);
  outline-offset: 2px;
}

.match-meta {
  min-width: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  overflow: hidden;
  padding: 8px 10px 7px;
  border-bottom: 1px solid color-mix(in srgb, var(--line) 80%, transparent);
  color: var(--muted);
  background: color-mix(in srgb, var(--card-soft) 72%, var(--card));
  font-size: 11px;
  font-weight: 750;
}

.match-meta > span,
.match-venue {
  min-width: 0;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.match-meta .material-symbols-outlined,
.match-venue .material-symbols-outlined {
  flex: 0 0 auto;
  color: color-mix(in srgb, var(--primary) 58%, var(--weak));
  font-size: 15px;
  line-height: 1;
}

.match-meta strong {
  flex: 0 0 auto;
  color: var(--primary);
  font-size: 11px;
  font-weight: 900;
}

.team-line {
  min-width: 0;
  min-height: 38px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  cursor: pointer;
  transition: background 120ms ease;
}

.team-line + .team-line {
  border-top: 1px solid color-mix(in srgb, var(--line) 76%, transparent);
}

.team-line:hover {
  background: var(--card-soft);
}

.team-line.win {
  background: color-mix(in srgb, var(--success) 11%, var(--card));
}

.team-name {
  min-width: 0;
  flex: 1;
  overflow: hidden;
  font-size: 14px;
  font-weight: 800;
  line-height: 1.2;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.game-score {
  min-width: 16px;
  color: var(--text);
  font-size: 16px;
  font-weight: 900;
  font-variant-numeric: tabular-nums;
  text-align: right;
}

.match-venue {
  padding: 7px 10px 8px;
  border-top: 1px solid color-mix(in srgb, var(--line) 72%, transparent);
  color: var(--muted);
  font-size: 11px;
  font-weight: 700;
}

.connector {
  position: relative;
  align-self: stretch;
  min-height: 92px;
}

.connector::before,
.connector::after,
.connector span::before,
.connector span::after {
  content: "";
  position: absolute;
  background: var(--bracket-line);
}

.connector::before {
  top: 34px;
  bottom: 34px;
  left: 12px;
  width: 2px;
  border-radius: 999px;
}

.connector::after {
  top: 50%;
  right: 0;
  left: 12px;
  height: 2px;
  border-radius: 999px;
  transform: translateY(-50%);
}

.connector span::before,
.connector span::after {
  left: 0;
  width: 14px;
  height: 2px;
  border-radius: 999px;
}

.connector span::before {
  top: 34px;
}

.connector span::after {
  bottom: 34px;
}

.bracket-path.single .connector::before,
.bracket-path.single .connector span::before,
.bracket-path.single .connector span::after {
  display: none;
}

.advance-node {
  width: 66px;
  min-height: 68px;
  display: grid;
  place-items: center;
  align-content: center;
  gap: 3px;
  border: 1px solid color-mix(in srgb, var(--primary) 28%, var(--line));
  border-radius: var(--radius-md);
  background:
    radial-gradient(90% 110% at 50% 0%, color-mix(in srgb, var(--accent) 34%, transparent), transparent 62%),
    var(--card);
  box-shadow: inset 0 0 0 1px color-mix(in srgb, #fff 45%, transparent);
}

.advance-node span,
.advance-node strong {
  max-width: 100%;
  overflow: hidden;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.advance-node span {
  color: var(--primary);
  font-size: 11px;
  font-weight: 900;
}

.advance-node strong {
  color: var(--text);
  font-size: 14px;
  font-weight: 900;
}

.state-text {
  padding: 42px 20px;
  border: 1px dashed var(--line);
  border-radius: var(--radius-lg);
  text-align: center;
  color: var(--muted);
  background: color-mix(in srgb, var(--card) 82%, transparent);
  font-size: 14px;
  font-weight: 700;
}

[data-theme='dark'] .bracket-hero,
[data-theme='dark'] .bracket-path {
  box-shadow: none;
}

[data-theme='dark'] .bracket-hero {
  background:
    radial-gradient(100% 120% at 0% 0%, color-mix(in srgb, var(--primary) 18%, transparent), transparent 48%),
    radial-gradient(84% 110% at 88% 4%, color-mix(in srgb, var(--accent) 13%, transparent), transparent 46%),
    color-mix(in srgb, var(--card) 94%, transparent);
}

@media (max-width: 390px) {
  .bracket-hero {
    padding: 12px;
  }

  .page-head h2 {
    font-size: 22px;
  }

  .bracket-overview {
    grid-template-columns: minmax(0, 1fr) auto;
  }

  .overview-stat.accent {
    display: none;
  }

  .path-body {
    grid-template-columns: minmax(0, 1fr) 22px 58px;
  }

  .bracket-path {
    padding: 10px;
  }

  .advance-node {
    width: 58px;
  }

  .advance-node strong {
    font-size: 13px;
  }

  .match-meta {
    align-items: flex-start;
    flex-direction: column;
    gap: 3px;
  }
}
</style>
