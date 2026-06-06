<script setup lang="ts">
import type { PostMatchSummary } from '@/types/postMatchSummary'
import AIThinking from './AIThinking.vue'

const props = defineProps<{
  summary: PostMatchSummary | null
  loading?: boolean
  error?: string
  canGenerate?: boolean
  spoilerSafe?: boolean
}>()

const emit = defineEmits<{
  (e: 'generate'): void
  (e: 'refresh'): void
}>()
</script>

<template>
  <article class="card post-summary-card">
    <div class="summary-head">
      <div>
        <span class="tag blue">AI</span>
        <h2>赛后补看摘要</h2>
      </div>
      <button
        v-if="summary"
        class="pill-btn"
        type="button"
        :disabled="loading"
        @click="emit('refresh')"
      >
        刷新
      </button>
    </div>

    <div v-if="loading" class="summary-state">
      <AIThinking />
      <span>正在生成赛后摘要...</span>
    </div>

    <div v-else-if="error" class="summary-state error">
      <span>{{ error }}</span>
      <button class="pill-btn primary" type="button" @click="emit('generate')">重试</button>
    </div>

    <div v-else-if="!summary" class="summary-empty">
      <span class="material-symbols-outlined">description</span>
      <p>比赛结束后可生成补看摘要。</p>
      <button
        class="pill-btn primary"
        type="button"
        :disabled="!canGenerate"
        @click="emit('generate')"
      >
        生成赛后补看摘要
      </button>
    </div>

    <div v-else class="summary-body">
      <div class="score-line">{{ summary.score_line }}</div>

      <div v-if="summary.worth_watching" class="worth-tag">
        <span class="material-symbols-outlined">visibility</span>
        {{ summary.worth_watching }}
      </div>

      <p class="summary-text">{{ summary.summary }}</p>

      <section v-if="summary.key_takeaways?.length" class="takeaways">
        <h3>关键要点</h3>
        <ul>
          <li v-for="item in summary.key_takeaways" :key="item">{{ item }}</li>
        </ul>
      </section>

      <div v-if="summary.qualification_impact" class="note-box">
        <span>出线影响</span>
        <p>{{ summary.qualification_impact }}</p>
      </div>

      <div v-if="summary.data_note" class="data-note">
        <small>{{ summary.data_note }}</small>
      </div>
    </div>
  </article>
</template>

<style scoped>
.post-summary-card {
  padding: 16px;
}

.summary-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.summary-head h2 {
  margin: 9px 0 0;
  font-size: 18px;
}

/* ── States ── */
.summary-state,
.summary-empty {
  min-height: 128px;
  display: grid;
  place-items: center;
  gap: 10px;
  text-align: center;
  color: var(--muted);
}

.summary-state.error {
  color: var(--primary);
}

.summary-empty .material-symbols-outlined {
  color: var(--primary);
  font-size: 32px;
}

.summary-empty p {
  margin: 0;
  max-width: 330px;
  line-height: 1.6;
}

/* ── Content ── */
.summary-body {
  display: grid;
  gap: 14px;
  margin-top: 16px;
}

.score-line {
  font-size: 20px;
  font-weight: 850;
  text-align: center;
  padding: 12px;
  border-radius: var(--radius-lg);
  background: color-mix(in srgb, var(--primary) 9%, transparent);
  color: var(--primary);
}

.worth-tag {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 13px;
  border-radius: var(--radius-md);
  background: var(--card-soft);
  font-size: 13px;
  font-weight: 700;
  color: var(--hot);
}

.worth-tag .material-symbols-outlined {
  font-size: 20px;
}

.summary-text {
  margin: 0;
  line-height: 1.7;
  font-size: 14px;
  color: var(--text);
}

.takeaways h3 {
  margin: 0 0 8px;
  font-size: 14px;
}

.takeaways ul {
  margin: 0;
  padding-left: 18px;
  color: var(--muted);
  line-height: 1.7;
  font-size: 13px;
}

.note-box {
  padding: 12px;
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
}

.note-box span {
  color: var(--primary);
  font-size: 12px;
  font-weight: 750;
}

.note-box p {
  margin: 7px 0 0;
  color: var(--muted);
  line-height: 1.55;
  font-size: 13px;
}

.data-note {
  color: var(--weak);
  font-size: 12px;
  text-align: center;
  line-height: 1.5;
}
</style>
