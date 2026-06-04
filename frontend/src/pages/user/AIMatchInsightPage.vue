<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import MatchInsightCard from '@/components/ai/MatchInsightCard.vue'
import TeamFlag from '@/components/common/TeamFlag.vue'
import { useAIStore } from '@/stores/useAIStore'
import { useMatchStore } from '@/stores/useMatchStore'
import type { Match } from '@/types/match'

const route = useRoute()
const ai = useAIStore()
const matchStore = useMatchStore()
const match = ref<Match | null>(null)

const matchId = Number(route.params.id)

function generate(forceRefresh = false) {
  if (!matchId) return
  ai.generateMatchInsight(matchId, forceRefresh).catch(() => {})
}

onMounted(async () => {
  match.value = await matchStore.fetchMatchDetail(matchId)
})
</script>

<template>
  <div class="ai-match-page">
    <article v-if="match" class="card match-summary">
      <div class="team">
        <TeamFlag :value="match.home_flag" :alt="match.home_team_name" :fallback="match.home_team_code" />
        <strong>{{ match.home_team_name }}</strong>
      </div>
      <span class="vs">VS</span>
      <div class="team away">
        <strong>{{ match.away_team_name }}</strong>
        <TeamFlag :value="match.away_flag" :alt="match.away_team_name" :fallback="match.away_team_code" />
      </div>
    </article>

    <MatchInsightCard
      :insight="ai.currentMatchInsight"
      :loading="ai.matchInsightLoading"
      :error="ai.matchInsightError"
      @generate="generate(false)"
      @refresh="generate(true)"
    />
  </div>
</template>

<style scoped>
.ai-match-page {
  display: grid;
  gap: 14px;
}

.match-summary {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  align-items: center;
  gap: 12px;
  padding: 15px;
}

.team {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 10px;
}

.team.away {
  justify-content: flex-end;
  text-align: right;
}

.team strong {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.vs {
  display: grid;
  place-items: center;
  min-width: 48px;
  height: 36px;
  border-radius: 999px;
  color: var(--primary);
  background: var(--card-soft);
  font-weight: 850;
}
</style>
