<script setup lang="ts">
import TeamFlag from '@/components/common/TeamFlag.vue'

defineProps<{
  homeTeam?: { name: string; code: string; flag: string }
  awayTeam?: { name: string; code: string; flag: string }
  homeScore?: number | null
  awayScore?: number | null
  homeWinner?: boolean
  awayWinner?: boolean
  empty?: boolean
  placeholder?: string
  matchId?: number
}>()

const emit = defineEmits<{
  (e: 'click', matchId?: number): void
}>()
</script>

<template>
  <div class="bracket-slot" :class="{ empty }" @click="emit('click', matchId)">
    <!-- Home team (top) -->
    <div class="team-row" :class="{ winner: homeWinner }">
      <template v-if="!empty && homeTeam">
        <TeamFlag :value="homeTeam.flag" :alt="homeTeam.name" :fallback="homeTeam.code" size="sm" />
        <span class="team-name">{{ homeTeam.name }}</span>
        <span class="slot-score">{{ homeScore ?? '-' }}</span>
      </template>
      <template v-else-if="placeholder">
        <span class="team-name placeholder">{{ placeholder }}</span>
      </template>
      <template v-else>
        <span class="team-name placeholder">TBD</span>
      </template>
    </div>

    <!-- Connector middle -->
    <div class="connector-line"></div>

    <!-- Away team (bottom) -->
    <div class="team-row" :class="{ winner: awayWinner }">
      <template v-if="!empty && awayTeam">
        <TeamFlag :value="awayTeam.flag" :alt="awayTeam.name" :fallback="awayTeam.code" size="sm" />
        <span class="team-name">{{ awayTeam.name }}</span>
        <span class="slot-score">{{ awayScore ?? '-' }}</span>
      </template>
      <template v-else-if="placeholder">
        <span class="team-name placeholder">{{ placeholder }}</span>
      </template>
      <template v-else>
        <span class="team-name placeholder">TBD</span>
      </template>
    </div>
  </div>
</template>

<style scoped>
.bracket-slot {
  position: relative;
  display: grid;
  gap: 0;
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  background: var(--card);
  overflow: hidden;
  cursor: pointer;
  transition: border-color 180ms ease;
}

.bracket-slot:hover {
  border-color: var(--primary);
}

.bracket-slot.empty {
  opacity: 0.5;
}

.team-row {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 7px 10px;
  min-height: 34px;
  font-size: 12px;
}

.team-row:first-child {
  border-bottom: 1px solid var(--line);
}

.team-row.winner {
  background: color-mix(in srgb, var(--success) 8%, transparent);
}

.team-name {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  font-weight: 650;
}

.team-name.placeholder {
  color: var(--weak);
}

.slot-score {
  font-weight: 850;
  font-variant-numeric: tabular-nums;
}

.connector-line {
  display: none;
}
</style>
