<script setup lang="ts">
import type { MatchInsight } from '@/types/ai'
import AIThinking from './AIThinking.vue'

const props = defineProps<{
  insight: MatchInsight | null
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  (e: 'generate'): void
  (e: 'refresh'): void
}>()

async function copySummary() {
  if (!props.insight || !navigator.clipboard) return
  await navigator.clipboard.writeText(props.insight.summary)
}
</script>

<template>
  <article class="ai-card">
    <div class="ai-head">
      <div>
        <span class="tag blue">AI 看点</span>
        <h2>这场值不值得看</h2>
      </div>
      <button
        class="pill-btn"
        type="button"
        :disabled="loading"
        @click="insight ? emit('refresh') : emit('generate')"
      >
        {{ insight ? '重新生成' : '生成看点' }}
      </button>
    </div>

    <div v-if="loading" class="state">
      <AIThinking />
      <span>正在整理比赛信息...</span>
    </div>

    <div v-else-if="error" class="state error">
      <span>{{ error }}</span>
      <button class="pill-btn primary" type="button" @click="emit('generate')">重试</button>
    </div>

    <div v-else-if="!insight" class="empty">
      <span class="material-symbols-outlined">auto_awesome</span>
      <p>基于赛程、球队和积分数据生成一份简短看球建议。</p>
    </div>

    <div v-else class="content">
      <div class="summary-row">
        <div>
          <strong>{{ insight.summary }}</strong>
          <span>看点指数 {{ insight.watch_rating }}/5</span>
        </div>
        <button class="icon-action" type="button" title="复制摘要" @click="copySummary">
          <span class="material-symbols-outlined">content_copy</span>
        </button>
      </div>

      <div class="rating">
        <i
          v-for="n in 5"
          :key="n"
          :class="{ active: n <= Math.round(insight.watch_rating || 0) }"
        ></i>
      </div>

      <section v-if="insight.reasons?.length">
        <h3>推荐理由</h3>
        <ul>
          <li v-for="item in insight.reasons" :key="item">{{ item }}</li>
        </ul>
      </section>

      <section v-if="insight.team_comparison?.length">
        <h3>双方对比</h3>
        <ul>
          <li v-for="item in insight.team_comparison" :key="item">{{ item }}</li>
        </ul>
      </section>

      <section v-if="insight.beginner_tips?.length">
        <h3>小白看点</h3>
        <ul>
          <li v-for="item in insight.beginner_tips" :key="item">{{ item }}</li>
        </ul>
      </section>

      <div class="note-grid">
        <div v-if="insight.qualification_impact">
          <span>出线影响</span>
          <p>{{ insight.qualification_impact }}</p>
        </div>
        <div v-if="insight.should_stay_up">
          <span>熬夜建议</span>
          <p>{{ insight.should_stay_up }}</p>
        </div>
      </div>

      <div v-if="insight.suitable_for?.length" class="chips">
        <span v-for="item in insight.suitable_for" :key="item">{{ item }}</span>
      </div>
    </div>
  </article>
</template>

<style scoped>
.ai-card {
  display: grid;
  gap: 16px;
  padding-top: 20px;
  border-top: 1px solid color-mix(in srgb, var(--line) 78%, transparent);
}

.ai-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.ai-head h2 {
  margin: 9px 0 0;
  font-size: 18px;
}

.state,
.empty {
  min-height: 128px;
  display: grid;
  place-items: center;
  gap: 10px;
  text-align: center;
  color: var(--muted);
}

.state.error {
  color: var(--primary);
}

.empty .material-symbols-outlined {
  color: var(--primary);
  font-size: 32px;
}

.empty p {
  margin: 0;
  max-width: 330px;
  line-height: 1.6;
}

.content {
  display: grid;
  gap: 14px;
}

.summary-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 0;
  border-top: 1px solid color-mix(in srgb, var(--line) 70%, transparent);
  border-bottom: 1px solid color-mix(in srgb, var(--line) 70%, transparent);
  background: transparent;
}

.summary-row strong {
  display: block;
  font-size: 16px;
  line-height: 1.45;
}

.summary-row span {
  display: block;
  margin-top: 6px;
  color: var(--muted);
  font-size: 12px;
}

.rating {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 6px;
}

.rating i {
  height: 7px;
  border-radius: 999px;
  background: var(--card-soft);
}

.rating i.active {
  background: linear-gradient(90deg, var(--primary), var(--gold));
}

section h3 {
  margin: 0 0 8px;
  font-size: 14px;
}

ul {
  margin: 0;
  padding-left: 18px;
  color: var(--muted);
  line-height: 1.7;
  font-size: 13px;
}

.note-grid {
  display: grid;
  gap: 10px;
}

.note-grid div {
  padding: 12px 0;
  border-top: 1px solid color-mix(in srgb, var(--line) 72%, transparent);
}

.note-grid span {
  color: var(--primary);
  font-size: 12px;
  font-weight: 750;
}

.note-grid p {
  margin: 7px 0 0;
  color: var(--muted);
  line-height: 1.55;
  font-size: 13px;
}

.chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.chips span {
  padding: 6px 9px;
  border-radius: 999px;
  color: var(--muted);
  background: var(--card-soft);
  font-size: 12px;
  font-weight: 700;
}
</style>
