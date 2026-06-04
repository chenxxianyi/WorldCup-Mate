<script setup lang="ts">
import type { ShareCopyResult } from '@/types/ai'
import AIThinking from './AIThinking.vue'

defineProps<{
  result: ShareCopyResult | null
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  (e: 'generate'): void
}>()

async function copyText(text?: string) {
  if (!text || !navigator.clipboard) return
  await navigator.clipboard.writeText(text)
}
</script>

<template>
  <article class="card share-card">
    <div class="card-head">
      <div>
        <span class="tag gold">分享文案</span>
        <h2>生成可复制的看球文案</h2>
      </div>
      <button class="pill-btn primary" type="button" :disabled="loading" @click="emit('generate')">
        生成
      </button>
    </div>

    <div v-if="loading" class="state">
      <AIThinking />
      <span>正在组织一段好发出去的话...</span>
    </div>
    <div v-else-if="error" class="state error">{{ error }}</div>
    <div v-else-if="!result" class="state">选择比赛和语气后生成文案。</div>
    <div v-else class="result">
      <strong v-if="result.title">{{ result.title }}</strong>
      <p>{{ result.content }}</p>
      <div class="actions">
        <button class="pill-btn" type="button" @click="copyText(result.content)">复制文案</button>
      </div>
      <ul v-if="result.tips?.length">
        <li v-for="item in result.tips" :key="item">{{ item }}</li>
      </ul>
    </div>
  </article>
</template>

<style scoped>
.share-card {
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
  min-height: 132px;
  display: grid;
  place-items: center;
  gap: 10px;
  text-align: center;
  color: var(--muted);
}

.state.error {
  color: var(--primary);
}

.result {
  display: grid;
  gap: 12px;
  margin-top: 14px;
}

.result strong {
  font-size: 16px;
}

.result p {
  margin: 0;
  padding: 13px;
  border-radius: var(--radius-md);
  white-space: pre-wrap;
  overflow-wrap: anywhere;
  color: var(--text);
  background: var(--card-soft);
  line-height: 1.7;
  font-size: 14px;
}

.actions {
  display: flex;
  justify-content: flex-end;
}

ul {
  margin: 0;
  padding-left: 18px;
  color: var(--muted);
  line-height: 1.6;
  font-size: 13px;
}
</style>
