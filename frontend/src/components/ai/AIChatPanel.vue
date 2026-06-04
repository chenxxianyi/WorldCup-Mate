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

    <AIInputBox class="chat-input" :loading="loading" @send="emit('send', $event)" />
  </section>
</template>

<style scoped>
.chat-panel {
  display: grid;
  gap: 10px;
}

.messages {
  min-height: 390px;
  display: grid;
  align-content: start;
  gap: 12px;
}

.empty {
  min-height: 210px;
  display: grid;
  place-items: center;
  gap: 8px;
  padding: 18px 8px;
  text-align: center;
  color: var(--muted);
}

.empty .material-symbols-outlined {
  color: var(--weak);
  font-size: 28px;
}

.empty p {
  max-width: 320px;
  margin: 0;
  line-height: 1.6;
}

.assistant-loading {
  width: fit-content;
  padding: 2px 0 2px 42px;
}

.error {
  margin: 0;
  color: var(--blue);
  font-size: 13px;
}

@media (max-width: 767px) {
  .chat-panel {
    padding-bottom: calc(var(--nav-h) + 118px + env(safe-area-inset-bottom));
  }

  .chat-input {
    position: fixed;
    z-index: 18;
    left: 50%;
    bottom: calc(var(--nav-h) + 36px + env(safe-area-inset-bottom));
    width: min(calc(100vw - 48px), 560px);
    transform: translateX(-50%);
  }
}
</style>
