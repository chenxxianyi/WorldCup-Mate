<script setup lang="ts">
import type { AIChatMessage } from '@/types/ai'
import AIInputBox from './AIInputBox.vue'
import AIMessageBubble from './AIMessageBubble.vue'
import AIThinking from './AIThinking.vue'

defineProps<{
  messages: AIChatMessage[]
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  (e: 'send', value: string): void
}>()
</script>

<template>
  <section class="chat-panel">
    <div class="messages">
      <div v-if="!messages.length" class="empty">
        <span class="material-symbols-outlined">forum</span>
        <p>可以问赛程、规则、小组形势，或者让它帮你选一场比赛。</p>
      </div>
      <AIMessageBubble v-for="(message, index) in messages" :key="message.id || index" :message="message" />
      <div v-if="loading" class="assistant-loading">
        <AIThinking />
      </div>
      <p v-if="error" class="error">{{ error }}</p>
    </div>

    <AIInputBox :loading="loading" @send="emit('send', $event)" />
  </section>
</template>

<style scoped>
.chat-panel {
  display: grid;
  gap: 12px;
}

.messages {
  min-height: 420px;
  display: grid;
  align-content: start;
  gap: 12px;
}

.empty {
  min-height: 260px;
  display: grid;
  place-items: center;
  gap: 10px;
  padding: 24px;
  border: 1px dashed var(--line);
  border-radius: var(--radius-xl);
  text-align: center;
  color: var(--muted);
}

.empty .material-symbols-outlined {
  color: var(--primary);
  font-size: 34px;
}

.empty p {
  max-width: 320px;
  margin: 0;
  line-height: 1.6;
}

.assistant-loading {
  width: fit-content;
  padding: 12px 14px;
  border: 1px solid var(--line);
  border-radius: 16px;
  background: var(--card);
}

.error {
  margin: 0;
  color: var(--primary);
  font-size: 13px;
}
</style>
