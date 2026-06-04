<script setup lang="ts">
import type { GroupAnalysis } from '@/types/ai'
import AIThinking from './AIThinking.vue'

defineProps<{
  analysis: GroupAnalysis | null
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  (e: 'generate'): void
  (e: 'refresh'): void
}>()
</script>

<template>
  <article class="card group-ai-card">
    <div class="card-head">
      <div>
        <span class="tag blue">AI 出线解读</span>
        <h2>小组形势</h2>
      </div>
      <button class="pill-btn" type="button" :disabled="loading" @click="analysis ? emit('refresh') : emit('generate')">
        {{ analysis ? '重新生成' : '生成解读' }}
      </button>
    </div>

    <div v-if="loading" class="state">
      <AIThinking />
      <span>正在阅读积分榜...</span>
    </div>
    <div v-else-if="error" class="state error">
      <span>{{ error }}</span>
      <button class="pill-btn primary" type="button" @click="emit('generate')">重试</button>
    </div>
    <div v-else-if="!analysis" class="state">
      <span>根据当前积分和剩余赛程生成小组出线解读。</span>
    </div>
    <div v-else class="content">
      <strong class="summary">{{ analysis.summary }}</strong>

      <ul v-if="analysis.key_points?.length">
        <li v-for="item in analysis.key_points" :key="item">{{ item }}</li>
      </ul>

      <div v-if="analysis.teams?.length" class="team-grid">
        <div v-for="team in analysis.teams" :key="team.team_id || team.team_name">
          <span>{{ team.team_name }}</span>
          <strong>{{ team.status || '待观察' }}</strong>
          <p v-if="team.note">{{ team.note }}</p>
        </div>
      </div>

      <p v-if="analysis.qualification_rules" class="rules">{{ analysis.qualification_rules }}</p>
      <p v-if="analysis.data_note" class="note">{{ analysis.data_note }}</p>
    </div>
  </article>
</template>

<style scoped>
.group-ai-card {
  padding: 16px;
}

.card-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.card-head h2 {
  margin: 9px 0 0;
  font-size: 18px;
}

.state {
  min-height: 104px;
  display: grid;
  place-items: center;
  gap: 10px;
  text-align: center;
  color: var(--muted);
}

.state.error {
  color: var(--primary);
}

.content {
  display: grid;
  gap: 12px;
  margin-top: 14px;
}

.summary {
  padding: 12px;
  border-radius: var(--radius-md);
  background: var(--card-soft);
  line-height: 1.5;
}

ul {
  margin: 0;
  padding-left: 18px;
  color: var(--muted);
  line-height: 1.65;
  font-size: 13px;
}

.team-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(128px, 1fr));
  gap: 9px;
}

.team-grid div {
  padding: 10px;
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
}

.team-grid span {
  display: block;
  color: var(--muted);
  font-size: 12px;
}

.team-grid strong {
  display: block;
  margin-top: 4px;
  color: var(--primary);
  font-size: 14px;
}

.team-grid p,
.rules,
.note {
  margin: 7px 0 0;
  color: var(--muted);
  line-height: 1.55;
  font-size: 13px;
}

.note {
  color: var(--weak);
}
</style>
