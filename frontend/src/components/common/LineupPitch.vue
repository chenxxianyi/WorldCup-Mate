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
  row: number
  rowSize: number
  labelPlacement: 'above' | 'below'
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
  return Math.max(12, Math.min(88, value))
}

function verticalPosition(side: 'home' | 'away', rowPosition: number, maxRowIndex: number) {
  const halfY = 7 + (rowPosition / maxRowIndex) * 38
  return side === 'home' ? 100 - halfY : halfY
}

function avoidCenterLabelZone(side: 'home' | 'away', y: number, rowPosition: number, maxRowIndex: number) {
  if (rowPosition !== maxRowIndex) return y
  return side === 'away' ? Math.min(y, 45) : Math.max(y, 55)
}

function labelPlacement(side: 'home' | 'away', rowPosition: number, maxRowIndex: number) {
  const isClosestToCenter = rowPosition === maxRowIndex
  if (side === 'home') return 'above'
  if (side === 'away') return isClosestToCenter ? 'above' : 'below'
  return 'below'
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

    return {
      key: `${team.side}-${item.player.playerId || item.player.name || item.index}`,
      player: item.player,
      side: team.side,
      x: team.side === 'away' ? 100 - x : x,
      y: avoidCenterLabelZone(team.side, verticalPosition(team.side, rowPosition, maxRowIndex), rowPosition, maxRowIndex),
      row: item.row,
      rowSize: rowPlayers.length,
      labelPlacement: labelPlacement(team.side, rowPosition, maxRowIndex),
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

      positioned.push({
        key: `${team.side}-${player.playerId || player.name || layerIndex}-${playerIndex}`,
        player,
        side: team.side,
        x: team.side === 'away' ? 100 - x : x,
        y: avoidCenterLabelZone(team.side, verticalPosition(team.side, layerIndex, maxLayerIndex), layerIndex, maxLayerIndex),
        row: layerIndex + 1,
        rowSize: rowPlayers.length,
        labelPlacement: labelPlacement(team.side, layerIndex, maxLayerIndex),
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
        :class="[
          item.side,
          `row-${item.row}`,
          `label-${item.labelPlacement}`,
          { dense: item.rowSize >= 4, crowded: item.rowSize >= 5 },
        ]"
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
  width: 100%;
  min-width: 0;
  display: grid;
  gap: 8px;
  overflow: hidden;
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
  justify-self: center;
  width: min(100%, 640px);
  min-height: clamp(560px, 116vw, 760px);
  aspect-ratio: 7 / 10;
  container-type: inline-size;
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
  --avatar-size: 34px;
  --label-gap: 5px;
  --name-height: 32px;
  --anchor-offset: calc(var(--avatar-size) / 2);

  position: absolute;
  z-index: 2;
  width: clamp(78px, 19cqw, 112px);
  height: calc(var(--avatar-size) + var(--label-gap) + var(--name-height));
  display: grid;
  grid-template-rows: var(--avatar-size) var(--label-gap) var(--name-height);
  justify-items: center;
  padding: 0;
  border: 0;
  color: #fff;
  background: transparent;
  overflow: visible;
}

.pitch-player.dense {
  width: clamp(64px, 16cqw, 92px);
}

.pitch-player.crowded {
  width: clamp(56px, 13cqw, 76px);
}

.pitch-player:focus-visible {
  outline: 2px solid rgba(255, 255, 255, 0.8);
  outline-offset: 4px;
  border-radius: 12px;
}

.pitch-avatar {
  position: relative;
  z-index: 2;
  width: var(--avatar-size);
  height: var(--avatar-size);
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
  width: 100%;
  height: var(--name-height);
  display: -webkit-box;
  padding: 4px 7px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 8px;
  color: #fff;
  background: rgba(0, 22, 16, 0.76);
  box-shadow: 0 6px 14px rgba(0, 0, 0, 0.18);
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow-wrap: anywhere;
  word-break: normal;
  text-align: center;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.35);
  white-space: normal;
  font-size: clamp(9px, 2cqw, 10.5px);
  font-weight: 800;
  line-height: 1.15;
}

.pitch-player.label-below .pitch-name {
  grid-row: 3;
}

.pitch-player.label-above .pitch-name {
  grid-row: 1;
}

.pitch-player.label-below .pitch-avatar {
  grid-row: 1;
}

.pitch-player.label-above .pitch-avatar {
  grid-row: 3;
}

.pitch-player.label-below {
  transform: translate(-50%, calc(-1 * var(--anchor-offset)));
}

.pitch-player.label-above {
  transform: translate(-50%, calc(-100% + var(--anchor-offset)));
}

@media (max-width: 720px) {
  .pitch {
    width: 100%;
    min-height: clamp(560px, 138vw, 720px);
  }

  .pitch-player {
    --avatar-size: 30px;
    --label-gap: 4px;
    --name-height: 30px;
    width: clamp(72px, 20cqw, 104px);
  }

  .pitch-name {
    padding: 3px 6px;
  }

  .pitch-avatar {
    font-size: 11px;
  }
}
</style>
