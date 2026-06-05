<script setup lang="ts">
import { computed } from 'vue'
import LineupPitch from '@/components/common/LineupPitch.vue'
import type { LineupPlayer, MatchLineups, TeamLineup } from '@/types/lineup'

const props = defineProps<{
  lineups: MatchLineups | null
  loading?: boolean
  error?: string
}>()

const teams = computed(() => {
  const values: TeamLineup[] = []
  if (props.lineups?.home) values.push(props.lineups.home)
  if (props.lineups?.away) values.push(props.lineups.away)
  return values
})

const hasLineupData = computed(() =>
  teams.value.some((team) => team.startXi.length || team.substitutes.length || team.formation || team.coachName),
)

const showPitch = computed(() => {
  const home = props.lineups?.home
  const away = props.lineups?.away
  if (!home || !away) return false

  const hasGrid = [...home.startXi, ...away.startXi].some((player) => player.grid)
  const hasFormation = Boolean(home.formation || away.formation)
  return Boolean((home.startXi.length || away.startXi.length) && (hasGrid || hasFormation))
})

const statusText = computed(() => {
  if (props.lineups?.status === 'available') return '已公布'
  if (props.lineups?.status === 'partial') return '更新中'
  return '未公布'
})

function sourceText(source?: string) {
  if (source === 'api-football') return 'API-Football'
  if (source === 'football-data') return 'football-data'
  return source || ''
}

function sideText(side: string) {
  return side === 'home' ? '主队' : '客队'
}

function playerNumber(player: LineupPlayer) {
  return player.shirtNumber ? `#${player.shirtNumber}` : '-'
}

function playerPosition(player: LineupPlayer) {
  return player.positionLabel || player.position
}
</script>

<template>
  <section class="lineups-section">
    <div class="lineups-head">
      <div>
        <h2>首发阵容</h2>
        <span>比赛双方阵型与替补席</span>
      </div>
      <strong class="lineups-status" :class="lineups?.status || 'unavailable'">{{ statusText }}</strong>
    </div>

    <div v-if="loading" class="lineups-state card">正在加载首发阵容...</div>
    <div v-else-if="error" class="lineups-state card error">{{ error }}</div>
    <div v-else-if="!lineups || !hasLineupData" class="lineups-state card">
      {{ lineups?.message || '首发阵容暂未公布' }}
    </div>

    <template v-else>
      <div v-if="lineups.source" class="lineups-source">
        <span>数据来源</span>
        <strong>{{ sourceText(lineups.source) }}</strong>
      </div>

      <LineupPitch
        v-if="showPitch && lineups.home && lineups.away"
        :home="lineups.home"
        :away="lineups.away"
      />

      <div class="lineup-teams">
        <article v-for="team in teams" :key="`${team.side}-${team.teamId}`" class="team-lineup card">
          <div class="team-lineup-head">
            <div>
              <span>{{ sideText(team.side) }}</span>
              <h3>{{ team.teamName || 'TBD' }}</h3>
            </div>
            <strong>{{ team.formation || '阵型待定' }}</strong>
          </div>

          <div v-if="team.coachName" class="coach-row">
            <span>主教练</span>
            <strong>{{ team.coachName }}</strong>
          </div>

          <div class="lineup-block">
            <div class="lineup-block-head">
              <h4>首发</h4>
              <span>{{ team.startXi.length }} 人</span>
            </div>
            <div v-if="team.startXi.length" class="lineup-list">
              <div
                v-for="(player, index) in team.startXi"
                :key="`${team.side}-starter-${player.playerId || player.name || index}`"
                class="lineup-player-row"
              >
                <span class="player-number">{{ playerNumber(player) }}</span>
                <div class="player-copy">
                  <strong>{{ player.name || '球员待定' }}</strong>
                  <span v-if="player.nameEn">{{ player.nameEn }}</span>
                </div>
                <span v-if="playerPosition(player)" class="position-pill">{{ playerPosition(player) }}</span>
              </div>
            </div>
            <div v-else class="mini-state">首发暂未公布</div>
          </div>

          <div class="lineup-block substitutes">
            <div class="lineup-block-head">
              <h4>替补席</h4>
              <span>{{ team.substitutes.length }} 人</span>
            </div>
            <div v-if="team.substitutes.length" class="lineup-list compact">
              <div
                v-for="(player, index) in team.substitutes"
                :key="`${team.side}-sub-${player.playerId || player.name || index}`"
                class="lineup-player-row"
              >
                <span class="player-number">{{ playerNumber(player) }}</span>
                <div class="player-copy">
                  <strong>{{ player.name || '球员待定' }}</strong>
                </div>
                <span v-if="playerPosition(player)" class="position-pill">{{ playerPosition(player) }}</span>
              </div>
            </div>
            <div v-else class="mini-state">替补席暂未公布</div>
          </div>
        </article>
      </div>
    </template>
  </section>
</template>

<style scoped>
.lineups-section {
  display: grid;
  gap: 12px;
}

.lineups-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 12px;
}

.lineups-head h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 750;
}

.lineups-head span {
  display: block;
  margin-top: 3px;
  color: var(--muted);
  font-size: 13px;
}

.lineups-status {
  min-height: 26px;
  display: inline-flex;
  align-items: center;
  padding: 0 10px;
  border-radius: 999px;
  color: var(--muted);
  background: var(--card-soft);
  font-size: 12px;
  white-space: nowrap;
}

.lineups-status.available {
  color: var(--success);
  background: color-mix(in srgb, var(--success) 12%, transparent);
}

.lineups-status.partial {
  color: var(--gold);
  background: color-mix(in srgb, var(--gold) 14%, transparent);
}

.lineups-state {
  min-height: 86px;
  display: grid;
  place-items: center;
  padding: 18px;
  color: var(--muted);
  font-size: 13px;
  text-align: center;
}

.lineups-state.error {
  color: var(--blue);
}

.lineups-source {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
  color: var(--muted);
  font-size: 12px;
}

.lineups-source strong {
  color: var(--text);
  font-weight: 750;
}

.lineup-teams {
  display: grid;
  grid-auto-flow: column;
  grid-auto-columns: 100%;
  gap: 12px;
  min-width: 0;
  overflow-x: auto;
  overflow-y: visible;
  scroll-padding-inline: 0;
  scroll-snap-type: x mandatory;
  overscroll-behavior-x: contain;
  -webkit-overflow-scrolling: touch;
}

.team-lineup {
  min-width: 0;
  display: grid;
  gap: 14px;
  padding: 14px;
  scroll-snap-align: start;
  scroll-snap-stop: always;
}

.team-lineup-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.team-lineup-head span,
.coach-row span,
.lineup-block-head span {
  color: var(--muted);
  font-size: 12px;
}

.team-lineup-head h3 {
  margin: 3px 0 0;
  overflow-wrap: anywhere;
  font-size: 16px;
  font-weight: 800;
}

.team-lineup-head strong {
  min-height: 28px;
  display: inline-flex;
  align-items: center;
  padding: 0 10px;
  border-radius: 999px;
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 9%, transparent);
  font-size: 13px;
  white-space: nowrap;
}

.coach-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 9px 10px;
  border-radius: var(--radius-sm);
  background: var(--card-soft);
}

.coach-row strong {
  min-width: 0;
  overflow-wrap: anywhere;
  text-align: right;
  font-size: 13px;
}

.lineup-block {
  display: grid;
  gap: 8px;
}

.lineup-block-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.lineup-block-head h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 800;
}

.lineup-list {
  display: grid;
  border-top: 1px solid var(--line);
}

.lineup-player-row {
  min-width: 0;
  min-height: 44px;
  display: grid;
  grid-template-columns: 42px minmax(0, 1fr) auto;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  border-bottom: 1px solid var(--line);
}

.lineup-list.compact .lineup-player-row {
  min-height: 38px;
  padding: 6px 0;
}

.player-number {
  min-width: 34px;
  height: 24px;
  display: inline-grid;
  place-items: center;
  border-radius: 999px;
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 9%, var(--card));
  font-size: 12px;
  font-weight: 850;
  white-space: nowrap;
}

.player-copy {
  min-width: 0;
  display: grid;
  gap: 2px;
}

.player-copy strong,
.player-copy span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.player-copy strong {
  font-size: 13px;
  line-height: 1.25;
}

.player-copy span {
  color: var(--muted);
  font-size: 11px;
  line-height: 1.25;
}

.position-pill {
  max-width: 76px;
  overflow: hidden;
  padding: 3px 7px;
  border-radius: 999px;
  color: var(--muted);
  background: var(--card-soft);
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 11px;
  font-weight: 750;
}

.mini-state {
  padding: 12px 0;
  color: var(--muted);
  font-size: 12px;
  text-align: center;
  border-top: 1px solid var(--line);
  border-bottom: 1px solid var(--line);
}

@media (max-width: 720px) {
  .lineups-source {
    justify-content: flex-start;
  }
}
</style>
