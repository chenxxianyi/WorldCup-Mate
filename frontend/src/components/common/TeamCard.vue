<script setup lang="ts">
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import type { Team } from '@/types/team'
import TeamFlag from '@/components/common/TeamFlag.vue'

defineProps<{ team: Team }>()

const fav = useFavoriteStore()
</script>

<template>
  <article class="card team-card">
    <button
      class="fav-mini"
      :class="{ active: fav.isTeamFollowed(team.id) }"
      @click="fav.toggleTeamFollow(team.id)"
    >
      {{ fav.isTeamFollowed(team.id) ? '★' : '☆' }}
    </button>
    <TeamFlag :value="team.flag" :alt="team.name" :fallback="team.code" size="lg" />
    <h3>{{ team.name }}</h3>
    <p>{{ team.name_en }} · {{ team.group_name }}</p>
    <span v-if="fav.isTeamFollowed(team.id)" class="tag green" style="margin-top: 12px">已关注</span>
    <span v-else class="tag" style="margin-top: 12px">{{ team.continent }}</span>
  </article>
</template>

<style scoped>
.team-card {
  position: relative;
  overflow: hidden;
  padding: 15px;
  min-height: 152px;
}

.team-card h3 {
  margin: 11px 0 2px;
  font-size: 16px;
  font-weight: 750;
}

.team-card p {
  margin: 0;
  color: var(--muted);
  font-size: 12px;
}

.fav-mini {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  border: 1px solid var(--line);
  border-radius: 999px;
  color: var(--weak);
  background: var(--card);
  font-size: 16px;
}

.fav-mini.active {
  color: #fff;
  border-color: transparent;
  background: var(--primary);
}
</style>
