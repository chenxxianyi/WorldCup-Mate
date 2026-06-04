<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import type { Team } from '@/types/team'
import TeamFlag from '@/components/common/TeamFlag.vue'

const props = defineProps<{ team: Team }>()

const router = useRouter()
const fav = useFavoriteStore()

function goDetail() {
  router.push(`/teams/${props.team.id}`)
}
</script>

<template>
  <article class="card team-card" tabindex="0" @click="goDetail" @keydown.enter="goDetail">
    <button
      class="fav-mini"
      :class="{ active: fav.isTeamFollowed(team.id) }"
      title="关注球队"
      @click.stop="fav.toggleTeamFollow(team.id)"
    >
      <span
        class="material-symbols-outlined"
        :style="fav.isTeamFollowed(team.id) ? 'font-variation-settings: \'FILL\' 1' : ''"
      >star</span>
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
  cursor: pointer;
}

.team-card:focus-visible {
  outline: 2px solid var(--primary);
  outline-offset: 2px;
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
  cursor: pointer;
}

.fav-mini .material-symbols-outlined {
  font-size: 18px;
}

.fav-mini.active {
  color: #fff;
  border-color: transparent;
  background: var(--primary);
}
</style>
