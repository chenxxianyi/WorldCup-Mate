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

const activeMatches = computed(() => (
  [...activeRound.value.matches].sort((a, b) => a.match_number - b.match_number)
))

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

onMounted(loadAll)
</script>

<template>
  <div class="bracket-page">
    <div class="bracket-nav-row">
      <button class="back-action" type="button" title="返回" aria-label="返回上一页" @click="goBack">
        <span class="material-symbols-outlined" aria-hidden="true">arrow_back</span>
      </button>
    </div>

    <div class="page-head">
      <div>
        <h2>淘汰赛对阵图</h2>
        <p>{{ activeRound.name }} · {{ activePaths.length || 0 }} 个晋级节点</p>
      </div>
      <span v-if="totalRounds > 0">{{ totalRounds }} 轮</span>
    </div>

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
        {{ round.name }}
        <span v-if="round.loading" class="tab-loading">...</span>
        <span v-else-if="round.matches.length" class="tab-count">{{ round.matches.length }}</span>
      </button>
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
            <span>路径 {{ path.index }}</span>
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
                  <span class="match-meta">{{ kickoffText(m) }}</span>

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
  display: grid;
  gap: 12px;
}

.bracket-nav-row {
  display: flex;
  align-items: center;
  min-height: 44px;
  margin-bottom: -2px;
}

.page-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.page-head > div {
  min-width: 0;
}

.back-action {
  width: 44px;
  height: 44px;
  display: inline-grid;
  flex: 0 0 auto;
  place-items: center;
  margin-left: -10px;
  border: 0;
  border-radius: 999px;
  color: var(--text);
  background: transparent;
  cursor: pointer;
  transition: color 160ms ease-out, background 160ms ease-out, transform 160ms ease-out;
}

.back-action:hover {
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 7%, transparent);
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

.page-head h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 800;
  letter-spacing: 0;
}

.page-head p {
  margin: 5px 0 0;
  color: var(--muted);
  font-size: 12px;
}

.page-head > span {
  margin-top: 4px;
  color: var(--muted);
  font-size: 13px;
  white-space: nowrap;
}

.round-tabs {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 4px;
  scrollbar-width: none;
}

.round-tabs::-webkit-scrollbar {
  display: none;
}

.round-tab {
  flex-shrink: 0;
  min-height: 36px;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 0 14px;
  border: 1px solid color-mix(in srgb, var(--primary) 22%, var(--line));
  border-radius: 999px;
  color: var(--text);
  background: var(--card);
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  transition: transform 150ms ease-out, color 150ms ease-out, background 150ms ease-out, border-color 150ms ease-out;
}

.round-tab:active {
  transform: scale(0.96);
}

.round-tab.active {
  color: #fff;
  border-color: var(--primary);
  background: var(--primary);
}

.round-tab.has:not(.active) {
  border-color: color-mix(in srgb, var(--primary) 36%, var(--line));
  background: color-mix(in srgb, var(--primary) 5%, var(--card));
}

.tab-count,
.tab-loading {
  font-size: 10px;
  opacity: 0.72;
}

.bracket-overview {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto auto;
  align-items: center;
  gap: 8px;
  padding: 10px;
  border: 1px solid color-mix(in srgb, var(--primary) 18%, var(--line));
  border-radius: var(--radius-md);
  background:
    linear-gradient(135deg, color-mix(in srgb, var(--primary) 7%, transparent), transparent 42%),
    var(--card);
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
  white-space: nowrap;
}

.overview-title strong {
  overflow: hidden;
  color: var(--text);
  font-size: 17px;
  font-weight: 850;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.overview-stat {
  min-width: 58px;
  justify-items: center;
  padding: 7px 9px;
  border-radius: var(--radius-sm);
  background: var(--card-soft);
}

.overview-stat strong {
  color: var(--text);
  font-size: 17px;
  font-weight: 850;
  line-height: 1;
}

.overview-stat.accent {
  background: color-mix(in srgb, var(--accent) 28%, var(--card));
}

.bracket-grid {
  display: grid;
  gap: 16px;
}

.bracket-path {
  display: grid;
  gap: 8px;
}

.path-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  min-width: 0;
}

.path-title span {
  color: var(--text);
  font-size: 13px;
  font-weight: 800;
}

.path-title small {
  min-width: 0;
  overflow: hidden;
  padding: 4px 8px;
  border-radius: 999px;
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 9%, var(--card));
  font-size: 11px;
  font-weight: 800;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.path-body {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 30px 66px;
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
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  color: var(--text);
  background: var(--card);
  text-align: left;
  transition: border-color 160ms ease-out, background 160ms ease-out;
}

.game-slot:hover {
  border-color: color-mix(in srgb, var(--primary) 48%, var(--line));
  background: color-mix(in srgb, var(--primary) 3%, var(--card));
}

.match-meta {
  min-width: 0;
  overflow: hidden;
  padding: 7px 10px 6px;
  border-bottom: 1px solid var(--line);
  color: var(--muted);
  background: color-mix(in srgb, var(--card-soft) 70%, var(--card));
  font-size: 11px;
  font-weight: 700;
  text-overflow: ellipsis;
  white-space: nowrap;
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
  border-top: 1px solid var(--line);
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
  font-weight: 750;
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
  background: color-mix(in srgb, var(--primary) 38%, var(--line));
}

.connector::before {
  top: 28px;
  bottom: 28px;
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
  top: 28px;
}

.connector span::after {
  bottom: 28px;
}

.bracket-path.single .connector::before,
.bracket-path.single .connector span::before,
.bracket-path.single .connector span::after {
  display: none;
}

.advance-node {
  width: 66px;
  min-height: 62px;
  display: grid;
  place-items: center;
  align-content: center;
  gap: 3px;
  border: 1px solid color-mix(in srgb, var(--primary) 34%, var(--line));
  border-radius: var(--radius-md);
  background:
    linear-gradient(180deg, color-mix(in srgb, var(--accent) 24%, transparent), transparent),
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
  font-weight: 850;
}

.advance-node strong {
  color: var(--text);
  font-size: 14px;
  font-weight: 900;
}

.state-text {
  padding: 42px 20px;
  border: 1px dashed var(--line);
  border-radius: var(--radius-md);
  text-align: center;
  color: var(--muted);
  background: var(--card);
  font-size: 14px;
}

@media (max-width: 390px) {
  .bracket-overview {
    grid-template-columns: minmax(0, 1fr) auto;
  }

  .overview-stat.accent {
    display: none;
  }

  .path-body {
    grid-template-columns: minmax(0, 1fr) 24px 58px;
  }

  .advance-node {
    width: 58px;
  }

  .advance-node strong {
    font-size: 13px;
  }
}
</style>
