<script setup lang="ts">
import type { AIChatMessage } from '@/types/ai'

defineProps<{
  message: AIChatMessage
}>()

async function copyText(text: string) {
  if (!navigator.clipboard) return
  await navigator.clipboard.writeText(text)
}
</script>

<template>
  <article class="message" :class="message.role">
    <div class="avatar">
      <span class="material-symbols-outlined">
        {{ message.role === 'user' ? 'person' : 'auto_awesome' }}
      </span>
    </div>
    <div class="bubble">
      <p>{{ message.content }}</p>
      <button
        v-if="message.role === 'assistant'"
        class="copy-btn"
        type="button"
        @click="copyText(message.content)"
      >
        复制
      </button>
    </div>
  </article>
</template>

<style scoped>
.message {
  display: flex;
  gap: 10px;
  align-items: flex-start;
}

.message.user {
  flex-direction: row-reverse;
}

.avatar {
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: 999px;
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 10%, transparent);
}

.avatar .material-symbols-outlined {
  font-size: 19px;
}

.bubble {
  max-width: min(78%, 680px);
  padding: 11px 13px;
  border: 1px solid var(--line);
  border-radius: 16px;
  color: var(--text);
  background: var(--card);
  box-shadow: var(--shadow);
}

.message.user .bubble {
  color: #fff;
  border-color: transparent;
  background: var(--primary);
}

.bubble p {
  margin: 0;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
  line-height: 1.65;
  font-size: 14px;
}

.copy-btn {
  min-height: 26px;
  margin-top: 8px;
  padding: 0 9px;
  border: 0;
  border-radius: 999px;
  color: var(--muted);
  background: var(--card-soft);
  font-size: 12px;
}
</style>
