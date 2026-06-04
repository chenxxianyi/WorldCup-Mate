<script setup lang="ts">
import { ref } from 'vue'
import type { AIChatMessage } from '@/types/ai'

defineProps<{
  message: AIChatMessage
}>()

const copyState = ref<'idle' | 'copied' | 'failed'>('idle')
let copyTimer: ReturnType<typeof window.setTimeout> | undefined

async function copyText(text: string) {
  const value = formatPlainText(text).trim()
  if (!value) return

  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(value)
    } else {
      copyByTextarea(value)
    }
    showCopyState('copied')
  } catch {
    try {
      copyByTextarea(value)
      showCopyState('copied')
    } catch {
      showCopyState('failed')
    }
  }
}

function copyByTextarea(text: string) {
  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.setAttribute('readonly', '')
  textarea.style.position = 'fixed'
  textarea.style.left = '-9999px'
  textarea.style.top = '0'
  document.body.appendChild(textarea)
  textarea.focus()
  textarea.select()

  const ok = document.execCommand('copy')
  document.body.removeChild(textarea)
  if (!ok) throw new Error('Copy failed')
}

function showCopyState(state: 'copied' | 'failed') {
  copyState.value = state
  if (copyTimer) window.clearTimeout(copyTimer)
  copyTimer = window.setTimeout(() => {
    copyState.value = 'idle'
  }, 1400)
}

function formatPlainText(text: string) {
  return text
    .replace(/\*\*([^*\n][\s\S]*?[^*\n])\*\*/g, '$1')
    .replace(/__([^_\n][\s\S]*?[^_\n])__/g, '$1')
    .replace(/\*([^*\n]+)\*/g, '$1')
    .replace(/_([^_\n]+)_/g, '$1')
    .replace(/\*+/g, '')
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
      <p>{{ formatPlainText(message.content) }}</p>
      <button
        v-if="message.role === 'assistant' && message.content.trim()"
        class="copy-btn"
        type="button"
        @click="copyText(message.content)"
      >
        {{ copyState === 'copied' ? '已复制' : copyState === 'failed' ? '复制失败' : '复制' }}
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
  width: 32px;
  height: 32px;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: 999px;
  color: var(--blue);
  background: color-mix(in srgb, var(--blue) 8%, transparent);
}

.avatar .material-symbols-outlined {
  font-size: 19px;
}

.bubble {
  max-width: min(78%, 680px);
  padding: 11px 13px;
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  color: var(--text);
  background: var(--card);
  box-shadow: none;
}

.message.user .bubble {
  color: #fff;
  border-color: transparent;
  background: var(--blue);
}

.bubble p {
  margin: 0;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
  line-height: 1.65;
  font-size: 14px;
}

.copy-btn {
  min-height: 32px;
  margin-top: 8px;
  padding: 0 10px;
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  color: var(--muted);
  background: transparent;
  font-size: 12px;
}

.copy-btn:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--blue) 55%, transparent);
  outline-offset: 2px;
}
</style>
