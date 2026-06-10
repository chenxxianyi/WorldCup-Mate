<script setup lang="ts">
import type { Match } from '@/types/match'
import FeaturedMatchPanel from '@/components/common/FeaturedMatchPanel.vue'
import MatchTicketRow from '@/components/common/MatchTicketRow.vue'

const props = defineProps<{
  match: Match
  featured?: boolean
}>()
</script>

<template>
  <FeaturedMatchPanel v-if="props.featured || props.match.is_featured" :match="props.match" />
  <MatchTicketRow v-else :match="props.match" />
</template>
