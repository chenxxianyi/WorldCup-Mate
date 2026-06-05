<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { LineupPlayer, TeamLineup } from '@/types/lineup'

const props = defineProps<{
  home: TeamLineup
  away: TeamLineup
}>()

interface PitchPlayer {
  key: string
  player: LineupPlayer
  side: 'home' | 'away'
  x: number
  y: number
}

const brokenPhotos = ref<Set<string>>(new Set())

function parseGrid(grid: string) {
  const [row, col] = grid.split(':').map((part) => Number(part))
  if (!row || !col) return null
  return { row, col }
}

function parseFormation(formation: string) {
  const layers = formation
    .split('-')
    .map((part) => Number(part.trim()))
    .filter((value) => Number.isFinite(value) && value > 0)

  return layers.length ? [1, ...layers] : []
}

function clampPercent(value: number) {
  return Math.max(7, Math.min(93, value))
}

function playersFromGrid(team: TeamLineup): PitchPlayer[] {
  const parsed = team.startXi
    .map((player, index) => {
      const grid = parseGrid(player.grid)
      return grid ? { player, index, ...grid } : null
    })
    .filter((value): value is { player: LineupPlayer; index: number; row: number; col: number } => Boolean(value))

  if (!parsed.length) return []

  const rows = [...new Set(parsed.map((item) => item.row))].sort((a, b) => a - b)
  const rowIndex = new Map(rows.map((row, index) => [row, index]))
  const maxRowIndex = Math.max(rows.length - 1, 1)

  return parsed.map((item) => {
    const rowPlayers = parsed.filter((candidate) => candidate.row === item.row)
    const maxCol = Math.max(...rowPlayers.map((candidate) => candidate.col), rowPlayers.length, 1)
    const rowPosition = rowIndex.get(item.row) || 0
    const x = clampPercent((item.col / (maxCol + 1)) * 100)
    const baseY = 10 + (rowPosition / maxRowIndex) * 80

    return {
      key: `${team.side}-${item.player.playerId || item.player.name || item.index}`,
      player: item.player,
      side: team.side,
      x: team.side === 'away' ? 100 - x : x,
      y: team.side === 'home' ? 100 - baseY : baseY,
    }
  })
}

function playersFromFormation(team: TeamLineup): PitchPlayer[] {
  const layers = parseFormation(team.formation)
  if (!layers.length || !team.startXi.length) return []

  const positioned: PitchPlayer[] = []
  let cursor = 0
  const maxLayerIndex = Math.max(layers.length - 1, 1)

  layers.forEach((count, layerIndex) => {
    const rowPlayers = team.startXi.slice(cursor, cursor + count)
    cursor += count

    rowPlayers.forEach((player, playerIndex) => {
      const x = clampPercent(((playerIndex + 1) / (rowPlayers.length + 1)) * 100)
      const baseY = 10 + (layerIndex / maxLayerIndex) * 80

      positioned.push({
        key: `${team.side}-${player.playerId || player.name || layerIndex}-${playerIndex}`,
        player,
        side: team.side,
        x: team.side === 'away' ? 100 - x : x,
        y: team.side === 'home' ? 100 - baseY : baseY,
      })
    })
  })

  return positioned
}

function positionTeam(team: TeamLineup) {
  const gridPlayers = playersFromGrid(team)
  if (gridPlayers.length) return gridPlayers
  return playersFromFormation(team)
}

const positionedPlayers = computed(() => [
  ...positionTeam(props.away),
  ...positionTeam(props.home),
])

function hasPhoto(item: PitchPlayer) {
  return Boolean(item.player.photoUrl) && !brokenPhotos.value.has(item.key)
}

function markPhotoBroken(item: PitchPlayer) {
  brokenPhotos.value = new Set([...brokenPhotos.value, item.key])
}

function playerMeta(player: LineupPlayer) {
  const bits = []
  if (player.shirtNumber) bits.push(`#${player.shirtNumber}`)
  if (player.positionLabel || player.position) bits.push(player.positionLabel || player.position)
  if (player.nameEn) bits.push(player.nameEn)
  return bits.join(' · ')
}

watch(
  () => [props.home.startXi, props.away.startXi],
  () => {
    brokenPhotos.value = new Set()
  },
)
</script>

<template>
  <div class="pitch-wrap" aria-label="首发阵容站位图">
    <div class="pitch-head">
      <span>{{ away.teamName }}</span>
      <strong>{{ away.formation || '阵型待定' }}</strong>
      <span>{{ home.teamName }}</span>
    </div>

    <div class="pitch">
      <div class="pitch-line halfway"></div>
      <div class="pitch-circle"></div>
      <div class="box box-top"></div>
      <div class="box box-bottom"></div>

      <button
        v-for="item in positionedPlayers"
        :key="item.key"
        class="pitch-player"
        :class="item.side"
        type="button"
        :style="{ left: `${item.x}%`, top: `${item.y}%` }"
        :title="`${item.player.name}${playerMeta(item.player) ? ' · ' + playerMeta(item.player) : ''}`"
      >
        <span class="pitch-avatar">
          <img
            v-if="hasPhoto(item)"
            :src="item.player.photoUrl"
            :alt="item.player.name"
            loading="lazy"
            @error="markPhotoBroken(item)"
          />
          <span v-else>{{ item.player.shirtNumber || '' }}</span>
        </span>
        <span class="pitch-name">{{ item.player.name }}</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
.pitch-wrap {
  display: grid;
  gap: 8px;
}

.pitch-head {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  align-items: center;
  gap: 8px;
  color: var(--muted);
  font-size: 12px;
}

.pitch-head span {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pitch-head span:last-child {
  text-align: right;
}

.pitch-head strong {
  color: var(--primary);
  font-size: 12px;
  white-space: nowrap;
}

.pitch {
  position: relative;
  min-height: 420px;
  aspect-ratio: 7 / 9;
  overflow: hidden;
  border: 1px solid color-mix(in srgb, var(--green) 30%, var(--line));
  border-radius: var(--radius-lg);
  background:
    linear-gradient(to bottom, rgba(255, 255, 255, 0.1), transparent 18%, transparent 82%, rgba(255, 255, 255, 0.1)),
    repeating-linear-gradient(
      to bottom,
      color-mix(in srgb, var(--green) 88%, #071f16) 0 46px,
      color-mix(in srgb, var(--green-2) 80%, #082016) 46px 92px
    );
}

.pitch::before {
  content: '';
  position: absolute;
  inset: 12px;
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 12px;
  pointer-events: none;
}

.pitch-line,
.pitch-circle,
.box {
  position: absolute;
  pointer-events: none;
}

.halfway {
  left: 12px;
  right: 12px;
  top: 50%;
  height: 1px;
  background: rgba(255, 255, 255, 0.5);
}

.pitch-circle {
  left: 50%;
  top: 50%;
  width: 86px;
  height: 86px;
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 999px;
  transform: translate(-50%, -50%);
}

.box {
  left: 22%;
  width: 56%;
  height: 84px;
  border: 1px solid rgba(255, 255, 255, 0.5);
}

.box-top {
  top: 12px;
  border-top: 0;
}

.box-bottom {
  bottom: 12px;
  border-bottom: 0;
}

.pitch-player {
  position: absolute;
  z-index: 2;
  width: 76px;
  min-height: 62px;
  display: grid;
  justify-items: center;
  gap: 4px;
  padding: 0;
  border: 0;
  color: #fff;
  background: transparent;
  transform: translate(-50%, -50%);
}

.pitch-player:focus-visible {
  outline: 2px solid rgba(255, 255, 255, 0.8);
  outline-offset: 4px;
  border-radius: 12px;
}

.pitch-avatar {
  width: 36px;
  height: 36px;
  display: grid;
  place-items: center;
  overflow: hidden;
  border: 2px solid rgba(255, 255, 255, 0.78);
  border-radius: 999px;
  background: var(--primary);
  box-shadow: 0 8px 18px rgba(0, 0, 0, 0.22);
  font-size: 12px;
  font-weight: 850;
}

.pitch-player.away .pitch-avatar {
  background: var(--blue);
}

.pitch-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.pitch-name {
  max-width: 76px;
  min-height: 18px;
  padding: 2px 5px;
  overflow: hidden;
  border-radius: 999px;
  color: #fff;
  background: rgba(0, 0, 0, 0.36);
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 10px;
  font-weight: 750;
  line-height: 1.2;
}

@media (max-width: 720px) {
  .pitch {
    min-height: 360px;
  }

  .pitch-player {
    width: 64px;
  }

  .pitch-name {
    max-width: 64px;
  }
}
</style>
