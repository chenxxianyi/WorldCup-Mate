<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { Player } from '@/types/player'

const props = defineProps<{
  players: Player[]
  loading?: boolean
  error?: string
}>()

const tabs = [
  { label: '全部', value: 'all' },
  { label: '门将', value: 'GK' },
  { label: '后卫', value: 'DF' },
  { label: '中场', value: 'MF' },
  { label: '前锋', value: 'FW' },
]

const activePosition = ref('all')
const brokenPhotos = ref<Set<number>>(new Set())

const sortedPlayers = computed(() =>
  [...props.players].sort((a, b) => {
    const an = a.shirtNumber ?? 999
    const bn = b.shirtNumber ?? 999
    if (an !== bn) return an - bn
    return a.name.localeCompare(b.name)
  }),
)

const filteredPlayers = computed(() => {
  if (activePosition.value === 'all') return sortedPlayers.value
  return sortedPlayers.value.filter((player) => player.position === activePosition.value)
})

function hasPhoto(player: Player) {
  return Boolean(player.photoUrl) && !brokenPhotos.value.has(player.id)
}

function markPhotoBroken(player: Player) {
  brokenPhotos.value = new Set([...brokenPhotos.value, player.id])
}

watch(
  () => props.players,
  () => {
    brokenPhotos.value = new Set()
    if (activePosition.value !== 'all' && !props.players.some((player) => player.position === activePosition.value)) {
      activePosition.value = 'all'
    }
  },
)
</script>

<template>
  <section class="roster-section">
    <div class="roster-head">
      <div>
        <h2>球员大名单</h2>
        <span>世界杯阵容名单</span>
      </div>
      <strong>{{ players.length }} 人</strong>
    </div>

    <div class="position-tabs" role="tablist" aria-label="按位置筛选球员">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        :class="{ active: activePosition === tab.value }"
        type="button"
        role="tab"
        :aria-selected="activePosition === tab.value"
        @click="activePosition = tab.value"
      >
        {{ tab.label }}
      </button>
    </div>

    <div v-if="loading" class="roster-state">正在加载球员名单...</div>
    <div v-else-if="error" class="roster-state error">{{ error }}</div>
    <div v-else-if="!players.length" class="roster-state">球员大名单暂未同步</div>
    <div v-else-if="!filteredPlayers.length" class="roster-state">这个位置暂无球员</div>

    <div v-else class="player-list">
      <article v-for="player in filteredPlayers" :key="player.id" class="player-row">
        <img
          v-if="hasPhoto(player)"
          class="player-photo"
          :src="player.photoUrl"
          :alt="player.name"
          loading="lazy"
          @error="markPhotoBroken(player)"
        />
        <span v-else class="player-photo-placeholder" aria-hidden="true"></span>

        <span class="player-number">#{{ player.shirtNumber || '-' }}</span>
        <div class="player-copy">
          <strong>{{ player.name }}</strong>
          <span v-if="player.positionLabel || player.position || player.club">
            {{ player.positionLabel || player.position }}
            <template v-if="player.club"> · {{ player.club }}</template>
          </span>
        </div>
      </article>
    </div>
  </section>
</template>

<style scoped>
.roster-section {
  display: grid;
  gap: 12px;
}

.roster-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 12px;
}

.roster-head h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 750;
}

.roster-head span {
  display: block;
  margin-top: 3px;
  color: var(--muted);
  font-size: 13px;
}

.roster-head strong {
  color: var(--muted);
  font-size: 13px;
  white-space: nowrap;
}

.position-tabs {
  display: flex;
  gap: 6px;
  overflow-x: auto;
  padding-bottom: 1px;
}

.position-tabs button {
  min-height: 32px;
  padding: 0 12px;
  border: 1px solid var(--line);
  border-radius: 999px;
  color: var(--muted);
  background: var(--card);
  font-size: 13px;
  font-weight: 750;
  white-space: nowrap;
  transition: color 160ms ease-out, background 160ms ease-out, border-color 160ms ease-out;
}

.position-tabs button.active {
  color: var(--blue);
  border-color: color-mix(in srgb, var(--blue) 32%, var(--line));
  background: color-mix(in srgb, var(--blue) 7%, var(--card));
}

.position-tabs button:focus-visible,
.player-row:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--blue) 55%, transparent);
  outline-offset: 2px;
}

.player-list {
  display: grid;
  border-top: 1px solid var(--line);
}

.player-row {
  min-width: 0;
  min-height: 62px;
  display: grid;
  grid-template-columns: 42px auto minmax(0, 1fr);
  align-items: center;
  gap: 10px;
  padding: 9px 0;
  border-bottom: 1px solid var(--line);
}

.player-photo,
.player-photo-placeholder {
  width: 42px;
  height: 42px;
  border-radius: 999px;
}

.player-photo {
  display: block;
  object-fit: cover;
  background: var(--card-soft);
}

.player-photo-placeholder {
  display: block;
  background: var(--card-soft);
}

.player-number {
  min-width: 36px;
  height: 24px;
  display: inline-grid;
  place-items: center;
  padding: 0 7px;
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
  gap: 3px;
}

.player-copy strong,
.player-copy span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.player-copy strong {
  font-size: 14px;
  line-height: 1.25;
}

.player-copy span {
  color: var(--muted);
  font-size: 12px;
  line-height: 1.3;
}

.roster-state {
  padding: 16px 0;
  color: var(--muted);
  font-size: 13px;
  text-align: center;
  border-top: 1px solid var(--line);
  border-bottom: 1px solid var(--line);
}

.roster-state.error {
  color: var(--blue);
}

@media (min-width: 768px) {
  .player-list {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    column-gap: 14px;
  }
}
</style>
