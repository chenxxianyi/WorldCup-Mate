<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { DemoTeam } from '@/data/leagueTheme'

const props = defineProps<{ team: DemoTeam; size?: 'small' | 'large' }>()
const failed = ref(false)
const crest = computed(() => (props.team[4] || '').trim())

watch(crest, () => { failed.value = false })
</script>

<template>
  <span
    class="team-badge"
    :class="[size, { 'has-crest': crest && !failed }]"
    :style="{ '--badge-a': team[3], '--badge-b': team[3] }"
    :aria-label="team[0]"
  >
    <img
      v-if="crest && !failed"
      class="team-badge-crest"
      :src="crest"
      alt=""
      loading="lazy"
      decoding="async"
      referrerpolicy="no-referrer"
      @error="failed = true"
    >
    <span
      v-else
      aria-hidden="true"
    >{{ team[1] }}</span>
  </span>
</template>

<style scoped>
.team-badge.has-crest {
  padding: 5px;
  overflow: hidden;
  border-color: color-mix(in srgb, var(--badge-a) 38%, rgba(255, 255, 255, 0.72));
  background: color-mix(in srgb, var(--badge-a) 7%, rgba(255, 255, 255, 0.96));
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.9), 0 10px 24px rgba(0, 0, 0, 0.14);
}

.team-badge.has-crest.small {
  padding: 4px;
}

.team-badge.has-crest.large {
  padding: 9px;
}

.team-badge-crest {
  position: relative;
  z-index: 1;
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
  filter: drop-shadow(0 2px 3px rgba(0, 0, 0, 0.16));
}
</style>
